package integration

import (
	"context"
	"errors"
	"testing"

	infra "github.com/britinogn/ctemzjournal/internal/cloudinary"
	"github.com/britinogn/ctemzjournal/internal/repository"
	"github.com/britinogn/ctemzjournal/internal/service"
	"github.com/britinogn/ctemzjournal/test/fixtures"
	"github.com/britinogn/ctemzjournal/test/helpers"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// fakeStorage stubs Cloudinary: no network, deterministic URLs.
type fakeStorage struct {
	destroyed []string
}

func (f *fakeStorage) SignUpload(userID, tradeID string) (infra.Signature, error) {
	return infra.Signature{
		UploadURL: "https://fake.cloudinary/upload",
		CloudName: "fake", APIKey: "fake",
		Timestamp: "123", Signature: "sig",
		Folder: "trade-images/" + userID + "/" + tradeID, Type: "authenticated",
	}, nil
}

func (f *fakeStorage) DeliveryURL(publicID string) (string, error) {
	return "https://fake.cloudinary/signed/" + publicID, nil
}

func (f *fakeStorage) Destroy(_ context.Context, publicID string) error {
	f.destroyed = append(f.destroyed, publicID)
	return nil
}

type imagesEnv struct {
	svc     *service.Images
	trades  *service.Trades
	fake    *fakeStorage
	userA   uuid.UUID
	userB   uuid.UUID
	account uuid.UUID
}

func newImagesEnv(t *testing.T) (*imagesEnv, context.Context) {
	t.Helper()
	database := helpers.MustConnect(t)
	helpers.Truncate(t, database)
	userA := fixtures.MustUser(t, database, "alice")
	userB := fixtures.MustUser(t, database, "bob")
	accounts := service.NewAccounts(repository.NewAccounts(database.Queries))
	ctx := context.Background()
	account, err := accounts.Create(ctx, userA, service.AccountCreate{Name: "Main", Type: "demo"})
	if err != nil {
		t.Fatalf("seed account: %v", err)
	}
	tradesRepo := repository.NewTrades(database.Queries)
	fake := &fakeStorage{}
	imagesSvc := service.NewImages(tradesRepo, repository.NewTradeImages(database.Queries), fake)
	tradesSvc := service.NewTrades(
		tradesRepo, repository.NewTradeTags(database.Queries),
		repository.NewAccounts(database.Queries), repository.NewSetups(database.Queries),
		repository.NewTags(database.Queries), database.Queries,
	)
	return &imagesEnv{svc: imagesSvc, trades: tradesSvc, fake: fake, userA: userA, userB: userB, account: account.ID}, ctx
}

func newOpenTrade(t *testing.T, env *imagesEnv, ctx context.Context) uuid.UUID {
	t.Helper()
	entry, sl, lots := 1.0850, 1.0800, 1.0
	trade, err := env.trades.Create(ctx, env.userA, service.TradeCreate{
		AccountID: env.account, Pair: "EUR/USD", Direction: "long",
		Entry: &entry, StopLoss: &sl, LotSize: &lots, Status: "open",
	})
	if err != nil {
		t.Fatalf("seed trade: %v", err)
	}
	return trade.ID
}

func TestImagesSignatureConfirmListDelete(t *testing.T) {
	env, ctx := newImagesEnv(t)
	tradeID := newOpenTrade(t, env, ctx)

	sig, err := env.svc.RequestSignature(ctx, env.userA, tradeID)
	if err != nil {
		t.Fatalf("signature: %v", err)
	}
	if sig.Type != "authenticated" || sig.Folder == "" {
		t.Fatalf("bad signature payload: %+v", sig)
	}
	if _, err := env.svc.RequestSignature(ctx, env.userB, tradeID); !errors.Is(err, pgx.ErrNoRows) {
		t.Fatalf("user B got signature for user A trade: %v", err)
	}

	kind := "entry"
	img, err := env.svc.Confirm(ctx, env.userA, tradeID, "trade-images/a/b/img1", &kind)
	if err != nil {
		t.Fatalf("confirm: %v", err)
	}
	if img.Position != 0 {
		t.Fatalf("first image position = %d, want 0", img.Position)
	}
	if _, err := env.svc.Confirm(ctx, env.userA, tradeID, "", &kind); err == nil {
		t.Fatalf("expected empty public_id error")
	}
	bad := "bogus"
	if _, err := env.svc.Confirm(ctx, env.userA, tradeID, "x", &bad); err == nil {
		t.Fatalf("expected kind validation error")
	}

	list, err := env.svc.List(ctx, env.userA, tradeID)
	if err != nil || len(list) != 1 {
		t.Fatalf("list: %v (n=%d)", err, len(list))
	}
	if list[0].URL == "" || list[0].URL == list[0].PublicID {
		t.Fatalf("list must return signed URLs, got %+v", list[0])
	}
	if _, err := env.svc.List(ctx, env.userB, tradeID); !errors.Is(err, pgx.ErrNoRows) {
		t.Fatalf("user B listed user A images: %v", err)
	}

	if err := env.svc.Delete(ctx, env.userA, tradeID, img.ID); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if len(env.fake.destroyed) != 1 || env.fake.destroyed[0] != "trade-images/a/b/img1" {
		t.Fatalf("destroy not called: %v", env.fake.destroyed)
	}
	list, err = env.svc.List(ctx, env.userA, tradeID)
	if err != nil || len(list) != 0 {
		t.Fatalf("after delete: %v (n=%d)", err, len(list))
	}
}

func TestImagesFiveCap(t *testing.T) {
	env, ctx := newImagesEnv(t)
	tradeID := newOpenTrade(t, env, ctx)

	for i := 0; i < infra.MaxImagesPerTrade; i++ {
		if _, err := env.svc.RequestSignature(ctx, env.userA, tradeID); err != nil {
			t.Fatalf("signature %d: %v", i, err)
		}
		if _, err := env.svc.Confirm(ctx, env.userA, tradeID, "img", nil); err != nil {
			t.Fatalf("confirm %d: %v", i, err)
		}
	}
	if _, err := env.svc.RequestSignature(ctx, env.userA, tradeID); err == nil {
		t.Fatalf("expected 5-image cap on signature")
	}
	if _, err := env.svc.Confirm(ctx, env.userA, tradeID, "img6", nil); err == nil {
		t.Fatalf("expected 5-image cap on confirm")
	}
}

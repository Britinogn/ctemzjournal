package service

import (
	"context"
	"errors"

	infra "github.com/britinogn/ctemzjournal/internal/cloudinary"
	sqlc "github.com/britinogn/ctemzjournal/internal/db/sqlc"
	"github.com/britinogn/ctemzjournal/internal/model"
	"github.com/britinogn/ctemzjournal/internal/repository"
	"github.com/google/uuid"
)

// Storage abstracts Cloudinary so tests can stub it.
// The production implementation is *cloudinary.Client.
type Storage interface {
	SignUpload(userID, tradeID string) (infra.Signature, error)
	DeliveryURL(publicID string) (string, error)
	Destroy(ctx context.Context, publicID string) error
}

// Images owns the upload flow: signature -> direct browser upload ->
// confirm, plus signed-URL reads and destroy-on-delete.
type Images struct {
	trades  *repository.Trades
	images  *repository.TradeImages
	storage Storage
}

func NewImages(trades *repository.Trades, images *repository.TradeImages, storage Storage) *Images {
	return &Images{trades: trades, images: images, storage: storage}
}

// RequestSignature checks ownership + the 5-image cap, then returns the
// signed upload payload for trade-images/{userID}/{tradeID}.
func (s *Images) RequestSignature(ctx context.Context, userID, tradeID uuid.UUID) (infra.Signature, error) {
	if _, err := s.trades.Get(ctx, tradeID, userID); err != nil {
		return infra.Signature{}, err
	}
	n, err := s.images.CountByTrade(ctx, tradeID)
	if err != nil {
		return infra.Signature{}, err
	}
	if n >= infra.MaxImagesPerTrade {
		return infra.Signature{}, errors.New("trade already has 5 images")
	}
	return s.storage.SignUpload(userID.String(), tradeID.String())
}

// Confirm saves the Cloudinary public_id after the browser upload.
func (s *Images) Confirm(ctx context.Context, userID, tradeID uuid.UUID, publicID string, kind *string) (sqlc.TradeImage, error) {
	if publicID == "" {
		return sqlc.TradeImage{}, errors.New("public_id is required")
	}
	if _, err := s.trades.Get(ctx, tradeID, userID); err != nil {
		return sqlc.TradeImage{}, err
	}
	if kind != nil && *kind != model.ImageEntry && *kind != model.ImageExit {
		return sqlc.TradeImage{}, errors.New("kind must be entry or exit")
	}
	n, err := s.images.CountByTrade(ctx, tradeID)
	if err != nil {
		return sqlc.TradeImage{}, err
	}
	if n >= infra.MaxImagesPerTrade {
		return sqlc.TradeImage{}, errors.New("trade already has 5 images")
	}
	return s.images.Create(ctx, tradeID, publicID, kind, int(n))
}

// ImageWithURL pairs a stored image with its short-lived signed URL.
type ImageWithURL struct {
	sqlc.TradeImage
	URL string `json:"url"`
}

// List returns the trade's images with signed delivery URLs (ownership-checked).
func (s *Images) List(ctx context.Context, userID, tradeID uuid.UUID) ([]ImageWithURL, error) {
	if _, err := s.trades.Get(ctx, tradeID, userID); err != nil {
		return nil, err
	}
	images, err := s.images.ListByTrade(ctx, tradeID)
	if err != nil {
		return nil, err
	}
	out := make([]ImageWithURL, 0, len(images))
	for _, img := range images {
		u, err := s.storage.DeliveryURL(img.PublicID)
		if err != nil {
			return nil, err
		}
		out = append(out, ImageWithURL{TradeImage: img, URL: u})
	}
	return out, nil
}

// Delete destroys the Cloudinary file (best-effort) and the row.
func (s *Images) Delete(ctx context.Context, userID, tradeID, imageID uuid.UUID) error {
	if _, err := s.trades.Get(ctx, tradeID, userID); err != nil {
		return err
	}
	img, err := s.images.Get(ctx, imageID)
	if err != nil {
		return err
	}
	if img.TradeID != tradeID {
		return errors.New("image does not belong to trade")
	}
	_ = s.storage.Destroy(ctx, img.PublicID) // best-effort: row delete is authoritative
	return s.images.Delete(ctx, imageID)
}

// DeleteByTrade destroys every Cloudinary file of a trade (best-effort).
// Call before deleting the trade row itself (FK cascade removes the rows).
func (s *Images) DeleteByTrade(ctx context.Context, userID, tradeID uuid.UUID) error {
	if _, err := s.trades.Get(ctx, tradeID, userID); err != nil {
		return err
	}
	images, err := s.images.ListByTrade(ctx, tradeID)
	if err != nil {
		return err
	}
	for _, img := range images {
		_ = s.storage.Destroy(ctx, img.PublicID)
	}
	return nil
}

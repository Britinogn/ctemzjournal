package service

import (
	"context"
	"fmt"
	"strings"

	sqlc "github.com/britinogn/ctemzjournal/internal/db/sqlc"
	"github.com/britinogn/ctemzjournal/internal/repository"
	"github.com/google/uuid"
)

// Me serves GET /me and PATCH /me.
// The frontend calls GET /me after login to redirect:
// admin -> /admin, user -> /dashboard.
type Me struct {
	profiles *repository.Profiles
}

func NewMe(profiles *repository.Profiles) *Me { return &Me{profiles: profiles} }

func (s *Me) Get(ctx context.Context, userID uuid.UUID) (sqlc.Profile, error) {
	return s.profiles.GetByID(ctx, userID)
}

type MeUpdate struct {
	DisplayName *string
	Timezone    *string
	AvatarPath  *string
}

func (s *Me) Update(ctx context.Context, userID uuid.UUID, in MeUpdate) (sqlc.Profile, error) {
	if in.Timezone != nil && !validateTimezone(strings.TrimSpace(*in.Timezone)) {
		return sqlc.Profile{}, fmt.Errorf("invalid timezone")
	}
	if in.AvatarPath != nil {
		v := strings.TrimSpace(*in.AvatarPath)
		if !validAvatarPath(v) {
			return sqlc.Profile{}, fmt.Errorf("invalid avatar")
		}
		trimmed := v
		in.AvatarPath = &trimmed
	}
	return s.profiles.Update(ctx, userID, in.DisplayName, in.Timezone, in.AvatarPath)
}

// validAvatarPath accepts '' (remove), legacy http(s) URLs, or a relative
// Supabase storage path like `{user_id}/avatar.webp`.
func validAvatarPath(v string) bool {
	if v == "" {
		return true
	}
	if strings.HasPrefix(v, "http://") || strings.HasPrefix(v, "https://") {
		return len(v) <= 2048
	}
	if len(v) > 512 || strings.Contains(v, "..") || strings.Contains(v, "\\") {
		return false
	}
	for _, r := range v {
		if !(r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' || r == '/' || r == '.' || r == '_' || r == '-') {
			return false
		}
	}
	lower := strings.ToLower(v)
	return strings.HasSuffix(lower, ".webp") || strings.HasSuffix(lower, ".png") || strings.HasSuffix(lower, ".jpg") || strings.HasSuffix(lower, ".jpeg")
}

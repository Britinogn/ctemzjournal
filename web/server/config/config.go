package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds everything the server needs.
// Values come from the environment (.env.sample documents every key).
// Removed Supabase vars stay removed: auth is verified via SUPABASE_JWKS_URL only.
type Config struct {
	AppEnv         string
	Port           string
	AppURL         string
	FrontendOrigin string

	DatabaseURL string

	SupabaseURL     string
	SupabaseJWKSURL string

	CloudinaryCloudName string
	CloudinaryAPIKey    string
	CloudinaryAPISecret string

	TwelveDataAPIKey    string
	RatesRefreshMinutes int

	BcryptRounds int

	JWTSecret  string
	JWTExpires string
	JWTReset   string

	CORSURL string
}

func Load() *Config {
	_ = godotenv.Load()

	return &Config{
		AppEnv:         getEnv("APP_ENV", "development"),
		Port:           getEnv("PORT", "8080"),
		AppURL:         getEnv("APP_URL", "http://localhost:3000"),
		FrontendOrigin: getEnv("FRONTEND_ORIGIN", "http://localhost:3000"),

		DatabaseURL: getEnv("DATABASE_URL", ""),

		SupabaseURL:     getEnv("SUPABASE_URL", ""),
		SupabaseJWKSURL: getEnv("SUPABASE_JWKS_URL", ""),

		CloudinaryCloudName: getEnv("CLOUDINARY_CLOUD_NAME", ""),
		CloudinaryAPIKey:    getEnv("CLOUDINARY_API_KEY", ""),
		CloudinaryAPISecret: getEnv("CLOUDINARY_API_SECRET", ""),

		TwelveDataAPIKey:    getEnv("TWELVE_DATA_API_KEY", ""),
		RatesRefreshMinutes: getEnvInt("RATES_REFRESH_MINUTES", 45),

		BcryptRounds: getEnvInt("BCRYPT_ROUNDS", 12),

		JWTSecret:  getEnv("JWT_SECRET", "change_me"),
		JWTExpires: getEnv("JWT_EXPIRES_IN", "1h"),
		JWTReset:   getEnv("JWT_RESET_AT", "7d"),

		CORSURL: getEnv("CORS_URL", "http://localhost:3000"),
	}
}

// AllowedOrigins dedupes the three origin vars (.env.sample has APP_URL,
// FRONTEND_ORIGIN and CORS_URL all pointing at the frontend).
func (c *Config) AllowedOrigins() []string {
	seen := map[string]struct{}{}
	out := []string{}
	for _, o := range []string{c.FrontendOrigin, c.CORSURL, c.AppURL} {
		if o == "" {
			continue
		}
		if _, ok := seen[o]; !ok {
			seen[o] = struct{}{}
			out = append(out, o)
		}
	}
	return out
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return fallback
	}
	return n
}

package service

import "time"

// DefaultTimezone is the fallback when the browser sends nothing usable.
const DefaultTimezone = "Africa/Lagos"

// normalizeTimezone keeps valid IANA names, falling back to Lagos.
// Used on signup/sync so a spoofed value never breaks profile creation.
func normalizeTimezone(tz string) string {
	if tz == "" {
		return DefaultTimezone
	}
	if len(tz) > 64 {
		return DefaultTimezone
	}
	if _, err := time.LoadLocation(tz); err != nil {
		return DefaultTimezone
	}
	return tz
}

// validateTimezone reports whether tz is a loadable IANA name.
func validateTimezone(tz string) bool {
	if tz == "" || len(tz) > 64 {
		return false
	}
	_, err := time.LoadLocation(tz)
	return err == nil
}

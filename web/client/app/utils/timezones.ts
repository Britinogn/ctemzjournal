/** Shared timezone helpers. DB default is Africa/Lagos; signup detects the browser zone. */

export const DEFAULT_TIMEZONE = 'Africa/Lagos'

export const TIMEZONES = [
  'Africa/Lagos',
  'UTC',
  'Europe/London',
  'Europe/Berlin',
  'America/New_York',
  'America/Chicago',
  'Asia/Dubai',
  'Asia/Tokyo',
  'Australia/Sydney',
]

/** Browser IANA zone (e.g. "America/Toronto"), falling back to Lagos. */
export function detectBrowserTimezone(): string {
  try {
    const tz = Intl.DateTimeFormat().resolvedOptions().timeZone
    if (typeof tz === 'string' && tz.length > 0 && tz.length <= 64)
      return tz
  }
  catch {
    // Intl unavailable — fall through to default.
  }
  return DEFAULT_TIMEZONE
}

/** Options for the settings dropdown — always includes the current value. */
export function timezoneOptions(current: string): string[] {
  if (TIMEZONES.includes(current))
    return TIMEZONES
  return [current, ...TIMEZONES]
}

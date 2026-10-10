/** Small display formatters. Money is quote-currency units unless noted. */

export function fmtSigned(n: number, digits = 2): string {
  const v = n.toFixed(digits);
  return n > 0 ? `+${v}` : v;
}

export function fmtMoney(n: number, currency = '$', digits = 0): string {
  return `${n < 0 ? '-' : '+'}${currency}${Math.abs(n).toLocaleString('en-US', {
    minimumFractionDigits: digits,
    maximumFractionDigits: digits,
  })}`;
}

export function fmtPct(n: number, digits = 1): string {
  return `${(n * 100).toFixed(digits)}%`;
}

export function fmtR(n: number, digits = 2): string {
  return `${n > 0 ? '+' : ''}${n.toFixed(digits)}R`;
}

/** "2 Oct" from an ISO date string. */
export function fmtDay(iso: string): string {
  const d = new Date(iso);
  return d.toLocaleDateString('en-GB', { day: 'numeric', month: 'short', timeZone: 'UTC' });
}

/** Fall back to UTC when a stored timezone is missing or invalid. */
export function safeTimezone(tz: string | null | undefined): string {
  if (!tz)
    return 'UTC';
  try {
    new Intl.DateTimeFormat('en-GB', { timeZone: tz });
    return tz;
  }
  catch {
    return 'UTC';
  }
}

/** "2 Oct" from an ISO instant, rendered in the given IANA timezone. */
export function fmtDayIn(iso: string, tz: string): string {
  const d = new Date(iso);
  return d.toLocaleDateString('en-GB', { day: 'numeric', month: 'short', timeZone: safeTimezone(tz) });
}

/** "2 Oct, 14:30" from an ISO instant, rendered in the given IANA timezone. */
export function fmtDateTimeIn(iso: string, tz: string): string {
  const d = new Date(iso);
  return d.toLocaleString('en-GB', {
    day: 'numeric', month: 'short', hour: '2-digit', minute: '2-digit', hour12: false,
    timeZone: safeTimezone(tz),
  });
}

/** Current hour (0-23) in the given IANA timezone — for greetings. */
export function hourInTimezone(tz: string, at: Date = new Date()): number {
  const zone = safeTimezone(tz);
  const hour = new Intl.DateTimeFormat('en-GB', { hour: 'numeric', hour12: false, timeZone: zone }).format(at);
  const n = Number.parseInt(hour, 10);
  return Number.isFinite(n) ? n % 24 : at.getHours();
}

/** ISO instant → `datetime-local` wall string in the given timezone. */
export function toDatetimeLocalIn(iso: string | null, tz: string): string {
  if (!iso)
    return '';
  const d = new Date(iso);
  if (Number.isNaN(d.getTime()))
    return '';
  const zone = safeTimezone(tz);
  const parts = new Intl.DateTimeFormat('en-CA', {
    year: 'numeric', month: '2-digit', day: '2-digit',
    hour: '2-digit', minute: '2-digit', hour12: false,
    timeZone: zone,
  }).formatToParts(d);
  const get = (t: string): string => parts.find(p => p.type === t)?.value ?? '';
  return `${get('year')}-${get('month')}-${get('day')}T${get('hour')}:${get('minute')}`;
}

/**
 * `datetime-local` wall string (as typed) in the given timezone → UTC ISO.
 * When the zone equals the browser zone this matches `new Date(v).toISOString()`.
 */
export function fromDatetimeLocalToISO(local: string, tz: string): string {
  if (!local)
    return '';
  const zone = safeTimezone(tz);
  const m = /^(\d{4})-(\d{2})-(\d{2})T(\d{2}):(\d{2})/.exec(local);
  if (!m)
    return new Date(local).toISOString();
  const target = Date.UTC(Number(m[1]), Number(m[2]) - 1, Number(m[3]), Number(m[4]), Number(m[5]));
  const fmt = new Intl.DateTimeFormat('en-CA', {
    year: 'numeric', month: '2-digit', day: '2-digit',
    hour: '2-digit', minute: '2-digit', hour12: false,
    timeZone: zone,
  });
  const wallMs = (instant: number): number | null => {
    const s = fmt.format(new Date(instant));
    const wm = /^(\d{4})-(\d{2})-(\d{2}), (\d{2}):(\d{2})/.exec(s)
      ?? /^(\d{4})-(\d{2})-(\d{2})T(\d{2}):(\d{2})/.exec(s);
    if (!wm)
      return null;
    return Date.UTC(Number(wm[1]), Number(wm[2]) - 1, Number(wm[3]), Number(wm[4]), Number(wm[5]));
  };
  // Fixed-point iteration: shift the guess by (target wall − actual wall).
  let guess = target;
  for (let i = 0; i < 3; i++) {
    const actual = wallMs(guess);
    if (actual === null)
      break;
    const diff = target - actual;
    guess += diff;
    if (diff === 0)
      break;
  }
  return new Date(guess).toISOString();
}

/** "AO" from "Ada Obi" (falls back to email prefix, then "?"). */
export function initials(name?: string | null, email?: string): string {
  const base = (name || '').trim() || (email || '').split('@')[0] || '';
  const parts = base.split(/[\s._-]+/).filter(Boolean);
  if (parts.length === 0)
    return '?';
  if (parts.length === 1)
    return parts[0].slice(0, 2).toUpperCase();
  return (parts[0][0] + parts[1][0]).toUpperCase();
}

/** "X min ago" / "2 h ago" / "3 d ago" for rates + journal freshness. */
export function timeAgo(iso: string | null | undefined): string {
  if (!iso)
    return 'never';
  const mins = Math.max(0, Math.round((Date.now() - new Date(iso).getTime()) / 60000));
  if (mins < 1)
    return 'just now';
  if (mins < 60)
    return `${mins} min ago`;
  const hours = Math.round(mins / 60);
  if (hours < 24)
    return `${hours} h ago`;
  return `${Math.round(hours / 24)} d ago`;
}

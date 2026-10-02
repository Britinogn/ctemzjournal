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

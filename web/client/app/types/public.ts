/**
 * Public home-page shapes. List journals carry ONLY safe fields —
 * no money amounts or lot sizes ever arrive here. Notes + all images
 * arrive only on GET /public/journals/{id} (public by owner choice).
 */

export interface PublicSettings {
  site_name: string;
  tagline: string | null;
  /** Full ready-to-use asset URLs (backend composes them). */
  logo_url: string | null;
  favicon_url: string | null;
  contact_email: string | null;
  footer_text: string | null;
  risk_disclaimer: string | null;
  social_links: Record<string, unknown>;
  allow_signups: boolean;
  maintenance_mode: boolean;
}

export type JournalResult = 'win' | 'loss';

export interface PublicJournal {
  id: string;
  display_name: string | null;
  pair: string;
  direction: string;
  timeframe: string | null;
  setup_name: string | null;
  /** Null for still-open public trades. */
  result: JournalResult | null;
  r_multiple: number | null;
  /** First image, signed URL. Null when the trade has no images. */
  image_url: string | null;
  date: string;
}

export interface PublicJournalImage {
  url: string;
  kind?: string | null;
}

/** GET /public/journals/{id} — public fields plus notes + every image. */
export interface PublicJournalDetail extends PublicJournal {
  notes: string | null;
  images: PublicJournalImage[] | null;
}

export interface Rate {
  pair: string;
  price: number;
  updated_at: string;
  stale: boolean;
}

/** GET /public/rates — cached prices; `stale` drives the warning badge. */
export interface RatesResponse {
  rates: Rate[];
  updated_at: string | null;
  stale: boolean;
}

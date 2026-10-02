/**
 * Trade screenshots. Three steps, all scoped to the owning trade:
 * 1. POST /trades/{id}/images/upload-signature → signed payload
 * 2. browser POSTs the file straight to Cloudinary (FormData)
 * 3. POST /trades/{id}/images with the returned public_id
 * Cap: 5 per trade. Compress to WebP ≤1600px / <2MB before upload.
 */
import type { ImageKind } from './trades';

export interface UploadSignature {
  upload_url: string;
  cloud_name: string;
  api_key: string;
  timestamp: string;
  signature: string;
  folder: string;
  type: string;
}

/** Confirm body after the direct browser upload. */
export interface ImageConfirm {
  public_id: string;
  kind?: ImageKind;
}

export interface TradeImage {
  ID: string;
  TradeID: string;
  PublicID: string;
  Kind: ImageKind | null;
  Position: number;
  CreatedAt: string;
}

/** List item: stored image + short-lived signed delivery URL. */
export interface TradeImageWithUrl extends TradeImage {
  URL: string;
}

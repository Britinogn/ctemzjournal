package cloudinary

import (
	"context"
	"net/url"
	"strconv"
	"time"

	cloudinary "github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/cloudinary/cloudinary-go/v2/asset"
)

// MaxImagesPerTrade and upload limits (docx §3, enforced server-side too).
const (
	MaxImagesPerTrade = 5
	MaxBytes          = 2 << 20 // 2 MB after frontend WebP compression
)

// AllowedTypes are the only file formats accepted (frontend compresses to webp).
var AllowedTypes = map[string]bool{"jpeg": true, "png": true, "webp": true}

// Client wraps Cloudinary: signed upload payloads, signed delivery URLs
// (type authenticated = private), and destroy. The API secret never leaves
// the backend.
type Client struct {
	cld       *cloudinary.Cloudinary
	cloudName string
	apiKey    string
}

// New builds the client. Empty credentials disable it (IsConfigured false);
// routes still mount but signature requests fail with a clear error.
func New(cloudName, apiKey, apiSecret string) (*Client, error) {
	if cloudName == "" || apiKey == "" || apiSecret == "" {
		return &Client{}, nil
	}
	cld, err := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	if err != nil {
		return nil, err
	}
	return &Client{cld: cld, cloudName: cloudName, apiKey: apiKey}, nil
}

// IsConfigured reports whether credentials were provided.
func (c *Client) IsConfigured() bool { return c.cld != nil }

// Signature is the signed upload payload the browser POSTs straight to
// Cloudinary with plain FormData (no extra package).
type Signature struct {
	UploadURL string `json:"upload_url"`
	CloudName string `json:"cloud_name"`
	APIKey    string `json:"api_key"`
	Timestamp string `json:"timestamp"`
	Signature string `json:"signature"`
	Folder    string `json:"folder"`
	Type      string `json:"type"`
}

// SignUpload returns the signed payload for trade-images/{userID}/{tradeID}.
// The browser uploads with type=authenticated so screenshots stay private.
func (c *Client) SignUpload(userID, tradeID string) (Signature, error) {
	if !c.IsConfigured() {
		return Signature{}, errNotConfigured
	}
	folder := "trade-images/" + userID + "/" + tradeID
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	params := url.Values{}
	params.Set("timestamp", timestamp)
	params.Set("folder", folder)
	params.Set("type", string(api.Authenticated))
	sig, err := api.SignParameters(params, c.cld.Config.Cloud.APISecret)
	if err != nil {
		return Signature{}, err
	}
	return Signature{
		UploadURL: "https://api.cloudinary.com/v1_1/" + c.cloudName + "/image/upload",
		CloudName: c.cloudName,
		APIKey:    c.apiKey,
		Timestamp: timestamp,
		Signature: sig,
		Folder:    folder,
		Type:      string(api.Authenticated),
	}, nil
}

// DeliveryURL creates a short-lived signed URL for a private (authenticated)
// image. Readers never see the raw public_id.
func (c *Client) DeliveryURL(publicID string) (string, error) {
	if !c.IsConfigured() {
		return "", errNotConfigured
	}
	cfg := c.cld.Config
	cfg.URL.SignURL = true // per-call copy: safe for concurrent workers
	img, err := asset.Image(publicID, &cfg)
	if err != nil {
		return "", err
	}
	img.DeliveryType = api.Authenticated
	return img.String()
}

// Destroy deletes an image from Cloudinary (called on image/trade delete).
func (c *Client) Destroy(ctx context.Context, publicID string) error {
	if !c.IsConfigured() {
		return errNotConfigured
	}
	_, err := c.cld.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID: publicID,
		Type:     string(api.Authenticated),
	})
	return err
}

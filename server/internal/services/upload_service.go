package services

import (
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"layanesia-server/internal/apperrors"

	"github.com/google/uuid"
)

// UploadPurpose menentukan aturan validasi & folder penyimpanan.
type UploadPurpose string

const (
	UploadPurposeKTP    UploadPurpose = "ktp"    // gambar JPEG/PNG
	UploadPurposeResume UploadPurpose = "resume" // PDF
)

type UploadResult struct {
	URL       string `json:"file_url"`
	PublicID  string `json:"public_id,omitempty"`
	Bytes     int64  `json:"bytes"`
	MIMEType  string `json:"mime_type"`
	Storage   string `json:"storage"` // cloudinary | local
	ExpiresIn string `json:"expires_in,omitempty"`
}

type UploadConfig struct {
	CloudName     string
	APIKey        string
	APISecret     string
	Folder        string
	LocalDir      string
	MaxImageBytes int64
	MaxDocBytes   int64
	AllowLocal    bool // fallback ke disk lokal (hanya dev)
	HTTPClient    *http.Client
}

type UploadService interface {
	Upload(ctx context.Context, ownerID uuid.UUID, purpose UploadPurpose, fh *multipart.FileHeader) (*UploadResult, error)
	LocalDir() string
	UsesLocalStorage() bool
}

type uploadService struct {
	cfg UploadConfig
}

func NewUploadService(cfg UploadConfig) UploadService {
	if cfg.HTTPClient == nil {
		cfg.HTTPClient = &http.Client{Timeout: 30 * time.Second}
	}
	if cfg.Folder == "" {
		cfg.Folder = "layanesia"
	}
	return &uploadService{cfg: cfg}
}

func (s *uploadService) cloudinaryConfigured() bool {
	return s.cfg.CloudName != "" && s.cfg.APIKey != "" && s.cfg.APISecret != ""
}

func (s *uploadService) UsesLocalStorage() bool { return !s.cloudinaryConfigured() && s.cfg.AllowLocal }
func (s *uploadService) LocalDir() string        { return s.cfg.LocalDir }

var allowedTypes = map[UploadPurpose]map[string]string{
	UploadPurposeKTP:    {"image/jpeg": ".jpg", "image/png": ".png"},
	UploadPurposeResume: {"application/pdf": ".pdf"},
}

// ParseUploadPurpose menerima "ktp" | "resume"/"cv" | "" (deteksi otomatis dari isi file).
func ParseUploadPurpose(s string) (UploadPurpose, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "ktp":
		return UploadPurposeKTP, nil
	case "resume", "cv":
		return UploadPurposeResume, nil
	case "":
		return "", nil
	}
	return "", apperrors.BadRequest("Parameter 'purpose' tidak valid: gunakan 'ktp' atau 'resume'.")
}

// Upload memvalidasi ukuran & tipe file berdasarkan isi (content sniffing, bukan
// ekstensi) lalu mengunggah ke Cloudinary (signed upload). Jika Cloudinary tidak
// dikonfigurasi dan AllowLocal aktif, file disimpan di disk lokal (dev only).
// purpose kosong = tentukan dari tipe file (gambar -> ktp, PDF -> resume).
func (s *uploadService) Upload(ctx context.Context, ownerID uuid.UUID, purpose UploadPurpose, fh *multipart.FileHeader) (*UploadResult, error) {
	if purpose != "" {
		if _, ok := allowedTypes[purpose]; !ok {
			return nil, apperrors.BadRequest("Jenis upload tidak dikenal.")
		}
	}

	limit := s.cfg.MaxDocBytes
	if s.cfg.MaxImageBytes > limit {
		limit = s.cfg.MaxImageBytes
	}
	if purpose == UploadPurposeKTP {
		limit = s.cfg.MaxImageBytes
	} else if purpose == UploadPurposeResume {
		limit = s.cfg.MaxDocBytes
	}
	if fh.Size <= 0 {
		return nil, apperrors.BadRequest("File kosong.")
	}
	if fh.Size > limit {
		return nil, apperrors.New(http.StatusRequestEntityTooLarge, "FILE_TOO_LARGE",
			fmt.Sprintf("Ukuran file maksimal %d MB.", limit/(1024*1024)))
	}

	f, err := fh.Open()
	if err != nil {
		return nil, apperrors.BadRequest("File tidak dapat dibaca.")
	}
	defer f.Close()

	data, err := io.ReadAll(io.LimitReader(f, limit+1))
	if err != nil {
		return nil, apperrors.BadRequest("File tidak dapat dibaca.")
	}
	if int64(len(data)) > limit {
		return nil, apperrors.New(http.StatusRequestEntityTooLarge, "FILE_TOO_LARGE",
			fmt.Sprintf("Ukuran file maksimal %d MB.", limit/(1024*1024)))
	}

	detected := sniffMIME(data)
	if purpose == "" {
		switch {
		case strings.HasPrefix(detected, "image/"):
			purpose = UploadPurposeKTP
		case detected == "application/pdf":
			purpose = UploadPurposeResume
		default:
			return nil, apperrors.New(http.StatusUnsupportedMediaType, "UNSUPPORTED_FILE_TYPE",
				fmt.Sprintf("Tipe file tidak didukung (%s). Diizinkan: image/jpeg, image/png, application/pdf.", detected))
		}
		if purpose == UploadPurposeKTP && int64(len(data)) > s.cfg.MaxImageBytes {
			return nil, apperrors.New(http.StatusRequestEntityTooLarge, "FILE_TOO_LARGE",
				fmt.Sprintf("Ukuran gambar maksimal %d MB.", s.cfg.MaxImageBytes/(1024*1024)))
		}
		if purpose == UploadPurposeResume && int64(len(data)) > s.cfg.MaxDocBytes {
			return nil, apperrors.New(http.StatusRequestEntityTooLarge, "FILE_TOO_LARGE",
				fmt.Sprintf("Ukuran dokumen maksimal %d MB.", s.cfg.MaxDocBytes/(1024*1024)))
		}
	}
	allowed := allowedTypes[purpose]
	ext, ok := allowed[detected]
	if !ok {
		names := make([]string, 0, len(allowed))
		for k := range allowed {
			names = append(names, k)
		}
		sort.Strings(names)
		return nil, apperrors.New(http.StatusUnsupportedMediaType, "UNSUPPORTED_FILE_TYPE",
			fmt.Sprintf("Tipe file tidak didukung (%s). Diizinkan: %s.", detected, strings.Join(names, ", ")))
	}

	publicID := fmt.Sprintf("%s/%s/%s_%d", s.cfg.Folder, purpose, ownerID.String(), time.Now().Unix())

	if s.cloudinaryConfigured() {
		return s.uploadToCloudinary(ctx, data, detected, publicID, purpose)
	}
	if s.cfg.AllowLocal {
		return s.saveLocal(data, detected, purpose, ext)
	}
	return nil, apperrors.Unavailable("Penyimpanan file belum dikonfigurasi.", errors.New("cloudinary not configured and local storage disabled"))
}

// sniffMIME mendeteksi tipe dari magic bytes. http.DetectContentType tidak mengenali
// PDF, jadi ditangani manual.
func sniffMIME(data []byte) string {
	if bytes.HasPrefix(data, []byte("%PDF-")) {
		return "application/pdf"
	}
	ct := http.DetectContentType(data)
	if i := strings.Index(ct, ";"); i >= 0 {
		ct = ct[:i]
	}
	return ct
}

func (s *uploadService) saveLocal(data []byte, mime string, purpose UploadPurpose, ext string) (*UploadResult, error) {
	dir := filepath.Join(s.cfg.LocalDir, string(purpose))
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, apperrors.Internal(err)
	}
	name := uuid.NewString() + ext
	if err := os.WriteFile(filepath.Join(dir, name), data, 0o644); err != nil {
		return nil, apperrors.Internal(err)
	}
	return &UploadResult{
		URL:      fmt.Sprintf("/uploads/%s/%s", purpose, name),
		Bytes:    int64(len(data)),
		MIMEType: mime,
		Storage:  "local",
	}, nil
}

// uploadToCloudinary melakukan signed upload via REST API (tanpa SDK).
// KTP diunggah dengan type=authenticated sehingga URL-nya tidak bisa diakses publik
// tanpa signed URL; resume memakai type=upload (dibagikan ke pemberi kerja setelah melamar).
func (s *uploadService) uploadToCloudinary(ctx context.Context, data []byte, mime, publicID string, purpose UploadPurpose) (*UploadResult, error) {
	resourceType := "raw"
	if strings.HasPrefix(mime, "image/") {
		resourceType = "image"
	}
	deliveryType := "upload"
	if purpose == UploadPurposeKTP {
		deliveryType = "authenticated"
	}

	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	params := map[string]string{
		"public_id": publicID,
		"timestamp": timestamp,
		"type":      deliveryType,
		"overwrite": "true",
	}
	signature := cloudinarySignature(params, s.cfg.APISecret)

	var body bytes.Buffer
	w := multipart.NewWriter(&body)
	for k, v := range params {
		_ = w.WriteField(k, v)
	}
	_ = w.WriteField("api_key", s.cfg.APIKey)
	_ = w.WriteField("signature", signature)
	part, err := w.CreateFormFile("file", "upload")
	if err != nil {
		return nil, apperrors.Internal(err)
	}
	if _, err := part.Write(data); err != nil {
		return nil, apperrors.Internal(err)
	}
	_ = w.Close()

	url := fmt.Sprintf("https://api.cloudinary.com/v1_1/%s/%s/upload", s.cfg.CloudName, resourceType)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, &body)
	if err != nil {
		return nil, apperrors.Internal(err)
	}
	req.Header.Set("Content-Type", w.FormDataContentType())

	resp, err := s.cfg.HTTPClient.Do(req)
	if err != nil {
		return nil, apperrors.Unavailable("Tidak dapat menghubungi layanan penyimpanan file.", err)
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, apperrors.Unavailable("Layanan penyimpanan file menolak unggahan.",
			fmt.Errorf("cloudinary %d: %s", resp.StatusCode, truncate(string(raw), 500)))
	}

	var out struct {
		SecureURL string `json:"secure_url"`
		PublicID  string `json:"public_id"`
		Bytes     int64  `json:"bytes"`
	}
	if err := json.Unmarshal(raw, &out); err != nil || out.SecureURL == "" {
		return nil, apperrors.Unavailable("Respons layanan penyimpanan tidak dapat dibaca.", err)
	}
	return &UploadResult{
		URL:      out.SecureURL,
		PublicID: out.PublicID,
		Bytes:    out.Bytes,
		MIMEType: mime,
		Storage:  "cloudinary",
	}, nil
}

// cloudinarySignature = SHA1(sorted "k=v" joined by "&" + api_secret).
func cloudinarySignature(params map[string]string, secret string) string {
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		parts = append(parts, k+"="+params[k])
	}
	sum := sha1.Sum([]byte(strings.Join(parts, "&") + secret))
	return hex.EncodeToString(sum[:])
}

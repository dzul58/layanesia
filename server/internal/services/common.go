package services

import (
	"strings"
	"time"

	"layanesia-server/internal/apperrors"
	"layanesia-server/internal/models"
)

// PageMeta adalah metadata pagination untuk endpoint daftar.
type PageMeta struct {
	Total  int64 `json:"total"`
	Limit  int   `json:"limit"`
	Offset int   `json:"offset"`
}

// postingLocation adalah lokasi hasil resolusi: input request diutamakan, jika kosong
// jatuh ke lokasi domisili user (ERD §2.3/§2.4).
type postingLocation struct {
	Country, Province, City, District string
	AddressDetail, FormattedAddress   *string
	OSMPlaceID                        *string
	Latitude, Longitude               *float64
}

type locationInput struct {
	Country, Province, City, District string
	AddressDetail, FormattedAddress   *string
	OSMPlaceID                        *string
	Latitude, Longitude               *float64
}

func resolveLocation(in locationInput, user *models.User) (postingLocation, error) {
	loc := postingLocation{
		Country:          firstNonBlank(in.Country, user.Country, "Indonesia"),
		Province:         firstNonBlank(in.Province, user.Province),
		City:             firstNonBlank(in.City, user.City),
		District:         firstNonBlank(in.District, user.District),
		AddressDetail:    in.AddressDetail,
		FormattedAddress: in.FormattedAddress,
		OSMPlaceID:       in.OSMPlaceID,
		Latitude:         in.Latitude,
		Longitude:        in.Longitude,
	}
	// Detail alamat/koordinat hanya diwariskan jika lokasi administratif juga diwariskan
	// (mencegah koordinat rumah user menempel pada posting di kota lain).
	inheritedAdmin := in.Province == "" && in.City == "" && in.District == ""
	if inheritedAdmin {
		if loc.AddressDetail == nil {
			loc.AddressDetail = user.AddressDetail
		}
		if loc.FormattedAddress == nil {
			loc.FormattedAddress = user.FormattedAddress
		}
		if loc.OSMPlaceID == nil {
			loc.OSMPlaceID = user.OSMPlaceID
		}
		if loc.Latitude == nil && loc.Longitude == nil {
			loc.Latitude, loc.Longitude = user.Latitude, user.Longitude
		}
	}
	if loc.Province == "" || loc.City == "" || loc.District == "" {
		return loc, apperrors.Validation(
			"Lokasi (provinsi, kota/kabupaten, kecamatan) wajib diisi atau diatur di profil domisili Anda.", nil)
	}
	if (loc.Latitude == nil) != (loc.Longitude == nil) {
		return loc, apperrors.Validation("Latitude dan longitude harus diisi bersamaan.", nil)
	}
	return loc, nil
}

func firstNonBlank(vals ...string) string {
	for _, v := range vals {
		if s := strings.TrimSpace(v); s != "" {
			return s
		}
	}
	return ""
}

// parseWorkDate menerima YYYY-MM-DD atau RFC3339 dan menolak tanggal lampau.
// Kosong berarti "besok".
func parseWorkDate(raw string, now time.Time) (time.Time, error) {
	raw = strings.TrimSpace(raw)
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	if raw == "" {
		return today.AddDate(0, 0, 1), nil
	}
	var t time.Time
	var err error
	if t, err = time.ParseInLocation("2006-01-02", raw, now.Location()); err != nil {
		if t, err = time.Parse(time.RFC3339, raw); err != nil {
			return time.Time{}, apperrors.Validation("Format tanggal kerja tidak valid.",
				map[string]any{"work_date": "gunakan format YYYY-MM-DD"})
		}
		t = t.In(now.Location())
	}
	t = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, now.Location())
	if t.Before(today) {
		return time.Time{}, apperrors.Validation("Tanggal kerja tidak boleh di masa lalu.",
			map[string]any{"work_date": "harus hari ini atau setelahnya"})
	}
	if t.After(today.AddDate(1, 0, 0)) {
		return time.Time{}, apperrors.Validation("Tanggal kerja terlalu jauh.",
			map[string]any{"work_date": "maksimal 1 tahun ke depan"})
	}
	return t, nil
}

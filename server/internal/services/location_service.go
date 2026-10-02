package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"layanesia-server/internal/apperrors"
	"layanesia-server/internal/models"
	"layanesia-server/internal/repositories"
)

type LocationService interface {
	GetProvinces() ([]string, error)
	GetCitiesByProvince(province string) ([]string, error)
	GetDistrictsByCity(city string) ([]models.Location, error)
	SearchLocations(query string) ([]models.Location, error)
	GeocodeNominatim(ctx context.Context, searchQuery string) ([]NominatimSearchResult, error)
}

type locationService struct {
	repo       repositories.LocationRepository
	httpClient *http.Client
	userAgent  string
}

type NominatimSearchResult struct {
	PlaceID     int64    `json:"place_id"`
	Licence     string   `json:"licence,omitempty"`
	OsmType     string   `json:"osm_type,omitempty"`
	OsmID       int64    `json:"osm_id,omitempty"`
	Boundingbox []string `json:"boundingbox,omitempty"`
	Lat         string   `json:"lat"`
	Lon         string   `json:"lon"`
	DisplayName string   `json:"display_name"`
	Class       string   `json:"class,omitempty"`
	Type        string   `json:"type,omitempty"`
	Importance  float64  `json:"importance,omitempty"`
}

func NewLocationService(repo repositories.LocationRepository) LocationService {
	return &locationService{
		repo:       repo,
		httpClient: &http.Client{Timeout: 10 * time.Second},
		// Wajib menurut Nominatim Usage Policy.
		userAgent: "LayanesiaMarketplace/1.0 (contact@layanesia.id)",
	}
}

func (s *locationService) GetProvinces() ([]string, error) {
	out, err := s.repo.GetProvinces()
	if err != nil {
		return nil, apperrors.Internal(err)
	}
	return out, nil
}

func (s *locationService) GetCitiesByProvince(province string) ([]string, error) {
	out, err := s.repo.GetCitiesByProvince(strings.TrimSpace(province))
	if err != nil {
		return nil, apperrors.Internal(err)
	}
	return out, nil
}

func (s *locationService) GetDistrictsByCity(city string) ([]models.Location, error) {
	out, err := s.repo.GetDistrictsByCity(strings.TrimSpace(city))
	if err != nil {
		return nil, apperrors.Internal(err)
	}
	return out, nil
}

func (s *locationService) SearchLocations(query string) ([]models.Location, error) {
	query = strings.TrimSpace(query)
	if len(query) < 2 {
		return nil, apperrors.Validation("Kata kunci minimal 2 karakter.", map[string]any{"q": "minimal 2 karakter"})
	}
	if len(query) > 100 {
		return nil, apperrors.Validation("Kata kunci terlalu panjang.", map[string]any{"q": "maksimal 100 karakter"})
	}
	out, err := s.repo.SearchLocations(query)
	if err != nil {
		return nil, apperrors.Internal(err)
	}
	return out, nil
}

func (s *locationService) GeocodeNominatim(ctx context.Context, searchQuery string) ([]NominatimSearchResult, error) {
	searchQuery = strings.TrimSpace(searchQuery)
	if len(searchQuery) < 3 {
		return nil, apperrors.Validation("Kata kunci minimal 3 karakter.", map[string]any{"q": "minimal 3 karakter"})
	}
	if len(searchQuery) > 200 {
		return nil, apperrors.Validation("Kata kunci terlalu panjang.", map[string]any{"q": "maksimal 200 karakter"})
	}

	endpoint := fmt.Sprintf("https://nominatim.openstreetmap.org/search?format=json&q=%s&countrycodes=id&limit=5",
		url.QueryEscape(searchQuery))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, apperrors.Internal(err)
	}
	req.Header.Set("User-Agent", s.userAgent)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, apperrors.Unavailable("Layanan geocoding tidak dapat dihubungi.", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return nil, apperrors.Unavailable("Respons layanan geocoding tidak dapat dibaca.", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, apperrors.Unavailable("Layanan geocoding menolak permintaan.",
			fmt.Errorf("nominatim %d: %s", resp.StatusCode, truncate(string(body), 300)))
	}

	var results []NominatimSearchResult
	if err := json.Unmarshal(body, &results); err != nil {
		return nil, apperrors.Unavailable("Respons layanan geocoding tidak valid.", err)
	}
	if results == nil {
		results = []NominatimSearchResult{}
	}
	return results, nil
}

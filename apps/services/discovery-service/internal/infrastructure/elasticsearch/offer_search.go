// Package elasticsearch implements the search infrastructure for the Discovery service.
package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/yousoon/discovery-service/internal/domain"
)

const (
	offersIndex = "offers"
)

// OfferSearchRepository implements domain.OfferSearchService using Elasticsearch.
type OfferSearchRepository struct {
	client *elasticsearch.Client
}

// NewOfferSearchRepository creates a new Elasticsearch search repository.
func NewOfferSearchRepository(client *elasticsearch.Client) *OfferSearchRepository {
	return &OfferSearchRepository{
		client: client,
	}
}

// offerDocument represents an offer document in Elasticsearch.
type offerDocument struct {
	ID                    string    `json:"id"`
	PartnerID             string    `json:"partner_id"`
	EstablishmentID       string    `json:"establishment_id"`
	Title                 string    `json:"title"`
	Description           string    `json:"description"`
	ShortDescription      string    `json:"short_description"`
	CategoryID            string    `json:"category_id"`
	Tags                  []string  `json:"tags"`
	DiscountType          string    `json:"discount_type"`
	DiscountValue         int       `json:"discount_value"`
	OriginalPrice         int64     `json:"original_price,omitempty"`
	OriginalPriceCurrency string    `json:"original_price_currency,omitempty"`
	DiscountedPrice       int64     `json:"discounted_price,omitempty"`
	Formula               string    `json:"formula,omitempty"`
	Status                string    `json:"status"`
	ValidityStartDate     time.Time `json:"validity_start_date"`
	ValidityEndDate       time.Time `json:"validity_end_date"`
	Location              GeoPoint  `json:"location"`
	PartnerName           string    `json:"partner_name"`
	EstablishmentName     string    `json:"establishment_name"`
	EstablishmentCity     string    `json:"establishment_city"`
	Views                 int64     `json:"views"`
	Clicks                int64     `json:"clicks"`
	Bookings              int64     `json:"bookings"`
	Favorites             int64     `json:"favorites"`
	AvgRating             float64   `json:"avg_rating"`
	ReviewCount           int       `json:"review_count"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
	PublishedAt           time.Time `json:"published_at,omitempty"`
}

// GeoPoint represents a geo point for Elasticsearch.
type GeoPoint struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

// searchResult represents the Elasticsearch search response.
type searchResult struct {
	Hits struct {
		Total struct {
			Value int64 `json:"value"`
		} `json:"total"`
		Hits []struct {
			ID     string        `json:"_id"`
			Score  float64       `json:"_score"`
			Source offerDocument `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
	Aggregations map[string]json.RawMessage `json:"aggregations,omitempty"`
}

// EnsureIndex creates the offers index with proper mappings if it doesn't exist.
func (r *OfferSearchRepository) EnsureIndex(ctx context.Context) error {
	// Check if index exists
	res, err := r.client.Indices.Exists([]string{offersIndex})
	if err != nil {
		return fmt.Errorf("failed to check index existence: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		return nil // Index already exists
	}

	// Create index with mappings
	mapping := `{
		"settings": {
			"number_of_shards": 2,
			"number_of_replicas": 1,
			"analysis": {
				"analyzer": {
					"french_analyzer": {
						"type": "custom",
						"tokenizer": "standard",
						"filter": ["lowercase", "french_elision", "french_stop", "french_stemmer"]
					}
				},
				"filter": {
					"french_elision": {
						"type": "elision",
						"articles_case": true,
						"articles": ["l", "m", "t", "qu", "n", "s", "j", "d", "c", "jusqu", "quoiqu", "lorsqu", "puisqu"]
					},
					"french_stop": {
						"type": "stop",
						"stopwords": "_french_"
					},
					"french_stemmer": {
						"type": "stemmer",
						"language": "light_french"
					}
				}
			}
		},
		"mappings": {
			"properties": {
				"id": { "type": "keyword" },
				"partner_id": { "type": "keyword" },
				"establishment_id": { "type": "keyword" },
				"title": {
					"type": "text",
					"analyzer": "french_analyzer",
					"fields": {
						"keyword": { "type": "keyword" },
						"autocomplete": {
							"type": "text",
							"analyzer": "simple"
						}
					}
				},
				"description": {
					"type": "text",
					"analyzer": "french_analyzer"
				},
				"short_description": {
					"type": "text",
					"analyzer": "french_analyzer"
				},
				"category_id": { "type": "keyword" },
				"tags": { "type": "keyword" },
				"discount_type": { "type": "keyword" },
				"discount_value": { "type": "integer" },
				"original_price": { "type": "long" },
				"original_price_currency": { "type": "keyword" },
				"discounted_price": { "type": "long" },
				"formula": { "type": "text" },
				"status": { "type": "keyword" },
				"validity_start_date": { "type": "date" },
				"validity_end_date": { "type": "date" },
				"location": { "type": "geo_point" },
				"partner_name": {
					"type": "text",
					"analyzer": "french_analyzer",
					"fields": {
						"keyword": { "type": "keyword" }
					}
				},
				"establishment_name": {
					"type": "text",
					"analyzer": "french_analyzer",
					"fields": {
						"keyword": { "type": "keyword" }
					}
				},
				"establishment_city": {
					"type": "text",
					"analyzer": "french_analyzer",
					"fields": {
						"keyword": { "type": "keyword" }
					}
				},
				"views": { "type": "long" },
				"clicks": { "type": "long" },
				"bookings": { "type": "long" },
				"favorites": { "type": "long" },
				"avg_rating": { "type": "float" },
				"review_count": { "type": "integer" },
				"created_at": { "type": "date" },
				"updated_at": { "type": "date" },
				"published_at": { "type": "date" }
			}
		}
	}`

	res, err = r.client.Indices.Create(
		offersIndex,
		r.client.Indices.Create.WithBody(strings.NewReader(mapping)),
		r.client.Indices.Create.WithContext(ctx),
	)
	if err != nil {
		return fmt.Errorf("failed to create index: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("failed to create index: %s", res.String())
	}

	return nil
}

// IndexOffer indexes an offer in Elasticsearch.
func (r *OfferSearchRepository) IndexOffer(ctx context.Context, offer *domain.Offer) error {
	doc := r.offerToDocument(offer)

	data, err := json.Marshal(doc)
	if err != nil {
		return fmt.Errorf("failed to marshal offer: %w", err)
	}

	req := esapi.IndexRequest{
		Index:      offersIndex,
		DocumentID: offer.ID().String(),
		Body:       bytes.NewReader(data),
		Refresh:    "true",
	}

	res, err := req.Do(ctx, r.client)
	if err != nil {
		return fmt.Errorf("failed to index offer: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("failed to index offer: %s", res.String())
	}

	return nil
}

// UpdateOffer updates an offer in Elasticsearch.
func (r *OfferSearchRepository) UpdateOffer(ctx context.Context, offer *domain.Offer) error {
	return r.IndexOffer(ctx, offer) // Reindex the document
}

// DeleteOffer removes an offer from Elasticsearch.
func (r *OfferSearchRepository) DeleteOffer(ctx context.Context, id domain.OfferID) error {
	req := esapi.DeleteRequest{
		Index:      offersIndex,
		DocumentID: id.String(),
		Refresh:    "true",
	}

	res, err := req.Do(ctx, r.client)
	if err != nil {
		return fmt.Errorf("failed to delete offer: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() && res.StatusCode != 404 {
		return fmt.Errorf("failed to delete offer: %s", res.String())
	}

	return nil
}

// Search performs a full-text search on offers.
func (r *OfferSearchRepository) Search(ctx context.Context, query string, filter domain.OfferFilter) ([]*domain.OfferSummary, int64, error) {
	searchQuery := r.buildSearchQuery(query, filter)

	data, err := json.Marshal(searchQuery)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to marshal query: %w", err)
	}

	res, err := r.client.Search(
		r.client.Search.WithContext(ctx),
		r.client.Search.WithIndex(offersIndex),
		r.client.Search.WithBody(bytes.NewReader(data)),
		r.client.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, 0, fmt.Errorf("search error: %s", res.String())
	}

	var result searchResult
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, 0, fmt.Errorf("failed to decode response: %w", err)
	}

	summaries := make([]*domain.OfferSummary, 0, len(result.Hits.Hits))
	for _, hit := range result.Hits.Hits {
		summary := r.documentToSummary(&hit.Source)
		summaries = append(summaries, summary)
	}

	return summaries, result.Hits.Total.Value, nil
}

// SearchNearby searches for offers near a location.
func (r *OfferSearchRepository) SearchNearby(ctx context.Context, lat, lng float64, radiusKm float64, filter domain.OfferFilter) ([]*domain.OfferSummary, int64, error) {
	searchQuery := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []interface{}{
					map[string]interface{}{
						"term": map[string]interface{}{
							"status": string(domain.OfferStatusActive),
						},
					},
				},
				"filter": []interface{}{
					map[string]interface{}{
						"geo_distance": map[string]interface{}{
							"distance": fmt.Sprintf("%.0fkm", radiusKm),
							"location": map[string]interface{}{
								"lat": lat,
								"lon": lng,
							},
						},
					},
					map[string]interface{}{
						"range": map[string]interface{}{
							"validity_end_date": map[string]interface{}{
								"gte": "now",
							},
						},
					},
				},
			},
		},
		"sort": []interface{}{
			map[string]interface{}{
				"_geo_distance": map[string]interface{}{
					"location": map[string]interface{}{
						"lat": lat,
						"lon": lng,
					},
					"order":         "asc",
					"unit":          "km",
					"distance_type": "arc",
				},
			},
		},
		"from": filter.Offset,
		"size": filter.Limit,
	}

	// Add category filter if specified
	if filter.CategoryID != nil {
		searchQuery["query"].(map[string]interface{})["bool"].(map[string]interface{})["must"] = append(
			searchQuery["query"].(map[string]interface{})["bool"].(map[string]interface{})["must"].([]interface{}),
			map[string]interface{}{
				"term": map[string]interface{}{
					"category_id": filter.CategoryID.String(),
				},
			},
		)
	}

	data, err := json.Marshal(searchQuery)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to marshal query: %w", err)
	}

	res, err := r.client.Search(
		r.client.Search.WithContext(ctx),
		r.client.Search.WithIndex(offersIndex),
		r.client.Search.WithBody(bytes.NewReader(data)),
		r.client.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, 0, fmt.Errorf("search error: %s", res.String())
	}

	var result searchResult
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, 0, fmt.Errorf("failed to decode response: %w", err)
	}

	summaries := make([]*domain.OfferSummary, 0, len(result.Hits.Hits))
	for _, hit := range result.Hits.Hits {
		summary := r.documentToSummary(&hit.Source)
		summaries = append(summaries, summary)
	}

	return summaries, result.Hits.Total.Value, nil
}

// GetAutocomplete returns autocomplete suggestions for a query.
func (r *OfferSearchRepository) GetAutocomplete(ctx context.Context, query string, limit int) ([]string, error) {
	searchQuery := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []interface{}{
					map[string]interface{}{
						"multi_match": map[string]interface{}{
							"query":  query,
							"type":   "bool_prefix",
							"fields": []string{"title.autocomplete", "partner_name.keyword", "establishment_city.keyword"},
						},
					},
				},
				"filter": []interface{}{
					map[string]interface{}{
						"term": map[string]interface{}{
							"status": string(domain.OfferStatusActive),
						},
					},
				},
			},
		},
		"_source": []string{"title"},
		"size":    limit,
	}

	data, err := json.Marshal(searchQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal query: %w", err)
	}

	res, err := r.client.Search(
		r.client.Search.WithContext(ctx),
		r.client.Search.WithIndex(offersIndex),
		r.client.Search.WithBody(bytes.NewReader(data)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to search: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("search error: %s", res.String())
	}

	var result searchResult
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	suggestions := make([]string, 0, len(result.Hits.Hits))
	seen := make(map[string]bool)
	for _, hit := range result.Hits.Hits {
		title := hit.Source.Title
		if !seen[title] {
			suggestions = append(suggestions, title)
			seen[title] = true
		}
	}

	return suggestions, nil
}

// buildSearchQuery builds the Elasticsearch query from domain filter.
func (r *OfferSearchRepository) buildSearchQuery(query string, filter domain.OfferFilter) map[string]interface{} {
	must := make([]interface{}, 0)
	filterClauses := make([]interface{}, 0)

	// Full-text search query
	if query != "" {
		must = append(must, map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":     query,
				"fields":    []string{"title^3", "description", "short_description^2", "partner_name^2", "establishment_name", "tags"},
				"type":      "best_fields",
				"fuzziness": "AUTO",
			},
		})
	}

	// Status filter
	if filter.Status != nil {
		filterClauses = append(filterClauses, map[string]interface{}{
			"term": map[string]interface{}{
				"status": string(*filter.Status),
			},
		})
	} else {
		// Default to active offers
		filterClauses = append(filterClauses, map[string]interface{}{
			"term": map[string]interface{}{
				"status": string(domain.OfferStatusActive),
			},
		})
	}

	// Category filter
	if filter.CategoryID != nil {
		filterClauses = append(filterClauses, map[string]interface{}{
			"term": map[string]interface{}{
				"category_id": filter.CategoryID.String(),
			},
		})
	}

	// Partner filter
	if filter.PartnerID != nil {
		filterClauses = append(filterClauses, map[string]interface{}{
			"term": map[string]interface{}{
				"partner_id": filter.PartnerID.String(),
			},
		})
	}

	// Establishment filter
	if filter.EstablishmentID != nil {
		filterClauses = append(filterClauses, map[string]interface{}{
			"term": map[string]interface{}{
				"establishment_id": filter.EstablishmentID.String(),
			},
		})
	}

	// Active only filter (validity dates)
	if filter.ActiveOnly {
		now := time.Now()
		filterClauses = append(filterClauses,
			map[string]interface{}{
				"range": map[string]interface{}{
					"validity_start_date": map[string]interface{}{
						"lte": now,
					},
				},
			},
			map[string]interface{}{
				"range": map[string]interface{}{
					"validity_end_date": map[string]interface{}{
						"gte": now,
					},
				},
			},
		)
	}

	// Min rating filter
	if filter.MinRating != nil && *filter.MinRating > 0 {
		filterClauses = append(filterClauses, map[string]interface{}{
			"range": map[string]interface{}{
				"avg_rating": map[string]interface{}{
					"gte": *filter.MinRating,
				},
			},
		})
	}

	// Discount type filter
	if filter.DiscountType != nil {
		filterClauses = append(filterClauses, map[string]interface{}{
			"term": map[string]interface{}{
				"discount_type": string(*filter.DiscountType),
			},
		})
	}

	// Build sort
	sort := make([]interface{}, 0)
	switch filter.SortBy {
	case "rating":
		sort = append(sort, map[string]interface{}{
			"avg_rating": map[string]interface{}{"order": "desc"},
		})
	case "newest":
		sort = append(sort, map[string]interface{}{
			"published_at": map[string]interface{}{"order": "desc"},
		})
	case "popular":
		sort = append(sort, map[string]interface{}{
			"views": map[string]interface{}{"order": "desc"},
		})
	case "discount":
		sort = append(sort, map[string]interface{}{
			"discount_value": map[string]interface{}{"order": "desc"},
		})
	default:
		// Default: relevance score + popularity
		sort = append(sort, "_score")
		sort = append(sort, map[string]interface{}{
			"views": map[string]interface{}{"order": "desc"},
		})
	}

	searchQuery := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must":   must,
				"filter": filterClauses,
			},
		},
		"sort": sort,
		"from": filter.Offset,
		"size": filter.Limit,
	}

	// Add geo distance sort if location provided
	if filter.Latitude != nil && filter.Longitude != nil {
		sort = append([]interface{}{
			map[string]interface{}{
				"_geo_distance": map[string]interface{}{
					"location": map[string]interface{}{
						"lat": *filter.Latitude,
						"lon": *filter.Longitude,
					},
					"order": "asc",
					"unit":  "km",
				},
			},
		}, sort...)
		searchQuery["sort"] = sort
	}

	return searchQuery
}

// offerToDocument converts a domain offer to an Elasticsearch document.
func (r *OfferSearchRepository) offerToDocument(offer *domain.Offer) *offerDocument {
	doc := &offerDocument{
		ID:                offer.ID().String(),
		PartnerID:         offer.PartnerID().String(),
		EstablishmentID:   offer.EstablishmentID().String(),
		Title:             offer.Title(),
		Description:       offer.Description(),
		ShortDescription:  offer.ShortDescription(),
		CategoryID:        offer.CategoryID().String(),
		Tags:              offer.Tags(),
		DiscountType:      string(offer.Discount().Type),
		DiscountValue:     offer.Discount().Value,
		Status:            string(offer.Status()),
		ValidityStartDate: offer.Validity().StartDate,
		ValidityEndDate:   offer.Validity().EndDate,
		Location: GeoPoint{
			Lat: offer.EstablishmentSnapshot().Location.Latitude(),
			Lon: offer.EstablishmentSnapshot().Location.Longitude(),
		},
		PartnerName:       offer.PartnerSnapshot().Name,
		EstablishmentName: offer.EstablishmentSnapshot().Name,
		EstablishmentCity: offer.EstablishmentSnapshot().City,
		Views:             int64(offer.Stats().Views),
		Clicks:            int64(offer.Stats().Clicks),
		Bookings:          int64(offer.Stats().Bookings),
		Favorites:         int64(offer.Stats().Favorites),
		AvgRating:         offer.Stats().AvgRating,
		ReviewCount:       offer.Stats().ReviewCount,
		CreatedAt:         offer.CreatedAt(),
		UpdatedAt:         offer.UpdatedAt(),
	}

	// Optional fields
	if offer.Discount().OriginalPrice != nil {
		doc.OriginalPrice = *offer.Discount().OriginalPrice
	}
	if offer.Discount().Formula != "" {
		doc.Formula = offer.Discount().Formula
	}
	if offer.PublishedAt() != nil {
		doc.PublishedAt = *offer.PublishedAt()
	}

	return doc
}

// documentToSummary converts an Elasticsearch document to a domain summary.
func (r *OfferSearchRepository) documentToSummary(doc *offerDocument) *domain.OfferSummary {
	offerID, _ := domain.ParseOfferID(doc.ID)
	categoryID, _ := domain.ParseCategoryID(doc.CategoryID)

	// Build GeoLocation
	location, _ := domain.NewGeoLocation(doc.Location.Lon, doc.Location.Lat)

	// Build Discount
	var originalPrice *int64
	if doc.OriginalPrice > 0 {
		originalPrice = &doc.OriginalPrice
	}
	discount := domain.Discount{
		Type:          domain.DiscountType(doc.DiscountType),
		Value:         doc.DiscountValue,
		OriginalPrice: originalPrice,
		Formula:       doc.Formula,
	}

	return &domain.OfferSummary{
		ID:                offerID,
		Title:             doc.Title,
		ShortDescription:  doc.ShortDescription,
		Discount:          discount,
		PrimaryImage:      "", // Not stored in ES, would need separate lookup
		PartnerName:       doc.PartnerName,
		EstablishmentName: doc.EstablishmentName,
		City:              doc.EstablishmentCity,
		Location:          location,
		CategoryID:        categoryID,
	}
}

// BulkIndex indexes multiple offers in a single bulk request.
func (r *OfferSearchRepository) BulkIndex(ctx context.Context, offers []*domain.Offer) error {
	if len(offers) == 0 {
		return nil
	}

	var buf bytes.Buffer
	for _, offer := range offers {
		meta := map[string]interface{}{
			"index": map[string]interface{}{
				"_index": offersIndex,
				"_id":    offer.ID().String(),
			},
		}

		if err := json.NewEncoder(&buf).Encode(meta); err != nil {
			return fmt.Errorf("failed to encode meta: %w", err)
		}

		doc := r.offerToDocument(offer)
		if err := json.NewEncoder(&buf).Encode(doc); err != nil {
			return fmt.Errorf("failed to encode document: %w", err)
		}
	}

	res, err := r.client.Bulk(
		bytes.NewReader(buf.Bytes()),
		r.client.Bulk.WithContext(ctx),
		r.client.Bulk.WithIndex(offersIndex),
		r.client.Bulk.WithRefresh("true"),
	)
	if err != nil {
		return fmt.Errorf("failed to bulk index: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("bulk index error: %s", res.String())
	}

	return nil
}

// GetOffersByIDs retrieves offers by their IDs from Elasticsearch.
func (r *OfferSearchRepository) GetOffersByIDs(ctx context.Context, ids []domain.OfferID) ([]*domain.OfferSummary, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	idStrings := make([]string, len(ids))
	for i, id := range ids {
		idStrings[i] = id.String()
	}

	searchQuery := map[string]interface{}{
		"query": map[string]interface{}{
			"ids": map[string]interface{}{
				"values": idStrings,
			},
		},
		"size": len(ids),
	}

	data, err := json.Marshal(searchQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal query: %w", err)
	}

	res, err := r.client.Search(
		r.client.Search.WithContext(ctx),
		r.client.Search.WithIndex(offersIndex),
		r.client.Search.WithBody(bytes.NewReader(data)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to search: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("search error: %s", res.String())
	}

	var result searchResult
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	summaries := make([]*domain.OfferSummary, 0, len(result.Hits.Hits))
	for _, hit := range result.Hits.Hits {
		summary := r.documentToSummary(&hit.Source)
		summaries = append(summaries, summary)
	}

	return summaries, nil
}

// Refresh forces a refresh of the offers index.
func (r *OfferSearchRepository) Refresh(ctx context.Context) error {
	res, err := r.client.Indices.Refresh(
		r.client.Indices.Refresh.WithContext(ctx),
		r.client.Indices.Refresh.WithIndex(offersIndex),
	)
	if err != nil {
		return fmt.Errorf("failed to refresh index: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("refresh error: %s", res.String())
	}

	return nil
}

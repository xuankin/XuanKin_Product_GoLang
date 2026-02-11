package repository

import (
	"Product_Mangement_Api/models"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"log"
	"strings"
)

type ElasticsearchRepository interface {
	IndexProduct(ctx context.Context, product *models.EsProductIndex) error
	DeleteProduct(ctx context.Context, productID string) error
	SearchProducts(ctx context.Context, params models.FilterParams) ([]models.EsProductIndex, int64, error)
	CreateIndexIfNotExists(ctx context.Context) error
}

type esRepository struct {
	client *elasticsearch.Client
	index  string
}

func NewElasticsearchRepository(client *elasticsearch.Client) ElasticsearchRepository {
	return &esRepository{client: client, index: "products"}
}

// 1. Khởi tạo Index với Mapping chuẩn (Quan trọng)
func (r *esRepository) CreateIndexIfNotExists(ctx context.Context) error {
	exists, _ := r.client.Indices.Exists([]string{r.index})
	if exists.StatusCode == 200 {
		return nil
	}

	mapping := `{
		"settings": { "number_of_shards": 1, "number_of_replicas": 0 },
		"mappings": {
			"properties": {
				"name": { "properties": { "vi": { "type": "text" }, "en": { "type": "text" } } },
				"attributes_summary": { 
					"type": "nested",
					"properties": {
						"name": { "type": "keyword" },
						"values": { "type": "keyword" }
					}
				},
				"min_price": { "type": "double" },
				"max_price": { "type": "double" },
				"created_at": { "type": "date" }
			}
		}
	}`

	res, err := r.client.Indices.Create(r.index, r.client.Indices.Create.WithBody(strings.NewReader(mapping)))
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.IsError() {
		return fmt.Errorf("error creating index: %s", res.String())
	}
	log.Println("Elasticsearch Index 'products' created successfully")
	return nil
}

func (r *esRepository) IndexProduct(ctx context.Context, product *models.EsProductIndex) error {
	data, err := json.Marshal(product)
	if err != nil {
		return err
	}

	req := esapi.IndexRequest{
		Index:      r.index,
		DocumentID: product.ID.String(),
		Body:       bytes.NewReader(data),
		Refresh:    "true",
	}

	res, err := req.Do(ctx, r.client)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error indexing product: %s", res.String())
	}
	return nil
}

func (r *esRepository) DeleteProduct(ctx context.Context, productID string) error {
	req := esapi.DeleteRequest{
		Index:      r.index,
		DocumentID: productID,
	}
	res, err := req.Do(ctx, r.client)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return nil
}

func (r *esRepository) SearchProducts(ctx context.Context, params models.FilterParams) ([]models.EsProductIndex, int64, error) {
	var buf bytes.Buffer

	query := map[string]interface{}{
		"from": (params.Page - 1) * params.Limit,
		"size": params.Limit,
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []interface{}{
					map[string]interface{}{"match_all": map[string]interface{}{}},
				},
			},
		},
	}

	if params.Search != "" {
		query["query"] = map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  params.Search,
				"fields": []string{"name.vi", "name.en", "description.vi", "description.en", "slug"},
			},
		}
	}

	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, 0, err
	}

	res, err := r.client.Search(
		r.client.Search.WithContext(ctx),
		r.client.Search.WithIndex(r.index),
		r.client.Search.WithBody(&buf),
	)
	if err != nil {
		return nil, 0, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, 0, fmt.Errorf("search error: %s", res.String())
	}

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, 0, err
	}

	hits := result["hits"].(map[string]interface{})
	total := int64(hits["total"].(map[string]interface{})["value"].(float64))
	hitList := hits["hits"].([]interface{})

	var products []models.EsProductIndex
	for _, hit := range hitList {
		source := hit.(map[string]interface{})["_source"]
		sourceBytes, _ := json.Marshal(source)
		var p models.EsProductIndex
		json.Unmarshal(sourceBytes, &p)
		products = append(products, p)
	}

	return products, total, nil
}

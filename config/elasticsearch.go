package config

import (
	"github.com/elastic/go-elasticsearch/v7"
	"log"
)

func ConnectElasticsearch(cfg *Config) *elasticsearch.Client {
	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{cfg.ElasticsearchAddress},
	})
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	res, err := es.Info()
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()
	log.Println("Successfully connected to Elasticsearch")
	return es
}

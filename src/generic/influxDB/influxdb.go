package genericInfluxDB

import (
	"context"
	"fmt"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
)

type InfluxConfig struct {
	URL    string
	Token  string
	Org    string
	Bucket string
}

type InfluxClient struct {
	Client   influxdb2.Client
	WriteAPI api.WriteAPIBlocking
	QueryAPI api.QueryAPI
	Org      string
	Bucket   string
}

func NewInfluxClient(cfg InfluxConfig) (*InfluxClient, error) {
	client := influxdb2.NewClient(cfg.URL, cfg.Token)

	// nil yerine context.Background() kullan
	ok, err := client.Ping(context.Background())
	if err != nil || !ok {
		return nil, fmt.Errorf("influxdb bağlantısı kurulamadı: %w", err)
	}

	writeAPI := client.WriteAPIBlocking(cfg.Org, cfg.Bucket)
	queryAPI := client.QueryAPI(cfg.Org)

	return &InfluxClient{
		Client:   client,
		WriteAPI: writeAPI,
		QueryAPI: queryAPI,
		Org:      cfg.Org,
		Bucket:   cfg.Bucket,
	}, nil
}

func (ic *InfluxClient) Close() {
	ic.Client.Close()
}

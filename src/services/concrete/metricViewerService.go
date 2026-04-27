package servicesConcrete

import (
	"context"
	"errors"
	"fmt"
	"os"
	genericInfluxDB "server-watcher-app/src/generic/influxDB"
	modelsDTOs "server-watcher-app/src/models/dtos"
	servicesAbstarct "server-watcher-app/src/services/abstract"
)

type metricViewerService struct {
	influxClient *genericInfluxDB.InfluxClient
}

func NewMetricViewerService(influxClient *genericInfluxDB.InfluxClient) servicesAbstarct.MetricViewerService {
	return &metricViewerService{
		influxClient: influxClient,
	}
}

// func (s *metricViewerService) GetLogsFromDB(ctx context.Context, dto *modelsDTOs.ContainerLogRequestDTO) (*[]modelsDTOs.LogResponseDTO, error) {

// 	bucket := os.Getenv("FLUX_BUCKET_NAME")

// 	if bucket == "" {
// 		return nil, errors.New("Bucket name is not applied")
// 	}

// 	hostID := dto.HostID
// 	containerName := dto.ContainerName
// 	limit := dto.Limit

// 	query := fmt.Sprintf(`
//     	from(bucket: "%s")
//     	    |> range(start: -1h)
//     	    |> filter(fn: (r) => r["_measurement"] == "container_logs")
//     	    |> filter(fn: (r) => r["host_id"] == "%s")
//     	    |> filter(fn: (r) => r["container_name"] == "%s")
//     	    |> filter(fn: (r) => r["_field"] == "message")
//     	    |> limit(n: %d)
//     	    |> sort(columns: ["_time"], desc: true)
// 		`, bucket, hostID, containerName, limit)

// 	result, err := s.influxClient.QueryAPI.Query(ctx, query)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var logs []modelsDTOs.LogResponseDTO
// 	for result.Next() {
// 		logs = append(logs, modelsDTOs.LogResponseDTO{
// 			Time:    result.Record().Time(),
// 			Message: fmt.Sprintf("%v", result.Record().Value()),
// 			HostID:  hostID,
// 		})
// 	}

// 	return &logs, nil

// }

func (s *metricViewerService) GetLogsFromDB(ctx context.Context, dto *modelsDTOs.ContainerLogRequestDTO) (*modelsDTOs.LogResponseDTO, error) {
	bucket := os.Getenv("FLUX_BUCKET_NAME")
	if bucket == "" {
		return nil, errors.New("bucket name is not applied")
	}

	hostID := dto.HostID
	containerName := dto.ContainerName
	limit := dto.Limit

	query := fmt.Sprintf(`
        from(bucket: "%s")
            |> range(start: -1h)
            |> filter(fn: (r) => r["_measurement"] == "container_logs")
            |> filter(fn: (r) => r["host_id"] == "%s")
            |> filter(fn: (r) => r["container_name"] == "%s")
            |> filter(fn: (r) => r["_field"] == "message")
            |> sort(columns: ["_time"], desc: true)
            |> limit(n: %d)
        `, bucket, hostID, containerName, limit)

	result, err := s.influxClient.QueryAPI.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	// Tek bir response objesi oluşturuyoruz
	response := &modelsDTOs.LogResponseDTO{
		HostID:        hostID,
		ContainerName: containerName,
		Messages:      []string{}, // Boş dizi olarak başlatıyoruz
	}

	for result.Next() {
		// Her bir log satırını diziye ekliyoruz
		logLine := fmt.Sprintf("%v", result.Record().Value())
		response.Messages = append(response.Messages, logLine)
	}

	return response, nil
}

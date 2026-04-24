package servicesConcrete

import (
	"context"
	"server-watcher-app/src/models"
	repositoryAbstarct "server-watcher-app/src/repository/abstract"
	repositoryConcrete "server-watcher-app/src/repository/concrete"
	servicesAbstarct "server-watcher-app/src/services/abstract"
)

type containerStatService struct {
	statRepo repositoryConcrete.ContainerStatLogRepo
}

func NewContainerStatService(statRepo repositoryConcrete.ContainerStatLogRepo) servicesAbstarct.ContainerStatsService {
	return &containerStatService{
		statRepo: statRepo,
	}
}

func (s *containerStatService) GetStatsByName(ctx context.Context, containerName string) ([]models.ContainerStatLog, error) {

	stats, err := s.statRepo.GetStatsWithContainerName(ctx, containerName)

	if err != nil {
		return nil, err
	}

	return stats, nil

}

func (s *containerStatService) GetStatsByNameP(ctx context.Context, p repositoryAbstarct.Pagination) (*repositoryAbstarct.PaginatedResult[models.ContainerStatLog], error) {

	result, err := s.statRepo.GetWithPaginate(ctx, p)

	if err != nil {
		return nil, err
	}

	responses := make([]models.ContainerStatLog, len(result.Data))

	// for i, cargo := range result.Data {

	// 	responses[i] = dtos.CargoResponseDTO{
	// 		ID:           cargo.ID,
	// 		Name:         cargo.Name,
	// 		SupplierType: cargo.SupplierType,
	// 		IsShared:     cargo.IsShared,
	// 		FamilyID:     cargo.FamilyID,
	// 	}

	// }

	return &repositoryAbstarct.PaginatedResult[models.ContainerStatLog]{
		Data:       responses,
		Page:       result.Page,
		Limit:      result.Limit,
		Total:      result.Total,
		TotalPages: result.TotalPages,
	}, nil

}

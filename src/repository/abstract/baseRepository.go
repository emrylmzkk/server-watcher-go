package repositoryAbstarct

import (
	"context"
	"errors"
	"math"
	"regexp"
	"strings"

	"gorm.io/gorm"
)

type Pagination struct {
	Page      int    `json:"page" query:"page"`
	Limit     int    `json:"limit" query:"limit"`
	SortBy    string `json:"sort_by" query:"sort_by"`
	SortOrder string `json:"sort_order" query:"sort_order"` // asc | desc
}

type PaginatedResult[T any] struct {
	Data       []T   `json:"data"`
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

type BaseRepository[T any] interface {
	Create(ctx context.Context, entity *T) error
	GetByID(ctx context.Context, id int) (*T, error)
	GetAll(ctx context.Context) ([]T, error)
	GetAllPaginated(ctx context.Context, p Pagination) (*PaginatedResult[T], error)
	Update(ctx context.Context, entity *T) error
	Delete(ctx context.Context, id int) error

	// Service içinde özel query yazmak için
	Query(ctx context.Context) *gorm.DB // deleted olanları filtreler
	Raw(ctx context.Context) *gorm.DB   // filtrelemez

	GetByUserIdPaginated(ctx context.Context, userId uint, p Pagination) (*PaginatedResult[T], error)
	GetByPagination(ctx context.Context, p Pagination) (*PaginatedResult[T], error)
	FindPaginated(ctx context.Context, p Pagination, query interface{}, args ...interface{}) (*PaginatedResult[T], error)
	PaginateWithUserFilter(ctx context.Context, userId uint, p Pagination) (*PaginatedResult[T], error)
	GetWithUserFilter(ctx context.Context, userId uint, entityId uint) (*T, error)
}

type GormBaseRepository[T any] struct {
	db *gorm.DB
}

func NewBaseRepository[T any](db *gorm.DB) BaseRepository[T] {
	return &GormBaseRepository[T]{db: db}
}

func (r *GormBaseRepository[T]) Create(ctx context.Context, entity *T) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

func (r *GormBaseRepository[T]) GetByID(ctx context.Context, id int) (*T, error) {
	var entity T
	err := r.baseQuery(ctx).Where("id = ?", id).First(&entity).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *GormBaseRepository[T]) GetAll(ctx context.Context) ([]T, error) {
	var entities []T
	if err := r.baseQuery(ctx).Find(&entities).Error; err != nil {
		return nil, err
	}
	return entities, nil
}

func (r *GormBaseRepository[T]) GetAllPaginated(ctx context.Context, p Pagination) (*PaginatedResult[T], error) {
	return r.paginate(r.baseQuery(ctx), p)
}

func (r *GormBaseRepository[T]) GetByPagination(ctx context.Context, p Pagination) (*PaginatedResult[T], error) {
	return r.paginate(r.baseQuery(ctx), p)
}

func (r *GormBaseRepository[T]) GetByUserIdPaginated(ctx context.Context, userId uint, p Pagination) (*PaginatedResult[T], error) {
	return r.paginate(r.baseQuery(ctx).Where("user_id = ?", userId), p)
}

func (r *GormBaseRepository[T]) FindPaginated(ctx context.Context, p Pagination, query interface{}, args ...interface{}) (*PaginatedResult[T], error) {
	return r.paginate(r.baseQuery(ctx).Where(query, args...), p)
}

func (r *GormBaseRepository[T]) paginate(db *gorm.DB, p Pagination) (*PaginatedResult[T], error) {
	p = normalizePagination(p)

	var total int64
	// Count işlemi için Model set edilmeli ki doğru tablo sayısını alsın
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}

	var entities []T
	offset := (p.Page - 1) * p.Limit
	orderClause := p.SortBy + " " + p.SortOrder

	if err := db.Order(orderClause).
		Offset(offset).
		Limit(p.Limit).
		Find(&entities).Error; err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(total) / float64(p.Limit)))
	if totalPages == 0 {
		totalPages = 1
	}

	return &PaginatedResult[T]{
		Data:       entities,
		Page:       p.Page,
		Limit:      p.Limit,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}

func (r *GormBaseRepository[T]) PaginateWithUserFilter(ctx context.Context, userId uint, p Pagination) (*PaginatedResult[T], error) {

	var total int64

	var entities []T

	p = normalizePagination(p)

	base := r.Query(ctx).Where("user_id = ?", userId)

	err := base.Count(&total).Error

	if err != nil {
		return nil, err
	}

	offset := (p.Page - 1) * p.Limit

	order := p.SortBy + " " + p.SortOrder

	if err := base.Order(order).Offset(offset).Limit(p.Limit).Find(&entities).Error; err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(total) / float64(p.Limit)))

	if totalPages == 0 {
		totalPages = 1
	}

	return &PaginatedResult[T]{
		Data:       entities,
		Page:       p.Page,
		Limit:      p.Limit,
		Total:      total,
		TotalPages: totalPages,
	}, nil

}

func (r *GormBaseRepository[T]) Update(ctx context.Context, entity *T) error {
	return r.db.WithContext(ctx).Save(entity).Error
}

// func (r *GormBaseRepository[T]) SoftDelete(ctx context.Context, id uint) error {
// 	now := time.Now()

// 	// model_status + deleted_at birlikte set ediyoruz
// 	return r.db.WithContext(ctx).
// 		Model(new(T)).
// 		Where("id = ? AND model_status <> ?", id, enum.Deleted).
// 		Updates(map[string]interface{}{
// 			"model_status": enum.Deleted,
// 			"deleted_at":   now,
// 		}).Error
// }

func (r *GormBaseRepository[T]) Delete(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).
		Unscoped().
		Where("id = ?", id).
		Delete(new(T)).Error
}

func (r *GormBaseRepository[T]) Query(ctx context.Context) *gorm.DB {
	return r.baseQuery(ctx)
}

func (r *GormBaseRepository[T]) Raw(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx).Model(new(T))
}

func (r *GormBaseRepository[T]) baseQuery(ctx context.Context) *gorm.DB {
	// deleted kayıtları otomatik dışarıda bırak
	return r.db.WithContext(ctx).
		Model(new(T))
}

func normalizePagination(p Pagination) Pagination {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.Limit <= 0 {
		p.Limit = 10
	}
	if p.Limit > 100 {
		p.Limit = 100
	}

	if !isSafeColumn(p.SortBy) {
		p.SortBy = "created_at"
	}

	order := strings.ToLower(strings.TrimSpace(p.SortOrder))
	if order != "asc" && order != "desc" {
		order = "desc"
	}
	p.SortOrder = order

	return p
}

func isSafeColumn(s string) bool {
	// SQL injection riskini azaltmak için basit whitelist formatı
	matched, _ := regexp.MatchString(`^[a-zA-Z_][a-zA-Z0-9_]*$`, s)
	return matched
}

func (r *GormBaseRepository[T]) GetWithUserFilter(ctx context.Context, userId uint, entityId uint) (*T, error) {

	model := new(T)

	query := r.Query(ctx).
		Where("user_id = ? AND id = ?", userId, entityId).
		First(model)

	if errors.Is(query.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("resource not found")
	}

	if query.Error != nil {
		return nil, query.Error
	}

	return model, nil
}

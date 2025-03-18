package pgsql

import (
	"backoffice/internal/entities"
	e "backoffice/internal/errors"
	"backoffice/internal/repositories"
	"context"
	"errors"
	"gorm.io/gorm"
)

type BaseRepository[T any] struct {
	conn *gorm.DB
}

func NewBaseRepository[T any](conn *gorm.DB) repositories.BaseRepository[T] {
	return &BaseRepository[T]{
		conn: conn,
	}
}

func (r *BaseRepository[T]) DB() *gorm.DB {
	return r.conn
}

func (r *BaseRepository[T]) Find(ctx context.Context, conditions map[string]interface{}) (data []*T, err error) {
	query := r.conn.WithContext(ctx).Where(conditions)

	err = query.Find(&data).Error

	return data, err
}

func (r *BaseRepository[T]) FindLimit(ctx context.Context, conditions map[string]interface{}, limit, offset int) (data []*T, total int64, err error) {
	query := r.conn.WithContext(ctx).Where(conditions)

	query.Count(&total)

	err = query.Limit(limit).Offset(offset).Find(&data).Error

	return data, total, err
}

func (r *BaseRepository[T]) FindBy(ctx context.Context, params map[string]interface{}) (entity *T, err error) {
	err = r.conn.WithContext(ctx).Where(params).First(&entity).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, e.ErrEntityNotFound
	}

	return entity, err
}

func (r *BaseRepository[T]) FindByWith(ctx context.Context, params map[string]interface{}, references ...string) (entity *T, err error) {
	query := r.conn.WithContext(ctx).Where(params)

	for _, reference := range references {
		query = query.Preload(reference)
	}

	err = query.First(&entity).Error

	return entity, err
}

func (r *BaseRepository[T]) Create(ctx context.Context, entity *T) (*T, error) {
	// TODO: add returning
	if err := r.conn.WithContext(ctx).Create(entity).Error; err != nil {
		return nil, err
	}

	return entity, nil
}

func (r *BaseRepository[T]) CreateNoReturn(ctx context.Context, entity *T) error {
	return r.conn.WithContext(ctx).Create(&entity).Error
}

func (r *BaseRepository[T]) Delete(ctx context.Context, entity *T, conditions ...interface{}) error {
	return r.conn.WithContext(ctx).Model(&entity).Delete(&entity, conditions...).Error
}

func (r *BaseRepository[T]) Save(ctx context.Context, entity *T) (*T, error) {
	if err := r.conn.WithContext(ctx).Save(&entity).Error; err != nil {
		return nil, err
	}

	return entity, nil
}

func (r *BaseRepository[T]) Update(ctx context.Context, entity *T, values interface{}, conditions map[string]interface{}) (*T, error) {
	if err := r.conn.WithContext(ctx).Model(&entity).Where(conditions).Updates(values).Error; err != nil {
		return nil, err
	}

	return entity, nil
}

func (r *BaseRepository[T]) Paginate(ctx context.Context, filters map[string]interface{}, order string, limit int, page int) (pagination entities.Pagination[T], err error) {
	conn := r.conn.WithContext(ctx).Where(filters).Order(order)

	items := make([]*T, 0)

	var entity T

	var total int64

	if err = conn.Model(&entity).Count(&total).Error; err != nil {
		return pagination, err
	}

	conn = conn.
		Limit(limit).
		Offset(limit * (page - 1))

	if err = conn.Find(&items).Error; err != nil {
		return
	}

	pagination.Total = int(total)
	pagination.Limit = limit
	pagination.CurrentPage = page
	pagination.Items = items

	return
}

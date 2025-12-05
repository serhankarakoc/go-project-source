package repositories

import (
	"context"
	"errors"
	"strings"

	"zatrano/pkg/currentuser"
	"zatrano/pkg/queryparams"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	ErrNotFound      = errors.New("kayıt bulunamadı")
	ErrMissingUserID = errors.New("context içinde geçerli user_id yok")
)

type IBaseRepository[T any] interface {
	GetAll(ctx context.Context, params queryparams.ListParams, scopes ...func(*gorm.DB) *gorm.DB) ([]T, int64, error)
	GetByID(ctx context.Context, id uint) (*T, error)
	Create(ctx context.Context, entity *T) error
	CreateWithRelations(ctx context.Context, entity *T) error
	BulkCreate(ctx context.Context, entities []T) error
	BulkCreateWithRelations(ctx context.Context, entities []T) error
	Update(ctx context.Context, id uint, data map[string]interface{}) error
	UpdateWithRelations(ctx context.Context, entity *T) error
	BulkUpdate(ctx context.Context, condition map[string]interface{}, data map[string]interface{}) error
	BulkUpdateWithRelations(ctx context.Context, entities []T) error
	Delete(ctx context.Context, id uint) error
	DeleteWithRelations(ctx context.Context, id uint) error
	BulkDelete(ctx context.Context, condition map[string]interface{}) error
	BulkDeleteWithRelations(ctx context.Context, ids []uint) error
	GetCount(ctx context.Context) (int64, error)
	CountByCondition(ctx context.Context, condition map[string]interface{}) (int64, error)
}

type BaseRepository[T any] struct {
	db                 *gorm.DB
	allowedSortColumns map[string]bool
	preloads           []string
}

func NewBaseRepository[T any](db *gorm.DB) *BaseRepository[T] {
	return &BaseRepository[T]{
		db: db,
		allowedSortColumns: map[string]bool{
			"id":         true,
			"created_at": true,
		},
	}
}

func (r *BaseRepository[T]) SetAllowedSortColumns(columns []string) {
	r.allowedSortColumns = make(map[string]bool)
	for _, col := range columns {
		r.allowedSortColumns[col] = true
	}
}

func (r *BaseRepository[T]) SetPreloads(preloads ...string) {
	r.preloads = preloads
}

func (r *BaseRepository[T]) GetAll(ctx context.Context, params queryparams.ListParams, scopes ...func(*gorm.DB) *gorm.DB) ([]T, int64, error) {
	var results []T
	var totalCount int64
	var t T
	query := r.db.WithContext(ctx).Model(&t)

	for _, preload := range r.preloads {
		query = query.Preload(preload)
	}
	for _, s := range scopes {
		query = s(query)
	}

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}
	if totalCount == 0 {
		return results, 0, nil
	}

	sortBy := params.SortBy
	orderBy := strings.ToLower(params.OrderBy)
	if orderBy != "asc" && orderBy != "desc" {
		orderBy = queryparams.DefaultOrderBy
	}
	if _, ok := r.allowedSortColumns[sortBy]; !ok {
		sortBy = queryparams.DefaultSortBy
	}

	query = query.Order(sortBy + " " + orderBy).Limit(params.PerPage).Offset(params.CalculateOffset())
	err := query.Find(&results).Error
	return results, totalCount, err
}

func (r *BaseRepository[T]) GetByID(ctx context.Context, id uint) (*T, error) {
	var result T
	query := r.db.WithContext(ctx)
	for _, preload := range r.preloads {
		query = query.Preload(preload)
	}
	err := query.First(&result, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	return &result, err
}

// Create — context’ten userID alıp created_by alanını set eder
func (r *BaseRepository[T]) Create(ctx context.Context, entity *T) error {
	if uid, ok := ctx.Value(currentuser.ContextUserIDKey).(uint); ok && uid > 0 {
		if err := r.db.WithContext(ctx).Model(entity).Set("gorm:created_by", uid).Create(entity).Error; err != nil {
			return err
		}
		return nil
	}
	return r.db.WithContext(ctx).Create(entity).Error
}

func (r *BaseRepository[T]) CreateWithRelations(ctx context.Context, entity *T) error {
	return r.db.WithContext(ctx).Session(&gorm.Session{FullSaveAssociations: true}).Create(entity).Error
}

func (r *BaseRepository[T]) BulkCreate(ctx context.Context, entities []T) error {
	return r.db.WithContext(ctx).Create(&entities).Error
}

func (r *BaseRepository[T]) BulkCreateWithRelations(ctx context.Context, entities []T) error {
	tx := r.db.WithContext(ctx).Session(&gorm.Session{FullSaveAssociations: true})
	return tx.Create(&entities).Error
}

// Update — updated_by alanını context’ten otomatik set eder
func (r *BaseRepository[T]) Update(ctx context.Context, id uint, data map[string]interface{}) error {
	if uid, ok := ctx.Value(currentuser.ContextUserIDKey).(uint); ok && uid > 0 {
		data["updated_by"] = uid
	}
	var t T
	result := r.db.WithContext(ctx).Model(&t).Where("id = ?", id).Updates(data)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *BaseRepository[T]) UpdateWithRelations(ctx context.Context, entity *T) error {
	return r.db.WithContext(ctx).Session(&gorm.Session{FullSaveAssociations: true}).Save(entity).Error
}

func (r *BaseRepository[T]) BulkUpdate(ctx context.Context, condition map[string]interface{}, data map[string]interface{}) error {
	if uid, ok := ctx.Value(currentuser.ContextUserIDKey).(uint); ok && uid > 0 {
		data["updated_by"] = uid
	}
	var t T
	tx := r.db.WithContext(ctx).Model(&t).Where(condition).Updates(data)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *BaseRepository[T]) BulkUpdateWithRelations(ctx context.Context, entities []T) error {
	tx := r.db.WithContext(ctx).Session(&gorm.Session{FullSaveAssociations: true})
	for _, entity := range entities {
		if err := tx.Save(&entity).Error; err != nil {
			return err
		}
	}
	return nil
}

// Delete — deleted_by alanını context’ten userID alarak set eder
func (r *BaseRepository[T]) Delete(ctx context.Context, id uint) error {
	var entity T
	userID, ok := ctx.Value(currentuser.ContextUserIDKey).(uint)
	if !ok || userID == 0 {
		return ErrMissingUserID
	}

	tx := r.db.WithContext(ctx)
	if err := tx.First(&entity, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}
	if err := tx.Model(&entity).Update("deleted_by", userID).Error; err != nil {
		return err
	}
	return tx.Delete(&entity).Error
}

func (r *BaseRepository[T]) DeleteWithRelations(ctx context.Context, id uint) error {
	var entity T
	userID, ok := ctx.Value(currentuser.ContextUserIDKey).(uint)
	if !ok || userID == 0 {
		return ErrMissingUserID
	}
	tx := r.db.WithContext(ctx)
	if err := tx.Preload(clause.Associations).First(&entity, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}
	if err := tx.Model(&entity).Update("deleted_by", userID).Error; err != nil {
		return err
	}
	return tx.Select(clause.Associations).Delete(&entity).Error
}

func (r *BaseRepository[T]) BulkDelete(ctx context.Context, condition map[string]interface{}) error {
	var entities []T
	userID, ok := ctx.Value(currentuser.ContextUserIDKey).(uint)
	if !ok || userID == 0 {
		return ErrMissingUserID
	}
	tx := r.db.WithContext(ctx)
	if err := tx.Where(condition).Find(&entities).Error; err != nil {
		return err
	}
	if len(entities) == 0 {
		return ErrNotFound
	}
	for _, entity := range entities {
		if err := tx.Model(&entity).Update("deleted_by", userID).Error; err != nil {
			return err
		}
		if err := tx.Delete(&entity).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *BaseRepository[T]) BulkDeleteWithRelations(ctx context.Context, ids []uint) error {
	var entities []T
	userID, ok := ctx.Value(currentuser.ContextUserIDKey).(uint)
	if !ok || userID == 0 {
		return ErrMissingUserID
	}
	tx := r.db.WithContext(ctx)
	if err := tx.Preload(clause.Associations).Find(&entities, ids).Error; err != nil {
		return err
	}
	if len(entities) == 0 {
		return ErrNotFound
	}
	for _, entity := range entities {
		if err := tx.Model(&entity).Update("deleted_by", userID).Error; err != nil {
			return err
		}
	}
	return tx.Select(clause.Associations).Delete(&entities).Error
}

func (r *BaseRepository[T]) GetCount(ctx context.Context) (int64, error) {
	var totalCount int64
	var t T
	err := r.db.WithContext(ctx).Model(&t).Count(&totalCount).Error
	return totalCount, err
}

func (r *BaseRepository[T]) CountByCondition(ctx context.Context, condition map[string]interface{}) (int64, error) {
	var count int64
	var t T
	err := r.db.WithContext(ctx).Model(&t).Where(condition).Count(&count).Error
	return count, err
}

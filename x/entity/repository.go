package entity

import (
	"context"
	"github.com/totegamma/concurrent/x/core"
	"gorm.io/gorm"
	"time"
)

// Repository is the interface for host repository
type Repository interface {
    Get(ctx context.Context, key string) (core.Entity, error)
    Create(ctx context.Context, entity *core.Entity) error
    Upsert(ctx context.Context, entity *core.Entity) error
    GetList(ctx context.Context) ([]SafeEntity, error)
    ListModified(ctx context.Context, modified time.Time) ([]SafeEntity, error)
    Delete(ctx context.Context, key string) error
    Update(ctx context.Context, entity *core.Entity) error
}

type repository struct {
	db *gorm.DB
}

// NewRepository creates a new host repository
func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// Get returns a entity by key
func (r *repository) Get(ctx context.Context, key string) (core.Entity, error) {
	ctx, span := tracer.Start(ctx, "RepositoryGet")
	defer span.End()

	var entity core.Entity
	err := r.db.WithContext(ctx).First(&entity, "id = ?", key).Error
	return entity, err
}

// Create creates new entity
func (r *repository) Create(ctx context.Context, entity *core.Entity) error {
	ctx, span := tracer.Start(ctx, "RepositoryCreate")
	defer span.End()

	return r.db.WithContext(ctx).Create(&entity).Error
}

// Upsert updates a entity
func (r *repository) Upsert(ctx context.Context, entity *core.Entity) error {
	ctx, span := tracer.Start(ctx, "RepositoryUpsert")
	defer span.End()

	return r.db.WithContext(ctx).Save(&entity).Error
}

// GetList returns all entities
func (r *repository) GetList(ctx context.Context) ([]SafeEntity, error) {
	ctx, span := tracer.Start(ctx, "RepositoryGetList")
	defer span.End()

	var entities []SafeEntity
	err := r.db.WithContext(ctx).Model(&core.Entity{}).Where("host IS NULL or host = ''").Find(&entities).Error
	return entities, err
}

// ListModified returns all entities which modified after given time
func (r *repository) ListModified(ctx context.Context, time time.Time) ([]SafeEntity, error) {
	ctx, span := tracer.Start(ctx, "RepositoryListModified")
	defer span.End()

	var entities []SafeEntity
	err := r.db.WithContext(ctx).Model(&core.Entity{}).Where("m_date > ?", time).Find(&entities).Error
	return entities, err
}

// Delete deletes a entity
func (r *repository) Delete(ctx context.Context, id string) error {
	ctx, span := tracer.Start(ctx, "RepositoryDelete")
	defer span.End()

	return r.db.WithContext(ctx).Delete(&core.Entity{}, "id = ?", id).Error
}

// Update updates a entity
func (r *repository) Update(ctx context.Context, entity *core.Entity) error {
	ctx, span := tracer.Start(ctx, "RepositoryUpdate")
	defer span.End()

	return r.db.WithContext(ctx).Where("id = ?", entity.ID).Updates(&entity).Error
}

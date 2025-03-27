package repository

import (
	"errors"
	"gorm.io/gorm"
	"sky_ISService/utils"
	"sky_ISService/utils/database"
	"time"
)

type BaseRepository[T any] struct {
	db *gorm.DB
}

// NewBaseRepository 创建一个新的 BaseRepository
func NewBaseRepository[T any](db *gorm.DB) *BaseRepository[T] {
	return &BaseRepository[T]{db: db}
}

// BaseCreate 插入新记录
func (repo *BaseRepository[T]) BaseCreate(entity *T) error {
	return repo.db.Create(entity).Error
}

// BaseGetByID 根据 ID 查询（排除已软删除数据）
func (repo *BaseRepository[T]) BaseGetByID(id int) (*T, error) {
	var entity T
	if err := repo.db.Where("id = ? AND is_deleted = false", id).First(&entity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // 未找到
		}
		return nil, err
	}
	return &entity, nil
}

// BaseGetWithPagination 通用分页查询
func (repo *BaseRepository[T]) BaseGetWithPagination(page, size int, conditions map[string]interface{}, orderBy string) (*utils.Pagination, error) {
	// 初始化查询
	db := repo.db.Model(new(T)).Where("is_deleted = false")

	// 应用动态查询条件
	db = database.ApplyConditions(db, conditions)

	// 执行查询并获取符合条件的数据总数
	var total int64
	var users []T

	if len(conditions) > 0 {
		// 查询数据总数
		if err := db.Find(&users).Count(&total).Error; err != nil {
			return nil, err
		}
	} else {
		// 执行分页查询
		if err := db.Offset((page - 1) * size).Limit(size).Find(&users).Error; err != nil {
			return nil, err
		}
	}

	// 计算分页结果
	pagination := &utils.Pagination{
		Total: total,
		Data:  users,
	}

	return pagination, nil
}

// BaseUpdate 更新记录（仅更新非软删除的数据）
func (repo *BaseRepository[T]) BaseUpdate(entity *T, id int) error {
	return repo.db.Model(entity).Where("id = ? AND is_deleted = false", id).Updates(entity).Error
}

// BaseSoftDelete 软删除记录
func (repo *BaseRepository[T]) BaseSoftDelete(id int) error {
	return repo.db.Model(new(T)).Where("id = ? AND is_deleted = false", id).
		Updates(map[string]interface{}{
			"is_deleted": true,
			"updated_at": time.Now(),
		}).Error
}

// BaseFetchAll 查询所有未删除的记录
func (repo *BaseRepository[T]) BaseFetchAll(orderBy string) ([]T, error) {
	var entities []T
	err := repo.db.Where("is_deleted = false").Order(orderBy).Find(&entities).Error
	if err != nil {
		return nil, err
	}
	return entities, nil
}

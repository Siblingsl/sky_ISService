package utils

import (
	"math"
	"reflect"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Pagination 结构体
type Pagination struct {
	Page       int         `json:"page"`        // 当前页码
	Limit      int         `json:"limit"`       // 每页数量
	Total      int64       `json:"total"`       // 数据总数
	TotalPages int         `json:"total_pages"` // 总页数
	Data       interface{} `json:"data"`        // 结果集
}

// PaginateFunc 分页函数签名
type PaginateFunc func(db *gorm.DB) *gorm.DB

// NewPagination 解析请求中的分页参数
func NewPagination(c *gin.Context) *Pagination {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	return &Pagination{
		Page:  page,
		Limit: limit,
	}
}

// Paginate GORM 分页封装
func Paginate(c *gin.Context, db *gorm.DB, model interface{}, filters ...PaginateFunc) (*Pagination, error) {
	p := NewPagination(c)

	// 应用可选的过滤器
	for _, filter := range filters {
		db = filter(db)
	}

	var total int64
	if err := db.Model(model).Count(&total).Error; err != nil {
		return nil, err
	}

	// 计算分页偏移量
	offset := (p.Page - 1) * p.Limit

	// 确保 `model` 是 slice 指针
	slicePtr := reflect.New(reflect.SliceOf(reflect.TypeOf(model).Elem())).Interface()

	if err := db.Limit(p.Limit).Offset(offset).Find(slicePtr).Error; err != nil {
		return nil, err
	}

	p.Total = total
	p.TotalPages = int(math.Ceil(float64(total) / float64(p.Limit)))
	p.Data = slicePtr

	return p, nil
}

// ResponseWithPagination 返回分页结果
func ResponseWithPagination(c *gin.Context, data *Pagination) {
	Success(c, data)
}

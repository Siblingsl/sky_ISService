package service

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"sky_ISService/proto/system"
	"sky_ISService/services/system/dto"
	"sky_ISService/services/system/repository"
	"sky_ISService/services/system/repository/models"
	"sky_ISService/shared/mq"
	"sky_ISService/utils"
	"time"
)

type AdminsService struct {
	adminsRepository *repository.AdminsRepository
	rabbitClient     *mq.RabbitMQClient
	system.UnimplementedSystemServiceServer
}

func NewUserService(adminsRepository *repository.AdminsRepository, rabbitClient *mq.RabbitMQClient) *AdminsService {
	return &AdminsService{
		adminsRepository: adminsRepository,
		rabbitClient:     rabbitClient,
	}
}

// CreateAdmin 添加管理员
func (s *AdminsService) CreateAdmin(ctx context.Context, req dto.CreateAdminsRequest) (*models.SkySystemAdmins, error) {
	// 1. 查询用户名是否已存在
	exists, err := s.adminsRepository.IsUsernameExists(req.Username)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("当前用户已存在")
	}

	admin := &models.SkySystemAdmins{
		Username:  req.Username,
		Password:  req.Password, // 这里需要 加密 处理
		FullName:  req.FullName,
		Email:     req.Email,
		Phone:     req.Phone,
		Notes:     req.Notes,
		Status:    req.Status,
		CreatedBy: req.CreatedBy,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	// 调用仓库层存储管理员信息
	if err := s.adminsRepository.PostCreateAdmin(admin); err != nil {
		return nil, err
	}

	// 发布消息到 MQ
	message := []byte(`{"username": "` + req.Username + `", "full_name": "` + req.FullName + `"}`)
	if err := s.rabbitClient.SendMessage("admin_created_queue", string(message)); err != nil {
		return nil, errors.New("管理员创建成功，但发布消息失败: " + err.Error())
	}

	// 可以把 admin 缓存到 redis 中
	return admin, nil
}

// GetAdminsByID 查询单个管理员用户
func (s *AdminsService) GetAdminsByID(id int) (*models.SkySystemAdmins, error) {
	admin, err := s.adminsRepository.GetAdminByID(id) // 调整参数
	if err != nil {
		return nil, err
	}
	return admin, nil
}

// GetUsersWithPagination 获取分页的管理员用户
func (s *AdminsService) GetUsersWithPagination(ctx *gin.Context, page int, size int, conditions map[string]interface{}) (*utils.Pagination, error) {
	// 调用 repository 层的分页查询并传递关键字
	pagination, err := s.adminsRepository.GetUsersWithPagination(ctx, page, size, conditions)
	if err != nil {
		return nil, err
	}
	return pagination, nil
}

// VerifyIsSystemAdmin 方法实现 (auth 子服务调用，不要动)
func (s *AdminsService) VerifyIsSystemAdmin(ctx context.Context, req *system.VerifyIsSystemAdminRequest) (*system.VerifyIsSystemAdminResponse, error) {
	// 这里实现你的业务逻辑
	return &system.VerifyIsSystemAdminResponse{
		IsAdmin: true, // 示例返回值
	}, nil
}

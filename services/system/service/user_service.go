package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sky_ISService/proto/system"
	"sky_ISService/services/system/dto"
	"sky_ISService/services/system/repository"
	"sky_ISService/services/system/repository/models"
	"sky_ISService/shared/mq"
	"sky_ISService/utils"
	"sky_ISService/utils/database"
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
func (s *AdminsService) CreateAdmin(req dto.CreateAdminsRequest) (*dto.SkySystemAdminsResponse, error) {
	// 1. 查询用户名是否已存在
	exists, err := s.adminsRepository.IsUsernameExists(req.Username)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("当前用户已存在")
	}

	if req.UserType == "00" {
		return nil, errors.New("无法创建顶级管理员账号")
	}

	admin := &models.SkySystemAdmins{
		Username: req.Username,
		Password: req.Password, // 这里需要 加密 处理
		FullName: req.FullName,
		UserType: req.UserType,
		Email:    req.Email,
		Phone:    req.Phone,
		CommonBase: database.CommonBase{
			Status:    req.Status,
			CreatedBy: req.CreatedBy,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Notes:     req.Notes,
		},
	}
	// 调用仓库层存储管理员信息
	if err := s.adminsRepository.BaseCreate(admin); err != nil {
		return nil, err
	}

	// 处理角色绑定
	roleIDs := req.RoleIDs
	// 过滤掉超级管理员角色 ID = 1
	var filteredRoleIDs []int
	for _, roleID := range roleIDs {
		if roleID == 1 {
			return nil, errors.New("不可以添加超级管理员角色")
		}
		filteredRoleIDs = append(filteredRoleIDs, roleID)
	}
	// 如果没有提供角色，自动分配默认角色（默认角色 ID 为 5 (普通用户)）
	if len(filteredRoleIDs) == 0 {
		defaultRoleID := int(5)
		filteredRoleIDs = append(filteredRoleIDs, defaultRoleID)
	}
	// 绑定角色
	err = s.BindRoles(admin.ID, filteredRoleIDs)
	if err != nil {
		return nil, errors.New("账号创建成功，但角色绑定失败: " + err.Error())
	}

	// 发布消息到 MQ
	messageData := map[string]string{
		"username":  req.Username,
		"full_name": req.FullName,
	}
	message, _ := json.Marshal(messageData)
	if err := s.rabbitClient.SendMessage("admin_created_queue", string(message)); err != nil {
		return nil, errors.New("管理员创建成功，但发布消息失败: " + err.Error())
	}

	adminResponse := &dto.SkySystemAdminsResponse{
		ID:        admin.ID,
		Username:  admin.Username,
		FullName:  admin.FullName,
		UserType:  admin.UserType,
		Email:     admin.Email,
		Phone:     admin.Phone,
		Status:    admin.Status,
		CreatedBy: admin.CreatedBy,
		UpdatedBy: admin.UpdatedBy,
		CreatedAt: admin.CreatedAt,
		UpdatedAt: admin.UpdatedAt,
		Notes:     admin.Notes,
		RolesID:   filteredRoleIDs,
	}

	// 可以把 admin 缓存到 redis 中
	return adminResponse, nil
}

// GetAdminsByID 查询单个管理员用户
func (s *AdminsService) GetAdminsByID(id int) (*dto.SkySystemAdminsResponse, error) {
	admin, err := s.adminsRepository.BaseGetByID(id)
	if err != nil {
		return nil, err
	}
	// 获取管理员角色信息
	roleIDs, err := s.adminsRepository.GetRoleIDsByAdminID(id)
	if err != nil {
		return nil, err
	}
	adminResponse := &dto.SkySystemAdminsResponse{
		ID:        admin.ID,
		Username:  admin.Username,
		Password:  admin.Password,
		FullName:  admin.FullName,
		UserType:  admin.UserType,
		Email:     admin.Email,
		Phone:     admin.Phone,
		Status:    admin.Status,
		CreatedBy: admin.CreatedBy,
		UpdatedBy: admin.UpdatedBy,
		CreatedAt: admin.CreatedAt,
		UpdatedAt: admin.UpdatedAt,
		Notes:     admin.Notes,
		RolesID:   roleIDs,
	}
	return adminResponse, nil
}

// GetUsersWithPagination 获取分页的管理员用户
func (s *AdminsService) GetUsersWithPagination(page int, size int, conditions map[string]interface{}) (*utils.Pagination, error) {
	return s.adminsRepository.BaseGetWithPagination(page, size, conditions, "id DESC")
}

// UpdateAdmin 修改管理员
func (s *AdminsService) UpdateAdmin(req dto.UpdateAdminsRequest) (*dto.SkySystemAdminsResponse, error) {
	// 检查管理员是否存在
	admin, err := s.adminsRepository.BaseGetByID(int(req.ID))
	if err != nil {
		return nil, err
	}

	// 验证是否为顶级管理员
	if admin.UserType == "00" {
		return nil, errors.New("无法修改顶级管理员账号")
	}

	// 更新管理员信息
	if req.Username != "" {
		admin.Username = req.Username
	}
	if req.Password != "" {
		admin.Password = req.Password // 这里应该对密码进行加密存储
	}
	if req.FullName != "" {
		admin.FullName = req.FullName
	}
	if req.UserType != "" {
		admin.UserType = req.UserType
	}
	if req.Email != "" {
		admin.Email = req.Email
	}
	if req.Phone != "" {
		admin.Phone = req.Phone
	}
	if req.Notes != "" {
		admin.Notes = req.Notes
	}
	if req.Status != false {
		admin.Status = req.Status
	}
	if req.UpdatedBy != 0 {
		admin.UpdatedBy = req.UpdatedBy
	}
	if req.UpdatedAt != "" {
		admin.UpdatedAt = time.Now()
	}

	err = s.adminsRepository.BaseUpdate(admin, int(req.ID))
	if err != nil {
		return nil, err
	}

	// 获取管理员当前角色
	currentRoleIDs, err := s.adminsRepository.GetRoleIDsByAdminID(int(req.ID))
	if err != nil {
		return nil, err
	}

	// 如果请求中包含角色信息，处理角色变动
	if len(req.RoleIDs) > 0 {
		// 删除不在新角色列表中的角色
		for _, roleID := range currentRoleIDs {
			if !utils.Contains(req.RoleIDs, roleID) {
				err := s.adminsRepository.RemoveAdminRole(req.ID, roleID)
				if err != nil {
					return nil, err
				}
			}
		}

		// 添加新角色
		for _, roleID := range req.RoleIDs {
			if !utils.Contains(currentRoleIDs, roleID) {
				err := s.adminsRepository.AddAdminRole(req.ID, roleID)
				if err != nil {
					return nil, err
				}
			}
		}
	} else {
		// 如果没有角色，默认绑定角色id = 5
		err := s.adminsRepository.AddAdminRole(req.ID, 5)
		if err != nil {
			return nil, err
		}
	}

	// 获取更新后的角色列表
	updatedRoleIDs, err := s.adminsRepository.GetRoleIDsByAdminID(int(req.ID))
	if err != nil {
		return nil, err
	}

	// 构建返回的响应数据
	adminResponse := &dto.SkySystemAdminsResponse{
		ID:        admin.ID,
		Username:  admin.Username,
		FullName:  admin.FullName,
		UserType:  admin.UserType,
		Email:     admin.Email,
		Phone:     admin.Phone,
		Status:    admin.Status,
		CreatedBy: admin.CreatedBy,
		UpdatedBy: admin.UpdatedBy,
		CreatedAt: admin.CreatedAt,
		UpdatedAt: admin.UpdatedAt,
		Notes:     admin.Notes,
		RolesID:   updatedRoleIDs,
	}

	return adminResponse, nil
}

// DeleteAdminByID 删除管理员
func (s *AdminsService) DeleteAdminByID(id int) (*models.SkySystemAdmins, error) {
	// 获取管理员信息，检查是否为顶级管理员
	admin, err := s.adminsRepository.BaseGetByID(id)
	if err != nil {
		return nil, fmt.Errorf("管理员不存在: %v", err)
	}
	// 如果是顶级管理员，则不能解绑角色
	if admin.UserType == "00" {
		return nil, errors.New("无法删除顶级管理员账号")
	}

	// 获取管理员绑定的角色
	roleIDs, err := s.adminsRepository.GetRoleIDsByAdminID(id)
	if err != nil {
		return nil, fmt.Errorf("获取管理员绑定的角色失败: %v", err)
	}

	// 执行角色解绑操作
	if len(roleIDs) > 0 {
		err = s.UnbindRoles(int(id), roleIDs)
		if err != nil {
			return nil, fmt.Errorf("删除角色与管理员绑定失败: %v", err)
		}
	}

	// 执行软删除
	if err := s.adminsRepository.BaseSoftDelete(id); err != nil {
		return nil, err
	}

	return admin, nil
}

// BindRoles 绑定角色
func (s *AdminsService) BindRoles(adminID int, roleIDs []int) error {
	if len(roleIDs) == 0 {
		return errors.New("角色 ID 不能为空")
	}
	// 先清空已有的角色绑定
	err := s.adminsRepository.AdminBindRoles(adminID)
	if err != nil {
		return err
	}
	// 重新插入新的角色绑定关系
	var adminRoles []models.AdminsRoles
	for _, roleID := range roleIDs {
		adminRoles = append(adminRoles, models.AdminsRoles{AdminID: adminID, RoleID: roleID})
	}
	// 避免空插入
	if len(adminRoles) == 0 {
		return nil
	}
	return s.adminsRepository.CreateAdminRoles(adminRoles)
}

// UnbindRoles 解绑角色
func (s *AdminsService) UnbindRoles(adminID int, roleIDs []int) error {
	// 获取管理员信息，检查是否为顶级管理员
	admin, err := s.adminsRepository.BaseGetByID(int(adminID))
	if err != nil {
		return err
	}
	// 如果是顶级管理员，则不能解绑角色
	if admin.UserType == "00" {
		return errors.New("无法解绑顶级管理员的角色")
	}
	// 遍历待解绑的角色
	for _, roleID := range roleIDs {
		exists, err := s.adminsRepository.CheckAdminRoleExist(adminID, roleID)
		if err != nil {
			return err
		}
		if !exists {
			continue // 如果角色未绑定，则跳过
		}
		// 解绑角色
		err = s.adminsRepository.RemoveAdminRole(adminID, roleID)
		if err != nil {
			return err
		}
	}
	// 解绑后检查剩余的角色
	remainingRoles, err := s.adminsRepository.GetRoleIDsByAdminID(int(adminID))
	if err != nil {
		return err
	}
	// 如果解绑后没有任何角色，则软删除该管理员
	if len(remainingRoles) == 0 {
		if err := s.adminsRepository.BaseSoftDelete(int(adminID)); err != nil {
			return err
		}
	}
	// 如果一切顺利，返回 nil
	return nil
}

// VerifyIsSystemAdmin 方法实现 (auth 子服务调用，不要动)
func (s *AdminsService) VerifyIsSystemAdmin(ctx context.Context, req *system.VerifyIsSystemAdminRequest) (*system.VerifyIsSystemAdminResponse, error) {
	// 这里实现你的业务逻辑
	return &system.VerifyIsSystemAdminResponse{
		IsAdmin: true, // 示例返回值
	}, nil
}

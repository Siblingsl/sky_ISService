package service

import (
	"errors"
	"fmt"
	"sky_ISService/services/system/dto"
	"sky_ISService/services/system/repository"
	"sky_ISService/services/system/repository/models"
	"sky_ISService/utils"
	"sky_ISService/utils/database"
	"time"
)

type RoleService struct {
	roleRepository *repository.RoleRepository
}

func NewRoleService(roleRepository *repository.RoleRepository) *RoleService {
	return &RoleService{roleRepository: roleRepository}
}

// CreateRole 添加角色
func (s *RoleService) CreateRole(req dto.CreateSkySystemRoleRequest) (*models.SkySystemRoles, error) {
	// 查询角色是否已存在
	isRole, err := s.roleRepository.IsRoleNameExists(req.RoleName)
	if err != nil {
		return nil, err
	}
	if isRole {
		return nil, errors.New("当前角色已存在")
	}

	role := &models.SkySystemRoles{
		RoleName:    req.RoleName,
		RoleKey:     req.RoleKey,
		RoleSort:    req.RoleSort,
		Description: req.Description,
		CommonBase: database.CommonBase{
			Status:    req.Status,
			CreatedBy: int(req.CreatedBy),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Notes:     req.Notes,
		},
	}
	if err := s.roleRepository.BaseCreate(role); err != nil {
		return nil, err
	}

	return role, nil
}

// GetRoleByID 查询单个角色
func (s *RoleService) GetRoleByID(id int) (*models.SkySystemRoles, error) {
	role, err := s.roleRepository.BaseGetByID(id)
	if err != nil {
		return nil, err
	}
	return role, nil
}

// GetRolesWithPagination 获取分页的全部角色
func (s *RoleService) GetRolesWithPagination(page int, size int, conditions map[string]interface{}) (*utils.Pagination, error) {
	return s.roleRepository.BaseGetWithPagination(page, size, conditions, "id DESC")
}

// UpdateRole 修改角色
func (s *RoleService) UpdateRole(req dto.UpdateSkySystemRoleRequest) (*models.SkySystemRoles, error) {
	// 检查角色是否存在
	role, err := s.GetRoleByID(int(req.ID))
	if err != nil {
		return nil, err
	}

	// 更新角色信息
	if req.RoleName != "" {
		role.RoleName = req.RoleName
	}
	if req.RoleKey != "" {
		role.RoleKey = req.RoleKey
	}
	if req.RoleSort != 0 {
		role.RoleSort = req.RoleSort
	}
	if req.Description != "" {
		role.Description = req.Description
	}
	// 更新公共字段
	role.CommonBase.Status = req.Status
	role.CommonBase.UpdatedBy = int(req.UpdatedBy)
	role.CommonBase.UpdatedAt = time.Now()

	err = s.roleRepository.BaseUpdate(role, int(req.ID))
	if err != nil {
		return nil, err
	}

	return role, nil
}

// DeleteRoleByID 删除角色
func (s *RoleService) DeleteRoleByID(id int) (*models.SkySystemRoles, error) {
	role, err := s.roleRepository.BaseGetByID(id)
	if err != nil {
		return nil, fmt.Errorf("管理员不存在: %v", err)
	}

	// 执行软删除
	if err := s.roleRepository.BaseSoftDelete(id); err != nil {
		return nil, err
	}

	return role, nil
}

// AssignMenusToRole 给角色分配可以打开菜单或者查看某些菜单中的部分数据还有可以读或写的权限
// AssignMenusToRole 给角色分配可以查看的菜单
func (s *RoleService) AssignMenusToRole(roleID int, menuIDs []int) (*models.SkySystemRoles, error) {
	// 检查传入的菜单列表是否为空
	if len(menuIDs) == 0 {
		return nil, errors.New("分配权限失败: 菜单ID列表为空")
	}

	// 获取角色信息
	role, err := s.roleRepository.BaseGetByID(int(roleID))
	if err != nil {
		return nil, err
	}

	// 删除角色当前所有的菜单权限
	if err := s.roleRepository.RemoveOldMenusFromRole(roleID); err != nil {
		return nil, err
	}

	// 批量为角色分配新的菜单
	if err := s.roleRepository.AssignMenusToRole(roleID, menuIDs); err != nil {
		return nil, err
	}

	return role, nil
}

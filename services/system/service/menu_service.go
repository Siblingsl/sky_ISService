package service

import (
	"errors"
	"fmt"
	"sky_ISService/services/system/dto"
	"sky_ISService/services/system/repository"
	"sky_ISService/services/system/repository/models"
	"sky_ISService/utils/database"
	"time"
)

type MenuService struct {
	menuRepository *repository.MenuRepository
}

func NewMenuService(menuRepository *repository.MenuRepository) *MenuService {
	return &MenuService{menuRepository: menuRepository}
}

func (s *MenuService) CreateMenu(req dto.CreateSkySystemMenuRequest) (*models.SkySystemMenus, error) {
	// 查寻菜单是否已经存在
	isMenu, err := s.menuRepository.IsMenuExist(req.MenuName)
	if err != nil {
		return nil, err
	}
	if isMenu {
		return nil, errors.New("当前菜单已存在")
	}

	menu := &models.SkySystemMenus{
		MenuName:    req.MenuName,
		MenuURL:     req.MenuURL,
		MenuType:    req.MenuType,
		ParentID:    req.ParentID,
		MenuSort:    req.MenuSort,
		MenuIcon:    req.MenuIcon,
		Description: req.Description,
		CommonBase: database.CommonBase{
			Status:    req.Status,
			CreatedBy: req.CreatedBy,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Notes:     req.Notes,
		},
	}
	if err := s.menuRepository.BaseCreate(menu); err != nil {
		return nil, err
	}

	return menu, nil
}

// GetMenuList 获取完整菜单列表
func (s *MenuService) GetMenuList() ([]dto.MenuItem, error) {
	return s.GetMenuTree()
}

// UpdateMenu 修改菜单
func (s *MenuService) UpdateMenu(req dto.UpdateSkySystemMenuRequest) (*models.SkySystemMenus, error) {
	// 检查菜单是否存在
	menu, err := s.menuRepository.BaseGetByID(int(req.ID))
	if err != nil {
		return nil, err
	}

	// 更新菜单信息（仅在请求参数非空时更新）
	if req.MenuName != "" {
		menu.MenuName = req.MenuName
	}
	if req.MenuURL != "" {
		menu.MenuURL = req.MenuURL
	}
	if req.MenuType != 0 {
		menu.MenuType = req.MenuType
	}
	if req.MenuIcon != "" {
		menu.MenuIcon = req.MenuIcon
	}
	if req.Description != "" {
		menu.Description = req.Description
	}
	if req.MenuSort != 0 {
		menu.MenuSort = req.MenuSort
	}
	if req.ParentID != 0 {
		menu.ParentID = req.ParentID
	}
	if req.Notes != "" {
		menu.Notes = req.Notes
	}
	// 更新公共字段
	menu.CommonBase.Status = req.Status
	menu.CommonBase.UpdatedBy = req.UpdatedBy
	menu.CommonBase.UpdatedAt = time.Now()

	if err := s.menuRepository.BaseUpdate(menu, int(req.ID)); err != nil {
		return nil, err
	}

	return menu, nil
}

// DeleteMenuByID 软删除菜单
func (s *MenuService) DeleteMenuByID(id int) (*models.SkySystemMenus, error) {
	// 获取菜单
	menu, err := s.menuRepository.BaseGetByID(id)
	if err != nil {
		return nil, fmt.Errorf("菜单不存在: %v", err)
	}

	// 判断是否可以删除
	if menu.ParentID == 0 || menu.ParentID == 1 {
		return nil, errors.New("无法删除当前目录或菜单")
	}

	// 软删除菜单
	if err := s.menuRepository.BaseSoftDelete(id); err != nil {
		return nil, fmt.Errorf("删除菜单失败: %v", err)
	}

	return menu, nil
}

// GetMenuTree 获取完整的菜单树
func (s *MenuService) GetMenuTree() ([]dto.MenuItem, error) {
	var menus []models.SkySystemMenus
	// 查询所有菜单，按父菜单分组
	menus, err := s.menuRepository.FetchAllMenus()
	if err != nil {
		return nil, err // 如果有错误，返回错误
	}

	// 构建菜单树
	menuTree := buildMenuTree(menus, 0)
	return menuTree, nil
}

// GetRoleMenusTreeByRoleId 获取根据角色得到的菜单树
func (s *MenuService) GetRoleMenusTreeByRoleId(roleID int) ([]dto.MenuItem, error) {
	// 查询角色所拥有的菜单权限
	var roleMenus []models.RolesMenus
	err := s.menuRepository.FetchMenusByRole(roleID, &roleMenus)
	if err != nil {
		return nil, err // 如果有错误，返回错误
	}
	// 获取所有菜单
	var menus []models.SkySystemMenus
	menus, err = s.menuRepository.FetchAllMenus()
	if err != nil {
		return nil, err
	}
	// 根据角色菜单权限过滤菜单
	var filteredMenus []models.SkySystemMenus
	for _, menu := range menus {
		for _, roleMenu := range roleMenus {
			if roleMenu.MenuID == menu.ID {
				filteredMenus = append(filteredMenus, menu)
				break
			}
		}
	}
	// 检查过滤后的菜单是否为空
	if len(filteredMenus) == 0 {
		return nil, errors.New("该角色没有分配任何菜单")
	}
	fmt.Println("filteredMenus", filteredMenus)

	// 构建菜单树
	menuTree := buildMenuTree(filteredMenus, roleID)
	fmt.Println("menuTree:", menuTree) // 打印菜单树以调试
	return menuTree, nil
}

// buildMenuTree 根据父菜单 ID 构建菜单树
func buildMenuTree(menus []models.SkySystemMenus, parentIDAndRoleID int) []dto.MenuItem {
	var result []dto.MenuItem
	for _, menu := range menus {
		if menu.ParentID == parentIDAndRoleID {
			// 查找子菜单
			children := buildMenuTree(menus, menu.ID)
			menuItem := dto.MenuItem{
				ID:       menu.ID,
				MenuName: menu.MenuName,
				MenuURL:  menu.MenuURL,
				ParentID: menu.ParentID,
				MenuSort: menu.MenuSort,
				MenuType: menu.MenuType,
				MenuIcon: menu.MenuIcon,
				Children: children,
			}
			result = append(result, menuItem)
		}
	}
	// 打印构建的菜单树，检查树是否正确
	fmt.Println("menuTree:", result)
	return result
}

package seeder

import (
	"gotribe/internal/model"

	"gorm.io/gorm"
)

// MenuSeeder 菜单种子
type MenuSeeder struct {
	*BaseSeeder
}

// NewMenuSeeder 创建菜单种子
func NewMenuSeeder() *MenuSeeder {
	return &MenuSeeder{
		BaseSeeder: NewBaseSeeder("menu"),
	}
}

// Run 执行菜单数据种子
func (s *MenuSeeder) Run(db *gorm.DB, syncExisting bool) error {
	// 获取角色
	var adminRole model.Role
	if err := db.First(&adminRole, 1).Error; err != nil {
		return err
	}

	var id0 int64 = 0
	var id1 int64 = 1
	var id8 int64 = 8
	var id10 int64 = 10
	var id18 int64 = 18
	componentStr := "component"
	systemAdminStr := "/system/admin"
	tableOfContents := "TableOfContents"
	briefcaseBusiness := "BriefcaseBusiness"
	alignHorizontalJustifyEnd := "AlignHorizontalJustifyEnd"
	lucideUsers := "LucideUsers"
	lucideUserRound := "LucideUserRound"
	menuSquare := "MenuSquare"
	network := "Network"
	lucideTableProperties := "LucideTableProperties"
	lucideUserCog2 := "LucideUserCog2"
	clipboardEdit := "ClipboardEdit"
	laptopMinimalCheck := "LaptopMinimalCheck"
	lucideHeadset := "LucideHeadset"
	boxes := "Boxes"
	wholeWord := "WholeWord"
	database := "Database"
	bookImage := "BookImage"
	lucideColumnsSettings := "LucideColumnsSettings"
	lucideTags := "LucideTags"
	lucideMartini := "LucideMartini"
	lucideDatabaseZap := "LucideDatabaseZap"
	lucideMessageSquareCode := "LucideMessageSquareCode"
	lucideTableConfig := "LucideTableConfig"

	menus := []model.Menu{
		{
			Model:     model.Model{ID: 34},
			Name:      "Dashboard",
			Title:     "控制台",
			Icon:      &alignHorizontalJustifyEnd,
			Path:      "/dashboard",
			Component: "Layout",
			Sort:      1,
			ParentID:  &id0,
			Roles:     []*model.Role{&adminRole},
			Creator:   "系统",
		},
		{
			Model:     model.Model{ID: 1},
			Name:      "System",
			Title:     "系统管理",
			Icon:      &componentStr,
			Path:      "/system",
			Component: "Layout",
			Redirect:  &systemAdminStr,
			Sort:      2,
			ParentID:  &id0,
			Roles:     []*model.Role{&adminRole},
			Creator:   "系统",
		},
		{
			Model:     model.Model{ID: 2},
			Name:      "Admin",
			Title:     "管理员",
			Icon:      &lucideUsers,
			Path:      "admin",
			Component: "/system/admin/index",
			Sort:      11,
			ParentID:  &id1,
			Roles:     []*model.Role{&adminRole},
			Creator:   "系统",
		},
		{
			Model:     model.Model{ID: 3},
			Name:      "Role",
			Title:     "角色权限",
			Icon:      &lucideUserRound,
			Path:      "role",
			Component: "/system/role/index",
			Sort:      12,
			ParentID:  &id1,
			Roles:     []*model.Role{&adminRole},
			Creator:   "系统",
		},
		{
			Model:     model.Model{ID: 4},
			Name:      "Menu",
			Title:     "菜单",
			Icon:      &menuSquare,
			Path:      "menu",
			Component: "/system/menu/index",
			Sort:      13,
			ParentID:  &id1,
			Roles:     []*model.Role{&adminRole},
			Creator:   "系统",
		},
		{
			Model:     model.Model{ID: 5},
			Name:      "Api",
			Title:     "接口权限",
			Icon:      &network,
			Path:      "api",
			Component: "/system/api/index",
			Sort:      14,
			ParentID:  &id1,
			Roles:     []*model.Role{&adminRole},
			Creator:   "系统",
		},
		{
			Model:     model.Model{ID: 7},
			Name:      "OperationLog",
			Title:     "操作日志",
			Icon:      &lucideMessageSquareCode,
			Path:      "operation-log",
			Component: "/system/operation-log/index",
			Sort:      21,
			ParentID:  &id1,
			Roles:     []*model.Role{&adminRole},
			Creator:   "系统",
		},
		{
			Model:     model.Model{ID: 8},
			Name:      "Business",
			Title:     "用户与项目",
			Icon:      &briefcaseBusiness,
			Path:      "/business",
			Component: "Layout",
			Sort:      3,
			ParentID:  &id0,
			Roles:     []*model.Role{&adminRole},
			Creator:   "系统",
		},
		{
			Model:     model.Model{ID: 9},
			Name:      "Project",
			Title:     "项目",
			Icon:      &lucideTableProperties,
			Path:      "/business/project",
			Component: "/business/project/index",
			Sort:      1,
			ParentID:  &id8,
			Roles:     []*model.Role{&adminRole},
			Creator:   "系统",
		},
		{
			Model:     model.Model{ID: 10},
			Name:      "Content",
			Title:     "内容管理",
			Icon:      &tableOfContents,
			Path:      "/content",
			Component: "Layout",
			Sort:      4,
			ParentID:  &id0,
			Roles:     []*model.Role{&adminRole},
			Creator:   "系统",
		},
		{
			Model:     model.Model{ID: 11},
			Name:      "Tag",
			Title:     "标签",
			Icon:      &lucideTags,
			Path:      "/content/tag",
			Component: "/content/tag/index",
			Sort:      3,
			ParentID:  &id10,
			Roles:     []*model.Role{&adminRole},
			Creator:   "系统",
		},
		{
			Model:     model.Model{ID: 12},
			Name:      "Category",
			Title:     "分类",
			Icon:      &boxes,
			Path:      "/content/category",
			Component: "/content/category/index",
			Sort:      2,
			ParentID:  &id10,
			Roles:     []*model.Role{&adminRole},
			Creator:   "系统",
		},
		{
			Model:     model.Model{ID: 13},
			Name:      "Article",
			Title:     "文章",
			Icon:      &wholeWord,
			Path:      "/content/article",
			Component: "/content/article/index",
			Sort:      1,
			ParentID:  &id10,
			Roles:     []*model.Role{&adminRole},
			Creator:   "系统",
		},
		{
			Model:     model.Model{ID: 14},
			Name:      "Config",
			Title:     "数据源",
			Icon:      &database,
			Path:      "/content/config",
			Component: "/content/config/index",
			Sort:      6,
			ParentID:  &id10,
			Roles:     []*model.Role{&adminRole},
			Creator:   "系统",
		},
		{
			Model:     model.Model{ID: 15},
			Name:      "Resource",
			Title:     "资源库",
			Icon:      &bookImage,
			Path:      "/content/resource",
			Component: "/content/resource/index",
			Sort:      5,
			ParentID:  &id10,
			Roles:     []*model.Role{&adminRole},
			Creator:   "系统",
		},
		{
			Model:     model.Model{ID: 16},
			Name:      "User",
			Title:     "用户",
			Icon:      &lucideUserCog2,
			Path:      "/business/user",
			Component: "/business/user/index",
			Sort:      2,
			ParentID:  &id8,
			Roles:     []*model.Role{&adminRole},
			Creator:   "系统",
		},
		{
			Model:     model.Model{ID: 17},
			Name:      "Column",
			Title:     "专栏",
			Icon:      &lucideColumnsSettings,
			Path:      "/content/column",
			Component: "/content/column/index",
			Sort:      4,
			ParentID:  &id10,
			Roles:     []*model.Role{&adminRole},
			Creator:   "系统",
		},
		{
			Model:     model.Model{ID: 18},
			Name:      "Operations",
			Title:     "运营管理",
			Icon:      &clipboardEdit,
			Path:      "/operations",
			Component: "Layout",
			Sort:      5,
			ParentID:  &id0,
			Roles:     []*model.Role{&adminRole},
			Creator:   "系统",
		},
		{
			Model:     model.Model{ID: 19},
			Name:      "Scene",
			Title:     "场景",
			Icon:      &laptopMinimalCheck,
			Path:      "/promotion/scene",
			Component: "/promotion/scene/index",
			Sort:      1,
			ParentID:  &id18,
			Roles:     []*model.Role{&adminRole},
			Creator:   "系统",
		},
		{
			Model:     model.Model{ID: 32},
			Name:      "Advertising",
			Title:     "广告",
			Icon:      &lucideHeadset,
			Path:      "/promotion/advertising",
			Component: "/promotion/advertising/index",
			Sort:      2,
			ParentID:  &id18,
			Roles:     []*model.Role{&adminRole},
			Creator:   "系统",
		},
		{
			Model:     model.Model{ID: 20},
			Name:      "comment",
			Title:     "评论",
			Icon:      &lucideMartini,
			Path:      "/operation/comment",
			Component: "/operation/comment/index",
			Sort:      3,
			ParentID:  &id18,
			Roles:     []*model.Role{&adminRole},
			Creator:   "系统",
		},
		{
			Model:     model.Model{ID: 21},
			Name:      "point",
			Title:     "积分",
			Icon:      &lucideDatabaseZap,
			Path:      "/operation/point",
			Component: "/operation/point/index",
			Sort:      4,
			ParentID:  &id18,
			Roles:     []*model.Role{&adminRole},
			Creator:   "系统",
		},
		{
			Model:     model.Model{ID: 29},
			Name:      "AdminConfig",
			Title:     "系统配置",
			Icon:      &lucideTableConfig,
			Path:      "/system/config",
			Component: "/system/config/index",
			Sort:      22,
			ParentID:  &id1,
			Roles:     []*model.Role{&adminRole},
			Creator:   "系统",
		},
	}

	for _, menu := range menus {
		if err := createOrSyncByID(db, &menu, menu.ID, syncExisting, []string{
			"name",
			"title",
			"icon",
			"path",
			"redirect",
			"component",
			"sort",
			"status",
			"hidden",
			"no_cache",
			"always_show",
			"breadcrumb",
			"active_menu",
			"parent_id",
			"creator",
		}); err != nil {
			return err
		}
		if err := ensureRoleMenu(db, adminRole.ID, menu.ID); err != nil {
			return err
		}
	}

	return nil
}

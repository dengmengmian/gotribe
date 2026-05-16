package seeder

import (
	"errors"
	"log"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Seeder 数据种子接口
type Seeder interface {
	Run(db *gorm.DB, syncExisting bool) error
	Name() string
}

// BaseSeeder 基础种子实现
type BaseSeeder struct {
	name string
}

// NewBaseSeeder 创建基础种子
func NewBaseSeeder(name string) *BaseSeeder {
	return &BaseSeeder{name: name}
}

// Name 返回种子名称
func (s *BaseSeeder) Name() string {
	return s.name
}

// SeedRegistry 种子注册表
type SeedRegistry struct {
	seeders []Seeder
}

// NewSeedRegistry 创建种子注册表
func NewSeedRegistry() *SeedRegistry {
	return &SeedRegistry{
		seeders: make([]Seeder, 0),
	}
}

// Register 注册种子
func (r *SeedRegistry) Register(seeder Seeder) {
	r.seeders = append(r.seeders, seeder)
}

// RunAll 运行所有种子
func (r *SeedRegistry) RunAll(db *gorm.DB, syncExisting bool) error {
	for _, seeder := range r.seeders {
		log.Printf("开始执行种子: %s", seeder.Name())
		if err := seeder.Run(db, syncExisting); err != nil {
			log.Printf("种子 %s 执行失败: %v", seeder.Name(), err)
			return err
		}
		log.Printf("种子 %s 执行完成", seeder.Name())
	}
	return nil
}

// 全局种子注册表
var GlobalRegistry = NewSeedRegistry()

// RegisterSeeder 注册种子的便捷函数
func RegisterSeeder(seeder Seeder) {
	GlobalRegistry.Register(seeder)
}

// RunSeeders 运行所有注册的种子
func RunSeeders(db *gorm.DB, syncExisting bool) error {
	return GlobalRegistry.RunAll(db, syncExisting)
}

// 通用创建函数
func createIfNotExists(db *gorm.DB, model interface{}, id int64) error {
	err := db.First(model, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return db.Create(model).Error
	}
	return err
}

func createOrSyncByID(db *gorm.DB, value interface{}, id int64, syncExisting bool, updateColumns []string) error {
	if syncExisting {
		return upsertByID(db, value, updateColumns)
	}
	return createIfNotExists(db, value, id)
}

// upsertByID 同步初始化参考数据。仅更新调用方显式列出的字段，避免覆盖业务运行时数据。
func upsertByID(db *gorm.DB, value interface{}, updateColumns []string) error {
	return db.
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			DoUpdates: clause.AssignmentColumns(updateColumns),
		}).
		Omit(clause.Associations).
		Create(value).
		Error
}

func ensureRoleMenu(db *gorm.DB, roleID, menuID int64) error {
	return db.
		Clauses(clause.OnConflict{DoNothing: true}).
		Table("role_menus").
		Create(map[string]interface{}{
			"role_id": roleID,
			"menu_id": menuID,
		}).
		Error
}

func ensureAdminRole(db *gorm.DB, adminID, roleID int64) error {
	return db.
		Clauses(clause.OnConflict{DoNothing: true}).
		Table("admin_roles").
		Create(map[string]interface{}{
			"admin_id": adminID,
			"role_id":  roleID,
		}).
		Error
}

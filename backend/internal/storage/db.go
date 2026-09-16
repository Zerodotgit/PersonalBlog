package storage

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"personal-blog/internal/config"
	"personal-blog/internal/model"

	"github.com/glebarez/sqlite"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Open 根据配置中的 driver 建立数据库连接（开发 sqlite / 生产 mysql），
// 随后执行迁移与初始管理员种子。
func Open(cfg *config.Config) (*gorm.DB, error) {
	dialector, driver, err := dialectorFor(cfg.Database)
	if err != nil {
		return nil, err
	}
	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("open %s: %w", driver, err)
	}

	if err := db.AutoMigrate(&model.User{}, &model.Article{}); err != nil {
		return nil, fmt.Errorf("auto migrate: %w", err)
	}

	if err := seedAdmin(db, cfg.Admin.Username, cfg.Admin.Password); err != nil {
		return nil, fmt.Errorf("seed admin: %w", err)
	}

	log.Printf("database initialized (driver=%s): tables migrated, admin seeded", driver)
	return db, nil
}

// dialectorFor 依据 driver 字段返回对应的 GORM Dialector。
// 未指定 driver 时默认 mysql，以保持与历史生产配置兼容。
func dialectorFor(dbc config.DatabaseConfig) (gorm.Dialector, string, error) {
	driver := strings.ToLower(strings.TrimSpace(dbc.Driver))
	switch driver {
	case "", "mysql":
		return mysql.Open(dbc.DSN), "mysql", nil
	case "sqlite":
		dsn := strings.TrimSpace(dbc.DSN)
		if dsn == "" {
			dsn = "blog.db"
		}
		// 为本地 sqlite 文件创建所在目录（忽略 query 串）
		file := dsn
		if i := strings.IndexAny(dsn, "?"); i >= 0 {
			file = dsn[:i]
		}
		if dir := filepath.Dir(file); dir != "." && dir != "" {
			if err := os.MkdirAll(dir, 0o755); err != nil {
				return nil, driver, fmt.Errorf("create sqlite dir: %w", err)
			}
		}
		return sqlite.Open(dsn), "sqlite", nil
	default:
		return nil, driver, fmt.Errorf("unsupported database driver: %s (only sqlite/mysql)", driver)
	}
}

// seedAdmin 在不存在管理员时写入默认账号，密码以 bcrypt 加密存储。
func seedAdmin(db *gorm.DB, username, password string) error {
	var count int64
	if err := db.Model(&model.User{}).Where("username = ?", username).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}
	return db.Create(&model.User{
		Username:     username,
		PasswordHash: string(hash),
	}).Error
}

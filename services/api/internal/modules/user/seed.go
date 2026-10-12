package user

import (
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/yanmai/api/pkg/hash"
)

type seedAccount struct {
	Username string
	Phone    string
	Nickname string
	Role     string
	Password string
}

// seedAccounts 本地联调种子账号（仅在 DB_SEED=true 时写入，且已存在则跳过）。
var seedAccounts = []seedAccount{
	{Username: "admin", Phone: "13800000001", Nickname: "平台管理员", Role: RoleAdmin, Password: "Admin@123456"},
	{Username: "lecturer", Phone: "13800000002", Nickname: "示范讲师", Role: RoleLecturer, Password: "Lecturer@123456"},
	{Username: "student", Phone: "13800000003", Nickname: "示范学员", Role: RoleStudent, Password: "Student@123456"},
}

// Seed 写入种子账号，幂等。
func Seed(db *gorm.DB, log *zap.Logger) error {
	created := 0
	for _, a := range seedAccounts {
		var count int64
		if err := db.Model(&User{}).Where("username = ?", a.Username).Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			continue
		}

		hashed, err := hash.Password(a.Password)
		if err != nil {
			return err
		}

		u := &User{
			Username:     a.Username,
			Phone:        a.Phone,
			PasswordHash: hashed,
			Nickname:     a.Nickname,
			Role:         a.Role,
			Status:       StatusActive,
		}
		if err := db.Create(u).Error; err != nil {
			return err
		}
		log.Info("已创建种子账号",
			zap.String("username", a.Username),
			zap.String("role", a.Role),
			zap.String("password", a.Password),
		)
		created++
	}

	if created > 0 {
		log.Warn("种子账号仅供本地联调使用，生产环境请保持 DB_SEED=false")
	} else {
		log.Info("种子账号已存在，跳过写入")
	}
	return nil
}

// SeedAccounts 返回种子账号清单（供 README / 验收脚本引用，避免文档与代码不一致）。
func SeedAccounts() []seedAccount { return seedAccounts }

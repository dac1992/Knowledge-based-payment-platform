package user

import (
	"context"
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"

	"github.com/yanmai/api/pkg/errs"
)

// Repository 用户数据访问接口，由本模块独占 users 表。
type Repository interface {
	Create(ctx context.Context, u *User) error
	FindByID(ctx context.Context, id uint64) (*User, error)
	FindByLogin(ctx context.Context, login string) (*User, error)
	ExistsByUsername(ctx context.Context, username string) (bool, error)
	ExistsByPhone(ctx context.Context, phone string) (bool, error)
	TouchLastLogin(ctx context.Context, id uint64, at time.Time) error
	UpdateFields(ctx context.Context, id uint64, fields map[string]any) error
	List(ctx context.Context, q ListQuery) ([]User, int64, error)
}

type repository struct{ db *gorm.DB }

// NewRepository 构造仓储实现。
func NewRepository(db *gorm.DB) Repository { return &repository{db: db} }

func (r *repository) Create(ctx context.Context, u *User) error {
	if err := r.db.WithContext(ctx).Create(u).Error; err != nil {
		if isDuplicateKey(err) {
			return errs.Wrap(errs.CodeUserExists, err).WithMessage("登录名或手机号已被占用")
		}
		return errs.Wrap(errs.CodeInternal, err).WithMessage("写入用户记录失败")
	}
	return nil
}

func (r *repository) FindByID(ctx context.Context, id uint64) (*User, error) {
	var u User
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&u).Error
	if err != nil {
		return nil, mapNotFound(err, errs.CodeUserNotFound)
	}
	return &u, nil
}

func (r *repository) FindByLogin(ctx context.Context, login string) (*User, error) {
	login = strings.TrimSpace(login)
	var u User
	err := r.db.WithContext(ctx).
		Where("username = ? OR (phone <> '' AND phone = ?)", login, login).
		First(&u).Error
	if err != nil {
		return nil, mapNotFound(err, errs.CodeUserNotFound)
	}
	return &u, nil
}

func (r *repository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	return r.exists(ctx, "username = ?", strings.TrimSpace(username))
}

func (r *repository) ExistsByPhone(ctx context.Context, phone string) (bool, error) {
	phone = strings.TrimSpace(phone)
	if phone == "" {
		return false, nil
	}
	return r.exists(ctx, "phone = ?", phone)
}

func (r *repository) exists(ctx context.Context, where string, args ...any) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&User{}).Where(where, args...).Count(&count).Error; err != nil {
		return false, errs.Wrap(errs.CodeInternal, err).WithMessage("查询用户是否存在失败")
	}
	return count > 0, nil
}

func (r *repository) TouchLastLogin(ctx context.Context, id uint64, at time.Time) error {
	err := r.db.WithContext(ctx).
		Model(&User{}).
		Where("id = ?", id).
		Update("last_login_at", at).Error
	if err != nil {
		return errs.Wrap(errs.CodeInternal, err).WithMessage("更新登录时间失败")
	}
	return nil
}

func (r *repository) UpdateFields(ctx context.Context, id uint64, fields map[string]any) error {
	if len(fields) == 0 {
		return nil
	}
	err := r.db.WithContext(ctx).Model(&User{}).Where("id = ?", id).Updates(fields).Error
	if err != nil {
		if isDuplicateKey(err) {
			return errs.Wrap(errs.CodeConflict, err).WithMessage("手机号已被其他账号占用")
		}
		return errs.Wrap(errs.CodeInternal, err).WithMessage("更新用户资料失败")
	}
	return nil
}

func (r *repository) List(ctx context.Context, q ListQuery) ([]User, int64, error) {
	base := func() *gorm.DB {
		tx := r.db.WithContext(ctx).Model(&User{})
		if kw := strings.TrimSpace(q.Keyword); kw != "" {
			like := "%" + kw + "%"
			tx = tx.Where("username ILIKE ? OR nickname ILIKE ? OR phone ILIKE ?", like, like, like)
		}
		if q.Role != "" {
			tx = tx.Where("role = ?", q.Role)
		}
		if q.Status != nil {
			tx = tx.Where("status = ?", *q.Status)
		}
		return tx
	}

	var total int64
	if err := base().Count(&total).Error; err != nil {
		return nil, 0, errs.Wrap(errs.CodeInternal, err).WithMessage("统计用户总数失败")
	}
	if total == 0 {
		return []User{}, 0, nil
	}

	page, size := NormalizePage(q.Page, q.PageSize)
	var list []User
	if err := base().Order("id DESC").Offset((page - 1) * size).Limit(size).Find(&list).Error; err != nil {
		return nil, 0, errs.Wrap(errs.CodeInternal, err).WithMessage("查询用户列表失败")
	}
	return list, total, nil
}

// NormalizePage 规范化分页参数。
func NormalizePage(page, size int) (int, int) {
	if page < 1 {
		page = 1
	}
	switch {
	case size < 1:
		size = 20
	case size > 100:
		size = 100
	}
	return page, size
}

func mapNotFound(err error, code errs.Code) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errs.New(code)
	}
	return errs.Wrap(errs.CodeInternal, err).WithMessage("查询用户失败")
}

func isDuplicateKey(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "duplicate key value") ||
		strings.Contains(msg, "SQLSTATE 23505") ||
		strings.Contains(msg, "UNIQUE constraint failed")
}

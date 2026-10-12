package user

import (
	"context"
	"strings"
	"time"

	"github.com/yanmai/api/pkg/errs"
	"github.com/yanmai/api/pkg/hash"
)

// Service 用户领域服务，是其他模块访问用户数据的唯一入口。
type Service interface {
	Create(ctx context.Context, in CreateInput) (*User, error)
	GetByID(ctx context.Context, id uint64) (*User, error)
	// Authenticate 校验凭证并返回用户；账号不存在与密码错误统一返回 CodeCredentialInvalid，避免账号枚举。
	Authenticate(ctx context.Context, login, password string) (*User, error)
	TouchLastLogin(ctx context.Context, id uint64) error
	Profile(ctx context.Context, id uint64) (Profile, error)
	UpdateProfile(ctx context.Context, id uint64, in UpdateInput) (Profile, error)
	List(ctx context.Context, q ListQuery) ([]Profile, int64, error)
}

type service struct{ repo Repository }

// NewService 构造用户服务。
func NewService(repo Repository) Service { return &service{repo: repo} }

func (s *service) Create(ctx context.Context, in CreateInput) (*User, error) {
	in.Username = strings.TrimSpace(in.Username)
	in.Phone = strings.TrimSpace(in.Phone)
	in.Email = strings.TrimSpace(in.Email)
	in.Nickname = strings.TrimSpace(in.Nickname)
	if in.Role == "" {
		in.Role = RoleStudent
	}

	if err := ValidateUsername(in.Username); err != nil {
		return nil, err
	}
	if err := ValidatePassword(in.Password); err != nil {
		return nil, err
	}
	if err := ValidateRole(in.Role); err != nil {
		return nil, err
	}

	taken, err := s.repo.ExistsByUsername(ctx, in.Username)
	if err != nil {
		return nil, err
	}
	if taken {
		return nil, errs.New(errs.CodeUserExists).WithMessage("登录名 %s 已被占用", in.Username)
	}

	if in.Phone != "" {
		taken, err = s.repo.ExistsByPhone(ctx, in.Phone)
		if err != nil {
			return nil, err
		}
		if taken {
			return nil, errs.New(errs.CodeUserExists).WithMessage("手机号 %s 已注册", maskPhone(in.Phone))
		}
	}

	hashed, err := hash.Password(in.Password)
	if err != nil {
		return nil, errs.Wrap(errs.CodeInternal, err).WithMessage("密码加密失败")
	}

	if in.Nickname == "" {
		in.Nickname = in.Username
	}

	u := &User{
		Username:     in.Username,
		Phone:        in.Phone,
		Email:        in.Email,
		PasswordHash: hashed,
		Nickname:     in.Nickname,
		Role:         in.Role,
		Status:       StatusActive,
	}
	if err := s.repo.Create(ctx, u); err != nil {
		return nil, err
	}
	return u, nil
}

func (s *service) GetByID(ctx context.Context, id uint64) (*User, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *service) Authenticate(ctx context.Context, login, password string) (*User, error) {
	u, err := s.repo.FindByLogin(ctx, login)
	if err != nil {
		if errs.Is(err, errs.CodeUserNotFound) {
			return nil, errs.New(errs.CodeCredentialInvalid)
		}
		return nil, err
	}
	if !hash.CheckPassword(u.PasswordHash, password) {
		return nil, errs.New(errs.CodeCredentialInvalid)
	}
	if !u.IsActive() {
		return nil, errs.New(errs.CodeUserDisabled)
	}
	return u, nil
}

func (s *service) TouchLastLogin(ctx context.Context, id uint64) error {
	return s.repo.TouchLastLogin(ctx, id, time.Now())
}

func (s *service) Profile(ctx context.Context, id uint64) (Profile, error) {
	u, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return Profile{}, err
	}
	return ToProfile(u), nil
}

func (s *service) UpdateProfile(ctx context.Context, id uint64, in UpdateInput) (Profile, error) {
	fields := make(map[string]any, 4)
	if in.Nickname != nil {
		v := strings.TrimSpace(*in.Nickname)
		if v == "" {
			return Profile{}, errs.New(errs.CodeInvalidParam).WithMessage("昵称不能为空")
		}
		fields["nickname"] = v
	}
	if in.Avatar != nil {
		fields["avatar"] = strings.TrimSpace(*in.Avatar)
	}
	if in.Email != nil {
		fields["email"] = strings.TrimSpace(*in.Email)
	}
	if in.Phone != nil {
		phone := strings.TrimSpace(*in.Phone)
		if phone != "" {
			taken, err := s.repo.ExistsByPhone(ctx, phone)
			if err != nil {
				return Profile{}, err
			}
			if taken {
				return Profile{}, errs.New(errs.CodeConflict).WithMessage("手机号 %s 已被其他账号占用", maskPhone(phone))
			}
		}
		fields["phone"] = phone
	}

	if err := s.repo.UpdateFields(ctx, id, fields); err != nil {
		return Profile{}, err
	}
	return s.Profile(ctx, id)
}

func (s *service) List(ctx context.Context, q ListQuery) ([]Profile, int64, error) {
	list, total, err := s.repo.List(ctx, q)
	if err != nil {
		return nil, 0, err
	}
	return ToProfiles(list), total, nil
}

// ValidateUsername 校验登录名规则：3-32 位，仅允许字母、数字与下划线。
func ValidateUsername(username string) error {
	if len(username) < 3 || len(username) > 32 {
		return errs.New(errs.CodeInvalidParam).WithMessage("登录名长度需在 3-32 个字符之间")
	}
	for _, r := range username {
		switch {
		case r >= 'a' && r <= 'z':
		case r >= 'A' && r <= 'Z':
		case r >= '0' && r <= '9':
		case r == '_':
		default:
			return errs.New(errs.CodeInvalidParam).WithMessage("登录名仅允许字母、数字与下划线")
		}
	}
	return nil
}

// ValidatePassword 校验密码强度。
func ValidatePassword(password string) error {
	if len(password) < 8 {
		return errs.New(errs.CodeInvalidParam).WithMessage("密码长度不得少于 8 位")
	}
	if len(password) > 64 {
		return errs.New(errs.CodeInvalidParam).WithMessage("密码长度不得超过 64 位")
	}
	return nil
}

// ValidateRole 校验角色取值。
func ValidateRole(role string) error {
	switch role {
	case RoleAdmin, RoleLecturer, RoleStudent:
		return nil
	default:
		return errs.New(errs.CodeInvalidParam).WithMessage("角色取值非法: %s", role)
	}
}

func maskPhone(phone string) string {
	if len(phone) < 7 {
		return phone
	}
	return phone[:3] + "****" + phone[len(phone)-4:]
}

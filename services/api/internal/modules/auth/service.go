// Package auth 认证模块：负责注册、登录、令牌签发与刷新。
//
// 职责边界：
//   - 本模块不直接访问 users 表，凭证校验委托给 user 模块的 Service
//   - 本模块只负责「身份凭证 → 令牌」的转换与令牌生命周期
package auth

import (
	"context"

	"github.com/yanmai/api/internal/modules/user"
	"github.com/yanmai/api/pkg/jwt"
)

// Result 登录态结果。
type Result struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int64
	User         user.Profile
}

// RegisterInput 注册入参。
type RegisterInput struct {
	Username string
	Password string
	Phone    string
	Email    string
	Nickname string
}

// LoginInput 登录入参。
type LoginInput struct {
	Login    string
	Password string
}

// userProvider 本模块所需的最小用户能力。
//
// 由消费方（本模块）定义接口、提供方（user 模块）实现，符合依赖倒置原则：
// 单元测试可注入 mock，未来把 user 拆成独立服务也只需更换实现。
type userProvider interface {
	Create(ctx context.Context, in user.CreateInput) (*user.User, error)
	Authenticate(ctx context.Context, login, password string) (*user.User, error)
	GetByID(ctx context.Context, id uint64) (*user.User, error)
	TouchLastLogin(ctx context.Context, id uint64) error
}

// Service 认证服务接口。
type Service interface {
	Register(ctx context.Context, in RegisterInput) (*Result, error)
	Login(ctx context.Context, in LoginInput) (*Result, error)
	Refresh(ctx context.Context, refreshToken string) (*Result, error)
	Profile(ctx context.Context, userID uint64) (user.Profile, error)
}

type service struct {
	users userProvider
	token *jwt.Manager
}

// NewService 构造认证服务。
func NewService(users userProvider, token *jwt.Manager) Service {
	return &service{users: users, token: token}
}

func (s *service) Register(ctx context.Context, in RegisterInput) (*Result, error) {
	u, err := s.users.Create(ctx, user.CreateInput{
		Username: in.Username,
		Password: in.Password,
		Phone:    in.Phone,
		Email:    in.Email,
		Nickname: in.Nickname,
		Role:     user.RoleStudent,
	})
	if err != nil {
		return nil, err
	}
	return s.issue(u)
}

func (s *service) Login(ctx context.Context, in LoginInput) (*Result, error) {
	u, err := s.users.Authenticate(ctx, in.Login, in.Password)
	if err != nil {
		return nil, err
	}
	// 登录时间失败不应阻断登录，忽略错误但保留实现位置便于后续接审计日志
	_ = s.users.TouchLastLogin(ctx, u.ID)

	return s.issue(u)
}

func (s *service) Refresh(ctx context.Context, refreshToken string) (*Result, error) {
	claims, err := s.token.ParseRefresh(refreshToken)
	if err != nil {
		return nil, err
	}

	// 回源校验账号仍然存在且可用，避免令牌有效但账号被禁用/删除
	u, err := s.users.GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, err
	}
	if !u.IsActive() {
		return nil, errRefreshForbidden()
	}
	return s.issue(u)
}

func (s *service) Profile(ctx context.Context, userID uint64) (user.Profile, error) {
	u, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return user.Profile{}, err
	}
	return user.ToProfile(u), nil
}

func (s *service) issue(u *user.User) (*Result, error) {
	access, err := s.token.GenerateAccess(u.ID, u.Username, u.Role)
	if err != nil {
		return nil, err
	}
	refresh, err := s.token.GenerateRefresh(u.ID, u.Username, u.Role)
	if err != nil {
		return nil, err
	}
	return &Result{
		AccessToken:  access.Value,
		RefreshToken: refresh.Value,
		ExpiresIn:    int64(s.token.AccessTTL().Seconds()),
		User:         user.ToProfile(u),
	}, nil
}

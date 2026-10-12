package user

import "time"

// Profile 对外暴露的用户资料（绝不包含密码哈希）。
type Profile struct {
	ID          uint64     `json:"id"`
	Username    string     `json:"username"`
	Phone       string     `json:"phone"`
	Email       string     `json:"email"`
	Nickname    string     `json:"nickname"`
	Avatar      string     `json:"avatar"`
	Role        string     `json:"role"`
	Status      int8       `json:"status"`
	LastLoginAt *time.Time `json:"lastLoginAt,omitempty"`
	CreatedAt   time.Time  `json:"createdAt"`
}

// ToProfile 把模型转换为对外资料。
func ToProfile(u *User) Profile {
	if u == nil {
		return Profile{}
	}
	return Profile{
		ID:          u.ID,
		Username:    u.Username,
		Phone:       u.Phone,
		Email:       u.Email,
		Nickname:    u.Nickname,
		Avatar:      u.Avatar,
		Role:        u.Role,
		Status:      u.Status,
		LastLoginAt: u.LastLoginAt,
		CreatedAt:   u.CreatedAt,
	}
}

// ToProfiles 批量转换。
func ToProfiles(us []User) []Profile {
	out := make([]Profile, 0, len(us))
	for i := range us {
		out = append(out, ToProfile(&us[i]))
	}
	return out
}

// CreateInput 创建用户入参（密码为明文，由本模块负责哈希）。
type CreateInput struct {
	Username string
	Password string
	Phone    string
	Email    string
	Nickname string
	Role     string
}

// UpdateInput 更新资料入参，nil 表示不修改。
type UpdateInput struct {
	Nickname *string
	Avatar   *string
	Email    *string
	Phone    *string
}

// ListQuery 用户列表查询条件。
type ListQuery struct {
	Page     int
	PageSize int
	Keyword  string
	Role     string
	Status   *int8
}

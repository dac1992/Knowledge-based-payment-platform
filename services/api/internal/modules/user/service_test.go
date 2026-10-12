package user

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/yanmai/api/pkg/errs"
)

// memRepo 内存仓储，用于在无数据库环境下验证领域逻辑。
type memRepo struct {
	seq  uint64
	rows map[uint64]*User
}

func newMemRepo() *memRepo { return &memRepo{rows: make(map[uint64]*User)} }

func (m *memRepo) Create(_ context.Context, u *User) error {
	for _, r := range m.rows {
		if r.Username == u.Username {
			return errs.New(errs.CodeUserExists)
		}
		if u.Phone != "" && r.Phone == u.Phone {
			return errs.New(errs.CodeUserExists)
		}
	}
	m.seq++
	u.ID = m.seq
	now := time.Now()
	u.CreatedAt, u.UpdatedAt = now, now
	cp := *u
	m.rows[u.ID] = &cp
	return nil
}

func (m *memRepo) FindByID(_ context.Context, id uint64) (*User, error) {
	if r, ok := m.rows[id]; ok {
		cp := *r
		return &cp, nil
	}
	return nil, errs.New(errs.CodeUserNotFound)
}

func (m *memRepo) FindByLogin(_ context.Context, login string) (*User, error) {
	login = strings.TrimSpace(login)
	for _, r := range m.rows {
		if r.Username == login || (r.Phone != "" && r.Phone == login) {
			cp := *r
			return &cp, nil
		}
	}
	return nil, errs.New(errs.CodeUserNotFound)
}

func (m *memRepo) ExistsByUsername(_ context.Context, username string) (bool, error) {
	for _, r := range m.rows {
		if r.Username == username {
			return true, nil
		}
	}
	return false, nil
}

func (m *memRepo) ExistsByPhone(_ context.Context, phone string) (bool, error) {
	if phone == "" {
		return false, nil
	}
	for _, r := range m.rows {
		if r.Phone == phone {
			return true, nil
		}
	}
	return false, nil
}

func (m *memRepo) TouchLastLogin(_ context.Context, id uint64, at time.Time) error {
	if r, ok := m.rows[id]; ok {
		t := at
		r.LastLoginAt = &t
		return nil
	}
	return errs.New(errs.CodeUserNotFound)
}

func (m *memRepo) UpdateFields(_ context.Context, id uint64, fields map[string]any) error {
	r, ok := m.rows[id]
	if !ok {
		return errs.New(errs.CodeUserNotFound)
	}
	for k, v := range fields {
		s, _ := v.(string)
		switch k {
		case "nickname":
			r.Nickname = s
		case "avatar":
			r.Avatar = s
		case "email":
			r.Email = s
		case "phone":
			r.Phone = s
		}
	}
	return nil
}

func (m *memRepo) List(_ context.Context, q ListQuery) ([]User, int64, error) {
	out := make([]User, 0, len(m.rows))
	for _, r := range m.rows {
		if q.Role != "" && r.Role != q.Role {
			continue
		}
		if q.Status != nil && r.Status != *q.Status {
			continue
		}
		if kw := strings.TrimSpace(q.Keyword); kw != "" && !strings.Contains(r.Username, kw) {
			continue
		}
		out = append(out, *r)
	}
	return out, int64(len(out)), nil
}

func newTestService() Service { return NewService(newMemRepo()) }

func TestCreateHashesPasswordAndAppliesDefaults(t *testing.T) {
	svc := newTestService()

	u, err := svc.Create(context.Background(), CreateInput{
		Username: "alice",
		Password: "Alice@2026",
	})
	if err != nil {
		t.Fatalf("创建用户失败: %v", err)
	}
	if u.ID == 0 {
		t.Fatal("应分配自增 ID")
	}
	if u.Role != RoleStudent {
		t.Fatalf("默认角色应为 %s，实际 %s", RoleStudent, u.Role)
	}
	if u.Nickname != "alice" {
		t.Fatalf("昵称应回退为登录名，实际 %q", u.Nickname)
	}
	if u.Status != StatusActive {
		t.Fatalf("新用户应为启用状态，实际 %d", u.Status)
	}
	if u.PasswordHash == "Alice@2026" || u.PasswordHash == "" {
		t.Fatal("密码必须哈希存储")
	}
}

func TestCreateRejectsDuplicateUsername(t *testing.T) {
	svc := newTestService()
	in := CreateInput{Username: "bob", Password: "Bob@20260718"}

	if _, err := svc.Create(context.Background(), in); err != nil {
		t.Fatalf("首次创建应成功: %v", err)
	}
	if _, err := svc.Create(context.Background(), in); !errs.Is(err, errs.CodeUserExists) {
		t.Fatalf("重复登录名应报 CodeUserExists，实际: %v", err)
	}
}

func TestCreateRejectsInvalidInput(t *testing.T) {
	svc := newTestService()

	cases := []struct {
		name string
		in   CreateInput
	}{
		{"登录名过短", CreateInput{Username: "ab", Password: "Abcdef@2026"}},
		{"登录名含非法字符", CreateInput{Username: "bad name", Password: "Abcdef@2026"}},
		{"密码过短", CreateInput{Username: "carol", Password: "short"}},
		{"角色非法", CreateInput{Username: "dave", Password: "Abcdef@2026", Role: "superman"}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if _, err := svc.Create(context.Background(), c.in); !errs.Is(err, errs.CodeInvalidParam) {
				t.Fatalf("应报 CodeInvalidParam，实际: %v", err)
			}
		})
	}
}

func TestAuthenticateScenarios(t *testing.T) {
	ctx := context.Background()
	svc := newTestService()

	if _, err := svc.Create(ctx, CreateInput{
		Username: "erin", Password: "Erin@2026", Phone: "13900000001",
	}); err != nil {
		t.Fatalf("准备数据失败: %v", err)
	}

	t.Run("登录名 + 正确密码", func(t *testing.T) {
		u, err := svc.Authenticate(ctx, "erin", "Erin@2026")
		if err != nil || u.Username != "erin" {
			t.Fatalf("应登录成功，err=%v", err)
		}
	})

	t.Run("手机号 + 正确密码", func(t *testing.T) {
		if _, err := svc.Authenticate(ctx, "13900000001", "Erin@2026"); err != nil {
			t.Fatalf("应支持手机号登录，err=%v", err)
		}
	})

	t.Run("密码错误", func(t *testing.T) {
		if _, err := svc.Authenticate(ctx, "erin", "wrong-password"); !errs.Is(err, errs.CodeCredentialInvalid) {
			t.Fatalf("应报 CodeCredentialInvalid，实际: %v", err)
		}
	})

	t.Run("账号不存在时不得泄露账号是否存在", func(t *testing.T) {
		if _, err := svc.Authenticate(ctx, "nobody", "whatever2026"); !errs.Is(err, errs.CodeCredentialInvalid) {
			t.Fatalf("应返回与密码错误一致的错误码，实际: %v", err)
		}
	})

	t.Run("账号被禁用", func(t *testing.T) {
		u, err := svc.Create(ctx, CreateInput{Username: "frank", Password: "Frank@2026"})
		if err != nil {
			t.Fatalf("准备数据失败: %v", err)
		}
		svc.(*service).repo.(*memRepo).rows[u.ID].Status = StatusDisabled

		if _, err := svc.Authenticate(ctx, "frank", "Frank@2026"); !errs.Is(err, errs.CodeUserDisabled) {
			t.Fatalf("应报 CodeUserDisabled，实际: %v", err)
		}
	})
}

func TestUpdateProfile(t *testing.T) {
	ctx := context.Background()
	svc := newTestService()

	u, err := svc.Create(ctx, CreateInput{Username: "grace", Password: "Grace@2026"})
	if err != nil {
		t.Fatalf("准备数据失败: %v", err)
	}

	nickname := "格蕾丝"
	email := "grace@yanmai.cn"
	profile, err := svc.UpdateProfile(ctx, u.ID, UpdateInput{Nickname: &nickname, Email: &email})
	if err != nil {
		t.Fatalf("更新资料失败: %v", err)
	}
	if profile.Nickname != nickname || profile.Email != email {
		t.Fatalf("资料未更新: %+v", profile)
	}

	blank := "   "
	if _, err := svc.UpdateProfile(ctx, u.ID, UpdateInput{Nickname: &blank}); !errs.Is(err, errs.CodeInvalidParam) {
		t.Fatalf("空白昵称应被拒绝，实际: %v", err)
	}
}

func TestProfileNeverExposesPasswordHash(t *testing.T) {
	ctx := context.Background()
	svc := newTestService()

	u, err := svc.Create(ctx, CreateInput{Username: "henry", Password: "Henry@2026"})
	if err != nil {
		t.Fatalf("准备数据失败: %v", err)
	}
	p, err := svc.Profile(ctx, u.ID)
	if err != nil {
		t.Fatalf("查询资料失败: %v", err)
	}
	if strings.Contains(strings.ToLower(p.Nickname+p.Username+p.Email), "hash") {
		t.Fatal("对外资料不应包含密码相关字段")
	}
}

func TestValidateUsernameAndPassword(t *testing.T) {
	for _, name := range []string{"abc", "a_b-c", "user name", "ab", ""} {
		ok := ValidateUsername(name) == nil
		want := name == "abc"
		if ok != want {
			t.Errorf("ValidateUsername(%q) = %v，期望 %v", name, ok, want)
		}
	}

	if err := ValidatePassword("1234567"); err == nil {
		t.Error("7 位密码应被拒绝")
	}
	if err := ValidatePassword("12345678"); err != nil {
		t.Errorf("8 位密码应被接受，实际: %v", err)
	}
	if err := ValidatePassword(strings.Repeat("a", 65)); err == nil {
		t.Error("65 位密码应被拒绝")
	}
}

func TestNormalizePage(t *testing.T) {
	cases := []struct{ inPage, inSize, wantPage, wantSize int }{
		{0, 0, 1, 20},
		{-1, -5, 1, 20},
		{2, 50, 2, 50},
		{3, 999, 3, 100},
	}
	for _, c := range cases {
		p, s := NormalizePage(c.inPage, c.inSize)
		if p != c.wantPage || s != c.wantSize {
			t.Errorf("NormalizePage(%d,%d) = (%d,%d)，期望 (%d,%d)", c.inPage, c.inSize, p, s, c.wantPage, c.wantSize)
		}
	}
}

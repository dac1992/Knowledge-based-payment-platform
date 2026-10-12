package hash

import "testing"

func TestPasswordAndCheck(t *testing.T) {
	const plain = "Yanmai@2026"

	hashed, err := Password(plain)
	if err != nil {
		t.Fatalf("哈希密码失败: %v", err)
	}
	if hashed == plain {
		t.Fatal("哈希结果不得等于明文")
	}
	if len(hashed) < 50 {
		t.Fatalf("bcrypt 哈希长度异常: %d", len(hashed))
	}

	if !CheckPassword(hashed, plain) {
		t.Fatal("正确密码校验应通过")
	}
	if CheckPassword(hashed, "wrong-password") {
		t.Fatal("错误密码校验不应通过")
	}
}

func TestPasswordIsSalted(t *testing.T) {
	a, err := Password("same-password")
	if err != nil {
		t.Fatalf("哈希失败: %v", err)
	}
	b, err := Password("same-password")
	if err != nil {
		t.Fatalf("哈希失败: %v", err)
	}
	if a == b {
		t.Fatal("相同明文两次哈希结果不应相同（缺少随机盐）")
	}
	if !CheckPassword(a, "same-password") || !CheckPassword(b, "same-password") {
		t.Fatal("两次哈希都应能通过校验")
	}
}

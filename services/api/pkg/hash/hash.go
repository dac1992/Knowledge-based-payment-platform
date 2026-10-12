// Package hash 封装密码哈希，避免业务层直接依赖具体算法。
//
// 现阶段使用 bcrypt（cost=10）。若未来更换算法，只需保证哈希串前缀可区分即可平滑升级。
package hash

import "golang.org/x/crypto/bcrypt"

// Cost 密码哈希强度。
const Cost = 10

// Password 对明文密码做哈希。
func Password(plain string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(plain), Cost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// CheckPassword 校验明文密码与哈希是否匹配。
func CheckPassword(hashed, plain string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain)) == nil
}

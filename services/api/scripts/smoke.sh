#!/usr/bin/env bash
#
# 端到端冒烟测试：健康检查 → 登录 → 鉴权 → RBAC → 令牌刷新 → 占位模块
#
# 前置条件：API 服务已在 $BASE_URL 上运行，且已连接可用数据库。
#   cd services/api && make run
#   bash scripts/smoke.sh
#
set -uo pipefail

# 响应断言依赖 python3 解析 JSON，缺失时直接给出明确提示而不是逐条静默失败
if ! command -v python3 >/dev/null 2>&1; then
  printf '\033[31m冒烟测试需要 python3 解析 JSON 响应，请先安装 python3。\033[0m\n' >&2
  exit 2
fi

BASE="${BASE_URL:-http://127.0.0.1:8080}"
API="$BASE/api/v1"
ADMIN_LOGIN="${SMOKE_ADMIN_LOGIN:-admin}"
ADMIN_PASSWORD="${SMOKE_ADMIN_PASSWORD:-Admin@123456}"
STUDENT_LOGIN="${SMOKE_STUDENT_LOGIN:-student}"
STUDENT_PASSWORD="${SMOKE_STUDENT_PASSWORD:-Student@123456}"
NEW_USER="smoke$RANDOM$RANDOM"
NEW_PASSWORD="Smoke@123456"

pass=0
fail=0

ok() { printf '  \033[32m✓\033[0m %s\n' "$1"; pass=$((pass + 1)); }
ng() { printf '  \033[31m✗\033[0m %s\n' "$1"; fail=$((fail + 1)); }
section() { printf '\n\033[1m==> %s\033[0m\n' "$1"; }

# jget <json> <python 表达式>
#
# 表达式经环境变量传入，而不是直接拼进 Python 源码：调用方传 "['code']" 时，
# 拼接会生成 eval('d['code']') 这种引号嵌套的非法语法，解析会静默失败。
# 显式按 UTF-8 解码 stdin，避免在非 UTF-8 locale（如 Windows 默认）下解析失败。
jget() {
  JGET_EXPR="$2" python3 -c '
import json, os, sys
d = json.loads(sys.stdin.buffer.read().decode("utf-8"))
print(eval("d" + os.environ["JGET_EXPR"]))
' <<<"$1"
}

expect() { # expect <实际> <期望> <说明>
  if [ "$1" = "$2" ]; then ok "$3"; else ng "$3（期望 $2，实际 ${1:-空}）"; fi
}

# ---------------------------------------------------------------- 健康检查
section "健康检查"

HEALTH=$(curl -s "$BASE/health")
expect "$(jget "$HEALTH" "['code']")" "0" "GET /health 返回成功信封"
expect "$(jget "$HEALTH" "['data']['service']")" "yanmai-api" "/health 返回服务名"
expect "$(jget "$HEALTH" "['data']['database']['status']")" "up" "数据库连通"

READY_CODE=$(curl -s -o /dev/null -w '%{http_code}' "$BASE/health/ready")
expect "$READY_CODE" "200" "GET /health/ready 就绪"

META=$(curl -s "$API/meta/modules")
expect "$(jget "$META" "['data']['modules'].__len__()")" "12" "已挂载 12 个业务模块"

# ---------------------------------------------------------------- 登录
section "登录与令牌"

LOGIN=$(curl -s -X POST "$API/auth/login" -H 'Content-Type: application/json' \
  -d "{\"login\":\"$ADMIN_LOGIN\",\"password\":\"$ADMIN_PASSWORD\"}")
ACCESS=$(jget "$LOGIN" "['data']['accessToken']")
REFRESH=$(jget "$LOGIN" "['data']['refreshToken']")
expect "$(jget "$LOGIN" "['data']['user']['role']")" "admin" "管理员登录成功"

if [ -n "$ACCESS" ]; then ok "签发 accessToken"; else ng "未签发 accessToken：$LOGIN"; fi
if [ -n "$REFRESH" ]; then ok "签发 refreshToken"; else ng "未签发 refreshToken"; fi

BAD_CODE=$(curl -s -o /dev/null -w '%{http_code}' -X POST "$API/auth/login" \
  -H 'Content-Type: application/json' \
  -d "{\"login\":\"$ADMIN_LOGIN\",\"password\":\"definitely-wrong\"}")
expect "$BAD_CODE" "401" "密码错误返回 401"

BAD_BODY=$(curl -s -X POST "$API/auth/login" -H 'Content-Type: application/json' \
  -d "{\"login\":\"$ADMIN_LOGIN\",\"password\":\"definitely-wrong\"}")
expect "$(jget "$BAD_BODY" "['code']")" "20003" "密码错误业务码 20003"

NOBODY=$(curl -s -X POST "$API/auth/login" -H 'Content-Type: application/json' -d "{\"login\":\"ab\"}")
expect "$(jget "$NOBODY" "['code']")" "10001" "参数校验失败业务码 10001"

# ---------------------------------------------------------------- 鉴权
section "鉴权与会话"

if [ -z "$ACCESS" ]; then
  ng "缺少 accessToken，跳过鉴权用例"
else
  ME=$(curl -s "$API/users/me" -H "Authorization: Bearer $ACCESS")
  expect "$(jget "$ME" "['data']['username']")" "$ADMIN_LOGIN" "GET /users/me 返回当前账号"

  if echo "$ME" | grep -q 'passwordHash'; then
    ng "响应中泄露了密码哈希"
  else
    ok "响应未泄露密码哈希"
  fi

  LIST=$(curl -s "$API/admin/users?page=1&pageSize=10" -H "Authorization: Bearer $ACCESS")
  if [ "$(jget "$LIST" "['data']['total']")" -ge 3 ] 2>/dev/null; then
    ok "GET /admin/users 返回种子账号（total ≥ 3）"
  else
    ng "GET /admin/users 数据异常：$LIST"
  fi

  MISMATCH_CODE=$(curl -s -o /dev/null -w '%{http_code}' "$API/users/me" \
    -H "Authorization: Bearer $REFRESH")
  expect "$MISMATCH_CODE" "401" "Refresh Token 不可访问受保护接口"

  MISMATCH_BODY=$(curl -s "$API/users/me" -H "Authorization: Bearer $REFRESH")
  expect "$(jget "$MISMATCH_BODY" "['code']")" "20008" "令牌类型不匹配业务码 20008"

  REFRESHED=$(curl -s -X POST "$API/auth/refresh" -H 'Content-Type: application/json' \
    -d "{\"refreshToken\":\"$REFRESH\"}")
  NEW_ACCESS=$(jget "$REFRESHED" "['data']['accessToken']")
  if [ -n "$NEW_ACCESS" ]; then ok "POST /auth/refresh 签发新令牌"; else ng "刷新令牌失败：$REFRESHED"; fi

  NEW_ME=$(curl -s "$API/auth/me" -H "Authorization: Bearer $NEW_ACCESS")
  expect "$(jget "$NEW_ME" "['data']['username']")" "$ADMIN_LOGIN" "新 Access Token 可访问 /auth/me"
fi

NO_TOKEN_CODE=$(curl -s -o /dev/null -w '%{http_code}' "$API/users/me")
expect "$NO_TOKEN_CODE" "401" "无令牌访问受保护接口返回 401"

NO_TOKEN_BODY=$(curl -s "$API/users/me")
expect "$(jget "$NO_TOKEN_BODY" "['code']")" "20005" "缺令牌业务码 20005"

# ---------------------------------------------------------------- 注册
section "注册"

REG=$(curl -s -X POST "$API/auth/register" -H 'Content-Type: application/json' \
  -d "{\"username\":\"$NEW_USER\",\"password\":\"$NEW_PASSWORD\",\"nickname\":\"冒烟用户\"}")
expect "$(jget "$REG" "['data']['user']['role']")" "student" "注册默认为学员角色"

REG_AGAIN=$(curl -s -X POST "$API/auth/register" -H 'Content-Type: application/json' \
  -d "{\"username\":\"$NEW_USER\",\"password\":\"$NEW_PASSWORD\"}")
expect "$(jget "$REG_AGAIN" "['code']")" "20002" "重复注册业务码 20002"

# ---------------------------------------------------------------- RBAC
section "RBAC 角色守卫"

STUDENT_LOGIN_RES=$(curl -s -X POST "$API/auth/login" -H 'Content-Type: application/json' \
  -d "{\"login\":\"$STUDENT_LOGIN\",\"password\":\"$STUDENT_PASSWORD\"}")
STUDENT_TOKEN=$(jget "$STUDENT_LOGIN_RES" "['data']['accessToken']")

if [ -n "$STUDENT_TOKEN" ]; then
  ok "学员账号登录成功"

  FORBIDDEN_CODE=$(curl -s -o /dev/null -w '%{http_code}' "$API/admin/users" \
    -H "Authorization: Bearer $STUDENT_TOKEN")
  expect "$FORBIDDEN_CODE" "403" "学员访问管理接口返回 403"

  FORBIDDEN_BODY=$(curl -s "$API/admin/users" -H "Authorization: Bearer $STUDENT_TOKEN")
  expect "$(jget "$FORBIDDEN_BODY" "['code']")" "20009" "权限不足业务码 20009"

  SELF=$(curl -s "$API/users/me" -H "Authorization: Bearer $STUDENT_TOKEN")
  expect "$(jget "$SELF" "['data']['username']")" "$STUDENT_LOGIN" "学员可访问自身资料"
else
  ng "学员账号登录失败：$STUDENT_LOGIN_RES"
fi

# ---------------------------------------------------------------- 错误处理
section "统一错误处理"

PLACEHOLDER=$(curl -s "$API/course")
expect "$(jget "$PLACEHOLDER" "['code']")" "10008" "占位模块业务码 10008"

PLACEHOLDER_CODE=$(curl -s -o /dev/null -w '%{http_code}' "$API/course")
expect "$PLACEHOLDER_CODE" "501" "占位模块返回 501"

NOT_FOUND=$(curl -s "$API/not-exist")
expect "$(jget "$NOT_FOUND" "['code']")" "10002" "未知路由业务码 10002"

METHOD_CODE=$(curl -s -o /dev/null -w '%{http_code}' -X DELETE "$API/auth/login")
expect "$METHOD_CODE" "405" "不支持的方法返回 405"

if [ -n "$(jget "$HEALTH" "['traceId']")" ]; then ok "响应携带 traceId"; else ng "响应缺少 traceId"; fi

# ---------------------------------------------------------------- 汇总
printf '\n\033[1m结果：%d 通过，%d 失败\033[0m\n' "$pass" "$fail"
[ "$fail" -eq 0 ] || exit 1

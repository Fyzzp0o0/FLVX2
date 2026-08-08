package api

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// JWT 手写实现,兼容原 Java JwtUtil:
//   - Header {"alg":"HmacSHA256","typ":"JWT"}(非标准 HS256 字符串,验证只重算签名)
//   - 三段 Base64URL(无 padding)拼接 header.payload.signature
//   - 签名 = HMAC-SHA256(secret, "header.payload")
//   - Claims: sub(userId 字符串) / user / name / role_id / iat / exp
//   - 有效期 90 天(原代码 90L*24*60*60*1000)
const jwtExpire = 90 * 24 * time.Hour

var b64 = base64.RawURLEncoding

type jwtHeader struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}

type jwtClaims struct {
	Sub    string `json:"sub"`
	User   string `json:"user"`
	Name   string `json:"name"`
	RoleID int64  `json:"role_id"`
	Iat    int64  `json:"iat"`
	Exp    int64  `json:"exp"`
}

// GenerateToken 签发 JWT(与原 Java 端逐字段对齐)
func GenerateToken(secret string, userID int64, username string, roleID int64) (string, error) {
	now := time.Now().Unix()
	headerJSON, _ := json.Marshal(jwtHeader{Alg: "HmacSHA256", Typ: "JWT"})
	claims := jwtClaims{
		Sub:    strconv.FormatInt(userID, 10),
		User:   username,
		Name:   username,
		RoleID: roleID,
		Iat:    now,
		Exp:    now + int64(jwtExpire.Seconds()),
	}
	claimsJSON, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}
	header := b64.EncodeToString(headerJSON)
	payload := b64.EncodeToString(claimsJSON)
	signingInput := header + "." + payload
	sig := sign(secret, signingInput)
	return signingInput + "." + sig, nil
}

// ValidateToken 校验 JWT:三段结构 + 签名相等 + exp 未过期
func ValidateToken(secret, token string) (userID int64, username string, roleID int64, ok bool) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return 0, "", 0, false
	}
	signingInput := parts[0] + "." + parts[1]
	if !hmac.Equal([]byte(sign(secret, signingInput)), []byte(parts[2])) {
		return 0, "", 0, false
	}
	payloadJSON, err := b64.DecodeString(parts[1])
	if err != nil {
		return 0, "", 0, false
	}
	var claims jwtClaims
	if err := json.Unmarshal(payloadJSON, &claims); err != nil {
		return 0, "", 0, false
	}
	if claims.Exp <= time.Now().Unix() {
		return 0, "", 0, false
	}
	uid, err := strconv.ParseInt(claims.Sub, 10, 64)
	if err != nil {
		return 0, "", 0, false
	}
	return uid, claims.User, claims.RoleID, true
}

func sign(secret, input string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(input))
	return b64.EncodeToString(mac.Sum(nil))
}

// md5Hex 无盐 MD5 十六进制小写(兼容原 Md5Util.md5)
func md5Hex(s string) string {
	sum := md5.Sum([]byte(s))
	return fmt.Sprintf("%x", sum)
}

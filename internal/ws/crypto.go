package ws

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"sync"
	"time"
)

// AES-256-GCM 加密,与 agent(go-gost/panel/crypto/aes.go)与 Java AESCrypto 双端对齐:
//   - 密钥 = SHA-256(secret) 取 32 字节
//   - nonce = 12 随机字节;ciphertext = GCM-Seal(nonce, plaintext, aad=nil)
//   - 载荷 = nonce‖ciphertext‖tag;data = Base64.StdEncoding(载荷)
//   - 包装: {"encrypted":true,"data":"<base64>","timestamp":...}(timestamp 不校验)
type encryptedEnvelope struct {
	Encrypted bool   `json:"encrypted"`
	Data      string `json:"data"`
	Timestamp any    `json:"timestamp,omitempty"`
}

// Crypto 按 secret 派生密钥的加解密器(带缓存)
type Crypto struct {
	mu     sync.Mutex
	cache  map[string]cipher.AEAD // secret → gcm
}

func NewCrypto() *Crypto {
	return &Crypto{cache: make(map[string]cipher.AEAD)}
}

func (c *Crypto) gcm(secret string) (cipher.AEAD, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if g, ok := c.cache[secret]; ok {
		return g, nil
	}
	key := sha256.Sum256([]byte(secret))
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}
	g, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	c.cache[secret] = g
	return g, nil
}

// Encrypt 加密明文,返回包装 JSON 字节
func (c *Crypto) Encrypt(secret string, plaintext []byte) ([]byte, error) {
	g, err := c.gcm(secret)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, g.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}
	sealed := g.Seal(nonce, nonce, plaintext, nil) // nonce‖ciphertext‖tag
	env := encryptedEnvelope{Encrypted: true, Data: base64.StdEncoding.EncodeToString(sealed), Timestamp: nowUnixMillis()}
	return json.Marshal(env)
}

// Decrypt 解密包装 JSON;非加密包装返回原始字节(明文)
func (c *Crypto) Decrypt(secret string, raw []byte) ([]byte, error) {
	var env encryptedEnvelope
	if err := json.Unmarshal(raw, &env); err != nil {
		return raw, nil // 非 JSON 或非包装 → 明文
	}
	if !env.Encrypted || env.Data == "" {
		return raw, nil // 明文
	}
	sealed, err := base64.StdEncoding.DecodeString(env.Data)
	if err != nil {
		return nil, err
	}
	g, err := c.gcm(secret)
	if err != nil {
		return nil, err
	}
	nonceSize := g.NonceSize()
	if len(sealed) < nonceSize {
		return nil, err
	}
	return g.Open(nil, sealed[:nonceSize], sealed[nonceSize:], nil)
}

func nowUnixMillis() int64 {
	return time.Now().UnixMilli()
}

package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
)

var (
	ErrEmptyJWTSigningKey   = errors.New("JWT_SIGNING_KEY is empty")
	ErrEmptyAccessTokenTTL  = errors.New("empty accessTokenTTL for JWT")
	ErrEmptyRefreshTokenTTL = errors.New("empty refreshTokenTTL for JWT")

	ErrEmptySub     = errors.New("empty sub")
	ErrEmptyRole    = errors.New("empty role")
	ErrEmptyExp     = errors.New("empty exp")
	ErrExpiredToken = errors.New("token has expired")
)

type TokenManager interface {
	NewJWT(userID int, role int) (string, error)
	NewRefreshToken() string
	Parse(token string) (int, int, error)
}

type Manager struct {
	jwtSigningKey   string
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func NewManager(jwtSigningKey string, accessTokenTTL, refreshTokenTTL time.Duration) (*Manager, error) {
	if jwtSigningKey == "" {
		return nil, ErrEmptyJWTSigningKey
	}

	if accessTokenTTL == 0 {
		return nil, ErrEmptyAccessTokenTTL
	}

	if refreshTokenTTL == 0 {
		return nil, ErrEmptyRefreshTokenTTL
	}

	return &Manager{
		jwtSigningKey:   jwtSigningKey,
		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,
	}, nil
}

func (m *Manager) NewJWT(userID int, role int) (string, error) {
	header := map[string]interface{}{
		"alg": "HS256",
		"typ": "JWT",
	}
	headerJSON, _ := json.Marshal(header)
	headerBase64 := base64.RawURLEncoding.EncodeToString(headerJSON)

	payload := map[string]interface{}{
		"sub":  fmt.Sprintf("%v", userID),
		"exp":  fmt.Sprintf("%v", time.Now().Add(m.accessTokenTTL).Unix()),
		"role": fmt.Sprintf("%v", role),
	}
	payloadJSON, _ := json.Marshal(payload)
	payloadBase64 := base64.RawURLEncoding.EncodeToString(payloadJSON)

	signature := sign(headerBase64+"."+payloadBase64, m.jwtSigningKey)
	return headerBase64 + "." + payloadBase64 + "." + signature, nil
}

func (m *Manager) NewRefreshToken() string {
	token := uuid.NewV4()
	return token.String()
}

func (m *Manager) Parse(token string) (int, int, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return -1, -1, errors.New("invalid token format")
	}

	headerBase64, payloadBase64, signature := parts[0], parts[1], parts[2]
	if !verify(headerBase64+"."+payloadBase64, signature, m.jwtSigningKey) {
		return -1, -1, errors.New("invalid signature")
	}

	payloadJSON, _ := base64.RawURLEncoding.DecodeString(payloadBase64)
	var payload map[string]interface{}
	json.Unmarshal(payloadJSON, &payload)

	sub, err := strconv.Atoi(fmt.Sprintf("%v", payload["sub"]))
	if err != nil {
		return -1, -1, err
	}

	role, err := strconv.Atoi(fmt.Sprintf("%v", payload["role"]))
	if err != nil {
		return -1, -1, err
	}

	exp, err := strconv.ParseInt(fmt.Sprintf("%v", payload["exp"]), 10, 64)
	if err != nil {
		return -1, -1, err
	}
	tm := time.Unix(exp, 0)

	if time.Now().After(tm) {
		return sub, role, ErrExpiredToken
	}

	return sub, role, nil
}

func sign(data, key string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(data))
	return base64.RawURLEncoding.EncodeToString(h.Sum(nil))
}

func verify(data, signature, key string) bool {
	expectedSignature := sign(data, key)
	return hmac.Equal([]byte(signature), []byte(expectedSignature))
}

package service

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/ucho456job/lgtmeme/config"
	"github.com/ucho456job/lgtmeme/internal/auth/model"
)

type JwtService interface {
	GetPublicKeys() ([]byte, error)
	GenerateAccessToken(userID *uuid.UUID, oauthClient *model.OauthClient, expiresIn time.Duration) (string, error)
	GenerateRefreshToken() (string, error)
	GenerateIDToken(oauthClient *model.OauthClient, user *model.User, nonce string) (string, error)
}

type jwtService struct{}

func NewJwtService() JwtService {
	return &jwtService{}
}

type CustomClaims struct {
	jwt.RegisteredClaims
	Scope string `json:"scope,omitempty"`
}

func (s *jwtService) loadPrivateKey() (*rsa.PrivateKey, error) {
	base64Key := os.Getenv("JWT_PRIVATE_KEY_BASE64")
	keyBytes, err := base64.StdEncoding.DecodeString(base64Key)
	if err != nil {
		return nil, err
	}
	return jwt.ParseRSAPrivateKeyFromPEM(keyBytes)
}

func (s *jwtService) loadPublicKey() (*rsa.PublicKey, error) {
	base64Key := os.Getenv("JWT_PUBLIC_KEY_BASE64")
	keyBytes, err := base64.StdEncoding.DecodeString(base64Key)
	if err != nil {
		return nil, err
	}
	return jwt.ParseRSAPublicKeyFromPEM(keyBytes)
}

func (s *jwtService) GetPublicKeys() ([]byte, error) {
	publicKey, err := s.loadPublicKey()
	if err != nil {
		return nil, err
	}

	keySet := jwk.NewSet()
	key, err := jwk.New(publicKey)
	if err != nil {
		return nil, err
	}

	if err := key.Set(jwk.KeyIDKey, "1"); err != nil {
		return nil, err
	}
	if err := key.Set(jwk.AlgorithmKey, "RS256"); err != nil {
		return nil, err
	}
	if err := key.Set(jwk.KeyUsageKey, "sig"); err != nil {
		return nil, err
	}

	keySet.Add(key)

	jsonKeySet, err := json.MarshalIndent(keySet, "", "  ")
	if err != nil {
		return nil, err
	}

	return jsonKeySet, nil
}

func (s *jwtService) GenerateAccessToken(userID *uuid.UUID, oauthClient *model.OauthClient, expiresIn time.Duration) (string, error) {
	scopes := []string{}
	for _, scope := range oauthClient.Scopes {
		scopes = append(scopes, scope.Code)
	}
	scopesStr := strings.Join(scopes, " ")

	claims := jwt.MapClaims{
		"aud":   oauthClient.ApplicationURL,
		"azp":   oauthClient.ClientID,
		"scope": scopesStr,
		"iss":   os.Getenv("BASE_URL"),
		"exp":   time.Now().Add(expiresIn).Unix(),
		"iat":   time.Now().Unix(),
		"jti":   uuid.New().String(),
	}

	if userID != nil {
		claims["sub"] = userID.String()
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	privateKey, err := s.loadPrivateKey()
	if err != nil {
		return "", err
	}

	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *jwtService) GenerateRefreshToken() (string, error) {
	bytes := make([]byte, config.REFRESH_TOKEN_SIZE)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (s *jwtService) GenerateIDToken(oauthClient *model.OauthClient, user *model.User, nonce string) (string, error) {
	claims := jwt.MapClaims{
		"aud":   oauthClient.ApplicationURL,
		"sub":   user.ID,
		"azp":   oauthClient.ClientID,
		"iss":   os.Getenv("BASE_URL"),
		"exp":   time.Now().Add(config.ID_TOKEN_EXPIRES_IN).Unix(),
		"iat":   time.Now().Unix(),
		"nonce": nonce,
		"name":  user.Name,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	privateKey, err := s.loadPrivateKey()
	if err != nil {
		return "", err
	}

	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

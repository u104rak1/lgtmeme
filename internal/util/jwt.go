package util

import (
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/ucho456job/my_authn_authz/internal/model"
)

type JwtService interface {
	GenerateAccessToken(userID model.UserID, oauthClient *model.OauthClient, expiresIn time.Duration) (string, error)
	GenerateRefreshToken(userID model.UserID) (string, error)
	GenerateIDToken(oauthClient *model.OauthClient, user *model.User, nonce string) (string, error)
}

type jwtService struct {
	jwtKey    []byte
	issuerURL string
}

func NewJwtService(jwtKey []byte, issuerURL string) JwtService {
	return &jwtService{
		jwtKey:    jwtKey,
		issuerURL: issuerURL,
	}
}

type CustomClaims struct {
	jwt.RegisteredClaims
	Scope string `json:"scope,omitempty"`
}

func (s *jwtService) GenerateAccessToken(userID model.UserID, oauthClient *model.OauthClient, expiresIn time.Duration) (string, error) {
	scopes := []string{}
	for _, scope := range oauthClient.Scopes {
		scopes = append(scopes, scope.Code)
	}
	scopesStr := strings.Join(scopes, " ")

	claims := jwt.MapClaims{
		"aud":   oauthClient.ApplicationURL,
		"sub":   userID,
		"azp":   oauthClient.ClientID,
		"scope": scopesStr,
		"iss":   s.issuerURL,
		"exp":   time.Now().Add(expiresIn).Unix(),
		"iat":   time.Now().Unix(),
		"jti":   uuid.New().String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(s.jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *jwtService) GenerateRefreshToken(userID model.UserID) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"iss": s.issuerURL,
		"iat": time.Now().Unix(),
		"jti": uuid.New().String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *jwtService) GenerateIDToken(oauthClient *model.OauthClient, user *model.User, nonce string) (string, error) {
	claims := jwt.MapClaims{
		"aud":   oauthClient.ApplicationURL,
		"sub":   user.ID,
		"azp":   oauthClient.ClientID,
		"iss":   s.issuerURL,
		"exp":   time.Now().Add(ID_TOKEN_EXPIRES_IN).Unix(),
		"iat":   time.Now().Unix(),
		"nonce": nonce,
		"name":  user.Name,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(s.jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

package jwx

import (
	"context"
	"time"

	"github.com/lestrrat-go/jwx/v3/jwa"
	"github.com/lestrrat-go/jwx/v3/jwk"
	"github.com/lestrrat-go/jwx/v3/jwt"
)

const (
	defaultIssuer      = "rcs"
	defaultSubject     = "user"
	tokenExpiry        = 24 * 30 * time.Hour
	refreshTokenExpiry = 24 * 35 * time.Hour
	jwksCacheDuration  = 1 * time.Hour
)

type SelfTokenClaims struct {
	Uid      uint64 `json:"uid"`
	UserType string `json:"user_type"`
	Role     string `json:"role"`
}

type JwtConfig struct {
	TokenSecret string `yaml:"token_secret"`
	OauthJwkUrl string `yaml:"oauth_jwk_url"`
}

type JwtProvider interface {
	SignSelfToken(uid uint64, userType, role string) (string, string, error)
	ValidateSelfToken(tokenString string) (*SelfTokenClaims, error)
	ValidateOauthToken(ctx context.Context, tokenString string) (string, error)
}

type jwtProvider struct {
	secret        string
	jwkUrl        string
	idpIdentifier string
}

func NewJwtProvider(cfg *JwtConfig) JwtProvider {
	return &jwtProvider{
		secret: cfg.TokenSecret,
		jwkUrl: cfg.OauthJwkUrl,
	}
}

func (p *jwtProvider) SignSelfToken(uid uint64, userType, role string) (string, string, error) {
	now := time.Now()
	token, err := jwt.NewBuilder().
		Issuer(defaultIssuer).
		Subject(defaultSubject).
		Expiration(now.Add(tokenExpiry)).
		NotBefore(now).
		IssuedAt(now).
		Claim("uid", uid).
		Claim("user_type", userType).
		Claim("role", role).
		Build()
	if err != nil {
		return "", "", err
	}

	refreshToken, err := jwt.NewBuilder().
		Issuer(defaultIssuer).
		Subject(defaultSubject).
		Expiration(now.Add(refreshTokenExpiry)).
		NotBefore(now).
		IssuedAt(now).
		Claim("uid", uid).
		Claim("user_type", userType).
		Claim("role", role).
		Build()
	if err != nil {
		return "", "", err
	}

	signedToken, err := jwt.Sign(token, jwt.WithKey(jwa.HS256(), []byte(p.secret)))
	if err != nil {
		return "", "", err
	}
	signedRefresh, err := jwt.Sign(refreshToken, jwt.WithKey(jwa.HS256(), []byte(p.secret)))
	if err != nil {
		return "", "", err
	}

	return string(signedToken), string(signedRefresh), nil
}

func (p *jwtProvider) ValidateSelfToken(tokenString string) (*SelfTokenClaims, error) {
	token, err := jwt.ParseString(tokenString,
		jwt.WithKey(jwa.HS256(), []byte(p.secret)),
		jwt.WithIssuer(defaultIssuer),
		jwt.WithRequiredClaim("uid"),
		jwt.WithRequiredClaim("user_type"),
		jwt.WithRequiredClaim("role"),
	)
	if err != nil {
		return nil, err
	}

	claims := &SelfTokenClaims{}
	if err := token.Get("uid", &claims.Uid); err != nil {
		return nil, err
	}
	if err := token.Get("user_type", &claims.UserType); err != nil {
		return nil, err
	}
	if err := token.Get("role", &claims.Role); err != nil {
		return nil, err
	}

	return claims, nil
}

func (p *jwtProvider) ValidateOauthToken(ctx context.Context, tokenString string) (string, error) {
	set, err := jwk.Fetch(ctx, p.jwkUrl)
	if err != nil {
		return "", err
	}

	token, err := jwt.Parse([]byte(tokenString), jwt.WithKeySet(set),
		jwt.WithRequiredClaim("sub"),
		jwt.WithRequiredClaim("iat"),
		jwt.WithClaimValue("idp_identifier", "xx"))
	if err != nil {
		return "", err
	}

	sub, _ := token.Subject()
	return sub, nil
}

package jwt

import (
	"backoffice/internal/transport/http"
	"backoffice/pkg/auth"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type provider struct {
	cfg *Config
}

func NewProvider(cfg *Config) auth.Authorizer {
	return &provider{
		cfg: cfg,
	}
}

func (p *provider) Token(opts ...auth.TokenOption) (*auth.Token, error) {
	options := auth.NewTokenOptions(opts...)
	expiredAt := time.Now().Add(options.Expiry)

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ID:        options.ID,
		Issuer:    p.cfg.Issuer,
		ExpiresAt: jwt.NewNumericDate(expiredAt),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	})

	tok, err := t.SignedString([]byte(p.cfg.Fingerprint))
	if err != nil {
		return nil, err
	}

	return &auth.Token{
		Token:     tok,
		ExpiredAt: expiredAt,
		CreatedAt: time.Now(),
	}, nil
}

func (p *provider) Refresh(opts ...auth.TokenOption) (*auth.Auth, error) {
	options := auth.NewTokenOptions(opts...)

	secret := options.RefreshToken
	if len(options.Secret) > 0 {
		secret = options.Secret
	}

	if _, err := p.Inspect(secret); err != nil {
		return nil, err
	}

	access, err := p.Token(auth.WithExpiry(p.cfg.AccessTokenLifetime), auth.WithTokenID(options.ID))
	if err != nil {
		return nil, err
	}

	refresh, err := p.Token(auth.WithExpiry(p.cfg.RefreshTokenLifetime), auth.WithTokenID(options.ID))
	if err != nil {
		return nil, err
	}

	return &auth.Auth{
		CreatedAt:    access.CreatedAt,
		ExpiredAt:    access.ExpiredAt,
		AccessToken:  access.Token,
		RefreshToken: refresh.Token,
	}, nil
}

func (p *provider) Generate(opts ...auth.GenerateOption) (*auth.Auth, error) {
	options := auth.NewGenerateOptions(opts...)

	access, err := p.Token(auth.WithExpiry(p.cfg.AccessTokenLifetime), auth.WithTokenID(options.ID))
	if err != nil {
		return nil, err
	}

	refresh, err := p.Token(auth.WithExpiry(p.cfg.RefreshTokenLifetime))
	if err != nil {
		return nil, err
	}

	return &auth.Auth{
		CreatedAt:    access.CreatedAt,
		ExpiredAt:    access.ExpiredAt,
		AccessToken:  access.Token,
		RefreshToken: refresh.Token,
	}, nil
}

func (p *provider) Verify(ctx *gin.Context) (*string, error) {
	t := strings.Replace(ctx.GetHeader(p.cfg.HeaderName), p.cfg.HeaderScheme, "", -1)

	if t != "" {
		return p.Inspect(t)
	}

	token, ok := ctx.GetQuery("token")
	if !ok {
		return nil, http.ErrAuthHeaderIsRequired
	}

	return p.Inspect(token)
}

func (p *provider) Inspect(t string) (*string, error) {
	token, err := p.parse(t)
	if token != nil && token.Valid {
		return &token.Claims.(*jwt.RegisteredClaims).ID, nil
	}

	if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			return nil, jwt.ErrTokenExpired
		}
	}

	return nil, jwt.ErrTokenUnverifiable
}

func (p *provider) parse(t string) (token *jwt.Token, err error) {
	token, err = jwt.ParseWithClaims(t, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrTokenUnverifiable
		}

		return []byte(p.cfg.Fingerprint), nil
	})

	return token, err
}

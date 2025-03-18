package pgsql

import (
	"backoffice/internal/entities"
	"context"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type tokenRepository struct {
	conn *gorm.DB
}

func NewTokenRepository(conn *gorm.DB) *tokenRepository {
	return &tokenRepository{
		conn: conn,
	}
}

func (r *tokenRepository) Create(ctx context.Context, token *entities.Token) error {
	return r.conn.WithContext(ctx).Create(&token).Error
}

func (r *tokenRepository) GetByRefresh(ctx context.Context, t string) (token *entities.Token, err error) {
	if err = r.conn.WithContext(ctx).Where("refresh_token = ?", t).First(&token).Error; err != nil {
		return nil, err
	}

	return token, nil
}

func (r *tokenRepository) GetByAccountID(ctx context.Context, accountID uuid.UUID) (tokens []*entities.Token, err error) {
	if err = r.conn.WithContext(ctx).Where("account_id = ?", accountID.String()).Find(&tokens).Error; err != nil {
		return nil, err
	}

	return tokens, nil
}

func (r *tokenRepository) DeleteByAccess(ctx context.Context, t string) error {
	return r.conn.WithContext(ctx).Where("access_token = ?", t).Delete(&entities.Token{}).Error
}

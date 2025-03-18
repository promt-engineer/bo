package services

import (
	"backoffice/internal/entities"
	e "backoffice/internal/errors"
	"backoffice/internal/repositories"
	"backoffice/internal/repositories/redis"
	"context"
	"errors"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"time"
)

type SessionService struct {
	sessionRepository repositories.SessionRepository
}

func NewSessionService(sessionRepository repositories.SessionRepository) *SessionService {
	return &SessionService{
		sessionRepository: sessionRepository,
	}
}

func (s *SessionService) Get(ctx context.Context, id uuid.UUID) (*entities.Session, error) {
	return s.sessionRepository.Get(ctx, id)
}

func (s *SessionService) GetKey(ctx context.Context, accountID uuid.UUID) ([]uuid.UUID, error) {
	return s.sessionRepository.GetKeys(ctx, accountID)
}

func (s *SessionService) Delete(ctx context.Context, sessionID uuid.UUID) error {
	return s.sessionRepository.Delete(ctx, redis.SessionCacheKeyPrefix, sessionID)
}

func (s *SessionService) Create(ctx context.Context, key uuid.UUID, session *entities.Session, expiration time.Time) error {
	return s.sessionRepository.Create(ctx, key, session, expiration)
}

func (s *SessionService) SwitchOrganization(ctx context.Context, session *entities.Session, sessionID, organizationID uuid.UUID) error {
	session.OrganizationID = organizationID

	return s.sessionRepository.Update(ctx, sessionID, session)
}

func (s *SessionService) UpdateAccountInfo(ctx context.Context, account *entities.Account) error {
	sessions, err := s.sessionRepository.Find(ctx, account.ID)
	if err != nil {
		return err
	}

	for _, session := range sessions {
		session.Account = account
		if err = s.sessionRepository.Update(ctx, session.ID, session); err != nil {
			return err
		}
	}

	return nil
}

func (s *SessionService) DeleteAll(ctx context.Context, account *entities.Account) error {
	sessions, err := s.sessionRepository.Find(ctx, account.ID)
	if err != nil {
		if errors.Is(err, e.ErrEntityNotFound) {
			zap.S().Warnf("%v, on %v, request", err, account.ID)
		} else {
			return err
		}
	}

	for _, session := range sessions {
		if err = s.sessionRepository.Delete(ctx, redis.SessionCacheKeyPrefix, session.ID); err != nil {
			if errors.Is(err, e.ErrEntityNotFound) {
				zap.S().Warnf("%v, on %v, request", err, session.ID)

				continue
			}

			return err
		}
	}

	if err = s.sessionRepository.Delete(ctx, redis.AccountsCacheKeyPrefix, account.ID); err != nil {
		if errors.Is(err, e.ErrEntityNotFound) {
			zap.S().Warnf("%v, on %v, request", err, account.ID)

			return nil
		}

		return err
	}

	return nil
}

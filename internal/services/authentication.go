package services

import (
	"backoffice/internal/entities"
	"backoffice/internal/repositories"
	"backoffice/pkg/auth"
	"context"
	"errors"
	"github.com/google/uuid"
)

var (
	ErrNotValidPassword = errors.New("not valid password")
)

type AuthenticationService struct {
	authProvider        auth.Authorizer
	authorizeService    *AuthorizationService
	accountService      *AccountService
	sessionService      *SessionService
	organizationService *OrganizationService
	tokenRepository     repositories.TokenRepository
}

func NewAuthenticationService(authProvider auth.Authorizer, authorizeService *AuthorizationService, accountService *AccountService, sessionService *SessionService, organizationService *OrganizationService, tokenRepository repositories.TokenRepository) *AuthenticationService {
	return &AuthenticationService{
		authProvider:        authProvider,
		authorizeService:    authorizeService,
		accountService:      accountService,
		sessionService:      sessionService,
		tokenRepository:     tokenRepository,
		organizationService: organizationService,
	}
}

func (s *AuthenticationService) Authenticate(ctx context.Context, account *entities.Account) (*auth.Auth, error) {
	jti := uuid.New()
	tokens, err := s.authProvider.Generate(auth.WithSubject(account.ID.String()), auth.WithID(jti.String()))
	if err != nil {
		return nil, err
	}

	t := &entities.Token{
		ID:           jti,
		AccountID:    account.ID,
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiredAt:    tokens.ExpiredAt,
	}

	if err = s.tokenRepository.Create(ctx, t); err != nil {
		return nil, err
	}

	session := &entities.Session{
		ID:             jti,
		Account:        account,
		OrganizationID: account.GetDefaultOrganizationID(),
	}

	if err = s.sessionService.Create(ctx, jti, session, t.ExpiredAt); err != nil {
		return nil, err
	}

	return tokens, nil
}

func (s *AuthenticationService) Refresh(ctx context.Context, rt string) (*auth.Auth, error) {
	token, err := s.tokenRepository.GetByRefresh(ctx, rt)
	if err != nil {
		return nil, err
	}

	jti := uuid.New()
	tokens, err := s.authProvider.Refresh(auth.WithRefreshToken(rt), auth.WithTokenID(jti.String()))
	if err != nil {
		return nil, err
	}

	if err = s.tokenRepository.DeleteByAccess(ctx, token.AccessToken); err != nil {
		return nil, err
	}

	t := &entities.Token{
		ID:           jti,
		AccountID:    token.AccountID,
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiredAt:    tokens.ExpiredAt,
	}

	if err = s.tokenRepository.Create(ctx, t); err != nil {
		return nil, err
	}

	return tokens, nil
}

func (s *AuthenticationService) Logout(ctx context.Context, sessionID uuid.UUID, token string) (err error) {
	if err = s.tokenRepository.DeleteByAccess(ctx, token); err != nil {
		return
	}

	if err = s.sessionService.Delete(ctx, sessionID); err != nil {
		return
	}

	return
}

func (s *AuthenticationService) LogoutAllByRole(ctx context.Context, roleID string) (errs []error) {
	accounts, err := s.accountService.FindAllByRoleID(ctx, roleID)
	if err != nil {
		errs = append(errs, err)

		return errs
	}

	for _, account := range accounts {
		tokens, err := s.tokenRepository.GetByAccountID(ctx, account.ID)
		if err != nil {
			errs = append(errs, err)

			continue
		}

		for _, token := range tokens {
			if err = s.tokenRepository.DeleteByAccess(ctx, token.AccessToken); err != nil {
				errs = append(errs, err)

				continue
			}
		}
	}

	for _, account := range accounts {
		sessionIDs, err := s.sessionService.GetKey(ctx, account.ID)
		if err != nil {
			errs = append(errs, err)

			continue
		}

		for _, sessionID := range sessionIDs {
			if err = s.sessionService.Delete(ctx, sessionID); err != nil {
				errs = append(errs, err)

				continue
			}
		}
	}

	return errs
}

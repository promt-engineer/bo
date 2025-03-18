package services

import (
	"backoffice/internal/entities"
	e "backoffice/internal/errors"
	"backoffice/internal/repositories"
	"backoffice/pkg/totp"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var (
	ErrCanNotRemoveAccount = errors.New("can't remove account")
)

type AccountService struct {
	accountRepository repositories.AccountRepository
	sessionService    *SessionService
	mailingService    *MailingService
}

func NewAccountService(accountRepository repositories.AccountRepository, sessionService *SessionService, mailingService *MailingService) *AccountService {
	return &AccountService{
		accountRepository: accountRepository,
		sessionService:    sessionService,
		mailingService:    mailingService,
	}
}

func (s *AccountService) FindBy(ctx context.Context, params map[string]interface{}) (*entities.Account, error) {
	return s.accountRepository.FindBy(ctx, params)
}

func (s *AccountService) FindAllByRoleID(ctx context.Context, roleID string) ([]*entities.Account, error) {
	return s.accountRepository.FindAllByRoleID(ctx, roleID)
}

func (s *AccountService) Update(ctx context.Context, id, firstName, lastName string) (*entities.Account, error) {
	ac, err := s.FindBy(ctx, map[string]interface{}{"id": id})
	if err != nil {
		return nil, err
	}

	account := &entities.Account{
		ID:                ac.ID,
		AuthProvider:      entities.AuthProviderEmail,
		AuthProviderID:    ac.AuthProviderID,
		AuthProviderToken: ac.AuthProviderToken,
		FirstName:         firstName,
		LastName:          lastName,
		Email:             ac.Email,
	}

	return s.accountRepository.Update(ctx, account)
}

func (s *AccountService) Paginate(ctx context.Context, organizationID uuid.UUID, filters map[string]interface{}, order string, limit int, offset int) ([]*entities.Account, int64, error) {
	p, total, err := s.accountRepository.Paginate(ctx, organizationID, repositories.Filters{Where: filters}, order, limit, offset)
	if err != nil {
		return p, total, err
	}

	p = lo.Filter(p, func(account *entities.Account, index int) bool {
		if lo.ContainsBy(account.Roles, func(role *entities.Role) bool {
			return role.Type == entities.RootRoleTypeName
		}) {
			total--
			return false
		}

		return true
	})

	return p, total, nil
}

func (s *AccountService) Create(ctx context.Context, authProviderID, authProviderToken, firstName, lastName string) (*entities.Account, error) {
	_, err := s.accountRepository.FindBy(ctx, map[string]interface{}{"auth_provider_id": authProviderID})
	if err == nil {
		return nil, e.ErrAccountAlreadyExists
	}

	if err != nil && !errors.Is(err, e.ErrEntityNotFound) {
		return nil, err
	}

	tokenHash, err := bcrypt.GenerateFromPassword([]byte(authProviderToken), 15)
	if err != nil {
		return nil, err
	}

	account := &entities.Account{
		ID:                uuid.New(),
		AuthProvider:      entities.AuthProviderEmail,
		AuthProviderID:    authProviderID,
		AuthProviderToken: string(tokenHash),
		FirstName:         firstName,
		LastName:          lastName,
		Email:             authProviderID,
	}

	account, err = s.accountRepository.Create(ctx, account)
	if err != nil {
		return account, err
	}

	return account, s.mailingService.NotifyUserEmail(authProviderID, authProviderID, authProviderToken)
}

func (s *AccountService) Delete(ctx context.Context, accountDeleter *entities.Account, accountID string) error {
	account, err := s.accountRepository.FindBy(ctx, map[string]interface{}{"id": accountID})
	if err != nil {
		return err
	}

	isRoot := accountDeleter.IsRoot()

	for _, role := range account.Roles {
		if entities.RootRoleTypeName == role.Type || (!isRoot && role.Type == entities.AdminRoleTypeName) {
			return fmt.Errorf("%s with role type: %v", ErrCanNotRemoveAccount, role.Type)
		}
	}

	if err = s.sessionService.DeleteAll(ctx, account); err != nil {
		return err
	}

	return s.accountRepository.Delete(ctx, account)
}

func (s *AccountService) UpdateTOTPSecret(ctx context.Context, account *entities.Account, secret, url string) (*entities.Account, error) {
	account.TOTPSecret = secret
	account.TOTPURL = url

	account, err := s.accountRepository.Update(ctx, &entities.Account{
		ID: account.ID,
		Params: &totp.Params{
			TOTPSecret: secret,
			TOTPURL:    url,
		},
	})

	if err != nil {
		return nil, err
	}

	if err = s.sessionService.UpdateAccountInfo(ctx, account); err != nil {
		return nil, err
	}

	return account, nil
}

func (s *AccountService) EnableTOTP(ctx context.Context, account *entities.Account) (*entities.Account, error) {
	account.TOTPEnabled = true

	account, err := s.accountRepository.Update(ctx, &entities.Account{
		ID: account.ID,
		Params: &totp.Params{
			TOTPEnabled: true,
		},
	})
	if err != nil {
		return nil, err
	}

	if err = s.sessionService.UpdateAccountInfo(ctx, account); err != nil {
		return nil, err
	}

	return account, nil
}

func (s *AccountService) DisableTOTP(ctx context.Context, account *entities.Account) (*entities.Account, error) {
	account.TOTPEnabled = false
	account.TOTPSecret = ""
	account.TOTPURL = ""

	account, err := s.accountRepository.DisableTOTP(ctx, account)
	if err != nil {
		return nil, err
	}

	if err = s.sessionService.UpdateAccountInfo(ctx, account); err != nil {
		return nil, err
	}

	return account, nil
}

func (s *AccountService) Auth(ctx context.Context, id string, token string) (*entities.Account, error) {
	account, err := s.FindBy(ctx, map[string]interface{}{"auth_provider_id": id})
	if err != nil {
		return nil, err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(account.AuthProviderToken), []byte(token)); err != nil {
		return nil, ErrNotValidPassword
	}

	return account, nil
}

func (s *AccountService) ChangePassword(ctx context.Context, account *entities.Account, password, newPassword string) error {
	if _, err := s.Auth(ctx, account.AuthProviderID, password); err != nil {
		return err
	}

	tokenHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), 15)
	if err != nil {
		return err
	}

	if _, err = s.accountRepository.Update(ctx, &entities.Account{ID: account.ID, AuthProviderToken: string(tokenHash)}); err != nil {
		return err
	}

	return nil
}

func (s *AccountService) ResetPasswordRequest(ctx context.Context, email string) error {
	account, err := s.FindBy(ctx, map[string]interface{}{"email": email})
	if err != nil {
		return err
	}

	resetToken := uuid.New().String()
	resetTokenExpiry := time.Now().Add(1 * time.Hour)

	account.ResetPasswordToken = resetToken
	account.ResetPasswordExpiresAt = &resetTokenExpiry

	_, err = s.accountRepository.Update(ctx, account)
	if err != nil {
		return err
	}

	return s.mailingService.ResetUserPassword(email, resetToken)
}

func (s *AccountService) ResetPassword(ctx context.Context, token, newPassword string) error {
	account, err := s.FindBy(ctx, map[string]interface{}{"reset_password_token": token})
	if err != nil {
		return err
	}

	if account.ResetPasswordExpiresAt != nil && account.ResetPasswordExpiresAt.After(time.Now()) {
		tokenHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), 15)
		if err != nil {
			return err
		}

		account.AuthProviderToken = string(tokenHash)
		account.ResetPasswordToken = ""
		account.ResetPasswordExpiresAt = nil

		_, err = s.accountRepository.Update(ctx, account)
		if err != nil {
			return err
		}

		return nil
	}

	return e.ErrAccountInvalidResetToken
}

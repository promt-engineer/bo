package errors

import (
	"errors"
	"fmt"
)

var (
	ErrCanNotRemoveOrganization    = errors.New("can not remove organization")
	ErrOrganizationIsNotIntegrator = errors.New("selected organization is not integrator")
	ErrOrganizationIsNotProvider   = errors.New("selected organization is not provider")
	ErrOrganizationIsNotOperator   = errors.New("selected organization is not operator")
	ErrDoesNotHavePermission       = errors.New("does not have permissions")

	ErrOrganizationNameMustBeUnique = errors.New("role slug must be unique")
	ErrOrganizationInUse            = errors.New("cannot delete assigned organization")
	ErrOrganizationAlreadyAssigned  = errors.New("this organization is already assigned to the user")
	ErrOperatorAlreadyAssigned      = errors.New("this operator is already assigned to the user")

	ErrAccountAlreadyExists     = errors.New("account with provided credentials already exists")
	ErrAccountInvalidResetToken = errors.New("invalid or expired password reset token")

	ErrRoleInUse           = errors.New("cannot delete assigned role")
	ErrRoleAlreadyAssigned = errors.New("this role is already assigned to the user")

	ErrNotAuthorized = errors.New("unauthorized")

	ErrDefaultWagerOutOfList = errors.New("default wager out of wager levels")
	ErrNegativeWager         = errors.New("wager must be positive")

	ErrInternal = errors.New("internal error")

	ErrEntityNotFound     = errors.New("not found")
	ErrEntityAlreadyExist = errors.New("entity already exist")

	ErrGameNotExist = errors.New("the integrator's game does not exist")

	ErrValidationFailed = func(param string) error {
		return fmt.Errorf("validation failed on parameter: %s", param)
	}
)

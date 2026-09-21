package pkg

import (
	"errors"

	"github.com/Freelance-launchpad/backend-interview/common/jerror/v2"
)

var (
	// ErrUserNotFound is returned when a user id is not found
	ErrUserNotFound = jerror.NewNotFoundError("USER_NOT_FOUND")

	// ErrContractNotFound is raised when no contract is found for a user.
	ErrContractNotFound = jerror.NewNotFoundError("CONTRACT_NOT_FOUND")

	ErrUserAlreadyHasLifeOffer = jerror.NewAlreadyExistsError("USER_ALREADY_HAS_LIVE_LIFE_OFFER")

	ErrUserNotInRegion = jerror.NewBusinessError("USER_NOT_IN_REGION")

	// ErrUnexpectedStatus is an error returned when an unexpected status is returned by an external service.
	ErrUnexpectedStatus = errors.New("UNEXPECTED_STATUS")
)

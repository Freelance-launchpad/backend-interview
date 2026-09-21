package jgin

import (
	"context"

	"github.com/Freelance-launchpad/backend-interview/common/jentity"
	"github.com/Freelance-launchpad/backend-interview/common/jgin/keys"
	pkgusers "github.com/Freelance-launchpad/backend-interview/services/users/pkg"
	"github.com/google/uuid"
)

// GetUserIDFromContext returns the user id stored in the context, if no userID is set returns an empty string
func GetUserIDFromContext(ctx context.Context) string {
	value, _ := ctx.Value(keys.UserIDKey).(string)
	return value
}

// ContextWithUserID returns the given context with a user id.
func ContextWithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, keys.UserIDKey, userID)
}

// GetEmailVerifiedFromContext returns the emailVerified value stored in the context
func GetEmailVerifiedFromContext(ctx context.Context) bool {
	value, _ := ctx.Value(keys.EmailVerifiedKey).(bool)
	return value
}

// ContextWithEmailVerified returns the given context with the emailVerified property.
func ContextWithEmailVerified(ctx context.Context, emailVerified bool) context.Context {
	return context.WithValue(ctx, keys.EmailVerifiedKey, emailVerified)
}

// GetUserEmailFromContext returns the user email stored in the context, if no user email is set returns an empty string
//
// Deprecated
func GetUserEmailFromContext(ctx context.Context) string {
	value, _ := ctx.Value(keys.UserEmailKey).(string)
	return value
}

// ContextWithUserEmail returns the given context with a user email.
//
// Deprecated
func ContextWithUserEmail(ctx context.Context, userEmail string) context.Context {
	return context.WithValue(ctx, keys.UserEmailKey, userEmail)
}

// GetEntityFromContext returns the user entity stored in the context.
func GetEntityFromContext(ctx context.Context) (jentity.Entity, error) {
	value, ok := ctx.Value(keys.EntityKey).(jentity.Entity)
	if !ok {
		return "", jentity.ErrInvalidEntity
	}
	return value, nil
}

// ContextWithEntity returns the given context with the entity of the requester.
func ContextWithEntity(ctx context.Context, userEntity jentity.Entity) context.Context {
	return context.WithValue(ctx, keys.EntityKey, userEntity)
}

// GetOfferIDFromContext returns the offerID stored in the context, if there is no offerID the value will be empty.
func GetOfferIDFromContext(ctx context.Context) uuid.UUID {
	value, _ := ctx.Value(keys.OfferIDKey).(uuid.UUID)
	return value
}

// ContextWithOfferID returns the given context with the value offerID.
func ContextWithOfferID(ctx context.Context, offerID uuid.UUID) context.Context {
	return context.WithValue(ctx, keys.OfferIDKey, offerID)
}

// GetOfferTypeFromContext returns the offerID stored in the context, if there is no offerID the value will be empty.
func GetOfferTypeFromContext(ctx context.Context) jentity.OfferType {
	value, _ := ctx.Value(keys.OfferTypeKey).(jentity.OfferType)
	return value
}

func ContextWithOfferType(ctx context.Context, offerType jentity.OfferType) context.Context {
	return context.WithValue(ctx, keys.OfferTypeKey, offerType)
}

func GetOfferContractFromContext(ctx context.Context) pkgusers.OfferWithContract {
	value, _ := ctx.Value(keys.OfferContractKey).(pkgusers.OfferWithContract)
	return value
}

func ContextWithOfferContract(ctx context.Context, oc pkgusers.OfferWithContract) context.Context {
	return context.WithValue(ctx, keys.OfferContractKey, oc)
}

// ServiceContext returns the given context with a service key.
func ServiceContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, keys.IsServiceKey, true)
}

// IsServiceContext returns true if the given context contains a service key.
func IsServiceContext(ctx context.Context) bool {
	value, _ := ctx.Value(keys.IsServiceKey).(bool)
	return value
}

func ContextWithPermissions(ctx context.Context, permissions []string) context.Context {
	return context.WithValue(ctx, keys.PermissionsKey, permissions)
}

func GetPermissionsFromContext(ctx context.Context) []string {
	value, _ := ctx.Value(keys.PermissionsKey).([]string)
	return value
}

func ContextWithRefreshTokenID(ctx context.Context, refreshToken uuid.UUID) context.Context {
	return context.WithValue(ctx, keys.RefreshTokenIDKey, refreshToken)
}

func GetRefreshTokenIDFromContext(ctx context.Context) uuid.UUID {
	value, _ := ctx.Value(keys.RefreshTokenIDKey).(uuid.UUID)
	return value
}

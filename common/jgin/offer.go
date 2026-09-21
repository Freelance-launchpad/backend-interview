package jgin

import (
	"context"
	"fmt"
	"slices"

	"github.com/Freelance-launchpad/backend-interview/common/jentity"
	"github.com/Freelance-launchpad/backend-interview/common/jerror/v2"
	"github.com/Freelance-launchpad/backend-interview/common/jgin/keys"
	"github.com/Freelance-launchpad/backend-interview/common/jhttp/v2"
	pkgusers "github.com/Freelance-launchpad/backend-interview/services/users/pkg"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var (
	ErrInvalidOfferID = jerror.NewValidationError("INVALID_OFFER_ID")

	ErrUnauthorizedUser     = jerror.NewPermissionDeniedError("UNAUTHORIZED_USER")
	ErrUserNotLinkedToOffer = jerror.NewPermissionDeniedError("USER_NOT_LINKED_TO_OFFER")
	ErrOfferNotFound        = jerror.NewPermissionDeniedError("OFFER_NOT_FOUND")
	ErrInvalidRegion        = jerror.NewPermissionDeniedError("INVALID_REGION")
	ErrInvalidOfferEntity   = jerror.NewPermissionDeniedError("INVALID_OFFER_ENTITY")
	ErrMissingOfferContract = jerror.NewPermissionDeniedError("MISSING_OFFER_CONTRACT")
	ErrInvalidOfferType     = jerror.NewPermissionDeniedError("INVALID_OFFER_TYPE")
	ErrOfferNotActive       = jerror.NewPermissionDeniedError("OFFER_NOT_ACTIVE")
	ErrEmailNotVerified     = jerror.NewPermissionDeniedError("EMAIL_NOT_VERIFIED")
)

// HasRegion checks if the user belongs to at least one of the allowed regions.
//
// TODO: add regions to the JWT token?
func HasRegion(client interface {
	GetUserByID(ctx context.Context, userID string) (pkgusers.User, error)
}, allowedRegions ...jentity.Region) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		userID := GetUserIDFromContext(ctx)
		if userID == "" {
			HandleError(c, ErrUnauthorizedUser)
			return
		}

		user, err := client.GetUserByID(ctx, userID)
		if err != nil {
			HandleError(c, err)
			return
		}

		for _, userRegion := range user.Regions {
			if slices.Contains(allowedRegions, userRegion) {
				return
			}
		}
		HandleError(c, ErrInvalidRegion)
	}
}

// EnforceHasOfferID read the header 'offer_id' checks whether the offer belongs to the user
// and set the offer_id in the query context.
// It also set the entity in the context.
// It also enforces that the user email address is verified.
func EnforceHasOfferID(client interface {
	GetOfferWithContract(ctx context.Context, offerID uuid.UUID) (pkgusers.OfferWithContract, error)
}) gin.HandlerFunc {
	return func(c *gin.Context) {
		const errMsg = "EnforceHasOfferID has failed"

		ctx := c.Request.Context()

		userID := GetUserIDFromContext(ctx)
		if userID == "" {
			HandleError(c, ErrUnauthorizedUser)
			return
		}

		emailVerified := GetEmailVerifiedFromContext(ctx)
		if !emailVerified {
			HandleError(c, ErrEmailNotVerified)
			return
		}

		offerID, err := uuid.Parse(c.GetHeader(keys.HeaderOfferID))
		if err != nil {
			HandleError(c, fmt.Errorf("%s: %w", err.Error(), ErrInvalidOfferID))
			return
		}

		offerWithContract, err := client.GetOfferWithContract(ctx, offerID)
		if err != nil {
			if jhttp.IsNotFoundAPIError(err) {
				err = fmt.Errorf("%s: %w", err, ErrOfferNotFound)
			}
			HandleError(c, fmt.Errorf("%s with offer %s: GetOfferByID err: %w", errMsg, offerID.String(), err))
			return
		}

		if offerWithContract.Offer.UserID != userID {
			HandleError(c, ErrUserNotLinkedToOffer)
			return
		}

		ctx = ContextWithOfferID(ctx, offerID)
		ctx = ContextWithOfferType(ctx, offerWithContract.Offer.Type)
		ctx = ContextWithEntity(ctx, offerWithContract.GetEntity())
		ctx = ContextWithOfferContract(ctx, offerWithContract)
		c.Request = c.Request.WithContext(ctx)
	}
}

type offerFilter struct {
	offerTypes      []jentity.OfferType
	entities        []jentity.Entity
	excludeEntities []jentity.Entity
	activeOnly      bool
}

// Note: the resulting handler must be used after [EnforceHasOfferID].
func AllowOffer() *offerFilter {
	return &offerFilter{}
}

func (f *offerFilter) WithOfferTypes(types ...jentity.OfferType) *offerFilter {
	f.offerTypes = types
	return f
}

func (f *offerFilter) WithEntities(entities ...jentity.Entity) *offerFilter {
	f.entities = entities
	return f
}

func (f *offerFilter) WithExcludeEntities(entities ...jentity.Entity) *offerFilter {
	f.excludeEntities = entities
	return f
}

func (f *offerFilter) WithActiveOnly() *offerFilter {
	f.activeOnly = true
	return f
}

// Handler returns a gin.HandlerFunc that enforces the filter criteria.
func (f *offerFilter) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		if GetOfferIDFromContext(ctx) == uuid.Nil {
			// EnforceHasOfferID was not called
			HandleError(c, ErrUnauthorizedUser)
			return
		}

		if len(f.offerTypes) > 0 {
			offerType := GetOfferTypeFromContext(ctx)
			if !slices.Contains(f.offerTypes, offerType) {
				HandleError(c, ErrInvalidOfferType)
				return
			}
		}

		if len(f.entities) > 0 || len(f.excludeEntities) > 0 {
			contractEntity, err := GetEntityFromContext(ctx)
			if err != nil {
				HandleError(c, ErrMissingOfferContract)
				return
			}
			if len(f.entities) > 0 && !slices.Contains(f.entities, contractEntity) {
				HandleError(c, ErrInvalidOfferEntity)
				return
			}
			if len(f.excludeEntities) > 0 && slices.Contains(f.excludeEntities, contractEntity) {
				HandleError(c, ErrInvalidOfferEntity)
				return
			}
		}

		if f.activeOnly {
			offerContract := GetOfferContractFromContext(ctx)
			if (offerContract == pkgusers.OfferWithContract{}) {
				HandleError(c, ErrMissingOfferContract)
				return
			}

			if !offerContract.IsActive() {
				HandleError(c, ErrOfferNotActive)
				return
			}
		}
	}
}

// IsLifeOffer checks if the current offer is a Life offer.
// If allowedEntities is provided, the contract entity must be one of the allowed entities. EntityOnboarding can be used to allow offers without a contract.
// If allowedEntities is not provided, any Life offer with a valid contract entity will be allowed. Offers without a contract (ie. EntityOnboarding) will be rejected.
// Note: IsLifeOffer must be called after [EnforceHasOfferID].
func IsLifeOffer(allowedEntities ...jentity.Entity) gin.HandlerFunc {
	f := AllowOffer().WithOfferTypes(jentity.OfferTypeLife)
	if len(allowedEntities) == 0 {
		return f.WithExcludeEntities(jentity.EntityOnboarding).Handler()
	}
	return f.WithEntities(allowedEntities...).Handler()
}

// AllowedEntities checks if the current offer's entity is allowed.
// Note: AllowedEntities must be called after [EnforceHasOfferID].
func AllowedEntities(allowedEntities ...jentity.Entity) gin.HandlerFunc {
	return AllowOffer().WithEntities(allowedEntities...).Handler()
}

// IsOpenOffer checks if the current offer is an Open offer.
// Note: IsOpenOffer must be called after [EnforceHasOfferID].
func IsOpenOffer() gin.HandlerFunc {
	return AllowOffer().WithOfferTypes(jentity.OfferTypeOpen).Handler()
}

// IsOfferActive checks if the current offer is active.
// Note: IsOfferActive must be called after [EnforceHasOfferID].
func IsOfferActive() gin.HandlerFunc {
	return AllowOffer().WithActiveOnly().Handler()
}

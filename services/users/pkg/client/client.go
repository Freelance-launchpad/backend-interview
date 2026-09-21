package client

import (
	"context"
	"iter"
	"net/url"

	"github.com/Freelance-launchpad/backend-interview/common/jdb"
	"github.com/Freelance-launchpad/backend-interview/common/jentity"
	"github.com/Freelance-launchpad/backend-interview/services/users/pkg"
	"github.com/google/uuid"
)

type Client interface {
	ListUsersInfo(ctx context.Context, filters pkg.ListUsersInfoParams) ([]pkg.UserInfo, uint64, error)
	IterateUsersInfo(ctx context.Context, filters pkg.ListUsersInfoParams) iter.Seq2[pkg.UserInfo, error]
	GetLifeUserInfo(ctx context.Context, offerID uuid.UUID) (pkg.LifeUserInfo, error)
	// GetUserByID return a user using its id.
	GetUserByID(ctx context.Context, userID string) (pkg.User, error)
	// GetUserByOfferID returns the user associated with the offer.
	GetUserByOfferID(ctx context.Context, offerID uuid.UUID) (pkg.User, error)
	// GetUserWithOffersAndContracts returns a user with their offers and contracts.
	GetUserWithOffersAndContracts(ctx context.Context, userID string) (pkg.UserWithOffersAndContracts, error)
	// UpdateUser updates an user.
	UpdateUser(ctx context.Context, userID string, updateUser pkg.RequestUpdateUser) error
	UpdateUserEmail(ctx context.Context, userID, email string) error
	// UpdateUserProfile update the user profile, Some information cannot be updated wether the user already has a life offer.
	UpdateUserProfile(ctx context.Context, userID string, profile pkg.UserPersonalInfo) error

	// ListOffersByUserID returns offers of a user.
	// This method should only be used in situations where no offer ID is known beforehand.
	ListOffersByUserID(ctx context.Context, userID string, active *bool) ([]pkg.Offer, error)
	// GetOfferByID returns an offer from its id.
	GetOfferByID(ctx context.Context, id uuid.UUID) (pkg.Offer, error)
	// GetOfferInfo returns the offer additional info.
	GetOfferInfo(ctx context.Context, offerID uuid.UUID) (pkg.OfferInfo, error)
	// GetOfferWithContract returns the offer and the contract attached to this offer if it exists.
	GetOfferWithContract(ctx context.Context, offerID uuid.UUID) (pkg.OfferWithContract, error)
	// ListLifeOfferIDsByOwner returns all offers IDs which belongs to an admin.
	ListLifeOfferIDsByOwner(ctx context.Context, salesOwnerEmail, csmOwnerEmail string) ([]uuid.UUID, error)
	// GetOfferLegalInfo returns the legal info of the offer.
	GetOfferLegalInfo(ctx context.Context, offerID uuid.UUID) (pkg.LegalInfo, error)
	// UpdateLifeOfferLegalInfo updates the legal info of a Life offer.
	UpdateLifeOfferLegalInfo(ctx context.Context, offerID uuid.UUID, newEntity jentity.Entity) error
	// GetLifeOfferAddress returns the personal address of a life offer.
	GetLifeOfferAddress(ctx context.Context, offerID uuid.UUID) (pkg.Address, error)
	// EnableOffer enables an offer.
	EnableOffer(ctx context.Context, offerID uuid.UUID) error
	// DisableOffer disables an offer.
	DisableOffer(ctx context.Context, offerID uuid.UUID) error

	// UploadOfferFile uploads a file for the offer.
	UploadOfferFile(ctx context.Context, offerID uuid.UUID, content []byte, filename string, fields map[string]string) (uuid.UUID, error)
	ListOfferFiles(ctx context.Context, offerID uuid.UUID, filters pkg.ListFileInfosFilters) ([]pkg.FileInfo, int64, error)
	DeleteFile(ctx context.Context, fileID uuid.UUID) error

	// GetOfferJobContract returns the job contract of a given offer.
	GetOfferJobContract(ctx context.Context, id uuid.UUID) (pkg.Contract, error)
	// CreateJobContract creates a job contract for an offer.
	CreateJobContract(ctx context.Context, offerID uuid.UUID, contract pkg.Contract, address pkg.Address) (uuid.UUID, error)
	DeleteJobContract(ctx context.Context, offerID uuid.UUID) error

	// ListHRInfo lists HR informations
	ListHRInfo(ctx context.Context, entity *jentity.Entity, active bool, pagination jdb.QueryPagination) ([]pkg.HRInfo, uint64, error)
	IterateHRInfo(ctx context.Context, entity *jentity.Entity, active bool) iter.Seq2[pkg.HRInfo, error]
	// GetHRInfo returns HR info from an HR_ID.
	GetHRInfo(ctx context.Context, hrID string) (pkg.HRInfo, error)
	GetHRIDs(ctx context.Context, ssn string, entity *jentity.Entity) ([]string, error)
	UpdateHRID(ctx context.Context, offerID uuid.UUID, hrID string) error

	// UpdateIntercomID updates the intercom ID of a user.
	UpdateIntercomID(ctx context.Context, userID string, intercomID string) error
	// DeleteCRMData deletes the CRM data of users using the contact ID.
	DeleteCRMData(ctx context.Context, crmContactID string) error

	GetChurnByOfferID(ctx context.Context, offerID uuid.UUID) (pkg.Churn, error)
}

func New(httpClient any, url url.URL) Client {
	return Client(nil)
}

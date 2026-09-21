package keys

type authContextKeys uint8

const (
	UserIDKey authContextKeys = iota
	UserEmailKey
	IsServiceKey
	OfferIDKey
	EntityKey
	OfferTypeKey
	OfferContractKey
	EmailVerifiedKey
	PermissionsKey
	RefreshTokenIDKey
)

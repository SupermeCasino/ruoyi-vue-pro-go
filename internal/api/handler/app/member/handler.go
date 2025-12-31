package member

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewAppMemberAddressHandler,
	NewAppAuthHandler,
	NewAppMemberPointRecordHandler,
	NewAppMemberSignInConfigHandler,
	NewAppMemberSignInRecordHandler,
	NewAppSocialUserHandler,
	NewAppMemberUserHandler,
	NewHandlers,
)

type Handlers struct {
	Address      *AppMemberAddressHandler
	Auth         *AppAuthHandler
	PointRecord  *AppMemberPointRecordHandler
	SignInConfig *AppMemberSignInConfigHandler
	SignInRecord *AppMemberSignInRecordHandler
	SocialUser   *AppSocialUserHandler
	User         *AppMemberUserHandler
}

func NewHandlers(
	address *AppMemberAddressHandler,
	auth *AppAuthHandler,
	pointRecord *AppMemberPointRecordHandler,
	signInConfig *AppMemberSignInConfigHandler,
	signInRecord *AppMemberSignInRecordHandler,
	socialUser *AppSocialUserHandler,
	user *AppMemberUserHandler,
) *Handlers {
	return &Handlers{
		Address:      address,
		Auth:         auth,
		PointRecord:  pointRecord,
		SignInConfig: signInConfig,
		SignInRecord: signInRecord,
		SocialUser:   socialUser,
		User:         user,
	}
}

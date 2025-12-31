package member

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewMemberConfigHandler,
	NewMemberGroupHandler,
	NewMemberLevelHandler,
	NewMemberPointRecordHandler,
	NewMemberSignInConfigHandler,
	NewMemberSignInRecordHandler,
	NewMemberTagHandler,
	NewMemberUserHandler,
	NewHandlers,
)

type Handlers struct {
	Config       *MemberConfigHandler
	Group        *MemberGroupHandler
	Level        *MemberLevelHandler
	PointRecord  *MemberPointRecordHandler
	SignInConfig *MemberSignInConfigHandler
	SignInRecord *MemberSignInRecordHandler
	Tag          *MemberTagHandler
	User         *MemberUserHandler
}

func NewHandlers(
	config *MemberConfigHandler,
	group *MemberGroupHandler,
	level *MemberLevelHandler,
	pointRecord *MemberPointRecordHandler,
	signInConfig *MemberSignInConfigHandler,
	signInRecord *MemberSignInRecordHandler,
	tag *MemberTagHandler,
	user *MemberUserHandler,
) *Handlers {
	return &Handlers{
		Config:       config,
		Group:        group,
		Level:        level,
		PointRecord:  pointRecord,
		SignInConfig: signInConfig,
		SignInRecord: signInRecord,
		Tag:          tag,
		User:         user,
	}
}

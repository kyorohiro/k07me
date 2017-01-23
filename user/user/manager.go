package user

import (
	p "github.com/kyorohiro/k07me/pointer"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

type UserManagerConfig struct {
	RootGroup       string
	UserKind        string
	UserPointerKind string
	LengthHash      int
	LimitOfFinding  int
}

type UserManager struct {
	config         UserManagerConfig
	pointerManager *p.PointerManager
}

func NewUserManager(config UserManagerConfig) *UserManager {
	obj := new(UserManager)
	if config.RootGroup == "" {
		config.RootGroup = "FFUser"
	}
	if config.UserKind == "" {
		config.UserKind = "FFUser"
	}
	if config.UserPointerKind == "" {
		config.UserPointerKind = config.UserKind + "-pointer"
	}
	if config.LimitOfFinding <= 0 {
		config.LimitOfFinding = 20
	}
	obj.config = config

	obj.pointerManager = p.NewPointerManager(p.PointerManagerConfig{
		RootGroup: config.RootGroup,
		Kind:      config.UserPointerKind,
	})

	return obj
}

func (obj *UserManager) GetUserKind() string {
	return obj.config.UserKind
}

func (obj *UserManager) NewNewUser(ctx context.Context, sign string) *User {
	return obj.newUserWithUserName(ctx, sign)
}

func (obj *UserManager) GetUserFromUserName(ctx context.Context, userName string, sign string) (*User, error) {
	userObj := obj.newUser(ctx, userName, sign)
	e := userObj.loadFromDB(ctx)
	return userObj, e
}

func (obj *UserManager) SaveUser(ctx context.Context, userObj *User) error {
	return userObj.pushToDB(ctx)
}

func (obj *UserManager) DeleteUser(ctx context.Context, userName string, sign string) error {
	gaeKey := obj.newUserGaeObjectKey(ctx, userName, sign)
	return datastore.Delete(ctx, gaeKey)
}

func (obj *UserManager) FindAuthPointer(ctx context.Context, userName string) p.FoundPointers {
	q := obj.pointerManager.NewQueryFromOwner(userName)
	return obj.pointerManager.FindPointerFromQueryAll(ctx, q)

	//for k := range founded.Keys {
	//	kInfo := obj.pointerManager.GetKeyInfoFromStringId(k)
	//	kInfo.IdentifyType
	//}
	//(obj *PointerManager) FindPointerFromQueryAll(ctx context.Context, q *datastore.Query) FoundPointers
}

//
func Debug(ctx context.Context, message string) {
	log.Infof(ctx, message)
}

package user

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	m "github.com/kyorohiro/k07me/prop"
)

type UserManagerConfig struct {
	UserKind        string
	UserPointerKind string
	LengthHash      int
	LimitOfFinding  int
}

type UserManager struct {
	config UserManagerConfig
}

func NewUserManager(config UserManagerConfig) *UserManager {
	obj := new(UserManager)
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

	return obj
}

func (obj *UserManager) GetUserKind() string {
	return obj.config.UserKind
}

func (obj *UserManager) NewNewUser(ctx context.Context) *User {
	return obj.newUserWithUserName(ctx)
}

func (obj *UserManager) GetUserFromUserName(ctx context.Context, userName string) (*User, error) {
	userObj := obj.newUser(ctx, userName)
	e := userObj.loadFromDB(ctx)
	return userObj, e
}

func (obj *UserManager) SaveUser(ctx context.Context, userObj *User) error {
	return userObj.pushToDB(ctx)
}

func (obj *UserManager) DeleteUser(ctx context.Context, userName string, sign string) error {
	gaeKey := obj.newUserGaeObjectKey(ctx, userName)
	return datastore.Delete(ctx, gaeKey)
}

//
func Debug(ctx context.Context, message string) {
	log.Infof(ctx, message)
}

func MakePropValue(name, v string) string {
	p := m.NewMiniProp()
	p.SetString(name, v)
	v = string(p.ToJson())
	return v
}
package user

import (
	m "github.com/kyorohiro/k07me/prop"
)

type UserKeyInfo struct {
	Kind      string
	RootGroup string
	UserName  string
	Sign      string
}

func (obj *UserManager) MakeUserGaeObjectKeyStringId(userName string, sign string) string {
	propObj := m.NewMiniProp()
	propObj.SetString("k", obj.config.UserKind)
	propObj.SetString("p", obj.config.RootGroup)
	propObj.SetString("n", userName)
	propObj.SetString("s", sign)
	return string(propObj.ToJson())
}

func (obj *UserManager) GetUserKeyInfo(stringId string) *UserKeyInfo {
	propObj := m.NewMiniPropFromJson([]byte(stringId))
	return &UserKeyInfo{
		Kind:      propObj.GetString("k", ""),
		RootGroup: propObj.GetString("p", ""),
		UserName:  propObj.GetString("n", ""),
		Sign:      propObj.GetString("s", ""),
	}
}

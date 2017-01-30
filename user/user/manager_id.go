package user

import (
	m "github.com/kyorohiro/k07me/prop"
)

type UserKeyInfo struct {
	UserName string
}

func (obj *UserManager) MakeUserGaeObjectKeyStringId(userName string) string {
	propObj := m.NewMiniProp()
	propObj.SetString("n", userName)
	return string(propObj.ToJson())
}

func (obj *UserManager) GetUserKeyInfo(stringId string) *UserKeyInfo {
	propObj := m.NewMiniPropFromJson([]byte(stringId))
	return &UserKeyInfo{
		UserName: propObj.GetString("n", ""),
	}
}

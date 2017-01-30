package handler

import (
	"net/http"

	"io/ioutil"

	miniprop "github.com/kyorohiro/k07me/prop"
	"google.golang.org/appengine"
)

func (obj *UserHandler) HandleUpdateInfo(w http.ResponseWriter, r *http.Request) {
	outputProp := miniprop.NewMiniProp()
	v, _ := ioutil.ReadAll(r.Body)
	inputProp := miniprop.NewMiniPropFromJson(v)
	ctx := appengine.NewContext(r)
	userName := inputProp.GetString("userName", "")
	displayName := inputProp.GetString("displayName", "")
	content := inputProp.GetString("content", "")
	token := inputProp.GetString("token", "")

	reqErr := obj.OnUpdateUserRequest(w, r, obj, inputProp, outputProp)
	if reqErr != nil {
		obj.OnUpdateUserFailed(w, r, obj, inputProp, outputProp)
		obj.HandleError(w, r, outputProp, 2001, reqErr.Error())
		return
	}
	//
	// check token
	{
		loginResult := obj.CheckLoginFromToken(r, token, false)
		if loginResult.IsLogin == false {
			obj.OnUpdateUserFailed(w, r, obj, inputProp, outputProp)
			obj.HandleError(w, r, miniprop.NewMiniProp(), 2001, "need to login")
			return
		}

		if userName == "" {
			userName = loginResult.AccessTokenObj.GetUserName()
		}
		if userName != loginResult.AccessTokenObj.GetUserName() {
			usrObj, userErr := obj.GetManager().GetUserFromUserName(ctx, loginResult.AccessTokenObj.GetUserName())
			if userErr != nil {
				obj.OnUpdateUserFailed(w, r, obj, inputProp, outputProp)
				obj.HandleError(w, r, outputProp, 2002, userErr.Error())
				return
			}
			if true == usrObj.IsMaster() {
				obj.OnUpdateUserFailed(w, r, obj, inputProp, outputProp)
				obj.HandleError(w, r, outputProp, 2002, "need to admin status ")
			}
		}
	}

	usrObj, userErr := obj.GetManager().GetUserFromUserName(ctx, userName)
	if userErr != nil {
		obj.OnUpdateUserFailed(w, r, obj, inputProp, outputProp)
		obj.HandleError(w, r, outputProp, 2002, userErr.Error())
		return
	}
	usrObj.SetDisplayName(displayName)
	usrObj.SetCont(content)
	defChec := obj.OnUpdateUserBeforeSave(w, r, obj, usrObj, inputProp, outputProp)
	if defChec != nil {
		obj.OnUpdateUserFailed(w, r, obj, inputProp, outputProp)
		obj.HandleError(w, r, outputProp, 2003, defChec.Error())
		return
	}
	nextUserObj, nextUserErr := obj.GetManager().SaveUserWithImmutable(ctx, usrObj)
	if nextUserErr != nil {
		obj.OnUpdateUserFailed(w, r, obj, inputProp, outputProp)
		obj.HandleError(w, r, outputProp, 2004, userErr.Error())
		return
	}
	//
	sucErr := obj.OnUpdateUserSuccess(w, r, obj, usrObj, inputProp, outputProp)
	if sucErr != nil {
		obj.OnUpdateUserFailed(w, r, obj, inputProp, outputProp)
		obj.HandleError(w, r, outputProp, 2005, sucErr.Error())
		return
	}
	outputProp.CopiedOver(miniprop.NewMiniPropFromMap(nextUserObj.ToMapPublic()))
	w.WriteHeader(http.StatusOK)
	w.Write(outputProp.ToJson())

}

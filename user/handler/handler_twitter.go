package handler

import (
	"net/http"

	"io/ioutil"

	miniprop "github.com/kyorohiro/k07me/prop"
	minisession "github.com/kyorohiro/k07me/session"
	"google.golang.org/appengine"
)

func (obj *UserHandler) HandleTwitterRequestToken(w http.ResponseWriter, r *http.Request) {
	obj.twitterHandler.HandleLoginEntry(w, r)
}

func (obj *UserHandler) HandleTwitterCallbackToken(w http.ResponseWriter, r *http.Request) {
	obj.twitterHandler.HandleLoginExit(w, r)
}

func (obj *UserHandler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	bodyBytes, _ := ioutil.ReadAll(r.Body)
	propObj := miniprop.NewMiniPropFromJson(bodyBytes)
	token := propObj.GetString("token", "")

	obj.sessionMgr.Logout(appengine.NewContext(r), token, minisession.MakeOptionInfo(r))
}

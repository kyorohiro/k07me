package handler

import (
	"net/http"

	"errors"

	"strconv"

	"github.com/kyorohiro/k07me/oauth/twitter"
	minipointer "github.com/kyorohiro/k07me/pointer"
	miniprop "github.com/kyorohiro/k07me/prop"
	minisession "github.com/kyorohiro/k07me/session"
	miniuser "github.com/kyorohiro/k07me/user/user"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
)

//
//
//
func (obj *UserHandler) NewTwitterHandlerObj(config twitter.TwitterOAuthConfig) *twitter.TwitterHandler {
	twitterHandlerObj := twitter.NewTwitterHandler( //
		config, twitter.TwitterHundlerOnEvent{
			OnFoundUser: func(w http.ResponseWriter, r *http.Request, handler *twitter.TwitterHandler, accesssToken *twitter.SendAccessTokenResult) map[string]string {
				ctx := appengine.NewContext(r)

				//
				//
				_, _, userObj, err1 := obj.LoginRegistFromTwitter(ctx, //
					accesssToken.GetScreenName(), //
					accesssToken.GetUserID(),     //
					accesssToken.GetOAuthToken()) //
				if err1 != nil {
					return map[string]string{"errcode": "2", "errindo": err1.Error()}
				}
				//
				//
				tokenObj, err := obj.sessionMgr.Login(ctx, //
					userObj.GetUserName(), //
					minisession.MakeAccessTokenConfigFromRequest(r))
				if err != nil {
					return map[string]string{"errcode": "1"}
				} else {
					return map[string]string{ //
						"token":    "" + tokenObj.GetLoginId(), //
						"userName": userObj.GetUserName(),
						"isMaster": strconv.Itoa(userObj.GetPermission())}
				}
			},
		})

	return twitterHandlerObj
}

//
//
//
func (obj *UserHandler) LoginRegistFromTwitter(ctx context.Context, screenName string, userId string, oauthToken string) (bool, *minipointer.Pointer, *miniuser.User, error) {
	return obj.LoginRegistFromSNS(ctx, screenName, userId, oauthToken, minipointer.TypeTwitter)
}

func (obj *UserHandler) LoginRegistFromFacebook(ctx context.Context, screenName string, userId string, oauthToken string) (bool, *minipointer.Pointer, *miniuser.User, error) {
	return obj.LoginRegistFromSNS(ctx, screenName, userId, oauthToken, minipointer.TypeFacebook)
}

func (obj *UserHandler) LoginRegistFromSNS(ctx context.Context, screenName string, userId string, oauthToken string, snsType string) (bool, *minipointer.Pointer, *miniuser.User, error) {

	snsIdProp := miniprop.NewMiniProp()
	snsIdProp.SetString("n", screenName)
	snsIdProp.SetString("i", userId)
	relayIdObj := obj.relayIdMgr.GetPointerWithNew(ctx, string(snsIdProp.ToJson()), snsType, map[string]string{"token": oauthToken})
	needMake := false

	//
	// new userObj
	var err error = nil
	var userObj *miniuser.User = nil
	var pointerObj *minipointer.Pointer = nil
	if relayIdObj.GetValue() != "" {
		needMake = true
		//		Debug(ctx, "LoginRegistFromTwitter (1) :"+relayIdObj.GetUserName())
		pointerObj = obj.relayIdMgr.GetPointerWithNewForRelayId(ctx, relayIdObj.GetValue())
		if pointerObj.GetValue() != "" {
			userObj, err = obj.GetManager().GetUserFromUserName(ctx, pointerObj.GetValue(), pointerObj.GetSign())
			if err != nil {
				userObj = nil
			}
		}
	}
	if userObj == nil {
		userObj = obj.GetManager().NewNewUser(ctx, "")
		userObj.SetDisplayName(screenName)
		//		Debug(ctx, "LoginRegistFromTwitter (2) :"+userObj.GetUserName())
		pointerObj = obj.relayIdMgr.GetPointerWithNewForRelayId(ctx, userObj.GetUserName())
		pointerObj.SetValue(userObj.GetUserName())
		pointerObj.SetOwner(userObj.GetUserName())
		pointerObj.SetSign("")
		//Debug(ctx, "LoginRegistFromTwitter :")
		err := obj.relayIdMgr.Save(ctx, pointerObj)
		if err != nil {
			return needMake, nil, nil, errors.New("failed to save pointreobj : " + err.Error())
		}
		//
		// set username
		relayIdObj.SetValue(pointerObj.GetValue())
		relayIdObj.SetOwner(userObj.GetUserName())
	}

	//
	// save relayId
	err = obj.relayIdMgr.Save(ctx, relayIdObj)
	if err != nil {
		return needMake, nil, nil, errors.New("failed to save sessionobj : " + err.Error())
	}

	//
	// save user
	err = obj.GetManager().SaveUser(ctx, userObj)
	if err != nil {
		return needMake, nil, nil, errors.New("failed to save userobj : " + err.Error())
	}
	return needMake, relayIdObj, userObj, nil
}

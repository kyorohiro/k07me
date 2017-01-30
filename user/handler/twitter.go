package handler

import (
	"net/http"

	//	"errors"

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
func (obj *UserHandler) NewTwitterHandlerObj(config twitter.TwitterOAuthConfig) *twitter.TwitterHandler {
	twitterHandlerObj := twitter.NewTwitterHandler( //
		config, twitter.TwitterHundlerOnEvent{
			OnFoundUser: func(w http.ResponseWriter, r *http.Request, handler *twitter.TwitterHandler, //
				accesssToken *twitter.SendAccessTokenResult) map[string]string {
				ctx := appengine.NewContext(r)

				//
				//
				userObj, err1 := obj.LoginRegistFromTwitter(ctx, //
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
					minisession.MakeOptionInfo(r))
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
/*
	_, _, userObj, err1 := obj.LoginRegistFromTwitter(ctx, //
		accesssToken.GetScreenName(), //
		accesssToken.GetUserID(),     //
		accesssToken.GetOAuthToken()) //
*/
func (obj *UserHandler) LoginRegistFromTwitter(ctx context.Context, screenName string, userId string, oauthToken string) (*miniuser.User, error) {
	return obj.LoginRegistFromSNS(ctx, screenName, userId, oauthToken, minipointer.TypeTwitter)
}

func (obj *UserHandler) LoginRegistFromSNS(ctx context.Context, screenName string, userId string, //
	oauthToken string, snsType string) (*miniuser.User, error) {

	fs := obj.GetManager().FindUserFromProp(ctx, snsType, screenName, "", false)
	var user *miniuser.User
	if len(fs.Users) <= 0 {
		user = obj.GetManager().NewNewUser(ctx)
	} else {
		user = fs.Users[0]
	}
	user.SetDisplayName(screenName)
	user.SetProp(snsType, screenName)
	privateProp := miniprop.NewMiniPropFromJson([]byte(user.GetPrivateInfo()))
	privateProp.SetString("n", screenName)
	privateProp.SetString("i", userId)
	privateProp.SetString("t", oauthToken)
	user.SetPrivateInfo(string(privateProp.ToJson()))
	obj.GetManager().SaveUser(ctx, user)

	return user, nil
}

package handler

import (
	"net/http"

	miniblob "github.com/kyorohiro/k07me/blob/blob"
	blobhandler "github.com/kyorohiro/k07me/blob/handler"
	"github.com/kyorohiro/k07me/oauth/twitter"
	//	minipointer "github.com/kyorohiro/k07me/pointer"
	miniprop "github.com/kyorohiro/k07me/prop"
	minisession "github.com/kyorohiro/k07me/session"
	miniuser "github.com/kyorohiro/k07me/user/user"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	//
	"crypto/sha1"
)

type UserHandler struct {
	manager *miniuser.UserManager
	//	relayIdMgr     *minipointer.PointerManager
	sessionMgr     *minisession.SessionManager
	blobHandler    *blobhandler.BlobHandler
	twitterHandler *twitter.TwitterHandler
	onEvents       UserHandlerOnEvent
	completeFunc   func(w http.ResponseWriter, r *http.Request, outputProp *miniprop.MiniProp, hh *blobhandler.BlobHandler, blobObj *miniblob.BlobItem) error
}

type UserHandlerManagerConfig struct {
	UserKind string
	//	RelayIdKind                string
	SessionKind                string
	BlobKind                   string
	BlobPointerKind            string
	BlobSign                   string
	MemcachedOnlyInBlobPointer bool
	LengthHash                 int
}

type UserHandlerOnEvent struct {
	OnGetUserRequestList       []func(w http.ResponseWriter, r *http.Request, h *UserHandler, o *miniprop.MiniProp) error
	OnGetUserFailedList        []func(w http.ResponseWriter, r *http.Request, h *UserHandler, o *miniprop.MiniProp)
	OnGetUserSuccessList       []func(w http.ResponseWriter, r *http.Request, h *UserHandler, i *miniuser.User, o *miniprop.MiniProp) error
	OnUpdateUserRequestList    []func(w http.ResponseWriter, r *http.Request, h *UserHandler, i *miniprop.MiniProp, o *miniprop.MiniProp) error
	OnUpdateUserFailedList     []func(w http.ResponseWriter, r *http.Request, h *UserHandler, i *miniprop.MiniProp, o *miniprop.MiniProp)
	OnUpdateUserBeforeSaveList []func(w http.ResponseWriter, r *http.Request, h *UserHandler, u *miniuser.User, i *miniprop.MiniProp, o *miniprop.MiniProp) error
	OnUpdateUserSuccessList    []func(w http.ResponseWriter, r *http.Request, h *UserHandler, u *miniuser.User, i *miniprop.MiniProp, o *miniprop.MiniProp) error
}

func NewUserHandler(callbackUrl string, //
	config UserHandlerManagerConfig) *UserHandler {
	if config.UserKind == "" {
		config.UserKind = "fu"
	}
	//	if config.RelayIdKind == "" {
	//		config.RelayIdKind = config.UserKind + "-pointer"
	//	}
	if config.SessionKind == "" {
		config.SessionKind = config.UserKind + "-session"
	}
	if config.BlobKind == "" {
		config.BlobKind = config.UserKind + "-blob"
	}
	if config.BlobPointerKind == "" {
		config.BlobPointerKind = config.UserKind + "-blob-pointer"
	}
	if config.BlobSign == "" {
		config.BlobSign = string(sha1.New().Sum([]byte("" + config.UserKind)))
	}
	//

	ret := &UserHandler{
		manager: miniuser.NewUserManager(miniuser.UserManagerConfig{
			UserKind: config.UserKind,
			//			UserPointerKind: config.RelayIdKind,
			LengthHash: config.LengthHash,
		}),
		//		relayIdMgr: minipointer.NewPointerManager( //
		//			minipointer.PointerManagerConfig{
		//				Kind:      config.RelayIdKind,
		//				RootGroup: config.RootGroup,
		//			}),
		sessionMgr: minisession.NewSessionManager(minisession.SessionManagerConfig{
			Kind: config.SessionKind,
		}),
		blobHandler: blobhandler.NewBlobHandler(callbackUrl, config.BlobSign, miniblob.BlobManagerConfig{
			Kind:                   config.BlobKind,
			PointerKind:            config.BlobPointerKind,
			CallbackUrl:            callbackUrl,
			MemcachedOnlyInPointer: config.MemcachedOnlyInBlobPointer,
			HashLength:             10,
		}),
		onEvents: UserHandlerOnEvent{},
	}

	ret.blobHandler.AddOnBlobComplete(ret.OnBlobComplete)
	return ret
}

//func (obj *UserHandler) GetPointerManager() *minipointer.PointerManager {
//	return obj.relayIdMgr
//}

func (obj *UserHandler) GetBlobHandler() *blobhandler.BlobHandler {
	return obj.blobHandler
}

func (obj *UserHandler) AddTwitterSession(twitterConfig twitter.TwitterOAuthConfig) {
	obj.twitterHandler = obj.NewTwitterHandlerObj(twitterConfig)
}

func (obj *UserHandler) GetSessionMgr() *minisession.SessionManager {
	return obj.sessionMgr
}

func (obj *UserHandler) GetManager() *miniuser.UserManager {
	return obj.manager
}

func Debug(ctx context.Context, message string) {
	log.Infof(ctx, message)
}

func (obj *UserHandler) CheckLoginFromToken(r *http.Request, token string, useIp bool) minisession.CheckResult {
	ctx := appengine.NewContext(r)
	return obj.GetSessionMgr().CheckAccessToken(ctx, token, minisession.MakeOptionInfo(r), useIp)
}

func (obj *UserHandler) HandleError(w http.ResponseWriter, r *http.Request, outputProp *miniprop.MiniProp, errorCode int, errorMessage string) {
	//
	//
	if outputProp == nil {
		outputProp = miniprop.NewMiniProp()
	}
	if errorCode != 0 {
		outputProp.SetInt("errorCode", errorCode)
	}
	if errorMessage != "" {
		outputProp.SetString("errorMessage", errorMessage)
	}
	w.WriteHeader(http.StatusBadRequest)
	w.Write(outputProp.ToJson())
}

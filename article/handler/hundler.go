package handler

import (
	"net/http"

	"io/ioutil"

	"github.com/kyorohiro/k07me/article/article"
	miniblob "github.com/kyorohiro/k07me/blob/blob"
	blobhandler "github.com/kyorohiro/k07me/blob/handler"
	miniprop "github.com/kyorohiro/k07me/prop"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

const (
	ErrorCodeFailedToSave                = 2001
	ErrorCodeFailedToCheckAboutGetCalled = 2002
	ErrorCodeNotFoundArticleId           = 2003
)

type ArticleHandler struct {
	projectId   string
	articleKind string
	blobKind    string
	pointerKind string
	artMana     *article.ArticleManager
	blobHundler *blobhandler.BlobHandler
	//	tagMana     *tag.TagManager
	//	onEvents ArticleHandlerOnEvent
}

type ArticleHandlerConfig struct {
	RootGroup       string
	ArticleKind     string
	PointerKind     string
	BlobKind        string
	BlobPointerKind string
	TagKind         string
	BlobCallbackUrl string
	BlobSign        string
	MemcachedOnly   bool
	LengthHash      int
}

func NewArtHandler(config ArticleHandlerConfig) *ArticleHandler {
	if config.RootGroup == "" {
		config.RootGroup = "ffstyle"
	}
	if config.ArticleKind == "" {
		config.ArticleKind = "ffart"
	}
	if config.PointerKind == "" {
		config.PointerKind = config.ArticleKind + "-pointer"
	}
	if config.BlobKind == "" {
		config.BlobKind = config.ArticleKind + "-blob"
	}
	if config.BlobPointerKind == "" {
		config.BlobPointerKind = config.ArticleKind + "-blob-pointer"
	}
	if config.TagKind == "" {
		config.TagKind = config.ArticleKind + "-tag"
	}
	artMana := article.NewArticleManager(article.ArticleManagerConfig{
		RootGroup:      config.RootGroup,
		KindArticle:    config.ArticleKind,
		KindPointer:    config.PointerKind,
		PrefixOfId:     "art",
		LimitOfFinding: 20,
		LengthHash:     config.LengthHash,
	})
	//	tagMana := tag.NewTagManager(config.TagKind, config.RootGroup)
	//
	//
	artHandlerObj := &ArticleHandler{
		projectId:   config.RootGroup,
		articleKind: config.ArticleKind,
		blobKind:    config.BlobKind,
		artMana:     artMana,
		//		tagMana:     tagMana,
		//		onEvents: ArticleHandlerOnEvent{},
	}

	//
	artHandlerObj.blobHundler = blobhandler.NewBlobHandler(config.BlobCallbackUrl, config.BlobSign,
		miniblob.BlobManagerConfig{
			RootGroup:              config.RootGroup,
			Kind:                   config.BlobKind,
			CallbackUrl:            config.BlobCallbackUrl,
			PointerKind:            config.BlobPointerKind,
			MemcachedOnlyInPointer: config.MemcachedOnly,
			HashLength:             10,
		})
	//	artHandlerObj.blobHundler.AddOnBlobBeforeSave(func(w http.ResponseWriter, r *http.Request, p *miniprop.MiniProp, h *blobhandler.BlobHandler, i *miniblob.BlobItem) error {
	//		dirSrc := r.URL.Query().Get("dir")
	//		articlId := artHandlerObj.GetArticleIdFromDir(dirSrc)
	//		i.SetOwner(articlId)
	//		return nil
	//	})
	artHandlerObj.blobHundler.AddOnBlobComplete(func(w http.ResponseWriter, r *http.Request, o *miniprop.MiniProp, hh *blobhandler.BlobHandler, i *miniblob.BlobItem) error {
		dirSrc := r.URL.Query().Get("dir")
		articlId := artHandlerObj.GetArticleIdFromDir(dirSrc)
		dir := artHandlerObj.GetDirFromDir(dirSrc)
		fileName := r.URL.Query().Get("file")
		//
		//
		ctx := appengine.NewContext(r)
		//Debug(ctx, "OnBlobComplete ::"+articlId+"::"+dir+"::"+fileName+"::")
		artObj, errGet := artHandlerObj.GetManager().GetArticleFromPointer(ctx, articlId)
		if errGet != nil {
			//Debug(ctx, "From Pointer GEt ER "+articlId)
			return errGet
		}
		//s	Debug(ctx, "=~====>> ICOM "+dir+"::"+fileName)
		if dir == "" && fileName == "icon" {
			artObj.SetIconUrl("key://" + i.GetBlobKey())
			// todo
			_, errSave := artHandlerObj.GetManager().SaveUsrWithImmutable(ctx, artObj)
			//			artHandlerObj.tagMana.DeleteTagsFromOwner(appengine.NewContext(r), nextArtObj.GetArticleId())
			//			artHandlerObj.tagMana.AddBasicTags(ctx, nextArtObj.GetTags(), "art://"+nextArtObj.GetGaeObjectKey().StringID(), artObj.GetArticleId(), "")

			if errSave != nil {
				return errSave
			}
		}
		return nil
	})
	return artHandlerObj
}

func (obj *ArticleHandler) GetManager() *article.ArticleManager {
	return obj.artMana
}

func (obj *ArticleHandler) GetBlobHandler() *blobhandler.BlobHandler {
	return obj.blobHundler
}

func (obj *ArticleHandler) HandleError(w http.ResponseWriter, r *http.Request, outputProp *miniprop.MiniProp, errorCode int, errorMessage string) {
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

func (obj *ArticleHandler) GetInputProp(w http.ResponseWriter, r *http.Request) *miniprop.MiniProp {
	v, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return miniprop.NewMiniProp()
	} else {
		return miniprop.NewMiniPropFromJson(v)
	}
}

//
//
//

// HandleBlobRequestTokenFromParams

func Debug(ctx context.Context, message string) {
	log.Infof(ctx, message)
}

///
//

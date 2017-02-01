package handler

import (
	"net/http"

	"github.com/kyorohiro/k07me/article/article"
	miniprop "github.com/kyorohiro/k07me/prop"
	"google.golang.org/appengine"
)

func (obj *ArticleHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	propObj := miniprop.NewMiniProp()
	ctx := appengine.NewContext(r)
	values := r.URL.Query()
	key := values.Get("key")
	articleId := values.Get("articleId")
	sign := values.Get("sign")
	//	mode := values.Get("m")
	//
	if key != "" {
		keyInfo := obj.GetManager().ExtractInfoFromStringId(key)
		articleId = keyInfo.ArticleId
		sign = keyInfo.Sign
	}
	var artObj *article.Article
	var err error
	//
	//
	errOnGAR := obj.OnGetArtRequest(w, r, obj, propObj)
	if errOnGAR != nil {
		obj.OnGetArtFailed(w, r, obj, propObj)
		obj.HandleError(w, r, propObj, ErrorCodeNotFoundArticleId, errOnGAR.Error())
		return
	}
	if sign != "" {
		artObj, err = obj.GetManager().GetArticleFromArticleId(ctx, articleId, sign)
	} else {
		artObj, _, err = obj.GetManager().GetArticleFromPointer(ctx, articleId)
	}
	if err != nil {
		obj.OnGetArtFailed(w, r, obj, propObj)
		obj.HandleError(w, r, propObj, ErrorCodeNotFoundArticleId, "not found article")
		return
	}
	if sign != "" {
		w.Header().Set("Cache-Control", "public, max-age=2592000")
	}
	errOnGAS := obj.OnGetArtSuccess(w, r, obj, artObj, propObj)
	if errOnGAS != nil {
		obj.OnGetArtFailed(w, r, obj, propObj)
		obj.HandleError(w, r, propObj, ErrorCodeNotFoundArticleId, errOnGAS.Error())
		return
	}
	propObj.CopiedOver(miniprop.NewMiniPropFromMap(artObj.ToMap()))
	w.Write(propObj.ToJson())
}

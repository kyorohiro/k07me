package handler

import (
	"net/http"

	ar "github.com/kyorohiro/k07me/article/article"
	miniprop "github.com/kyorohiro/k07me/prop"
	"google.golang.org/appengine"
)

func (obj *ArticleHandler) HandleNew(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	propObj := miniprop.NewMiniProp()
	//
	// load param from json
	inputProp := obj.GetInputProp(w, r)
	title := inputProp.GetString("title", "")
	//z	target := inputProp.GetString("target", "")
	content := inputProp.GetString("content", "")
	ownerName := inputProp.GetString("ownerName", "")
	tags := inputProp.GetPropStringList("", "tags", nil)
	articleId := inputProp.GetPropString("", "articleId", "")
	//
	//
	propKeys := inputProp.GetPropStringList("", "propKeys", make([]string, 0))
	propValues := inputProp.GetPropStringList("", "propValues", make([]string, 0))
	lat := inputProp.GetFloat("lat", -999.0)
	lng := inputProp.GetFloat("lng", -999.0)
	//
	//
	outputProp := miniprop.NewMiniProp()
	errCall := obj.OnNewRequest(w, r, obj, inputProp, outputProp)
	if nil != errCall {
		obj.OnNewArtFailed(w, r, obj, inputProp, outputProp)
		obj.HandleError(w, r, outputProp, ErrorCodeFailedToCheckAboutGetCalled, errCall.Error())
		return
	}
	//
	var artObj *ar.Article
	if articleId != "" {
		var artErr error
		artObj, artErr = obj.GetManager().NewArticleFromArticleId(ctx, articleId)
		if artErr != nil {
			obj.OnNewArtFailed(w, r, obj, inputProp, outputProp)
			obj.HandleError(w, r, outputProp, ErrorCodeFailedToCheckAboutGetCalled, artErr.Error())
			return
		}
	} else {
		artObj = obj.GetManager().NewArticle(ctx)
	}
	artObj.SetTitle(title)
	artObj.SetCont(content)
	artObj.SetUserName(ownerName)
	artObj.SetTags(tags)
	artObj.SetLat(lat)
	artObj.SetLng(lng)
	//
	//
	if len(propKeys) == len(propValues) {
		for i, kv := range propKeys {
			artObj.SetProp(kv, propValues[i])
		}
	}
	//
	errNew := obj.OnNewBeforeSave(w, r, obj, artObj, inputProp, outputProp)
	if nil != errNew {
		obj.OnNewArtFailed(w, r, obj, inputProp, outputProp)
		obj.HandleError(w, r, outputProp, ErrorCodeFailedToCheckAboutGetCalled, errNew.Error())
		return
	}
	//
	nextArtObj, errSave := obj.GetManager().SaveUsrWithImmutable(ctx, artObj)
	if errSave != nil {
		obj.OnNewArtFailed(w, r, obj, inputProp, outputProp)
		obj.HandleError(w, r, outputProp, ErrorCodeFailedToSave, errSave.Error())
		return
	}
	propObj.SetPropString("", "articleId", nextArtObj.GetArticleId())
	errOnSc := obj.OnNewArtSuccess(w, r, obj, nextArtObj, inputProp, outputProp)
	if nil != errOnSc {
		if nil != obj.GetManager().DeleteFromArticleId(ctx, artObj.GetArticleId(), "") {
			Debug(ctx, "<GOMIDATA>articleId="+artObj.GetArticleId())
		}
		obj.OnNewArtFailed(w, r, obj, inputProp, outputProp)
		obj.HandleError(w, r, outputProp, ErrorCodeFailedToCheckAboutGetCalled, errOnSc.Error())
		return
	}
	///
	// add tag
	///
	//	obj.tagMana.AddBasicTags(ctx, tags, "art://"+nextArtObj.GetGaeObjectKey().StringID(), nextArtObj.GetArticleId(), "")
	w.WriteHeader(http.StatusOK)
	w.Write(propObj.ToJson())

}

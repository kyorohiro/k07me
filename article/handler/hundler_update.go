package handler

import (
	"net/http"

	miniprop "github.com/kyorohiro/k07me/prop"
	"google.golang.org/appengine"
)

func (obj *ArticleHandler) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	propObj := miniprop.NewMiniProp()
	//
	// load param from json
	inputProp := obj.GetInputProp(w, r)
	articleId := inputProp.GetString("articleId", "")
	ownerName := inputProp.GetString("userName", "")
	title := inputProp.GetString("title", "")
	//target := inputProp.GetString("target", "")
	content := inputProp.GetString("content", "")
	tags := inputProp.GetPropStringList("", "tags", make([]string, 0))
	//
	//
	propKeys := inputProp.GetPropStringList("", "propKeys", make([]string, 0))
	propValues := inputProp.GetPropStringList("", "propValues", make([]string, 0))
	lat := inputProp.GetFloat("lat", -999.0)
	lng := inputProp.GetFloat("lng", -999.0)
	//
	//
	outputProp := miniprop.NewMiniProp()

	//
	if articleId == "" {
		obj.HandleError(w, r, outputProp, ErrorCodeNotFoundArticleId, "Not Found Article")
		return
	}

	artObj, errGetArt := obj.GetManager().GetArticleFromPointer(ctx, articleId)
	if errGetArt != nil {
		obj.HandleError(w, r, outputProp, ErrorCodeNotFoundArticleId, "Not Found Article")
		return
	}
	//

	artObj.SetTitle(title)
	artObj.SetUserName(ownerName)
	//	artObj.SetProp("target", target)
	artObj.SetCont(content)
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
	//
	_, errSave := obj.GetManager().SaveArticleWithImmutable(ctx, artObj)

	if errSave != nil {
		obj.HandleError(w, r, outputProp, ErrorCodeFailedToSave, errSave.Error())
		return
	} else {
		propObj.SetPropString("", "articleId", artObj.GetArticleId())
		w.WriteHeader(http.StatusOK)
		w.Write(propObj.ToJson())
	}
}

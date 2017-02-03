package article

import (
	"time"

	"errors"

	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/memcache"
)

func (obj *ArticleManager) GetArticleFromArticleId(ctx context.Context, articleId string, sign string) (*Article, error) {
	return obj.NewArticleFromGaeObjectKey(ctx, obj.NewGaeObjectKey(ctx, articleId, sign, ""))
}

func (obj *ArticleManager) NewArticleFromGaeObjectKey(ctx context.Context, key *datastore.Key) (*Article, error) {
	k := key
	//
	//
	artObjFMem, errNewFMem := obj.NewArticleFromMemcache(ctx, k.StringID())
	if errNewFMem == nil {
		log.Infof(ctx, ">>>> new article Obj from memcache")
		return artObjFMem, nil
	}
	//
	//
	var gaeObj GaeObjectArticle
	err := datastore.Get(ctx, k, &gaeObj)
	if err != nil {
		return nil, err
	}
	//
	//
	return obj.NewArticleFromGaeObject(ctx, k, &gaeObj), nil
}

func (obj *ArticleManager) NewArticleFromMemcache(ctx context.Context, stringId string) (*Article, error) {
	ret := new(Article)
	ret.gaeObject = new(GaeObjectArticle)
	idInfo := obj.ExtractInfoFromStringId(stringId)
	ret.gaeObjectKey = obj.NewGaeObjectKey(ctx, idInfo.ArticleId, idInfo.Sign, "")
	ret.kind = obj.config.KindArticle
	artObjSource, errGetFMem := memcache.Get(ctx, ret.gaeObjectKey.StringID())
	if errGetFMem != nil {
		return nil, errGetFMem
	}
	errSetParam := ret.SetParamFromsJson(ctx, string(artObjSource.Value))

	return ret, errSetParam
}

func (obj *ArticleManager) NewArticleFromGaeObject(ctx context.Context, gaeKey *datastore.Key, gaeObj *GaeObjectArticle) *Article {
	ret := new(Article)
	ret.gaeObject = gaeObj
	ret.gaeObjectKey = gaeKey
	ret.kind = obj.config.KindArticle
	//
	//

	return ret
}

func (obj *ArticleManager) NewArticleFromArticle(ctx context.Context, baseArtObj *Article, sign string) *Article {
	//
	ret := new(Article)
	ret.kind = obj.config.KindArticle
	ret.gaeObject = &GaeObjectArticle{}
	ret.gaeObjectKey = obj.NewGaeObjectKey(ctx, baseArtObj.GetArticleId(), sign, "")

	//
	baseArtData := baseArtObj.ToMap()
	baseArtData[TypeSign] = sign
	ret.SetParamFromsMap(baseArtData)
	return ret
}

func (obj *ArticleManager) NewArticle(ctx context.Context) *Article {
	created := time.Now()
	var secretKey string
	var articleId string
	var key *datastore.Key
	var art GaeObjectArticle
	sign := "0"
	for {
		secretKey = obj.makeRandomId() + obj.makeRandomId()
		articleId = obj.makeArticleId(created, secretKey)
		//stringId := obj.makeStringId(articleId, sign)
		//
		Debug(ctx, "<NewArticle>"+articleId)
		key = obj.NewGaeObjectKey(ctx, articleId, sign, "")
		err := datastore.Get(ctx, key, &art)
		if err != nil {
			break
		}
	}
	//
	ret := new(Article)
	ret.kind = obj.config.KindArticle
	ret.gaeObject = &art
	ret.gaeObjectKey = key
	ret.gaeObject.Sign = sign
	ret.gaeObject.Created = created
	ret.gaeObject.Updated = created
	ret.gaeObject.SecretKey = secretKey
	ret.gaeObject.ArticleId = articleId
	//
	return ret
}

func (obj *ArticleManager) NewArticleFromArticleId(ctx context.Context, articleId string) (*Article, error) {
	created := time.Now()
	secretKey := ""
	var key *datastore.Key
	var art GaeObjectArticle
	sign := "0"

	//articleId = obj.makeArticleId(created, secretKey)
	key = obj.NewGaeObjectKey(ctx, articleId, sign, "")
	err := datastore.Get(ctx, key, &art)
	if err == nil {
		return nil, errors.New("already found")
	}
	//
	ret := new(Article)
	ret.kind = obj.config.KindArticle
	ret.gaeObject = &art
	ret.gaeObjectKey = key
	ret.gaeObject.Sign = sign
	ret.gaeObject.Created = created
	ret.gaeObject.Updated = created
	ret.gaeObject.SecretKey = secretKey
	ret.gaeObject.ArticleId = articleId
	//
	return ret, nil
}

func (obj *ArticleManager) NewGaeObjectKey(ctx context.Context, articleId string, sign string, kind string) *datastore.Key {
	if kind == "" {
		kind = obj.config.KindArticle
	}
	return datastore.NewKey(ctx, kind, obj.makeStringId(articleId, sign), 0, nil)
}

func (obj *ArticleManager) NewGaeObjectKeyFromKeyString(ctx context.Context, key string) *datastore.Key {
	return datastore.NewKey(ctx, obj.config.KindArticle, key, 0, nil)
}

//
//
//
func (obj *Article) saveOnDB(ctx context.Context) error {
	_, err := datastore.Put(ctx, obj.gaeObjectKey, obj.gaeObject)
	obj.updateMemcache(ctx)
	return err
}

func (mgrObj *ArticleManager) SaveOnOtherDB(ctx context.Context, obj *Article, kind string) error {
	_, err := datastore.Put(ctx, mgrObj.NewGaeObjectKey(ctx, obj.GetArticleId(), obj.gaeObject.Sign, kind), obj.gaeObject)
	return err
}

func (mgrObj *ArticleManager) DeleteFromArticleId(ctx context.Context, articleId, sign string) error {
	key := mgrObj.NewGaeObjectKey(ctx, articleId, sign, mgrObj.GetKind())
	memcache.Delete(ctx, key.StringID())
	return datastore.Delete(ctx, mgrObj.NewGaeObjectKey(ctx, articleId, sign, mgrObj.GetKind()))
}

func (obj *ArticleManager) DeleteFromArticleIdWithPointer(ctx context.Context, articleId string) error {
	//
	key, e := obj.GetArticleKeyFromPointer(ctx, articleId)
	if e == nil {
		memcache.Delete(ctx, key)
		return datastore.Delete(ctx, obj.NewGaeObjectKeyFromKeyString(ctx, key))
	}
	return nil
}

func (obj *ArticleManager) GetArticleKeyFromPointer(ctx context.Context, articleId string) (string, error) {
	founded := obj.FindArticleFromArticleId(ctx, articleId, "", true)
	if len(founded.ArticleKeys) <= 0 {
		return "", errors.New("not found")
	}
	return founded.ArticleKeys[0], nil
}

func (obj *ArticleManager) GetArticleFromPointer(ctx context.Context, articleId string) (*Article, error) {
	Debug(ctx, "articleId = "+articleId)
	founded := obj.FindArticleFromArticleId(ctx, articleId, "", false)
	if len(founded.Articles) <= 0 {
		Debug(ctx, "articleId e1")
		return nil, errors.New("not found")
	}
	Debug(ctx, "articleId a1")
	return founded.Articles[0], nil
}

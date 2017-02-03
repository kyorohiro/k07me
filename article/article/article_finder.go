package article

import (
	miniprop "github.com/kyorohiro/k07me/prop"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

// https://cloud.google.com/appengine/docs/go/config/indexconfig#updating_indexes

type FoundArticles struct {
	Articles    []*Article
	ArticleKeys []string
	CursorOne   string
	CursorNext  string
}

/*
func (obj *ArticleManager) FindArticleFromUserName(ctx context.Context, userName string, cursorSrc string, keyOnly bool) *FoundArticles {
	q := datastore.NewQuery(obj.config.KindArticle)
	q = q.Filter("UserName =", userName) ////
	q = q.Order("-Updated").Limit(obj.config.LimitOfFinding)
	return obj.FindArticleFromQuery(ctx, q, cursorSrc, keyOnly)
}*/

func (obj *ArticleManager) FindArticleFromArticleId(ctx context.Context, articleId string, cursorSrc string, keyOnly bool) *FoundArticles {
	q := datastore.NewQuery(obj.config.KindArticle)
	q = q.Filter("ArticleId =", articleId)
	return obj.FindArticleFromQuery(ctx, q, cursorSrc, keyOnly)
}

func (obj *ArticleManager) FindArticleWithNewOrder(ctx context.Context, cursorSrc string, keyOnly bool) *FoundArticles {
	q := datastore.NewQuery(obj.config.KindArticle)
	return obj.FindArticleFromQuery(ctx, q, cursorSrc, keyOnly)
}

func (obj *ArticleManager) FindArticleFromQuery(ctx context.Context, q *datastore.Query, cursorSrc string, keyOnly bool) *FoundArticles {
	cursor := obj.newCursorFromSrc(cursorSrc)
	if cursor != nil {
		q = q.Start(*cursor)
	}
	q = q.KeysOnly()
	founds := q.Run(ctx)

	var retUser []*Article
	var articleIds []string = make([]string, 0)

	var cursorNext string = ""
	var cursorOne string = ""
	for i := 0; ; i++ {
		key, err := founds.Next(nil)

		if err != nil || err == datastore.Done {
			break
		} else {
			articleIds = append(articleIds, key.StringID())
			if keyOnly == false {
				userObj, errNewUserObj := obj.NewArticleFromGaeObjectKey(ctx, key)
				if errNewUserObj == nil {
					retUser = append(retUser, userObj)
				}
			}
		}
		if i == 0 {
			cursorOne = obj.makeCursorSrc(founds)
		}
	}
	cursorNext = obj.makeCursorSrc(founds)
	return &FoundArticles{
		Articles:    retUser,
		ArticleKeys: articleIds,
		CursorNext:  cursorNext,
		CursorOne:   cursorOne,
	}
}

func (obj *ArticleManager) newCursorFromSrc(cursorSrc string) *datastore.Cursor {
	c1, e := datastore.DecodeCursor(cursorSrc)
	if e != nil {
		return nil
	} else {
		return &c1
	}
}

func (obj *ArticleManager) makeCursorSrc(founds *datastore.Iterator) string {
	c, e := founds.Cursor()
	if e == nil {
		return c.String()
	} else {
		return ""
	}
}

//
//
//
type ArtQuery struct {
	q *datastore.Query
}

func (obj *ArticleManager) NewArtQuery() *ArtQuery {
	return &ArtQuery{
		q: datastore.NewQuery(obj.config.KindArticle),
	}

}

func (obj *ArtQuery) GetQuery() *datastore.Query {
	return obj.q
}

func (obj *ArtQuery) WithProp(ctx context.Context, props map[string]string) *ArtQuery {
	for k, v := range props {
		p := miniprop.NewMiniProp()
		p.SetString(k, v)
		v := string(p.ToJson())
		obj.q = obj.q.Filter("Props.Value =", v) ////
	}
	return obj
}

func (obj *ArtQuery) WithTags(ctx context.Context, tags []string) *ArtQuery {
	for _, tag := range tags {
		obj.q = obj.q.Filter("Tags.Tag =", tag) ////
	}
	return obj
}

func (obj *ArtQuery) WithUserName(ctx context.Context, userName string) *ArtQuery {
	obj.q = obj.q.Filter("UserName =", userName)
	return obj
}

func (obj *ArtQuery) WithUpdateMinus(ctx context.Context) *ArtQuery {
	obj.q = obj.q.Order("-Updated")
	return obj
}

func (obj *ArtQuery) WithLimitOfFinding(ctx context.Context, limitOfFinding int) *ArtQuery {
	obj.q = obj.q.Limit(limitOfFinding)
	return obj
}

package pointer

import (
	"time"

	m "github.com/kyorohiro/k07me/prop"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

type PointerKeyInfo struct {
	IdentifyType string
	Identify     string
	//Kind         string
	RootGroup string
}

func (obj *PointerManager) NewPointer(ctx context.Context, //screenName string, //
	userId string, identifyType string, infos map[string]string) *Pointer {
	gaeKey := obj.NewPointerGaeKey(ctx, userId, identifyType)
	gaeObj := GaePointerItem{
		//		PointerName: screenName,
		PointerId:   userId,
		PointerType: identifyType,
		RootGroup:   obj.rootGroup,
	}
	propObj := m.NewMiniPropFromJson([]byte(gaeObj.Info))
	for k, v := range infos {
		propObj.SetString(k, v)
	}
	gaeObj.Info = string(propObj.ToJson())
	gaeObj.Update = time.Now()
	return &Pointer{
		gaeObj: &gaeObj,
		gaeKey: gaeKey,
		kind:   obj.kind,
	}
}

func (obj *PointerManager) NewPointerGaeKey(ctx context.Context, identify string, identifyType string) *datastore.Key {
	return obj.NewPointerGaeKeyFromStringId(ctx, obj.MakePointerStringId(identify, identifyType))
}

func (obj *PointerManager) NewPointerGaeKeyFromStringId(ctx context.Context, stringId string) *datastore.Key {
	return datastore.NewKey(ctx, obj.kind, stringId, 0, nil)
}

func (obj *PointerManager) MakePointerStringId(identify string, identifyType string) string {
	prop := m.NewMiniProp()
	//prop.SetString("k", obj.kind)
	prop.SetString("g", obj.rootGroup)
	prop.SetString("i", identify)
	prop.SetString("t", identifyType)
	return string(prop.ToJson())
}

func (obj *PointerManager) GetKeyInfoFromStringId(stringId string) PointerKeyInfo {
	prop := m.NewMiniPropFromJson([]byte(stringId))
	return PointerKeyInfo{
		//Kind:         prop.GetString("k", ""),
		RootGroup:    prop.GetString("g", ""),
		Identify:     prop.GetString("i", ""),
		IdentifyType: prop.GetString("t", ""),
	}
}

//
//
func (obj *PointerManager) DeletePointerFromObj(ctx context.Context, item *Pointer) error {
	return obj.DeletePointer(ctx, item.GetId(), item.GetType())
}

func (obj *PointerManager) DeletePointer(ctx context.Context, userId, identifyType string) error {
	gaeKey := obj.NewPointerGaeKey(ctx, userId, identifyType)
	ret := datastore.Delete(ctx, gaeKey)
	obj.DeleteMemcache(ctx, gaeKey.StringID())
	return ret
}

func (obj *PointerManager) DeletePointerFromStringId(ctx context.Context, stringId string) error {
	gaeKey := obj.NewPointerGaeKeyFromStringId(ctx, stringId)
	ret := datastore.Delete(ctx, gaeKey)
	obj.DeleteMemcache(ctx, gaeKey.StringID())
	return ret
}

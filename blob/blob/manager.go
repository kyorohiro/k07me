package blob

import (
	"golang.org/x/net/context"

	p "github.com/kyorohiro/k07me/pointer"

	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

type BlobManager struct {
	config     BlobManagerConfig
	pointerMgr *p.PointerManager
}

type BlobManagerConfig struct {
	RootGroup              string
	Kind                   string
	PointerKind            string
	CallbackUrl            string
	MemcachedOnlyInPointer bool
	HashLength             int
}

func NewBlobManager(config BlobManagerConfig) *BlobManager {
	ret := new(BlobManager)
	ret.config = config
	ret.pointerMgr = p.NewPointerManager(p.PointerManagerConfig{
		RootGroup:     config.RootGroup,
		Kind:          config.PointerKind,
		MemcachedOnly: config.MemcachedOnlyInPointer, // todo
	})
	return ret
}

func (obj *BlobManager) GetPointerMgr() *p.PointerManager {
	return obj.pointerMgr
}

func (obj *BlobManager) GetPointer(ctx context.Context, parent, name string) (*p.Pointer, error) {
	return obj.pointerMgr.GetPointer(ctx, obj.MakeBlobId(parent, name), p.TypePointer)
}

func (obj *BlobManager) GetPointerGaeKey(ctx context.Context, parent, name string) *datastore.Key {
	return obj.pointerMgr.NewPointerGaeKey(ctx, obj.MakeBlobId(parent, name), p.TypePointer)
}

func Debug(ctx context.Context, message string) {
	log.Infof(ctx, message)
}

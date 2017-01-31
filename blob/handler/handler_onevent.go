package handler

import (
	"net/http"

	miniblob "github.com/kyorohiro/k07me/blob/blob"
	mm "github.com/kyorohiro/k07me/prop"
)

func (obj *BlobHandler) AddOnBlobRequest(f func(w http.ResponseWriter, r *http.Request, input *mm.MiniProp, output *mm.MiniProp, h *BlobHandler) (map[string]string, error)) {
	obj.onEvent.OnBlobRequestList = append(obj.onEvent.OnBlobRequestList, f)
}

func (obj *BlobHandler) AddOnBlobBeforeSave(f func(http.ResponseWriter, *http.Request, *mm.MiniProp, *BlobHandler, *miniblob.BlobItem) error) {
	obj.onEvent.OnBlobBeforeSaveList = append(obj.onEvent.OnBlobBeforeSaveList, f)
}

func (obj *BlobHandler) AddOnBlobComplete(f func(http.ResponseWriter, *http.Request, *mm.MiniProp, *BlobHandler, *miniblob.BlobItem) error) {
	obj.onEvent.OnBlobCompleteList = append(obj.onEvent.OnBlobCompleteList, f)
}

func (obj *BlobHandler) AddOnBlobFailed(f func(http.ResponseWriter, *http.Request, *mm.MiniProp, *BlobHandler, *miniblob.BlobItem)) {
	obj.onEvent.OnBlobFailedList = append(obj.onEvent.OnBlobFailedList, f)
}

func (obj *BlobHandler) OnBlobRequestList(w http.ResponseWriter, r *http.Request, i *mm.MiniProp, o *mm.MiniProp, h *BlobHandler) (map[string]string, error) {
	ret := map[string]string{}
	for _, f := range obj.onEvent.OnBlobRequestList {
		vsTmp, err := f(w, r, i, o, h)
		for k, v := range vsTmp {
			ret[k] = v
		}
		if err != nil {
			return ret, err
		}
	}
	return ret, nil
}

func (obj *BlobHandler) OnBlobBeforeSave(w http.ResponseWriter, r *http.Request, o *mm.MiniProp, h *BlobHandler, i *miniblob.BlobItem) error {
	for _, f := range obj.onEvent.OnBlobBeforeSaveList {
		err := f(w, r, o, h, i)
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *BlobHandler) OnBlobComplete(w http.ResponseWriter, r *http.Request, o *mm.MiniProp, h *BlobHandler, i *miniblob.BlobItem) error {
	for _, f := range obj.onEvent.OnBlobCompleteList {
		err := f(w, r, o, h, i)
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *BlobHandler) OnBlobFailed(w http.ResponseWriter, r *http.Request, o *mm.MiniProp, h *BlobHandler, i *miniblob.BlobItem) {
	for _, f := range obj.onEvent.OnBlobFailedList {
		f(w, r, o, h, i)
	}
}

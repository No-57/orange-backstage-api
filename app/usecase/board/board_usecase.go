package board

import (
	"context"
	"errors"
	"io"
	"mime/multipart"
	"orange-backstage-api/app/model"
	"orange-backstage-api/app/store"
	"orange-backstage-api/infra/api"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

type Usecase struct {
	store *store.Store

	uploadPath string
}

type Param struct {
	UploadPath string
}

func New(store *store.Store, param Param) *Usecase {
	return &Usecase{
		store:      store,
		uploadPath: param.UploadPath,
	}
}

func (u *Usecase) List(ctx context.Context) ([]model.Board, error) {
	records, err := u.store.Board.Select(ctx)
	if err != nil {
		return nil, api.NewStoreErr(err)
	}

	return records, nil
}

type CreateParam struct {
	Image *multipart.FileHeader
	Board *model.Board
}

func (u *Usecase) Create(ctx context.Context, param CreateParam) error {
	dst := filepath.Join(u.uploadPath, uuid.NewString()+filepath.Ext(param.Image.Filename))

	src, err := param.Image.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	if err = os.MkdirAll(filepath.Dir(dst), 0750); err != nil {
		return err
	}

	// #nosec G304: only create file from trusted source
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err = io.Copy(out, src); err != nil {
		return err
	}

	param.Board.ImageURL = dst
	if err := u.store.Board.Create(ctx, param.Board); err != nil {
		return api.NewStoreErr(err)
	}

	return nil
}

func (u *Usecase) Delete(ctx context.Context, id uint64) error {
	if id == 0 {
		return api.NewParamErr(errors.New("id is required"))
	}

	record, err := u.store.Board.SeleteByID(ctx, id)
	if err != nil {
		if errors.Is(err, store.ErrRecordNotFound) {
			return api.NewParamErr(err)
		}
		return api.NewStoreErr(err)
	}
	defer func() {
		os.Remove(record.ImageURL)
	}()

	if err := u.store.Board.Delete(ctx, id); err != nil {
		return api.NewStoreErr(err)
	}

	return nil
}

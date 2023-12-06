package board

import (
	"mime/multipart"
	"net/url"
	"orange-backstage-api/app/model"
	"orange-backstage-api/app/usecase/board"
	"orange-backstage-api/infra/api"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

type Router struct {
	usecase *board.Usecase
}

func New(usecase *board.Usecase) *Router {
	return &Router{usecase: usecase}
}

func (r *Router) Register(ginR gin.IRouter) {
	ginR.GET("/boards", r.List)
	ginR.POST("/boards", r.Create)
	ginR.DELETE("/boards/:id", r.Delete)
}

type ListItem struct {
	ID          uint64    `json:"id"`
	Code        string    `json:"code"`
	ImageURL    string    `json:"image_url"`
	ActionType  string    `json:"action_type"`
	Action      string    `json:"action"`
	CreatedDate time.Time `json:"created_date"`
	UpdatedDate time.Time `json:"updated_date"`
}

// List
//
//	@Summary		List boards
//	@Description	List boards
//	@Tags			board
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	api.DataResp{data=[]ListItem}
//	@Failure		400	{object}	api.ErrResp
//	@Failure		500	{object}	api.ErrResp
//	@Router			/boards [get]
func (r *Router) List(c *gin.Context) {
	ctx := api.NewCtx(c)

	boards, err := r.usecase.List(ctx)
	if err != nil {
		ctx.Resp().HandleErr(err)
		return
	}

	payload := make([]ListItem, 0, len(boards))
	for _, board := range boards {
		base := filepath.Base(board.ImageURL)
		imageURL, err := url.JoinPath("/images", base)
		if err != nil {
			imageURL = ""
		}

		payload = append(payload, ListItem{
			ID:          board.ID,
			Code:        board.Code,
			ImageURL:    c.Request.Host + imageURL,
			ActionType:  board.ActionType,
			Action:      board.Action,
			CreatedDate: board.CreatedAt,
			UpdatedDate: board.UpdatedAt,
		})
	}

	ctx.Resp().Data(payload)
}

type CreateJSON struct {
	Code       string `json:"code" form:"code" binding:"required"`
	ActionType string `json:"action_type" form:"action_type" binding:"required"`
	Action     string `json:"action" form:"action" binding:"required"`
}

type CreateReq struct {
	File *multipart.FileHeader `form:"file" binding:"required"`
	Data CreateJSON            `form:"data" binding:"required"`
}

// Create
//
//	@Summary		Create board
//	@Description	Create board with image
//	@Tags			board
//	@Accept			multipart/form-data
//	@Produce		json
//	@Security		BearerAuth
//	@Param			file	formData	file		true	"image file"
//	@Param			data	formData	CreateJSON	true	"payload"
//	@Success		200		{object}	api.CodeResp
//	@Failure		400		{object}	api.ErrResp
//	@Failure		500		{object}	api.ErrResp
//	@Router			/boards [post]
func (r *Router) Create(c *gin.Context) {
	ctx := api.NewCtx(c)

	var req CreateReq
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.Resp().InvalidParam(err)
		return
	}

	param := board.CreateParam{
		Image: req.File,
		Board: &model.Board{
			Code:       req.Data.Code,
			ActionType: req.Data.ActionType,
			Action:     req.Data.Action,
		},
	}

	if err := r.usecase.Create(ctx, param); err != nil {
		ctx.Resp().HandleErr(err)
		return
	}

	ctx.Resp().OK()
}

// Delete
//
//	@Summary		Delete board
//	@Description	Delete board
//	@Tags			board
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		uint64	true	"board id"
//	@Success		200	{object}	api.CodeResp
//	@Failure		400	{object}	api.ErrResp
//	@Failure		500	{object}	api.ErrResp
//	@Router			/boards/{id} [delete]
func (r *Router) Delete(c *gin.Context) {
	ctx := api.NewCtx(c)

	id := ctx.Param().ID()
	if id == 0 {
		ctx.Resp().NotFound()
		return
	}

	if err := r.usecase.Delete(ctx, id); err != nil {
		ctx.Resp().HandleErr(err)
		return
	}

	ctx.Resp().OK()
}

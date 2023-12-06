package theme

import (
	"orange-backstage-api/app/usecase/theme"
	"orange-backstage-api/infra/api"
	"time"

	"github.com/gin-gonic/gin"
)

type Router struct {
	usecase *theme.Usecase
}

func New(usecase *theme.Usecase) *Router {
	return &Router{usecase: usecase}
}

func (r *Router) Register(router *gin.RouterGroup) {
	router.GET("/themes", r.List)
	router.POST("/themes", r.Create)
	router.DELETE("/themes/:id", r.Delete)
}

type ListItem struct {
	ID          uint64    `json:"id"`
	Code        string    `json:"code"`
	Type        string    `json:"type"`
	Disable     bool      `json:"disable"`
	CreatedDate time.Time `json:"created_date"`
	UpdatedDate time.Time `json:"updated_date"`
}

// List
//
//	@Summary		List themes
//	@Description	List themes
//	@Tags			theme
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	api.DataResp{data=[]ListItem}
//	@Failure		400	{object}	api.ErrResp
//	@Failure		500	{object}	api.ErrResp
//	@Router			/themes [get]
func (r *Router) List(c *gin.Context) {
	ctx := api.NewCtx(c)

	themes, err := r.usecase.List(ctx)
	if err != nil {
		ctx.Resp().HandleErr(err)
		return
	}

	payload := make([]ListItem, 0, len(themes))
	for _, theme := range themes {
		payload = append(payload, ListItem{
			ID:          theme.ID,
			Code:        theme.Code,
			Type:        theme.Type,
			Disable:     theme.Disable,
			CreatedDate: theme.CreatedAt,
			UpdatedDate: theme.UpdatedAt,
		})
	}

	ctx.Resp().Data(payload)
}

type CreateReq struct {
	Code    string `json:"code" binding:"required"`
	Type    string `json:"type" binding:"required"`
	Disable bool   `json:"disable"`
}

// Create
//
//	@Summary		Create theme
//	@Description	Create theme
//	@Tags			theme
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			body	body		CreateReq	true	"body"
//	@Success		200		{object}	api.CodeResp
//	@Failure		400		{object}	api.ErrResp
//	@Failure		500		{object}	api.ErrResp
//	@Router			/themes [post]
func (r *Router) Create(c *gin.Context) {
	ctx := api.NewCtx(c)

	var req CreateReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Resp().InvalidParam(err)
		return
	}

	if err := r.usecase.Create(ctx, theme.CreateParam{
		Code:    req.Code,
		Type:    req.Type,
		Disable: req.Disable,
	}); err != nil {
		ctx.Resp().HandleErr(err)
		return
	}

	ctx.Resp().Data(nil)
}

// Delete
//
//	@Summary		Delete theme
//	@Description	Delete theme
//	@Tags			theme
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		uint64	true	"theme id"
//	@Success		200	{object}	api.CodeResp
//	@Failure		400	{object}	api.ErrResp
//	@Failure		500	{object}	api.ErrResp
//	@Router			/themes/{id} [delete]
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

	ctx.Resp().Data(nil)
}

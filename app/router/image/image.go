package image

import "github.com/gin-gonic/gin"

type Router struct {
	uploadPath string
}

type Param struct {
	UploadPath string
}

func New(param Param) *Router {
	return &Router{uploadPath: param.UploadPath}
}

func (r *Router) Register(ginR gin.IRouter) {
	ginR.GET("/images/:name", r.Get)
}

// Get
//
//	@Summary		Get image
//	@Description	Get image by name
//	@Tags			image
//	@Accept			json
//	@Produce		json
//	@Param			name	path		string	true	"image name"
//	@Success		200		{file}		file
//	@Failure		404		{object}	string	"404 page not found"
//	@Router			/images/{name} [get]
func (r *Router) Get(c *gin.Context) {
	c.File(r.uploadPath + "/" + c.Param("name"))
}

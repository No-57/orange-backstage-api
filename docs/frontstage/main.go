package frontstage

//	@title			Orange Frontstage API Document
//	@version		0.1.0
//	@description	For Orange Frontstage API Document
//	@host			localhost:8080
//	@BasePath		/api/v1

//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization
//	@description				Enter the access token with the `Bearer ` prefix, e.g. "Bearer abcde12345".

type ProductItem struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`

	// Type: seller type
	//	- phone
	//	- laptop
	//	- desktop
	//	- audio
	//	- tablet
	//	- earphone
	Type string `json:"type"`
}

// GetProducts
//
//	@Summary		Get products
//	@Description	Get products with pagination
//	@Tags			product
//	@Accept			json
//	@Produce		json
//	@Param			page		query		int		false	"page. Default: 1 "
//	@Param			page_size	query		int		false	"page size. Default: 10"
//	@Param			name		query		string	false	"name. Split by comma. Default: empty"
//	@Param			sort_by		query		string	false	"order by. Default: id"
//	@Param			order_by	query		string	false	"order. desc or asc. Default: asc"
//	@Param			fields		query		string	false	"fields to show, split by comma. Default: all fields"
//	@Success		200			{object}	api.DataResp{data=[]ProductItem}
//	@Failure		400			{object}	api.ErrResp
//	@Failure		500			{object}	api.ErrResp
//	@Router			/products [get]
func GetProducts() {}

type ProductDetail struct {
	ID          uint64 `json:"id"`
	Name        string `json:"name"`
	Price       string `json:"price"`
	Description string `json:"description"`

	// Like: is this product liked by current user
	// Use `SetFavorite` API to set this field.
	Like      bool   `json:"like"`
	ImageURL  string `json:"image_url"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// GetProductDetail
//
//	@Summary		Get product detail
//	@Description	Get product detail by specific id
//	@Tags			product
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"product id"
//	@Success		200	{object}	api.DataResp{data=ProductDetail}
//	@Failure		400	{object}	api.ErrResp
//	@Failure		500	{object}	api.ErrResp
//	@Router			/products/{id} [get]
func GetProductDetail() {}

type SetFavoriteReq struct {
	IDs      []uint64 `json:"ids"`
	Favorite bool     `json:"favorite"`
}

// SetFavorite
//
//	@Summary		Set favorite prodcuts
//	@Description	Set favorite products
//	@Tags			product
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			data	body		SetFavoriteReq	true	"payload"
//	@Success		200		{object}	api.DataResp
//	@Failure		400		{object}	api.ErrResp
//	@Failure		500		{object}	api.ErrResp
//	@Router			/products/favorite [put]
func SetFavorite() {}

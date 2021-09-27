package http

type ErrorResponse struct {
	Message string `json:"message" example:"error message"`
	Code    string `json:"code,omitempty" example:"IR001"`
}

type PageResponse struct {
	Data   interface{}        `json:"data" example:"{'item1', 'item2', 'item3'}" swaggertype:"ArrayOfString"`
	Paging PagePagingResponse `json:"paging"`
}

type PagePagingResponse struct {
	Previous string `json:"previous,omitempty" example:"https://quorum-key-manager.com/stores/your-store/secrets?page=1"`
	Next     string `json:"next,omitempty" example:"https://quorum-key-manager.com/stores/your-store/secrets?page=3"`
}

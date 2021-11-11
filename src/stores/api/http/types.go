package http

// @TODO Remove duplicated types once make gen-swagger issue is identified

// https://github.com/ConsenSys/quorum-key-manager/actions/runs/1170918392
type ErrorResponse struct {
	Message string `json:"message" example:"error message"`
	Code    string `json:"code,omitempty" example:"IR001"`
}

type PageResponse struct {
	Data   []string           `json:"data" example:"item1,item2,item3"`
	Paging PagePagingResponse `json:"paging"`
}

type PagePagingResponse struct {
	Previous string `json:"previous,omitempty" example:"http://quorum-key-manager.com/stores/your-store/secrets?page=1"`
	Next     string `json:"next,omitempty" example:"http://quorum-key-manager.com/stores/your-store/secrets?page=3"`
}

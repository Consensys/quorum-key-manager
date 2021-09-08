package types

type PageResponse struct {
	Data   []interface{} `json:"data" example:"{'item1', 'item2', 'item3'}"`
	Paging PagePagingResponse `json:"paging"`
}

type PagePagingResponse struct {
	Previous string `json:"previous,omitempty" example:"http://quorum-key-manager.com/stores/your-store/secrets?page=1"`
	Next     string `json:"next,omitempty" example:"http://quorum-key-manager.com/stores/your-store/secrets?page=3"`
}

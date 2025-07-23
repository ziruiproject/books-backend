package pagination

import "math"

type Page[T any] struct {
	Items       []T `json:"items"`
	CurrentPage int `json:"current_page"`
	PerPage     int `json:"per_page"`
	TotalPages  int `json:"total_pages"`
	TotalItems  int `json:"total_items"`
	NextPage    int `json:"next_page"`
	PrevPage    int `json:"prev_page"`
}

type Request struct {
	SortBy  string `json:"sort_by"`
	OrderBy string `json:"order_by"`
	Page    int    `json:"page"`
	Limit   int    `json:"limit"`
}

func NewPagination(request *Request) *Request {
	if request.Page < 1 {
		request.Page = 1
	}
	if request.Limit <= 0 {
		request.Limit = 10
	}
	if request.OrderBy == "" {
		request.OrderBy = "created_at"
	}
	if request.SortBy == "" {
		request.SortBy = "ASC"
	}

	return &Request{
		SortBy:  request.SortBy,
		OrderBy: request.OrderBy,
		Page:    request.Page,
		Limit:   request.Limit,
	}
}

func NewPage[T any](req Request, totalItems int64, items []T) *Page[T] {
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}

	totalPages := int(math.Ceil(float64(totalItems) / float64(req.Limit)))

	nextPage := 0
	if req.Page < totalPages {
		nextPage = req.Page + 1
	}

	prevPage := 0
	if req.Page > 1 {
		prevPage = req.Page - 1
	}

	return &Page[T]{
		Items:       items,
		CurrentPage: req.Page,
		PerPage:     req.Limit,
		TotalPages:  totalPages,
		TotalItems:  int(totalItems),
		NextPage:    nextPage,
		PrevPage:    prevPage,
	}
}

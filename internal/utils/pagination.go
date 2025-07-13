package utils

type PaginationParams struct {
	Limit  int
	Cursor *int
}

type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	NextCursor *int        `json:"next_cursor"`
	HasMore    bool        `json:"has_more"`
}

const DefaultLimit = 20
const MaxLimit = 100

type Identifiable interface {
	GetID() int
}

func CreatePaginatedResponse[T Identifiable](items []T, limit int) PaginatedResponse {
	response := PaginatedResponse{
		Data:    items,
		HasMore: false,
	}

	if len(items) == limit {
		response.HasMore = true
		if len(items) > 0 {
			cursor := items[len(items)-1].GetID()
			response.NextCursor = &cursor
		}
	}

	return response
}

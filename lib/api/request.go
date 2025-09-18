package api

type PaginationRequest struct {
	Limit  uint `query:"limit"`
	Offset uint `query:"offset"`
}

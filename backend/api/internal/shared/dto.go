package shared

type GenericResponse[T any] struct {
	Data T `json:"data"`
}

type Error struct {
	Field string `json:"field"`
	Tag   string `json:"tag"`
}

type ErrorResponse struct {
	Error  string `json:"error"`
	Status int    `json:"status"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

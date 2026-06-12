package response

type ErrorDTO struct {
	Data    any    `json:"data"`
	Message string `json:"message"`
}
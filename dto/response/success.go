package response

type JSONResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
}
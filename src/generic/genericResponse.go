package generic

type ApiGenericResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

func NewSuccessResponse(data interface{}) ApiGenericResponse {
	return ApiGenericResponse{
		Success: true,
		Data:    data,
	}
}

func NewErrorResponse(message string, err interface{}) ApiGenericResponse {
	return ApiGenericResponse{
		Success: false,
		Message: message,
		Error:   err,
	}
}

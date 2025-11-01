package domain

type Response struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Errors  interface{} `json:"errors,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func ErrorResponse(message string) Response {
	return Response{
		Status:  false,
		Message: message,
	}
}

func ErrorResponseWithData(message string, data interface{}) Response {
	return Response{
		Status:  false,
		Message: message,
		Data:    data,
	}
}

func SuccessResponse(message string, data interface{}) Response {
	return Response{
		Status:  true,
		Message: message,
		Data:    data,
	}
}

type PaginationResponse struct {
	Page      int         `json:"page"`
	Limit     int         `json:"limit"`
	TotalPage int         `json:"total_page,omitempty"`
	TotalData int         `json:"total_data,omitempty"`
	Data      interface{} `json:"data,omitempty"`
}
package response

type SuccessResponse struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func Success(args ...interface{}) map[string]interface{} {
	switch len(args) {
	case 1:
		return map[string]interface{}{
			"status": "success",
			"data":   args[0],
		}
	case 2:
		return map[string]interface{}{
			"message": args[0],
			"data":    args[1],
		}
	default:
		return map[string]interface{}{"status": "success"}
	}
}

func Error(message string) ErrorResponse {
	return ErrorResponse{
		Status:  "error",
		Message: message,
	}
}

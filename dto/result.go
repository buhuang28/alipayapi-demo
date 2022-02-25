package dto

type Result struct {
	Status  int64       `json:"status"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func (r *Result) Success(message string, data interface{}) Result {
	r.Status = 1
	if message == "" {
		r.Message = "success"
	} else {
		r.Message = message
	}
	r.Data = data
	return *r
}

func (r *Result) Error(status int64, message string) Result {
	if status >= 0 {
		status = -1
	}
	r.Status = status
	r.Message = message
	return *r
}

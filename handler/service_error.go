package handler

type ServiceError struct {
	StatusCode int
	Message    string
	Detail     string
}

func (e *ServiceError) Error() string {
	return e.Message
}

func (e *ServiceError) DetailError() string {
	return e.Detail
}

func NewServiceError(statusCode int, message, detail string) *ServiceError {
	return &ServiceError{
		StatusCode: statusCode,
		Message:    message,
		Detail:     detail,
	}
}

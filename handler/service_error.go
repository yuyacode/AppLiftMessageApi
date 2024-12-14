package handler

type ServiceError struct {
	StatusCode int
	Message    string
	Detail     string
}

func (se *ServiceError) Error() string {
	return se.Message
}

func (se *ServiceError) DetailError() string {
	return se.Detail
}

func NewServiceError(statusCode int, message, detail string) *ServiceError {
	return &ServiceError{
		StatusCode: statusCode,
		Message:    message,
		Detail:     detail,
	}
}

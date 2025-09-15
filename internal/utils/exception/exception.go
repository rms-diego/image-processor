package exception

type AppError struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

type Error interface {
	Error() string
}

func (e *AppError) Error() string {

	return e.Message
}

func New(message string, code int) Error {
	return &AppError{Message: message, Code: code}
}

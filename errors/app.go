package errors

type AppCode int

const (
	Conflict AppCode = iota
	NotFound
	InvalidInput
	Generic
)

type AppError struct {
	Code    AppCode
	Message string
}

func (err *AppError) Error() string {
	return err.Message
}

func ConflictError(message string) *AppError {
	return &AppError{
		Code:    Conflict,
		Message: message,
	}
}

func NotFoundError(message string) *AppError {
	return &AppError{
		Code:    NotFound,
		Message: message,
	}
}

func InvalidInputError(message string) *AppError {
	return &AppError{
		Code:    InvalidInput,
		Message: message,
	}
}

func GenericError(message string) *AppError {
	return &AppError{
		Code:    Generic,
		Message: message,
	}
}

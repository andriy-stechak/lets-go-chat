package errors

func ToHttpError(err *AppError) *AppHttpError {
	switch err.Code {
	case Conflict:
		return HttpConflict(err.Message)
	case NotFound:
		return HttpNotFound(err.Message)
	case InvalidInput:
		return HttpBadRequest(err.Message)
	case Generic:
		return HttpInternalError(err.Message)
	}
	return HttpInternalError(err.Message)
}

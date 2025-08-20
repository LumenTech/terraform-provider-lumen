package helper

func IsHttpClientError(statusCode int) bool {
	return statusCode >= 400 && statusCode < 500
}

func IsHttpServerError(statusCode int) bool {
	return statusCode >= 500
}

package models

type HttpError struct {
	ErrorCode 		string `json:"err_code"`
	ErrorMessage	interface{} `json:"err_msg"`
}

// TODO refactor errors according to this: https://blog.golang.org/go1.13-errors

type HTTPErrorLong struct {
	Error 			HttpError
	StatusCode		int // an HTTP error code
}

func NewHTTPErrorLong(errorCode string, errorMsg interface{}, statusCode int) HTTPErrorLong {
	errBody := HttpError{
		ErrorCode:    errorCode,
		ErrorMessage: errorMsg,
	}
	longErr := HTTPErrorLong{
		Error:      errBody,
		StatusCode: statusCode,
	}
	return longErr
}
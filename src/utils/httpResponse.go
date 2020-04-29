package utils

type HttpError struct {
	ErrorCode 		string `json:"err_code"`
	ErrorMessage	interface{} `json:"err_msg"`
}

type HTTPErrorLong struct {
	Error 			HttpError
	StatusCode		int // an HTTP error code
}
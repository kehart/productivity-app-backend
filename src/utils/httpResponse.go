package utils

type Metadata interface{} // placeholder for better def

type HTTPResponseObject struct {
	Meta Metadata `json:"_meta"`
	Data interface{} `json:"data"`
}

type HTTPResponseCollection struct {
	Meta Metadata `json:"_meta"`
	Items []interface{} `json:"items"`
}

type HttpError struct {
	ErrorCode 		string `json:"err_code"`
	ErrorMessage	interface{}`json:"err_msg"`
}

type HTTPErrorLong struct {
	Error 			HttpError
	StatusCode		int // an HTTP error code
}
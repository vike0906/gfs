package common

type Response struct {
	Code    uint16      `json:"code"`
	Message string      `json:"message,omitempty"`
	Content interface{} `json:"content,omitempty"`
}

const (
	SuccessCode    = 0
	FailCode       = 200
	SuccessMessage = "success"
)

func ResponseInstance() *Response {
	return new(Response)
}

func (r *Response) Success() Response {
	return Response{Code: SuccessCode, Message: SuccessMessage}
}

func (r *Response) SuccessWithMessage(message string) Response {
	return Response{Code: SuccessCode, Message: message}
}

func (r *Response) SuccessWithContent(content interface{}) Response {
	return Response{Code: SuccessCode, Message: SuccessMessage, Content: content}
}

func (r *Response) Fail(message string) Response {
	return Response{Code: FailCode, Message: message}
}

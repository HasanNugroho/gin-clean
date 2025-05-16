package response

type Meta struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type Response struct {
	Data interface{} `json:"data"`
	Meta Meta        `json:"meta"`
}

func Ok(data interface{}) Response {
	return Response{
		Data: data,
		Meta: Meta{
			Code:    "ok",
			Message: "Ok",
		},
	}
}

func Err(code string, message string) Response {
	return Response{
		Data: nil,
		Meta: Meta{
			Code:    code,
			Message: message,
		},
	}
}

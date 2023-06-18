package rest

type Response struct {
	code int
	raw  string
}

func (r *Response) GetRaw() string {
	return r.raw
}

func (r *Response) GetCode() int {
	return r.code
}

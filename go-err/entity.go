package go_err

type ResponseError struct {
	Status  string     `json:"status,omitempty"`
	Message string     `json:"message,omitempty"`
	Meta    *ErrorMeta `json:"meta,omitempty"`
}

type ErrorMeta struct {
	Code   int        `json:"code,omitempty"`
	Name   string     `json:"name,omitempty"`
	Fields ErrorField `json:"fields,omitempty"`
}

func NewResponseError(c CustomError) *ResponseError {
	resp := &ResponseError{
		Status:  "error",
		Message: c.Error(),
	}
	code, name := c.GetCode(), c.GetName()
	if code != 0 || name != "" {
		resp.Meta = &ErrorMeta{Code: code, Name: name}
	}
	return resp
}

func NewResponseErrorWithFields(c CustomError, fields map[string]string) *ResponseError {
	resp := NewResponseError(c)
	resp.SetField(fields)
	return resp
}

func (r *ResponseError) SetField(fields map[string]string) {
	if r.Meta == nil {
		r.Meta = &ErrorMeta{}
	}
	r.Meta.Fields = fields
}

func (r *ResponseError) AddField(key, value string) {
	if r.Meta == nil {
		r.Meta = &ErrorMeta{}
	}
	if r.Meta.Fields == nil {
		r.Meta.Fields = make(ErrorField)
	}
	r.Meta.Fields[key] = value
}

type ErrorField map[string]string

func (f ErrorField) GetFirstErrorAndOtherTotal() (string, int) {
	total := len(f)
	if total > 0 {
		total--
	}
	for k := range f {
		return f[k], total
	}
	return "", 0
}

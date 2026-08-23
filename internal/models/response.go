package models

type APIResponse[R any] struct {
	Error string `json:"error,omitempty"`
	Data  R      `json:"data,omitempty"`
}

func NewDataResponse[R any](data R) APIResponse[R] {
	res := APIResponse[R]{
		Data: data,
	}
	return res
}

func NewErrorResponse(err error) APIResponse[any] {
	if err == nil {
		return APIResponse[any]{
			Error: "Unknown error",
		}
	}
	res := APIResponse[any]{
		Error: err.Error(),
	}
	return res
}

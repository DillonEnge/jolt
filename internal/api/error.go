package api

type ApiError struct {
	Status int
	Err    error
}

func (a *ApiError) Error() string {
	if a.Err == nil {
		return "error not provided"
	}

	return a.Err.Error()
}

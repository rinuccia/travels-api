package handler

type errResponse struct {
	Error string `json:"error"`
}

func newErrResponse(errString string) *errResponse {
	return &errResponse{
		Error: errString,
	}
}

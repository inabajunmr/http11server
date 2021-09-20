package http

type HTTPError struct {
	Msg    string
	Status int
}

func (err *HTTPError) Error() string {
	return err.Msg
}

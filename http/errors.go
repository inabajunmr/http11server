package http

type HTTPError struct {
	Msg    string
	Status int
}

func (err *HTTPError) Error() string {
	return err.Msg
}

type WaitRequestError struct {
	Msg string
}

func (err *WaitRequestError) Error() string {
	return err.Msg
}

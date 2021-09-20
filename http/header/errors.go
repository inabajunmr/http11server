package header

type HeaderParserError struct {
	Msg string
}

func (err *HeaderParserError) Error() string {
	return err.Msg
}

package http

type HTTPVersion int

const (
	HTTP11 = iota
)

func (v HTTPVersion) ToString() string {
	return "HTTP/1.1"
}

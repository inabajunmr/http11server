package header

type ContentCoding int

const (
	CONTENT_CODING_COMPRESS ContentCoding = iota
	CONTENT_CODING_DEFLATE
	CONTENT_CODING_GZIP
	CONTENT_CODING_IDENTITY
)

func getContentCoding(v string) ContentCoding {
	switch v {
	case "compress":
		return CONTENT_CODING_COMPRESS
	case "deflate":
		return CONTENT_CODING_DEFLATE
	case "identity":
		return CONTENT_CODING_IDENTITY
	case "gzip":
		return CONTENT_CODING_GZIP
	}
	return 0 // TODO
}

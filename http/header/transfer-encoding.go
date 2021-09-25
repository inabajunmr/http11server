package header

type TransferEncoding int

const (
	TRANSFER_ENCODING_CHUNKED TransferEncoding = iota
	TRANSFER_ENCODING_COMPRESS
	TRANSFER_ENCODING_DEFLATE
	TRANSFER_ENCODING_GZIP
	TRANSFER_ENCODING_IDENTITY
)

func getTransferEncoding(v string) TransferEncoding {
	switch v {
	case "chunked":
		return TRANSFER_ENCODING_CHUNKED
	case "compress":
		return TRANSFER_ENCODING_COMPRESS
	case "deflate":
		return TRANSFER_ENCODING_DEFLATE
	case "gzip":
		return TRANSFER_ENCODING_GZIP
	case "identity":
		return TRANSFER_ENCODING_IDENTITY
	}
	return 0 // TODO
}

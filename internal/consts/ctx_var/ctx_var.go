package ctx_var

// Define the context keys as constants of type CtxKey
type CtxKey string

const (
	REQUEST_CONTENT CtxKey = "request_content"
	MERCHANT_ID     CtxKey = "merchant_id"
	TRACE_ID        CtxKey = "trace_id"
)

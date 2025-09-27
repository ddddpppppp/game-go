package middleware

import (
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"
)

// ResponseHandler is a middleware that handles responses without logging errors to console.
// It's a replacement for ghttp.MiddlewareHandlerResponse.
func ResponseHandler(r *ghttp.Request) {
	r.Middleware.Next()

	// There's already a response.
	if r.Response.BufferLength() > 0 {
		return
	}

	// No response, it gives the handler result as the response.
	var (
		msg     string
		err     error
		res     interface{}
		codeNum int
	)
	res = r.GetHandlerResponse()
	// Only use error when it's not nil.
	if err = r.GetError(); err != nil {
		// Instead of passing the original error directly, we extract the information
		// to avoid printing the stack trace
		errCode := gerror.Code(err)
		codeNum = errCode.Code()
		msg = err.Error()
		// Do not log the error here, just pass the message to the response
		r.SetError(nil) // Clear the error to prevent it from being logged elsewhere
	} else {
		codeNum = gcode.CodeOK.Code()
	}

	r.Response.WriteJson(ghttp.DefaultHandlerResponse{
		Code:    codeNum,
		Message: msg,
		Data:    res,
	})
}

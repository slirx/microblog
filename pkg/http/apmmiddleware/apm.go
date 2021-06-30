package apmmiddleware

import (
	"bytes"
	"io"
	"net/http"

	"go.elastic.co/apm"
	"go.elastic.co/apm/module/apmhttp"
)

// todo refactor this file. maybe write unit tests if possible

type multiResponseWriter struct {
	responseWriter http.ResponseWriter
	multiWriter    io.Writer
}

func (m multiResponseWriter) Header() http.Header {
	return m.responseWriter.Header()
}

func (m multiResponseWriter) Write(b []byte) (int, error) {
	return m.multiWriter.Write(b)
}

func (m multiResponseWriter) WriteHeader(statusCode int) {
	m.responseWriter.WriteHeader(statusCode)
}

func MultiWriter(l io.Writer, w http.ResponseWriter) http.ResponseWriter {
	multi := io.MultiWriter(l, w)

	return &multiResponseWriter{
		responseWriter: w,
		multiWriter:    multi,
	}
}

// Wrap wraps h such that it will report requests as transactions
// to Elastic APM, using route in the transaction name.
//
// By default, the returned Handle will use apm.DefaultTracer.
// Use WithTracer to specify an alternative tracer.
//
// By default, the returned Handle will recover panics, reporting
// them to the configured tracer. To override this behaviour, use
// WithRecovery.
func Wrap(h http.HandlerFunc, route string, o ...Option) http.HandlerFunc {
	opts := gatherOptions(o...)

	return func(w http.ResponseWriter, req *http.Request) {
		if !opts.tracer.Recording() || opts.requestIgnorer(req) {
			h(w, req)
			return
		}

		var tx *apm.Transaction

		tx, req = apmhttp.StartTransaction(opts.tracer, req.Method+" "+route, req)
		defer tx.End()

		var httpResponse *apmhttp.Response
		w, httpResponse = apmhttp.WrapResponseWriter(w)

		var buf bytes.Buffer
		w = MultiWriter(&buf, w)

		body := opts.tracer.CaptureHTTPRequestBody(req)

		// todo maybe move recover to separate middleware?
		defer func() {
			if v := recover(); v != nil {
				opts.recovery(w, req, httpResponse, body, tx, v)
			}
			apmhttp.SetTransactionContext(tx, req, httpResponse, body)
			body.Discard()
		}()

		h(w, req)

		if httpResponse.StatusCode == 0 {
			httpResponse.StatusCode = http.StatusOK
		}

		// todo will it write panic response to apm?
		tx.Context.SetCustom("response_body", buf.String())
	}
}

//// WrapNotFoundHandler wraps h so that it is traced. If h is nil, then http.NotFoundHandler() will be used.
//func WrapNotFoundHandler(h http.Handler, o ...Option) http.Handler {
//	if h == nil {
//		h = http.NotFoundHandler()
//	}
//	return wrapHandlerUnknownRoute(h, o...)
//}
//
//// WrapMethodNotAllowedHandler wraps h so that it is traced. If h is nil, then a default handler
//// will be used that returns status code 405.
//func WrapMethodNotAllowedHandler(h http.Handler, o ...Option) http.Handler {
//	if h == nil {
//		h = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
//			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
//		})
//	}
//	return wrapHandlerUnknownRoute(h, o...)
//}
//
//func wrapHandlerUnknownRoute(h http.Handler, o ...Option) http.Handler {
//	opts := gatherOptions(o...)
//	return apmhttp.Wrap(
//		h,
//		apmhttp.WithTracer(opts.tracer),
//		apmhttp.WithRecovery(opts.recovery),
//		apmhttp.WithServerRequestName(apmhttp.UnknownRouteRequestName),
//		apmhttp.WithServerRequestIgnorer(opts.requestIgnorer),
//	)
//}

func gatherOptions(o ...Option) options {
	opts := options{
		tracer: apm.DefaultTracer,
	}
	for _, o := range o {
		o(&opts)
	}
	if opts.requestIgnorer == nil {
		opts.requestIgnorer = apmhttp.NewDynamicServerRequestIgnorer(opts.tracer)
	}
	if opts.recovery == nil {
		opts.recovery = apmhttp.NewTraceRecovery(opts.tracer)
	}
	return opts
}

type options struct {
	tracer         *apm.Tracer
	recovery       apmhttp.RecoveryFunc
	requestIgnorer apmhttp.RequestIgnorerFunc
}

// Option sets options for tracing.
type Option func(*options)

// WithTracer returns an Option which sets t as the tracer
// to use for tracing server requests.
func WithTracer(t *apm.Tracer) Option {
	if t == nil {
		panic("t == nil")
	}
	return func(o *options) {
		o.tracer = t
	}
}

// WithRecovery returns an Option which sets r as the recovery
// function to use for tracing server requests.
func WithRecovery(r apmhttp.RecoveryFunc) Option {
	if r == nil {
		panic("r == nil")
	}
	return func(o *options) {
		o.recovery = r
	}
}

// todo I need this method
// WithRequestIgnorer returns a Option which sets r as the
// function to use to determine whether or not a request should
// be ignored. If r is nil, all requests will be reported.
//func WithRequestIgnorer(r apmhttp.RequestIgnorerFunc) Option {
//	if r == nil {
//		r = apmhttp.IgnoreNone
//	}
//	return func(o *options) {
//		o.requestIgnorer = r
//	}
//}

// Router wraps an httprouter.Router, instrumenting all added routes
// except static content served with ServeFiles.
//type Router struct {
//	*httprouter.Router
//	opts []Option
//}

// NewTracer returns a new Router which will instrument all added routes
// except static content served with ServeFiles.
//
// Router.NotFound and Router.MethodNotAllowed will be set, and will
// report transactions with the name "<METHOD> unknown route".
//func New(o ...Option) *Router {
//	router := httprouter.New()
//	router.NotFound = WrapNotFoundHandler(router.NotFound, o...)
//	router.MethodNotAllowed = WrapMethodNotAllowedHandler(router.MethodNotAllowed, o...)
//	return &Router{
//		Router: router,
//		opts:   o,
//	}
////}
//
//// DELETE calls r.Router.DELETE with a wrapped handler.
//func (r *Router) DELETE(path string, handle httprouter.Handle) {
//	r.Router.DELETE(path, Wrap(handle, path, r.opts...))
//}
//
//// GET calls r.Router.GET with a wrapped handler.
//func (r *Router) GET(path string, handle httprouter.Handle) {
//	r.Router.GET(path, Wrap(handle, path, r.opts...))
//}
//
//// HEAD calls r.Router.HEAD with a wrapped handler.
//func (r *Router) HEAD(path string, handle httprouter.Handle) {
//	r.Router.HEAD(path, Wrap(handle, path, r.opts...))
//}
//
//// Handle calls r.Router.Handle with a wrapped handler.
//func (r *Router) Handle(method, path string, handle httprouter.Handle) {
//	r.Router.Handle(method, path, Wrap(handle, path, r.opts...))
//}
//
//// HandlerFunc is equivalent to r.Router.HandlerFunc, but traces requests.
//func (r *Router) HandlerFunc(method, path string, handler http.HandlerFunc) {
//	r.Handler(method, path, handler)
//}
//
//// Handler is equivalent to r.Router.Handler, but traces requests.
//func (r *Router) Handler(method, path string, handler http.Handler) {
//	r.Handle(method, path, func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
//		ctx := req.Context()
//		ctx = context.WithValue(ctx, httprouter.ParamsKey, p)
//		req = req.WithContext(ctx)
//		handler.ServeHTTP(w, req)
//	})
//}
//
//// OPTIONS is equivalent to r.Router.OPTIONS, but traces requests.
//func (r *Router) OPTIONS(path string, handle httprouter.Handle) {
//	r.Router.OPTIONS(path, Wrap(handle, path, r.opts...))
//}
//
//// PATCH is equivalent to r.Router.PATCH, but traces requests.
//func (r *Router) PATCH(path string, handle httprouter.Handle) {
//	r.Router.PATCH(path, Wrap(handle, path, r.opts...))
//}
//
//// POST is equivalent to r.Router.POST, but traces requests.
//func (r *Router) POST(path string, handle httprouter.Handle) {
//	r.Router.POST(path, Wrap(handle, path, r.opts...))
//}
//
//// PUT is equivalent to r.Router.PUT, but traces requests.
//func (r *Router) PUT(path string, handle httprouter.Handle) {
//	r.Router.PUT(path, Wrap(handle, path, r.opts...))
//}

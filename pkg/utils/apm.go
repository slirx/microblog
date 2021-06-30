package utils

import (
	"net/http"

	"github.com/pkg/errors"
	"go.elastic.co/apm"
	"go.elastic.co/apm/module/apmhttp"

	"gitlab.com/slirx/newproj/pkg/api"
	"gitlab.com/slirx/newproj/pkg/logger"
)

func NewRecoveryFunc(zapLogger logger.Logger, responseBuilder api.ResponseBuilder) func(
	w http.ResponseWriter,
	req *http.Request,
	resp *apmhttp.Response,
	body *apm.BodyCapturer,
	tx *apm.Transaction,
	recovered interface{},
) {
	return func(
		w http.ResponseWriter,
		req *http.Request,
		resp *apmhttp.Response,
		body *apm.BodyCapturer,
		tx *apm.Transaction,
		recovered interface{},
	) {
		if err, ok := recovered.(error); ok {
			zapLogger.Error(err)
		} else if err, ok := recovered.(string); ok {
			zapLogger.Error(errors.New(err))
		} // todo handle else too

		responseBuilder.ErrorResponse(req.Context(), w, api.InternalError)
	}
}

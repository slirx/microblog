package media

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/pkg/errors"
	"go.elastic.co/apm/module/apmzap"

	"gitlab.com/slirx/newproj/pkg/api"
	"gitlab.com/slirx/newproj/pkg/http/jwtmiddleware"
	"gitlab.com/slirx/newproj/pkg/logger"
)

const (
	MaxFileSize = 10 << 20
)

var (
	ErrFileTooLarge    = errors.New("file is too large")
	ErrInvalidFileType = errors.New("file type is forbidden")
	ErrInvalidService  = errors.New("invalid service")
	ErrInvalidItemID   = errors.New("invalid item id")
)

var (
	mediaPath       = "web/images/media"
	mediaPublicPath = "media/image"
	allowedServices = map[string]struct{}{
		"user": {},
	}
	notFoundImages = map[string]string{
		"user": "user/not-found.png",
	}
)

type Handler interface {
	// Images returns list of images. It's used to fetch images by their IDs.
	Images(w http.ResponseWriter, r *http.Request)
	// UploadImage uploads one image and returns URL to it.
	UploadImage(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	Logger          logger.Logger
	ResponseBuilder api.ResponseBuilder
}

func (h handler) Images(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	serviceName := r.URL.Query().Get("service")
	idsString := strings.Split(r.URL.Query().Get("ids"), ",")

	var s string
	var publicName string
	var fileName string
	var err error

	images := make(map[string]string)
	host := "http://microblog.local:8080/" // todo move host to config

	for _, s = range idsString {
		publicName = path.Join(serviceName, fmt.Sprintf("%s.png", s))

		fileName = path.Join(mediaPath, publicName)
		if _, err = os.Stat(fileName); err != nil {
			images[s] = host + path.Join(mediaPublicPath, notFoundImages[serviceName])
			continue
		}

		images[s] = host + path.Join(mediaPublicPath, publicName) // todo uncomment
	}

	response := ImagesResponse{
		Images: images,
	}

	h.ResponseBuilder.DataResponse(ctx, w, response)
}

func (h handler) UploadImage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	uid, err := jwtmiddleware.UID(ctx)
	if err != nil {
		h.Logger.Error(err, apmzap.TraceContext(ctx)...)
		h.ResponseBuilder.ErrorResponse(ctx, w, api.NewRequestError(err))
		return
	}

	err = r.ParseMultipartForm(10 << 20)
	if err != nil {
		h.Logger.Error(err, apmzap.TraceContext(ctx)...)
		h.ResponseBuilder.ErrorResponse(ctx, w, api.InternalError)
		return
	}

	serviceName := r.FormValue("service")

	var itemID int
	if serviceName == "user" {
		itemID = uid
	}

	if itemID == 0 {
		h.Logger.Error(ErrInvalidItemID, apmzap.TraceContext(ctx)...)
		h.ResponseBuilder.ErrorResponse(ctx, w, api.NewRequestError(ErrInvalidItemID))
		return
	}

	if _, ok := allowedServices[serviceName]; !ok {
		h.Logger.Error(ErrInvalidService, apmzap.TraceContext(ctx)...)
		h.ResponseBuilder.ErrorResponse(ctx, w, api.NewRequestError(ErrInvalidService))
		return
	}

	file, fileHeader, err := r.FormFile("image")
	if err != nil {
		h.Logger.Error(err, apmzap.TraceContext(ctx)...)
		h.ResponseBuilder.ErrorResponse(ctx, w, api.RequestError)
		return
	}

	defer file.Close()

	if fileHeader.Size > MaxFileSize {
		h.Logger.Error(ErrFileTooLarge, apmzap.TraceContext(ctx)...)
		h.ResponseBuilder.ErrorResponse(ctx, w, api.NewRequestError(ErrFileTooLarge))
		return
	}

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		h.Logger.Error(err, apmzap.TraceContext(ctx)...)
		h.ResponseBuilder.ErrorResponse(ctx, w, api.InternalError)
		return
	}

	fileType := http.DetectContentType(fileBytes[:512])

	switch fileType {
	case "image/jpeg", "image/jpg", "image/png":
	default:
		h.Logger.Error(ErrInvalidFileType, apmzap.TraceContext(ctx)...)
		h.ResponseBuilder.ErrorResponse(ctx, w, api.NewRequestError(ErrInvalidFileType))
		return
	}

	fileName := fmt.Sprintf("%d.png", itemID)

	tempFile, err := os.Create(path.Join(mediaPath, serviceName, fileName))
	if err != nil {
		h.Logger.Error(err, apmzap.TraceContext(ctx)...)
		h.ResponseBuilder.ErrorResponse(ctx, w, api.InternalError)
		return
	}

	_, err = tempFile.Write(fileBytes)
	if err != nil {
		_ = tempFile.Close()
		h.Logger.Error(err, apmzap.TraceContext(ctx)...)
		h.ResponseBuilder.ErrorResponse(ctx, w, api.InternalError)
		return
	}

	err = tempFile.Close()
	if err != nil {
		h.Logger.Error(err, apmzap.TraceContext(ctx)...)
		h.ResponseBuilder.ErrorResponse(ctx, w, api.InternalError)
		return
	}

	response := UploadImageResponse{
		URL: path.Join(mediaPath, serviceName, fileName),
	}

	h.ResponseBuilder.DataResponse(ctx, w, response)
}

// NewHandler returns instance of implemented Handler interface.
func NewHandler(l logger.Logger, rb api.ResponseBuilder) Handler {
	return handler{Logger: l, ResponseBuilder: rb}
}

package jwtmiddleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"go.elastic.co/apm/module/apmzap"

	"gitlab.com/slirx/newproj/pkg/api"
	"gitlab.com/slirx/newproj/pkg/logger"
)

const ContextKeyUserID = "uid"
const ContextKeyLogin = "login"

func Wrap(h http.HandlerFunc, rb api.ResponseBuilder, l logger.Logger, secret []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")

		token, err := jwt.Parse(header, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, api.NewAccessError(fmt.Errorf("invalid JWT"))
			}

			return secret, nil
		})
		if err != nil {
			vErr, ok := err.(*jwt.ValidationError)
			if ok && vErr.Errors&jwt.ValidationErrorExpired == jwt.ValidationErrorExpired {
				err = api.NewAccessError(errors.New("token is expired"))
				l.Error(err, apmzap.TraceContext(r.Context())...)
				rb.ErrorResponse(r.Context(), w, err)
				return
			}

			l.Error(err, apmzap.TraceContext(r.Context())...)
			rb.ErrorResponse(r.Context(), w, err)
			return
		}

		if token.Valid {
			var ok bool
			var uid float64
			var uidTmp interface{}
			claims := make(map[string]interface{})

			claims, ok = token.Claims.(jwt.MapClaims)
			if !ok {
				l.Error(err, apmzap.TraceContext(r.Context())...)
				rb.ErrorResponse(r.Context(), w, err)
				return
			}

			uidTmp, ok = claims["uid"]
			if !ok {
				err = fmt.Errorf("invalid JWT")
				l.Error(err, apmzap.TraceContext(r.Context())...)
				rb.ErrorResponse(r.Context(), w, api.NewAccessError(err))
				return
			}

			uid, ok = uidTmp.(float64)
			if !ok {
				err = fmt.Errorf("invalid JWT")
				l.Error(err, apmzap.TraceContext(r.Context())...)
				rb.ErrorResponse(r.Context(), w, api.NewAccessError(err))
				return
			}

			ctx := r.Context()
			ctx = context.WithValue(ctx, ContextKeyUserID, int(uid))
			r = r.WithContext(ctx)

			h(w, r)

			return
		}

		err = fmt.Errorf("invalid JWT")
		l.Error(err, apmzap.TraceContext(r.Context())...)
		rb.ErrorResponse(r.Context(), w, api.NewAccessError(err))
	}
}

func WrapInternal(h http.HandlerFunc, rb api.ResponseBuilder, l logger.Logger, secrets map[string][]byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		claims := make(map[string]interface{})

		var login string

		token, err := jwt.Parse(header, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, api.NewAccessError(fmt.Errorf("invalid JWT"))
			}

			var ok bool
			var tmpLogin interface{}

			claims, ok = token.Claims.(jwt.MapClaims)
			if !ok {
				return nil, api.NewAccessError(fmt.Errorf("invalid JWT (invalid claims)"))
			}

			tmpLogin, ok = claims["login"]
			if !ok {
				return nil, api.NewAccessError(fmt.Errorf("invalid JWT (no login)"))
			}

			login, ok = tmpLogin.(string)
			if !ok {
				return nil, api.NewAccessError(fmt.Errorf("invalid JWT (invalid login)"))
			}

			secret, ok := secrets[login]
			if !ok {
				return nil, api.NewAccessError(fmt.Errorf("invalid JWT (no login)"))
			}

			return secret, nil
		})
		if err != nil {
			vErr, ok := err.(*jwt.ValidationError)
			if ok && vErr.Errors&jwt.ValidationErrorExpired == jwt.ValidationErrorExpired {
				err = api.NewAccessError(errors.New("token is expired"))
				l.Error(err, apmzap.TraceContext(r.Context())...)
				rb.ErrorResponse(r.Context(), w, err)
				return
			}

			l.Error(err, apmzap.TraceContext(r.Context())...)
			rb.ErrorResponse(r.Context(), w, err)
			return
		}

		if token.Valid {
			ctx := r.Context()
			ctx = context.WithValue(ctx, ContextKeyLogin, login)
			r = r.WithContext(ctx)

			h(w, r)

			return
		}

		err = fmt.Errorf("invalid JWT")
		l.Error(err, apmzap.TraceContext(r.Context())...)
		rb.ErrorResponse(r.Context(), w, api.NewAccessError(err))
	}
}

// WrapSetUIDContext sets uid to context in case used is authorized.
// It doesn't report an error in case user in not authorized. It should be used when authorization is optional.
//func WrapSetUIDContext(h http.HandlerFunc, rb api.ResponseBuilder, l logger.Logger, secret []byte) http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		header := r.Header.Get("Authorization")
//		if header == "" {
//			h(w, r)
//			return
//		}
//
//		token, err := jwt.Parse(header, func(token *jwt.Token) (interface{}, error) {
//			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
//				return nil, api.NewAccessError(fmt.Errorf("invalid JWT"))
//			}
//
//			return secret, nil
//		})
//		if err != nil {
//			vErr, ok := err.(*jwt.ValidationError)
//			if ok && vErr.Errors&jwt.ValidationErrorExpired == jwt.ValidationErrorExpired {
//				err = api.NewAccessError(errors.New("token is expired"))
//				l.Error(err, apmzap.TraceContext(r.Context())...)
//				rb.ErrorResponse(r.Context(), w, err)
//				return
//			}
//
//			l.Error(err, apmzap.TraceContext(r.Context())...)
//			rb.ErrorResponse(r.Context(), w, err)
//			return
//		}
//
//		if token.Valid {
//			var ok bool
//			var uid float64
//			var uidTmp interface{}
//			claims := make(map[string]interface{})
//
//			claims, ok = token.Claims.(jwt.MapClaims)
//			if !ok {
//				l.Error(err, apmzap.TraceContext(r.Context())...)
//				rb.ErrorResponse(r.Context(), w, err)
//				return
//			}
//
//			uidTmp, ok = claims["uid"]
//			if !ok {
//				err = fmt.Errorf("invalid JWT")
//				l.Error(err, apmzap.TraceContext(r.Context())...)
//				rb.ErrorResponse(r.Context(), w, api.NewAccessError(err))
//				return
//			}
//
//			uid, ok = uidTmp.(float64)
//			if !ok {
//				err = fmt.Errorf("invalid JWT")
//				l.Error(err, apmzap.TraceContext(r.Context())...)
//				rb.ErrorResponse(r.Context(), w, api.NewAccessError(err))
//				return
//			}
//
//			ctx := r.Context()
//			ctx = context.WithValue(ctx, ContextKeyUserID, uint(uid))
//			r = r.WithContext(ctx)
//
//			h(w, r)
//
//			return
//		}
//
//		err = fmt.Errorf("invalid JWT")
//		l.Error(err, apmzap.TraceContext(r.Context())...)
//		rb.ErrorResponse(r.Context(), w, api.NewAccessError(err))
//	}
//}

// UID returns user id from context.
func UID(ctx context.Context) (int, error) {
	var uid int
	var ok bool

	val := ctx.Value(ContextKeyUserID)
	if val == nil {
		return 0, errors.WithStack(fmt.Errorf("%s is not in context", ContextKeyUserID))
	}

	uid, ok = val.(int)
	if !ok {
		return 0, errors.WithStack(fmt.Errorf("invalid %s type", ContextKeyUserID))
	}

	if uid <= 0 {
		return 0, errors.WithStack(fmt.Errorf("invalid %s value", ContextKeyUserID))
	}

	return uid, nil
}

// OptionalUID returns user id from context, without any errors in case there is no UID in the context.
//func OptionalUID(ctx context.Context) uint {
//	var uid uint
//	var ok bool
//
//	val := ctx.Value(ContextKeyUserID)
//	if val == nil {
//		return 0
//	}
//
//	uid, ok = val.(uint)
//	if !ok {
//		return 0
//	}
//
//	if uid <= 0 {
//		return 0
//	}
//
//	return uid
//}

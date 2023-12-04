package api

import (
	"errors"
	"net/http"
	"sync/atomic"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type CodeResp struct {
	Code Code `json:"code" example:"101001"` // e.CodeSuccess
}

type DataResp struct {
	CodeResp
	Data any `json:"data,omitempty"`
}

type MessageResp struct {
	CodeResp
	Message string `json:"message"`
}

type ListDataResp struct {
	DataResp

	Total int64 `json:"total,omitempty"`
}

type ErrResp struct {
	MessageResp
	Extra string `json:"extra,omitempty"`
}

type resp struct {
	c *gin.Context

	done   atomic.Bool
	logger *zerolog.Logger
}

func newResp(c *gin.Context, log zerolog.Logger) *resp {
	return &resp{
		c:      c,
		logger: &log,
	}
}

func (r *resp) OK() {
	r.send(http.StatusOK, CodeResp{Code: CodeOK})
}

func (r *resp) Data(data any) {
	r.send(http.StatusOK, DataResp{
		CodeResp: CodeResp{Code: CodeOK},
		Data:     data,
	})
}

func (r *resp) DataWithTotal(data []any, total uint64) {
	r.send(http.StatusOK, ListDataResp{
		DataResp: DataResp{
			CodeResp: CodeResp{Code: CodeOK},
			Data:     data,
		},
		Total: int64(total),
	})
}

func (r *resp) Err(err error) {
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		r.sendErr(http.StatusBadRequest, CodeInvalidParam, err)
		return
	}

	r.sendErr(http.StatusInternalServerError, CodeUnknown, err)
}

func (r *resp) NotFound() {
	r.sendErr(http.StatusNotFound, CodeNotFound, nil)
}

func (r *resp) Forbidden(err error) {
	r.sendErr(http.StatusForbidden, CodeForbidden, err)
}

func (r *resp) Unauthorized(err error) {
	r.sendErr(http.StatusUnauthorized, CodeInvalidParam, err)
}

func (r *resp) ExpiredToken(err error) {
	r.sendErr(http.StatusUnauthorized, CodeTokenExpired, err)
}

func (r *resp) dbErr(err error) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		r.sendErr(http.StatusBadRequest, CodeInvalidParam, err)
		return
	}

	r.sendErr(http.StatusInternalServerError, CodeUnknown, err)
}

func (r *resp) InvalidParam(err error) {
	r.sendErr(http.StatusBadRequest, CodeInvalidParam, err)
}

func (r *resp) Unknown(err error) {
	r.sendErr(http.StatusInternalServerError, CodeUnknown, err)
}

func (r *resp) HandleErr(err error) {
	var (
		storeErr *StoreErr
		paramErr *ParamErr
		httpErr  *HTTPErr
	)
	switch {
	case errors.As(err, &storeErr):
		r.dbErr(storeErr)
	case errors.As(err, &paramErr):
		r.sendErr(http.StatusBadRequest, CodeInvalidParam, paramErr)
	case errors.As(err, &httpErr):
		r.sendErr(httpErr.HTTPCode, httpErr.Code, httpErr)
	default:
		r.Unknown(err)
	}
}

//
// Finisher
//

func (r *resp) sendErr(httpCode int, code Code, err error) {
	if err != nil && httpCode >= http.StatusBadRequest {
		r.c.Errors = append(r.c.Errors, &gin.Error{
			Err:  err,
			Type: gin.ErrorTypePrivate,
		})

		r.logger.
			Debug().
			CallerSkipFrame(2). // skip to gin handler layer
			Caller().
			Err(err).
			Msg("response error")
	}

	r.send(httpCode, ErrResp{
		MessageResp: MessageResp{
			CodeResp: CodeResp{Code: code},
			Message:  Message(code),
		},
	})
}

func (r *resp) send(httpCode int, payload any) {
	if r.done.Swap(true) {
		r.logger.Warn().
			Int("http_code", httpCode).
			Interface("payload", payload).
			Msg("response already done, skip")
		return
	}

	if payload == nil {
		if httpCode >= http.StatusBadRequest {
			r.c.AbortWithStatus(httpCode)
			return
		}

		r.c.AbortWithStatus(httpCode)
		return
	}

	r.c.AbortWithStatusJSON(httpCode, payload)
}

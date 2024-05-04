package router

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	se "github.com/usk81/ashihara/shared/domain/errors"
	"github.com/usk81/ashihara/shared/interface/transport/http/render"

	"github.com/usk81/ashihara/services/holidays/core/domain/errors"
)

// BindFromJSON ...
func BindFromJSON(body io.ReadCloser, v interface{}) (err error) {
	if err := json.NewDecoder(body).Decode(v); err != nil {
		return errors.NewCause(err, errors.CaseBadRequest)
	}
	return nil
}

// RenderJSON ...
func RenderJSON(w http.ResponseWriter, s int, v interface{}) {
	if err := render.JSON(w, s, v); err != nil {
		err := errors.NewCause(fmt.Errorf("fail to encode response : %w", err), errors.CaseBackendError)
		render.JSON(w, http.StatusInternalServerError, err) // nolint: errcheck
		return
	}
}

// RenderError ...
func RenderError(w http.ResponseWriter, err error) {
	ec, ok := err.(*se.Cause)
	if !ok {
		ec = errors.NewCause(err, errors.CaseBackendError).(*se.Cause)
	}
	render.JSON(w, ec.Code(), ec) // nolint: errcheck
}

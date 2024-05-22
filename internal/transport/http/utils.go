package http

import (
	"encoding/json"
	"errors"
	"mini-godis/internal/errs"
	"mini-godis/internal/transport/http/response"
	"net/http"
	"strconv"

	"go.uber.org/zap"
)

func (t transport) handleHTTPError(w http.ResponseWriter, err error, method string, status int) {
	w.Header().Set("Content-Type", "application/json")

	logError := err
	err = errs.Handle(err)
	if errors.Is(err, errs.NoRowsInResult) ||
		errors.Is(err, errs.NoReference) {
		status = http.StatusNotFound
	}

	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(response.Error{Error: err.Error()}); err != nil {
		t.l.Error(method, zap.Error(logError))
	}
	t.l.Error(method, zap.Error(logError))
}

func (t transport) handleHTTPOk(w http.ResponseWriter, resp interface{}, method string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if resp != nil {
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			t.l.Error(method, zap.Error(err))
		}
	}
	t.l.Info(method)
}

func IsJSON(str string) bool {
	var js json.RawMessage
	return json.Unmarshal([]byte(str), &js) == nil
}

func getIDFromURL(r *http.Request) (int, error) {
	strID := r.PathValue("id")
	if strID == "" {
		return -1, errs.NilID
	}
	ID, err := strconv.Atoi(strID)
	if err != nil {
		return -1, errs.WrongID
	}
	return ID, nil
}

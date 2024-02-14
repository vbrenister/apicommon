package apicommon

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type ServerConfig struct{}

type JsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func (app *ServerConfig) ReadJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytesSize := 1048576

	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytesSize))

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(data)

	if err != nil {
		return err
	}

	err = decoder.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("request body must only have a single JSON value")
	}

	return nil
}

func (app *ServerConfig) WriteJSON(w http.ResponseWriter, status int, data any, headers ...http.Header) error {
	out, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for k, v := range headers[0] {
			w.Header()[k] = v
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(out)

	return err
}

func (app *ServerConfig) ErroJSON(w http.ResponseWriter, err error, status ...int) error {
	statusCode := http.StatusBadRequest

	if len(status) > 0 {
		statusCode = status[0]
	}

	return app.WriteJSON(w, statusCode, JsonResponse{
		Error:   true,
		Message: err.Error(),
	})
}

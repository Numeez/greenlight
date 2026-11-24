package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/Numeez/greenlight/internal/validator"
	"github.com/julienschmidt/httprouter"
)

type envelope map[string]interface{}

func (app *application) readIDParam(r *http.Request) (int64, error) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}

	return id, nil

}
func (app *application) readJson(w http.ResponseWriter, r *http.Request, dst any) error {

	max_bytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(max_bytes))

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err := dec.Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly formatted json at character: %d", syntaxError.Offset)
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect json type for field: %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect json type at character %d", unmarshalTypeError.Offset)
		case errors.As(err, &invalidUnmarshalError):
			panic(err)
		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body contains unknown key: %s", fieldName)
		case err.Error() == "http request body too large":
			return fmt.Errorf("body must not be larger than %d bytes", max_bytes)
		default:
			return err
		}
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only contain a single JSON value")
	}

	return nil

}
func (app *application) writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}
	js = append(js, '\n')

	for k, v := range headers {
		w.Header()[k] = v
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(js))

	return nil

}

func (app *application) readString(qs url.Values, key string, defaultValue string) string {
	value := qs.Get(key)
	if value == "" {
		value = defaultValue
	}
	return value

}
func (app *application) readCSV(qs url.Values, key string, defaultValue []string) []string {
	value := qs.Get(key)
	if value == "" {
		return defaultValue
	}
	return strings.Split(value, ",")

}

func (app *application) readInt(v *validator.Validator, qs url.Values, key string, defaultValue int) int {
	value := qs.Get(key)
	if value == "" {
		return defaultValue
	}
	val, err := strconv.Atoi(value)
	if err != nil {
		v.AddError(key, "must be an integer value")
	}
	return val

}

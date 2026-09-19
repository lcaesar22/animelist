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

	"github.com/go-chi/chi/v5"
	"github.com/konaha-gakure/otaku/internal/validator"
)

type envelope map[string]any

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.errorResponse(w, r, http.StatusBadRequest, err.Error())
}

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, def any) error {
	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(def)
	if err != nil {
		// if there is an error during decoding
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError

		var maxBytesError *http.MaxBytesError

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly formated json at  character %d", syntaxError.Offset)
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("badly formated json")
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect json type for field %q", unmarshalTypeError.Offset)
			}
			return fmt.Errorf("body contains incorrect json type at character %d", unmarshalTypeError.Offset)
		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body contains unknown field %q", fieldName)
		case errors.As(err, &maxBytesError):
			return fmt.Errorf("body must not be larger than %d bytes", maxBytesError.Limit)
		case errors.As(err, &invalidUnmarshalError):
			panic(err)
		default:
			return err
		}

	}

	err = dec.Decode(&def)
	if !errors.Is(err, io.EOF) {
		return errors.New("body must contain a single json object")
	}
	return nil
}

// converts params to integer
func (app *application) readIDParam(r *http.Request) (int, error) {
	params := chi.URLParam(r, "id")
	id, err := strconv.Atoi(params)
	if err != nil || id < 1 {
		return 0, errors.New("id must be an integer")
	}
	return id, nil
}

// json write helper method
func (app *application) writeJSON(w http.ResponseWriter, status int, data any, headers http.Header) error {
	js, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	js = append(js, '\n')

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)
	return nil
}

/*
 * helpers for filtering
 *readString() returns a value from a query string if a match is found
 *readCSV() reads and splits a query string into two
 *readInt() reads a string value from the query and converts it to an integer1`q
 */

func (app *application) readString(qs url.Values, key string, def string) string {
	s := qs.Get(key)
	if s == "" {
		return def
	}
	return s
}

func (app *application) readCSV(qs url.Values, key string, def []string) []string {
	csv := qs.Get(key)
	if csv == "" {
		return def
	}
	return strings.Split(csv, ",")
}

func (app *application) readInt(qs url.Values, key string, def int, v *validator.Validator) int {
	s := qs.Get(key)
	if s == "" {
		return def
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		v.AddError(key, "must be an integer")
		return def
	}
	return i
}

// The background() helper accepts an arbitrary function as a parameter.
func (app *application) background(fn func()) {
	app.wg.Add(1)
	// Launch a background goroutine.
	go func() {
		defer app.wg.Done()
		// Recover any panic.
		defer func() {
			if err := recover(); err != nil {
				app.logger.Error(fmt.Sprintf("%v", err))
			}
		}()
		// Execute the arbitrary function that we passed as the parameter.
		fn()
	}()
}

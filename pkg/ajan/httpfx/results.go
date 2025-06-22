package httpfx

import (
	"encoding/json"
	"net/http"

	"github.com/eser/aya.is-services/pkg/ajan/results"
)

var (
	okResult = results.Define( //nolint:gochecknoglobals
		results.ResultKindSuccess,
		"OK",
		"OK",
	)
	errResult = results.Define( //nolint:gochecknoglobals
		results.ResultKindError,
		"ERR",
		"Error",
	)
)

// Result Options.
type ResultOption func(*Result)

func WithBody(body []byte) ResultOption {
	return func(result *Result) {
		result.InnerBody = body
	}
}

func WithPlainText(body string) ResultOption {
	return func(result *Result) {
		result.InnerBody = []byte(body)
	}
}

func WithJSON(body any) ResultOption {
	return func(result *Result) {
		encoded, err := json.Marshal(body)
		if err != nil {
			result.InnerBody = []byte("Failed to encode JSON")
			result.InnerStatusCode = http.StatusInternalServerError

			return
		}

		result.InnerBody = encoded
	}
}

// Results With Options.
type Results struct{}

func (r *Results) Ok(options ...ResultOption) Result {
	result := Result{
		Result: okResult.New(),

		InnerStatusCode:    http.StatusNoContent,
		InnerRedirectToURI: "",
		InnerBody:          make([]byte, 0),
	}

	for _, option := range options {
		option(&result)
	}

	return result
}

func (r *Results) Accepted(options ...ResultOption) Result {
	result := Result{
		Result: okResult.New(),

		InnerStatusCode:    http.StatusAccepted,
		InnerRedirectToURI: "",
		InnerBody:          make([]byte, 0),
	}

	for _, option := range options {
		option(&result)
	}

	return result
}

func (r *Results) NotFound(options ...ResultOption) Result {
	result := Result{
		Result: okResult.New(),

		InnerStatusCode:    http.StatusNotFound,
		InnerRedirectToURI: "",
		InnerBody:          []byte("Not Found"),
	}

	for _, option := range options {
		option(&result)
	}

	return result
}

func (r *Results) Unauthorized(options ...ResultOption) Result {
	result := Result{
		Result: errResult.New(),

		InnerStatusCode:    http.StatusUnauthorized,
		InnerRedirectToURI: "",
		InnerBody:          make([]byte, 0),
	}

	for _, option := range options {
		option(&result)
	}

	return result
}

func (r *Results) BadRequest(options ...ResultOption) Result {
	result := Result{
		Result: errResult.New(),

		InnerStatusCode:    http.StatusBadRequest,
		InnerRedirectToURI: "",
		InnerBody:          []byte("Bad Request"),
	}

	for _, option := range options {
		option(&result)
	}

	return result
}

func (r *Results) Error(statusCode int, options ...ResultOption) Result {
	result := Result{
		Result: errResult.New(),

		InnerStatusCode:    statusCode,
		InnerRedirectToURI: "",
		InnerBody:          make([]byte, 0),
	}

	for _, option := range options {
		option(&result)
	}

	return result
}

// Results Without Options.
func (r *Results) Bytes(body []byte) Result {
	return Result{
		Result: okResult.New(),

		InnerStatusCode:    http.StatusOK,
		InnerRedirectToURI: "",
		InnerBody:          body,
	}
}

func (r *Results) PlainText(body []byte) Result {
	return Result{
		Result: okResult.New(),

		InnerStatusCode:    http.StatusOK,
		InnerRedirectToURI: "",
		InnerBody:          body,
	}
}

func (r *Results) JSON(body any) Result {
	encoded, err := json.Marshal(body)
	if err != nil {
		// TODO(@eser) Log error
		return Result{
			Result: errResult.New(),

			InnerStatusCode:    http.StatusInternalServerError,
			InnerRedirectToURI: "",
			InnerBody:          []byte("Failed to encode JSON"),
		}
	}

	return Result{
		Result: okResult.New(),

		InnerStatusCode:    http.StatusOK,
		InnerRedirectToURI: "",
		InnerBody:          encoded,
	}
}

func (r *Results) Redirect(uri string) Result {
	return Result{
		Result: okResult.New(),

		InnerStatusCode:    http.StatusTemporaryRedirect,
		InnerRedirectToURI: uri,
		InnerBody:          make([]byte, 0),
	}
}

func (r *Results) Abort() Result {
	// TODO(@eser) implement this
	return Result{
		Result: errResult.New(),

		InnerStatusCode:    http.StatusNotImplemented,
		InnerRedirectToURI: "",
		InnerBody:          []byte("Not Implemented"),
	}
}

package shared

import (
	"encoding/json"
	"log"
	"net/http"
)

type Response struct {
	Data       interface{}    `json:"data"`
	StatusCode int            `json:"status_code"`
	Message    string         `json:"message"`
	Cookies    []*http.Cookie `json:"-"`
}

type AppHandler func(*http.Request) (*Response, *CustomError)

func (h AppHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	resp, cerr := h(r)
	if cerr != nil {
		errResponse := writeErrorResponse(cerr, w)
		if _, err := w.Write(errResponse); err != nil {
			log.Printf("error writing error response: %v", err)
		}
		return
	}

	// Set cookies if any are provided
	if resp.Cookies != nil {
		for _, cookie := range resp.Cookies {
			http.SetCookie(w, cookie)
		}
	}

	httpResponse := make(map[string]interface{})
	httpResponse["data"] = resp.Data
	httpResponse["status"] = resp.StatusCode

	response, err := json.Marshal(httpResponse)
	if err != nil {
		log.Printf("failed to marshal response: %v", err)
		errResponse := writeErrorResponse(NewCustomErr(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), err.Error()), w)
		if _, err := w.Write(errResponse); err != nil {
			log.Printf("failed to marshal response: %v", err)
		}
		return
	}
	w.WriteHeader(resp.StatusCode)
	if _, err := w.Write(response); err != nil {
		log.Printf("failed to write response: %v", err)
		return
	}
}

func writeErrorResponse(err *CustomError, w http.ResponseWriter) []byte {
	// Check if the header has already been written
	if w.Header().Get("Content-Type") == "" {
		w.WriteHeader(err.StatusCode) // Only set the header if it hasn't been set
	}

	errResponse := make(map[string]interface{})
	errResponse["message"] = err.Message
	errResponse["error_message"] = err.ErrorMessage
	errResponse["code"] = err.StatusCode

	response, er := json.Marshal(errResponse)
	if er != nil {
		return []byte(http.StatusText(http.StatusInternalServerError))
	}
	return response
}

type CustomError struct {
	StatusCode   int
	Message      string
	ErrorMessage string
}

func NewCustomErr(statusCode int, errMsg, msg string) *CustomError {
	return &CustomError{
		StatusCode:   statusCode,
		ErrorMessage: errMsg,
		Message:      msg,
	}
}

func (ce *CustomError) Error() string {
	return ce.ErrorMessage
}

func BadRequest(msg string) *CustomError {
	return NewCustomErr(http.StatusBadRequest, "bad_request", msg)
}

func InternalError(msg string) *CustomError {
	return NewCustomErr(http.StatusInternalServerError, "internal_server_error", msg)
}

func Unauthorized(msg string) *CustomError {
	return NewCustomErr(http.StatusUnauthorized, "unauthorized", msg)
}

func JSON(_ http.ResponseWriter, statusCode int, data interface{}) (*Response, *CustomError) {
	return &Response{
		Data:       data,
		StatusCode: statusCode,
	}, nil
}

package proxy_error

import (
	"fmt"
	"net/http"
)

const (
	BadRequest = iota
	InternalError
)

type ProxyError int

func (e ProxyError) String() string {
	switch e {
	case BadRequest:
		return "BAD_REQUEST"
	case InternalError:
		return "INTERNAL_ERROR"
	}

	panic(fmt.Sprintf("ProxyError unhandled: %d", e))
}

func (e ProxyError) HttpStatusCode() int {
	switch e {
	case BadRequest:
		return http.StatusBadRequest
	case InternalError:
		return http.StatusInternalServerError
	}

	panic(fmt.Sprintf("ProxyError unhandled: %d", e))
}

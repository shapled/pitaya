package pitaya

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type Request interface {
	Context() echo.Context
	Request()
}

type Response interface {
	Context() echo.Context
	Response()
}

type BaseRequest struct {
	ctx echo.Context
}

func (br *BaseRequest) Context() echo.Context {
	return br.ctx
}

func (br *BaseRequest) Request() {

}

type BaseResponse struct {
	ctx echo.Context
}

func (br *BaseResponse) Context() echo.Context {
	return br.ctx
}

func (br *BaseResponse) Response() {

}

type Validator struct {
	validator *validator.Validate
}

func (v *Validator) Validate(i interface{}) error {
	if err := v.validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return nil
}

type Server struct {
	*echo.Echo
	badRequestStatus int
	responseWrapper  func(Response) interface{}
	errorWrapper     func(error) interface{}
}

func NewServer() *Server {
	return NewServerWithArgs(0, nil, nil)
}

func NewServerWithArgs(badRequestStatus int, responseWrapper func(Response) interface{}, errorWrapper func(error) interface{}) *Server {
	if badRequestStatus == 0 {
		badRequestStatus = http.StatusBadRequest
	}
	if responseWrapper == nil {
		responseWrapper = func(x Response) interface{} { return x }
	}
	if errorWrapper == nil {
		errorWrapper = func(x error) interface{} { return x }
	}
	e := echo.New()
	e.Validator = &Validator{validator: validator.New()}
	return &Server{
		Echo:             e,
		badRequestStatus: badRequestStatus,
		responseWrapper:  responseWrapper,
		errorWrapper:     errorWrapper,
	}
}

func (server *Server) HandlerWrapper(handler func(Request) (Response, error), req Request) func(echo.Context) error {
	return func(ctx echo.Context) error {
		if req != nil {
			if err := ctx.Bind(req); err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, server.errorWrapper(err))
			}
		}
		resp, err := handler(req)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, server.errorWrapper(err))
		}
		return ctx.JSON(http.StatusOK, server.responseWrapper(resp))
	}
}

func (server *Server) GET(path string, handler func(Request) (Response, error), req Request) {
	server.Echo.GET(path, server.HandlerWrapper(handler, req))
}

func (server *Server) POST(path string, handler func(Request) (Response, error), req Request) {
	server.Echo.POST(path, server.HandlerWrapper(handler, req))
}

func (server *Server) PUT(path string, handler func(Request) (Response, error), req Request) {
	server.Echo.PUT(path, server.HandlerWrapper(handler, req))
}

func (server *Server) DELETE(path string, handler func(Request) (Response, error), req Request) {
	server.Echo.DELETE(path, server.HandlerWrapper(handler, req))
}

func (server *Server) OPTIONS(path string, handler func(Request) (Response, error), req Request) {
	server.Echo.OPTIONS(path, server.HandlerWrapper(handler, req))
}

func (server *Server) Start(address string) error {
	return server.Echo.Start(address)
}

func (server *Server) Stop(address string) error {
	return server.Echo.Close()
}

package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	"github.com/mopeneko/echo-msgpack/echomsgpack"
)

// Request of POST /hello
type Request struct {
	Name string `msgpack:"name"`
}

// Response of POST /hello
type Response struct {
	Message string `msgpack:"message"`
}

func main() {
	const addr = ":1350"

	e := echo.New()
	e.Binder = new(echomsgpack.Binder)
	e.Use(echomsgpack.ContextWrapper)

	e.POST("/hello", func(c echo.Context) error {
		cc, ok := c.(echomsgpack.Context)
		if !ok {
			return c.NoContent(http.StatusInternalServerError)
		}

		req := new(Request)
		res := new(Response)

		if err := cc.Bind(req); err != nil {
			res.Message = "failed to read a request."
			return cc.MessagePack(http.StatusBadRequest, res)
		}

		res.Message = fmt.Sprintf("Hello, %s.", req.Name)
		return cc.MessagePack(http.StatusOK, res)
	})

	e.Logger.Fatal(e.Start(addr))
}

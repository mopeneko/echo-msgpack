package echomsgpack

import (
	"log"
	"net/http"
	"strings"

	"github.com/labstack/echo"
	"github.com/vmihailenco/msgpack/v5"
)

const (
	charsetUTF8 = "charset=UTF-8"

	mime            = "application/x-msgpack"
	mimeCharsetUTF8 = mime + "; " + charsetUTF8
)

type (
	// Context adds a function to sending MessgePack response.
	Context interface {
		echo.Context
		MessagePack(code int, i interface{}) error
	}

	context struct {
		echo.Context
	}
)

// MessagePack sends a MessagePack response with status code.
func (c context) MessagePack(code int, i interface{}) error {
	// encode as MessagePack
	b, err := msgpack.Marshal(i)
	if err != nil {
		return err
	}

	return c.Blob(code, mimeCharsetUTF8, b)
}

// Binder for MessagePack.
type Binder struct{}

// Bind binds a request body to given interface.
func (b *Binder) Bind(i interface{}, c echo.Context) error {
	// validate with default binder
	db := new(echo.DefaultBinder)
	if err := db.Bind(i, c); err != echo.ErrUnsupportedMediaType {
		return err
	}

	req := c.Request()

	// check content type
	if !strings.HasPrefix(req.Header.Get(echo.HeaderContentType), mime) {
		log.Println(req.Header.Get(echo.HeaderContentType))
		return echo.ErrUnsupportedMediaType
	}

	// decode MessagePack body
	if err := msgpack.NewDecoder(req.Body).Decode(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error()).SetInternal(err)
	}

	return nil
}

// ContextOverrider wraps context for sending a MessagePack response.
func ContextWrapper(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c = context{c}
		return next(c)
	}
}

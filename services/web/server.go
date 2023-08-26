package web

import (
	"log"

	jwt "github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"deviloza.com.mx/misc"
	"deviloza.com.mx/auth"
)

type Server struct {
	Config      *misc.Secrets
	RunningPort string
}

func NewServer(config *misc.Secrets, runningPort string) *Server {
	s := Server{config, runningPort}
	return &s
}

func (s *Server) Start(routes Routes) error {
	e := echo.New()
	log.Println(s.Config)

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	// CORS default configuration
	e.Use(middleware.CORS())

	// Restricted endpoints
	r := e.Group("/restricted")
	// Configure middleware with the custom claims type
	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(auth.JWTCustomClaims)
		},
		SigningKey: []byte("secret"),
	}
	r.Use(echojwt.WithConfig(config))

	for _, route := range routes {
		switch method := route.Method; method {
		case "GET":
			if route.IsRestricted {
				r.GET(route.Url, route.Handler)
			} else {
				e.GET(route.Url, route.Handler)
			}
		case "POST":
			if route.IsRestricted {
				r.POST(route.Url, route.Handler)
			} else {
				e.POST(route.Url, route.Handler)
			}
		case "DELETE":
			if route.IsRestricted {
				r.DELETE(route.Url, route.Handler)
			} else {
				e.DELETE(route.Url, route.Handler)
			}
		case "PUT":
			if route.IsRestricted {
				r.PUT(route.Url, route.Handler)
			} else {
				e.PUT(route.Url, route.Handler)
			}
		default:
			log.Println("Unknown method for %s", method)
		}
	}

	e.Logger.Fatal(e.Start(":" + s.RunningPort))
	return nil
}

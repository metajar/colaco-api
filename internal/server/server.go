package server

import (
	"colaco-api/internal/api/v1"
	"colaco-api/internal/jwt"
	"colaco-api/svc"
	"fmt"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/labstack/echo/v4"
	emiddle "github.com/labstack/echo/v4/middleware"
	middleware "github.com/oapi-codegen/echo-middleware"
	"io/fs"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"
)

func s2ptr(s string) *string {
	return &s
}

func genErrorResponse(error string) v1.ErrorResp {
	return v1.ErrorResp{Error: s2ptr(error)}
}

func genMessageResponse(s string) v1.MessageResponse {
	return v1.MessageResponse{Message: &s}
}

var _ v1.ServerInterface = (*VendingMachine)(nil)

type VendingMachine struct {
	m           sync.RWMutex
	port        string
	SlotStorage svc.VendingStorageInterface
}

func authenticateUser(username, password string) bool {
	// For demonstration purposes only. We would actually call another method to verify
	// a username and password but this will be fine for now.
	return username == "admin" && password == "password"
}

func WithPort(port string) func(machine *VendingMachine) {
	return func(vm *VendingMachine) {
		vm.port = port
	}
}

func WithStorage(s svc.VendingStorageInterface) func(machine *VendingMachine) {
	return func(vm *VendingMachine) {
		vm.SlotStorage = s
	}
}

// WithStartingSodas sets the initial names of the vending slots in the
// VendingMachine based on the given list of sodas.
//
// The function takes in a slice of VendingSlot objects representing the slots in
// the vending machine, and returns a function that modifies the provided
// VendingMachine by setting the
func WithStartingSodas(sodas []v1.VendingSlot) func(machine *VendingMachine) {
	return func(vm *VendingMachine) {
		if vm.SlotStorage == nil {
			log.Fatalln("please initialize storage first via WithStorage option.")
		}
		for _, soda := range sodas {
			if soda.OccupiedSoda.Name != nil {
				// We will just set the name of the vending slot to the name of the Soda
				// until we need to support multiple slots with the same soda. Cola Co
				// hasn't stated they needed this.
				vm.SlotStorage.AddSlot(strings.ToLower(*soda.OccupiedSoda.Name), soda)
			}
		}
	}
}

// CreateMiddleware takes a JWSValidator and returns a slice of echo.MiddlewareFunc
// and an error. The function first tries to load the Swagger specification using the
// GetSwagger function. If there is an error loading the spec, it returns an error.
// Next, it creates a validator middleware using the OapiRequestValidatorWithOptions
// function from the "github.com/deepmap/oapi-codegen/v2/pkg/middleware" package. The
// validator middleware is configured with options to silence warning messages and
// set the authentication function using the NewAuthenticator function. Then, a custom
// skipAuthMiddleware is defined as a function that checks if the request path is
// "/auth/login", "/openapi.yaml", or "/docs". If the path matches any of these, it
// skips the validator middleware and proceeds to the next handler. Otherwise, it
// applies the validator middleware using the validator function returned by the
// oapi-codegen library. Finally, the skipAuthMiddleware is returned as the only
// element in the middleware slice.
func CreateMiddleware(v jwt.JWSValidator) ([]echo.MiddlewareFunc, error) {
	spec, err := v1.GetSwagger()
	if err != nil {
		return nil, fmt.Errorf("loading spec: %w", err)
	}

	validator := middleware.OapiRequestValidatorWithOptions(spec,
		&middleware.Options{
			SilenceServersWarning: true,
			Options: openapi3filter.Options{
				AuthenticationFunc: jwt.NewAuthenticator(v),
			},
		})

	// Wrap the validator in a custom middleware to exclude the /auth/login path
	skipAuthMiddleware := func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Check if the request path is for the authentication endpoint
			if c.Path() == "/auth/login" {
				// If so, skip the validator middleware and continue to the next handler
				return next(c)
			}
			if c.Path() == "/openapi.yaml" {
				// If so, skip the validator middleware and continue to the next handler
				return next(c)
			}
			if c.Path() == "/docs" {
				// If so, skip the validator middleware and continue to the next handler
				return next(c)
			}
			// For all other paths, apply the validator middleware
			return validator(next)(c)
		}
	}

	return []echo.MiddlewareFunc{skipAuthMiddleware}, nil
}

// NewVendingMachine creates a new instance of the VendingMachine struct with the
// provided options applied. The options parameter is a variadic function that
// takes in functions with
func NewVendingMachine(options ...func(machine *VendingMachine)) *VendingMachine {
	vm := &VendingMachine{}
	for _, option := range options {
		option(vm)
	}
	if vm.port == "" {
		// sets a default in case
		// port is not set via func opt.
		vm.port = "8080"
	}
	return vm
}

func (v *VendingMachine) Run() {
	e := echo.New()
	fa, err := jwt.NewFakeAuthenticator()
	if err != nil {
		log.Fatalln("error creating the authenticator and can't move forward:", err.Error())
	}
	mw, err := CreateMiddleware(fa)
	e.Use(emiddle.Logger())
	e.Use(mw...)
	e.GET("/openapi.yaml", func(c echo.Context) error {
		file, err := v1.Content.Open("api.yml")
		if err != nil {
			return err
		}
		defer func(file fs.File) {
			err := file.Close()
			if err != nil {
				log.Fatalln(err)
			}
		}(file)
		f, err := fs.ReadFile(v1.Content, "api.yml")
		return c.Blob(http.StatusOK, "application/x-yaml", f)
	})

	e.GET("/docs", func(c echo.Context) error {
		htmlContent, err := fs.ReadFile(v1.Content, "redoc.html")
		if err != nil {
			return err // Handle error
		}
		return c.HTMLBlob(http.StatusOK, htmlContent)
	})
	v1.RegisterHandlers(e, v)
	e.Logger.Fatal(e.Start(net.JoinHostPort("0.0.0.0", v.port)))
}

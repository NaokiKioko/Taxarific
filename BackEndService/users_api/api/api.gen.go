// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.3.0 DO NOT EDIT.
package api

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gin-gonic/gin"
	"github.com/oapi-codegen/runtime"
)

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Get all admins
	// (GET /admin)
	GetAdmin(c *gin.Context)
	// Create an admin
	// (POST /admin)
	PostAdmin(c *gin.Context)
	// Create an employee
	// (POST /admin/employee)
	PostAdminEmployee(c *gin.Context)
	// Get all cases
	// (GET /case)
	GetCase(c *gin.Context)
	// Get all pending cases
	// (GET /case/pending)
	GetCasePending(c *gin.Context)
	// Get all employees
	// (GET /employee)
	GetEmployee(c *gin.Context)
	// Allows Employees to update certain information
	// (PUT /employee/addcase/{caseid})
	PutEmployeeAddcaseCaseid(c *gin.Context, caseid string)
	// Login a user, employee, or admin
	// (POST /login)
	PostLogin(c *gin.Context)
	// Get all users
	// (GET /user)
	GetUser(c *gin.Context)
	// Create a user
	// (POST /user)
	PostUser(c *gin.Context)
	// Allows users to submit a case
	// (PUT /user/case)
	PutUserCase(c *gin.Context)
	// Get a user
	// (GET /user/profile)
	GetUserProfile(c *gin.Context)
	// Update a user
	// (PUT /user/profile)
	PutUserProfile(c *gin.Context)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandler       func(*gin.Context, error, int)
}

type MiddlewareFunc func(c *gin.Context)

// GetAdmin operation middleware
func (siw *ServerInterfaceWrapper) GetAdmin(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.GetAdmin(c)
}

// PostAdmin operation middleware
func (siw *ServerInterfaceWrapper) PostAdmin(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.PostAdmin(c)
}

// PostAdminEmployee operation middleware
func (siw *ServerInterfaceWrapper) PostAdminEmployee(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.PostAdminEmployee(c)
}

// GetCase operation middleware
func (siw *ServerInterfaceWrapper) GetCase(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.GetCase(c)
}

// GetCasePending operation middleware
func (siw *ServerInterfaceWrapper) GetCasePending(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.GetCasePending(c)
}

// GetEmployee operation middleware
func (siw *ServerInterfaceWrapper) GetEmployee(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.GetEmployee(c)
}

// PutEmployeeAddcaseCaseid operation middleware
func (siw *ServerInterfaceWrapper) PutEmployeeAddcaseCaseid(c *gin.Context) {

	var err error

	// ------------- Path parameter "caseid" -------------
	var caseid string

	err = runtime.BindStyledParameterWithOptions("simple", "caseid", c.Param("caseid"), &caseid, runtime.BindStyledParameterOptions{Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter caseid: %w", err), http.StatusBadRequest)
		return
	}

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.PutEmployeeAddcaseCaseid(c, caseid)
}

// PostLogin operation middleware
func (siw *ServerInterfaceWrapper) PostLogin(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.PostLogin(c)
}

// GetUser operation middleware
func (siw *ServerInterfaceWrapper) GetUser(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.GetUser(c)
}

// PostUser operation middleware
func (siw *ServerInterfaceWrapper) PostUser(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.PostUser(c)
}

// PutUserCase operation middleware
func (siw *ServerInterfaceWrapper) PutUserCase(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.PutUserCase(c)
}

// GetUserProfile operation middleware
func (siw *ServerInterfaceWrapper) GetUserProfile(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.GetUserProfile(c)
}

// PutUserProfile operation middleware
func (siw *ServerInterfaceWrapper) PutUserProfile(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.PutUserProfile(c)
}

// GinServerOptions provides options for the Gin server.
type GinServerOptions struct {
	BaseURL      string
	Middlewares  []MiddlewareFunc
	ErrorHandler func(*gin.Context, error, int)
}

// RegisterHandlers creates http.Handler with routing matching OpenAPI spec.
func RegisterHandlers(router gin.IRouter, si ServerInterface) {
	RegisterHandlersWithOptions(router, si, GinServerOptions{})
}

// RegisterHandlersWithOptions creates http.Handler with additional options
func RegisterHandlersWithOptions(router gin.IRouter, si ServerInterface, options GinServerOptions) {
	errorHandler := options.ErrorHandler
	if errorHandler == nil {
		errorHandler = func(c *gin.Context, err error, statusCode int) {
			c.JSON(statusCode, gin.H{"msg": err.Error()})
		}
	}

	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandler:       errorHandler,
	}

	router.GET(options.BaseURL+"/admin", wrapper.GetAdmin)
	router.POST(options.BaseURL+"/admin", wrapper.PostAdmin)
	router.POST(options.BaseURL+"/admin/employee", wrapper.PostAdminEmployee)
	router.GET(options.BaseURL+"/case", wrapper.GetCase)
	router.GET(options.BaseURL+"/case/pending", wrapper.GetCasePending)
	router.GET(options.BaseURL+"/employee", wrapper.GetEmployee)
	router.PUT(options.BaseURL+"/employee/addcase/:caseid", wrapper.PutEmployeeAddcaseCaseid)
	router.POST(options.BaseURL+"/login", wrapper.PostLogin)
	router.GET(options.BaseURL+"/user", wrapper.GetUser)
	router.POST(options.BaseURL+"/user", wrapper.PostUser)
	router.PUT(options.BaseURL+"/user/case", wrapper.PutUserCase)
	router.GET(options.BaseURL+"/user/profile", wrapper.GetUserProfile)
	router.PUT(options.BaseURL+"/user/profile", wrapper.PutUserProfile)
}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"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",
	"/wvvzpgvIwAA",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %w", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	res := make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	resolvePath := PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		pathToFile := url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}

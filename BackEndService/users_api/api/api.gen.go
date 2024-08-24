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
	// Update a user
	// (PUT /user/{userid})
	PutUserUserid(c *gin.Context, userid string)
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

// PutUserUserid operation middleware
func (siw *ServerInterfaceWrapper) PutUserUserid(c *gin.Context) {

	var err error

	// ------------- Path parameter "userid" -------------
	var userid string

	err = runtime.BindStyledParameterWithOptions("simple", "userid", c.Param("userid"), &userid, runtime.BindStyledParameterOptions{Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter userid: %w", err), http.StatusBadRequest)
		return
	}

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.PutUserUserid(c, userid)
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
	router.GET(options.BaseURL+"/employee", wrapper.GetEmployee)
	router.PUT(options.BaseURL+"/employee/addcase/:caseid", wrapper.PutEmployeeAddcaseCaseid)
	router.POST(options.BaseURL+"/login", wrapper.PostLogin)
	router.GET(options.BaseURL+"/user", wrapper.GetUser)
	router.POST(options.BaseURL+"/user", wrapper.PostUser)
	router.PUT(options.BaseURL+"/user/:userid", wrapper.PutUserUserid)
}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/+xYTW/bRhD9K4ttj5SotDkEurlOYLgo0Bxi9BAYwYocUZuQO8zurGzB4H8vdvhhSqRk",
	"yXbiHHwxaHKGs/Pemw/qTiZYlGjAkJPzO+mSFRSKL8/SQptwUVoswZIGvg2F0nm4WKItFMl5cyeStClB",
	"zqUjq00mq0jqNNht347k7STDSXOztLrQpNcw/XfxFRK6fN83mOiiREt8BkUrOZcZTgs0GaaLKdos5utJ",
	"avUabLxwaOLuhSH+7QRVqScJppCBmcAtWTUhlXEewVzO5RedRlhogqKkjayqSBpVwPDcVSRL5dwN2j1J",
	"HQj1tQ7VD1NF0sJ3ry2kcv45INUEjjo4u3DXHbTIIIWznCsHQ24S5eBLCi6xuiSNZjQNNqqp6Sj0nk8w",
	"butIkXej7/IO7HHv2km4dRzL7UNR5riBPfnxRQCSL363sJRz+Vt8L+O40XDMGFVdAGWt2oT/XxX8Kyj4",
	"yoEdZ/hYXl95fHkeg7s2S+STaMrDs0/qVlm91IkIHDtx9vFSRnIN1nFLkm+ms+kspIElGFVqOZd/8q2I",
	"weWDxqqdPhkw7FtdTV4ACZXngq2coJUisVJrEEsf7iYJODeVHMGq4HOZ1l71UAspuxJN003+mM1Ye2gI",
	"DEdTZZnrhD3jGrN2NB7dfVjfg+4T8NpO5Uzk2pHAZS+hKdPvfFEouxmkGwTOfH6WNUrXQRLoRnA6t6AI",
	"hDK1awuUO4zTR3Q9oL57cPQXppuTMHr00vAI0e9qeK98d9RO1kM10MKbk/I8JIEawzHOmYyEyUmF88xD",
	"oGSzy/wOgSPUV1FTLDH05+YDcmhta0VoJ9QiB0EoSH0DgUaEAS143h7QRzepX3XyFJ10MI5IpX12olrg",
	"npkxwfSlcrDBtoZNj1UWTlLKBdCWRn50071H8tTG2yW6r/d2Bj1IOxi3UY1VmgZA4rvwV6cVC92PwHxV",
	"ps9SkL6D+awOfc6BWdRWFUBgw4F3w1++D+kH+zpsd4gb7QLlXfQQMozjevtpi08mbZTtWol6tO1W4PVz",
	"tYrjN8WdqmfH40p89nNL3LMaBiUeybezt0PtdG4GSSzRm3RXumd5jjdOfLivYWxiiAQsKW1E2N1Cpw2v",
	"3KvrHLPml4DRqfJPeCxUrdBuw3A6M5AKX4olWkHtTjg+TfgVLzFFDgyLSFrM4eEp8sjxMXtCgoTfwAyJ",
	"+Pu/T6J+NPz8HtFd2FBFjlkGqdDmSNWx017F9aUQde0kEmgH+0swcY3AfPMteHAUsUNv1T9aYBdAvIz/",
	"uls/p7Zv8NRADYF7aOd/Sj12eL0udcfnqce7yDMW6yk7INN/oN7iu/D3mOXk0ULyrKMrDnPcKkIraKK1",
	"k2rP5uHbl77A5vGj9fyTvz7qlrVHcaetJIeHw5aexpQZrMGuW4F4m/NISVS+Qkfzd7N3MxlobBwHjc87",
	"wuLoIdFTkpNVdDf2jX7gp6XGvR5rI+5713ehTCqWOs8FemLF95YwsVAO0uBwEwJ39eD8otAUmFhrJUqd",
	"kLfQPwfcf/kMcEFDWhknCiCrE8dQlBYLoBX4/ksaA1ldV/8HAAD//3iGHwMZGgAA",
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

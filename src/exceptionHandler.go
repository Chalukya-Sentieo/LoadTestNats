package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"
	"time"

	"github.com/ansel1/merry"
	"github.com/gin-gonic/gin"
)

type ExceptionHandler interface {
	getExceptionHandler() ExceptionHandler
	captureException(errMsg string)
	handleError(c *gin.Context)
}

/**
Bad Request Exception and it's associated methods
*/
type BadRequestException struct {
	err        merry.Error
	statusCode int
}

func (exp *BadRequestException) getExceptionHandler() ExceptionHandler {
	return exp
}

func (exp *BadRequestException) captureException(errMsg string) {
	exp.err = merry.New(errMsg)
	exp.statusCode = 400
	panic(exp)
}

func (exp *BadRequestException) handleError(c *gin.Context) {
	errorBody := appendStacktrace(exp.err)
	c.AbortWithStatusJSON(exp.statusCode, errorBody)
}

/**
Unauthorized Request Exception and it's associated methods
*/
type UnauthorizedRequestException struct {
	err        merry.Error
	statusCode int
}

func (exp *UnauthorizedRequestException) getExceptionHandler() ExceptionHandler {
	return exp
}

func (exp *UnauthorizedRequestException) captureException(errMsg string) {
	exp.err = merry.New(errMsg)
	exp.statusCode = 401
	panic(exp)
}

func (exp *UnauthorizedRequestException) handleError(c *gin.Context) {
	errorBody := appendStacktrace(exp.err)
	c.AbortWithStatusJSON(exp.statusCode, errorBody)
}

/**
Unauthorized Request Exception and it's associated methods
*/
type InternalServerError struct {
	err        merry.Error
	statusCode int
}

func (exp *InternalServerError) getExceptionHandler() ExceptionHandler {
	return exp
}

func (exp *InternalServerError) captureException(errMsg string) {
	exp.err = merry.New(errMsg)
	exp.statusCode = 500
	//panic(err) //This always gets thrown from inside the exceptionhandler, so panic is redundant here
}

func (exp *InternalServerError) handleError(c *gin.Context) {
	errorBody := appendStacktrace(exp.err)
	c.AbortWithStatusJSON(exp.statusCode, errorBody)
}

/**
Recovery middleware for handling all application panic and responding with appropriate status
*/

func getCustomRecoveryMiddleware() gin.HandlerFunc {
	return gin.CustomRecovery(func(context *gin.Context, exp interface{}) {
		log.Printf("Custom Recovery Handler Invoked for (%s) type\n", reflect.TypeOf(exp).String())
		errHandler, ok := exp.(ExceptionHandler)

		//Wrap all unhandled system errors in InternalServerError
		if !ok {
			errObj := InternalServerError{}
			errHandler = errObj.getExceptionHandler()

			errHandler.captureException(fmt.Sprint(exp))
		}
		errHandler.handleError(context)
	})
}

func getErrBody(err merry.Error) gin.H {
	return gin.H{
		"error": fmt.Sprint(err),
	}
}

func appendStacktrace(err merry.Error) gin.H {
	errorMsg := getErrBody(err)

	if os.Getenv("APP_DEBUG") == "true" {
		errorMsg["stackTrace"] = merry.Stacktrace(err)
	}
	return errorMsg
}

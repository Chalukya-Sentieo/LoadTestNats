package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"
	"time"

	"github.com/ansel1/merry"
	"github.com/getsentry/sentry-go"
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
	sentry.CaptureException(exp.err)

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
	sentry.CaptureException(exp.err)

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
	sentry.CaptureException(exp.err)

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

/**
Initialize Sentry SDK
*/

func prettyPrint(v interface{}) string {
	pp, _ := json.MarshalIndent(v, "", "  ")
	return string(pp)
}

func initializeSentry() {
	APP_DEBUG := false
	if os.Getenv("APP_DEBUG") == "true" {
		// SENTRY_DEBUG_MODE set this to true in your env to make sure sentry can report exceptions even in your debug environments
		if os.Getenv("SENTRY_DEBUG_MODE") != "true" {
			return
		}
		APP_DEBUG = true
	}
	err := sentry.Init(sentry.ClientOptions{
		// Either set your DSN here or set the SENTRY_DSN environment variable.
		Dsn: os.Getenv("SENTRY_DSN"),
		// Either set environment and release here or set the SENTRY_ENVIRONMENT
		// and SENTRY_RELEASE environment variables.
		Environment: os.Getenv("APP_ENV"),
		Release:     os.Getenv("RELEASE"),
		// Enable printing of SDK debug messages.
		// Useful when getting started or trying to figure something out.
		Debug: APP_DEBUG,
		BeforeSend: func(event *sentry.Event, hint *sentry.EventHint) *sentry.Event {
			// if ex, ok := hint.OriginalException.(ExceptionHandler); ok {
			// 	event.Message = event.Message + " - " + fmt.Sprint(ex)
			// 	log.Println("Exception caught. Sending to sentry", ex)
			// }

			//fmt.Printf("Before send %s\n\n", prettyPrint(event))
			return event
		},
		BeforeBreadcrumb: func(breadcrumb *sentry.Breadcrumb, _ *sentry.BreadcrumbHint) *sentry.Breadcrumb {
			// if breadcrumb.Message == "Random breadcrumb 3" {
			// 	breadcrumb.Message = "Not so random breadcrumb 3"
			// }

			// log.Printf("Breadcrumbs %s\n\n", prettyPrint(breadcrumb))

			return breadcrumb
		},
	})
	if err != nil {
		log.Panicf("Error in sentry.Init: %s", err)
	}
	// Flush buffered events before the program terminates.
	// Set the timeout to the maximum duration the program can afford to wait.
	defer sentry.Flush(2 * time.Second)
	defer sentry.Recover()
}

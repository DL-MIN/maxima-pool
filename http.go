/*******************************************************************************
 * HTTP server
 *
 * @author     Lars Thoms <lars@thoms.io>
 * @date       2023-04-17
 ******************************************************************************/

package main

import (
	"Moodle_Maxima_Pool/controller"
	"Moodle_Maxima_Pool/models"
	"Moodle_Maxima_Pool/services"
	"context"
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
	"path"
	"strings"
)

var (
	router              *gin.Engine
	errUnauthenticated  = &models.ErrorResponseJSON{Status: http.StatusUnauthorized, Code: "unauthorized", Title: "Unauthorized", Details: "The request misses a valid API key."}
	errUndefinedRequest = &models.ErrorResponseJSON{Status: http.StatusRequestedRangeNotSatisfiable, Code: "undefined_request", Title: "Undefined request", Details: "The type of request is undefined."}
)

func initHTTPConfig() {
	// Set GIN mode
	if logger.Level() == Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// Check API key length
	if len(viper.GetString("server.api_key")) < 16 {
		logger.Warn("API key is very short")
	}
}

func initHTTPRoutes() {
	router = gin.Default()

	router.Use(gin.CustomRecovery(errorHandlerGin))

	router.Use(globalHeader())

	router.GET("/openapi.json", controller.GetOpenAPI)

	authorized := router.Group(path.Clean(viper.GetString("server.base_path")), validateAPIKey())

	// Job
	authorized.POST("/MaximaPool", controller.PostJob)
}

func startHTTPServer() {
	waitGroup.Add(1)
	defer waitGroup.Done()

	logger.Debug("configure web server")
	initHTTPConfig()

	logger.Debug("create routes")
	initHTTPRoutes()

	logger.Debug("create web server")
	server := &http.Server{Addr: fmt.Sprintf("%s:%d", viper.GetString("server.listen"), viper.GetInt("server.port")), Handler: router}

	go func() {
		logger.Infof("start web server and listen to http://%s:%d", viper.GetString("server.listen"), viper.GetInt("server.port"))
		if err := server.ListenAndServe(); err == http.ErrServerClosed {
			logger.Info(err)
		} else {
			logger.Fatal(err)
		}
	}()

	<-terminator

	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Warn(err)
	}
}

func validateAPIKey() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Header.Get("X-API-Key") == viper.GetString("server.api_key") {
			return
		}

		basicAuthHeader := strings.SplitN(c.Request.Header.Get("Authorization"), " ", 2)
		if len(basicAuthHeader) == 2 && strings.EqualFold(basicAuthHeader[0], "Basic") {
			if basicAuthPayload, err := base64.StdEncoding.DecodeString(basicAuthHeader[1]); err == nil {
				if basicAuthPair := strings.SplitN(string(basicAuthPayload), ":", 2); len(basicAuthPair) == 2 && basicAuthPair[1] == viper.GetString("server.api_key") {
					return
				}
			}
		}

		c.AbortWithStatusJSON(services.Error(errUnauthenticated))
	}
}

func globalHeader() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Link", "</openapi.json>; rel=\"service-desc\"")
	}
}

func errorHandlerGin(c *gin.Context, err any) {
	logger.Warn(err)
	c.AbortWithStatusJSON(services.Error(errUndefinedRequest))
}

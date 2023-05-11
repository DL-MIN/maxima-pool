/*******************************************************************************
 * Controller: GET OpenAPI
 *
 * @author     Lars Thoms <lars@thoms.io>
 * @date       2023-04-14
 ******************************************************************************/

package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

import _ "embed"

//go:embed openapi.json
var data []byte

func GetOpenAPI(c *gin.Context) {
	c.Data(http.StatusOK, gin.MIMEJSON, data)
}

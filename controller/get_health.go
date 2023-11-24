/*******************************************************************************
 * Controller: GET Health
 *
 * @author     Lars Thoms <lars@thoms.io>
 * @date       2023-04-14
 ******************************************************************************/

package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetHealth(c *gin.Context) {
	c.AbortWithStatus(http.StatusOK)
}

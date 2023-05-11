/*******************************************************************************
 * Controller: POST job
 *
 * @author     Lars Thoms <lars@thoms.io>
 * @date       2023-04-14
 ******************************************************************************/

package controller

import (
	"Moodle_Maxima_Pool/models"
	"Moodle_Maxima_Pool/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

func PostJob(c *gin.Context) {
	reqQuery := &models.JobRequestQuery{
		Timeout:     30000,
		PlotURLBase: "!ploturl!",
	}

	if err := c.ShouldBindQuery(reqQuery); err != nil {
		c.AbortWithStatusJSON(services.Error(errRequestInvalid))
		return
	}

	if resp, err := services.JobCreate(reqQuery); err != nil {
		c.JSON(http.StatusRequestedRangeNotSatisfiable, err)
	} else if resp.IsZIP {
		c.DataFromReader(http.StatusOK, int64(resp.Output.Len()), "application/zip", resp.Output, map[string]string{
			"Content-Disposition": `attachment; filename="output.zip"`,
		})
	} else {
		c.DataFromReader(http.StatusOK, int64(resp.Output.Len()), gin.MIMEPlain, resp.Output, nil)
	}
}

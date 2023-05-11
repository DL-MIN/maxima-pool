/*******************************************************************************
 * Model: job
 *
 * @author     Lars Thoms <lars@thoms.io>
 * @date       2023-04-24
 ******************************************************************************/

package models

import (
	"bytes"
)

type JobRequestQuery struct {
	Input       string `form:"input" binding:"required"`
	Timeout     int    `form:"timeout" binding:"omitempty"`
	PlotURLBase string `form:"ploturlbase" binding:"omitempty"`
	Version     string `form:"version" binding:"omitempty"`
}

type JobResponse struct {
	Output *bytes.Buffer
	IsZIP  bool
}

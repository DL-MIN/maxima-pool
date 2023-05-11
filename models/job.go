/*******************************************************************************
 * Model: job
 *
 * @author     Lars Thoms <lars@thoms.io>
 * @date       2023-04-24
 ******************************************************************************/

package models

import (
	"bytes"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/hashicorp/go-version"
	"regexp"
)

var versionRegex = regexp.MustCompile(version.VersionRegexpRaw)

type JobRequestQuery struct {
	Input       string `form:"input" binding:"required"`
	Timeout     int    `form:"timeout" binding:"omitempty"`
	PlotURLBase string `form:"ploturlbase" binding:"omitempty"`
	Version     string `form:"version" binding:"omitempty,version"`
}

type JobResponse struct {
	Output *bytes.Buffer
	IsZIP  bool
}

func init() {
	versionValidator := func(fl validator.FieldLevel) bool {
		return versionRegex.MatchString(fl.Field().String())
	}

	if validate, ok := binding.Validator.Engine().(*validator.Validate); ok {
		if err := validate.RegisterValidation("version", versionValidator); err != nil {
			panic("could not register validator")
		}
	} else {
		panic("could not handle validator engine")
	}
}

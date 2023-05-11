/*******************************************************************************
 * Service: error
 *
 * @author     Lars Thoms <lars@thoms.io>
 * @date       2023-04-14
 ******************************************************************************/

package services

import (
	"Moodle_Maxima_Pool/models"
)

func Error(error *models.ErrorResponseJSON) (status int, errors *models.ErrorsResponseJSON) {
	errors = &models.ErrorsResponseJSON{Errors: []*models.ErrorResponseJSON{error}}
	status = errors.Status()
	return
}

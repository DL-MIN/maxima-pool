/*******************************************************************************
 * Test: Service: error
 *
 * @author     Lars Thoms <lars@thoms.io>
 * @date       2023-04-21
 ******************************************************************************/

package services

import (
	"Moodle_Maxima_Pool/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestError(t *testing.T) {
	type args struct {
		error *models.ErrorResponseJSON
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErrors *models.ErrorsResponseJSON
	}{
		{"valid 400 error", args{error: &models.ErrorResponseJSON{Status: 400}}, 400, &models.ErrorsResponseJSON{Errors: []*models.ErrorResponseJSON{{Status: 400}}}},
		{"valid 404 error", args{error: &models.ErrorResponseJSON{Status: 404}}, 404, &models.ErrorsResponseJSON{Errors: []*models.ErrorResponseJSON{{Status: 404}}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotStatus, gotErrors := Error(tt.args.error)
			assert.Equal(t, tt.wantStatus, gotStatus)
			assert.Equal(t, tt.wantErrors, gotErrors)
		})
	}
}

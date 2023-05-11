/*******************************************************************************
 * Controller
 *
 * @author     Lars Thoms <lars@thoms.io>
 * @date       2023-04-14
 ******************************************************************************/

package controller

import (
	"Moodle_Maxima_Pool/models"
	"net/http"
)

var (
	errRequestInvalid  = &models.ErrorResponseJSON{Status: http.StatusBadRequest, Code: "invalid_input", Title: "Invalid input", Details: "The request is invalid."}
	errNotImplemented  = &models.ErrorResponseJSON{Status: http.StatusBadRequest, Code: "not_implemented", Title: "Not implemented", Details: "This action is not implemented."}
	errLanguageMissing = &models.ErrorResponseJSON{Status: http.StatusBadRequest, Code: "language_not_found", Title: "Language not found", Details: "The requested language does not exist."}
	errFileMissing     = &models.ErrorResponseJSON{Status: http.StatusNotFound, Code: "file_not_found", Title: "File not found", Details: "The requested file does not exist."}
	errFileCreation    = &models.ErrorResponseJSON{Status: http.StatusBadRequest, Code: "file_creation", Title: "File not created", Details: "The sent file could not be created."}
	errFileHash        = &models.ErrorResponseJSON{Status: http.StatusBadRequest, Code: "file_hash", Title: "File hash", Details: "The sent file hash does not equal to our hash calculation."}
)

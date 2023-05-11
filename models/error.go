/*******************************************************************************
 * Model: error
 *
 * @author     Lars Thoms <lars@thoms.io>
 * @date       2023-04-24
 ******************************************************************************/

package models

type ErrorsResponseJSON struct {
	Errors []*ErrorResponseJSON `json:"errors"`
}

type ErrorResponseJSON struct {
	Status  int    `json:"status"`
	Code    string `json:"code"`
	Title   string `json:"title"`
	Details string `json:"details"`
}

func (e *ErrorsResponseJSON) Status() int {
	return e.Errors[0].Status
}

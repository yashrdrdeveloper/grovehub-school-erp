// Package controller provides a common GoMock controller factory.
// It is intentionally domain-neutral and does not contact any external system.
package controller

import "go.uber.org/mock/gomock"

func New(reporter gomock.TestReporter) *gomock.Controller {
    return gomock.NewController(reporter)
}

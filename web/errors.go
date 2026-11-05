package web

import (
	"github.com/bigphu/goat/errs"
	"net/http"
)

var (
	ErrMalformedJSON = errs.BadRequest(
		"The request body could not be parsed as JSON.",
		errs.WithType("/errors/malformed-json"), errs.WithTitle("Malformed JSON body"))

	ErrRequestTooLarge = errs.New("", errs.WithStatus(http.StatusRequestEntityTooLarge),
		errs.WithType("/errors/request-too-large"), errs.WithTitle("Request body too large"))

	ErrValidation = errs.Unprocessable("",
		errs.WithType("/errors/validation-failed"), errs.WithTitle("Validation failed"))

	ErrInvalidRequest = errs.BadRequest("",
		errs.WithType("/errors/invalid-request"), errs.WithTitle("Invalid request"))

	ErrRouteNotFound = errs.NotFound("",
		errs.WithType("/errors/route-not-found"), errs.WithTitle("Route not found"))

	ErrMethodNotAllowed = errs.New("", errs.WithStatus(http.StatusMethodNotAllowed),
		errs.WithType("/errors/method-not-allowed"), errs.WithTitle("Method not allowed"))

	ErrUnsupportedMediaType = errs.New("", errs.WithStatus(http.StatusUnsupportedMediaType),
		errs.WithType("/errors/unsupported-media-type"), errs.WithTitle("Unsupported media type"))
)

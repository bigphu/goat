package apicommon_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bigphu/goat/errs"
	"github.com/bigphu/goat/web"
	"github.com/bigphu/goat/web/apicommon"
)

// Body lỗi server thật sự gửi (web.RenderProblem) phải khớp schema Problem trong
// api/common.yaml, vì frontend sinh type từ spec. Field lạ = spec lệch code.
func TestProblemSchemaMatchesErrs(t *testing.T) {
	e := errs.Unprocessable("The request body has 1 invalid field.",
		errs.WithType("/errors/validation-failed"), errs.WithTitle("Validation failed"),
		errs.WithFields(errs.NewField("email", "is required")))
	rec := httptest.NewRecorder()
	web.RenderProblem(rec, httptest.NewRequest(http.MethodGet, "/", nil), e)
	b := rec.Body.Bytes()

	dec := json.NewDecoder(bytes.NewReader(b))
	dec.DisallowUnknownFields()
	var p apicommon.Problem
	if err := dec.Decode(&p); err != nil {
		t.Fatalf("problem body %s does not fit apicommon.Problem: %v", b, err)
	}

	if p.Type != "/errors/validation-failed" || p.Title != "Validation failed" || p.Status != http.StatusUnprocessableEntity {
		t.Errorf("problem = %+v", p)
	}
	if p.Detail == nil || *p.Detail != "The request body has 1 invalid field." {
		t.Errorf("detail = %v", p.Detail)
	}
	if p.Errors == nil || len(*p.Errors) != 1 || (*p.Errors)[0] != (apicommon.FieldError{Field: "email", Detail: "is required"}) {
		t.Errorf("errors = %v", p.Errors)
	}
}

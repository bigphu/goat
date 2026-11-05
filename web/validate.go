package web

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/bigphu/goat/errs"
	"github.com/go-playground/validator/v10"
)

// Một validator dùng chung cho cả backend
var validate = newValidator()

func newValidator() *validator.Validate {
	v := validator.New(validator.WithRequiredStructEnabled())

	// Lỗi báo theo tên trong JSON (email), không theo tên field Go (Email)
	v.RegisterTagNameFunc(func(f reflect.StructField) string {
		name, _, _ := strings.Cut(f.Tag.Get("json"), ",")
		if name == "" || name == "-" {
			return f.Name
		}
		return name
	})
	return v
}

func Validate(v any) error {
	err := validate.Struct(v)
	if err == nil {
		return nil
	}

	var verrs validator.ValidationErrors
	if !errors.As(err, &verrs) {
		// Không phải lỗi dữ liệu của client mà là lỗi code => 500
		return errs.ErrInternal.With(
			errs.WithCause(fmt.Errorf("validate: %w", err)))
	}

	fields := make([]*errs.FieldError, 0, len(verrs))
	for _, fe := range verrs {
		fields = append(fields, fieldError(fieldPath(fe), ruleMessage(fe)))
	}

	return validationError(fields)
}

// detailRequired là câu báo thiếu field, xem fieldError
const detailRequired = "is required"

// fieldError dựng lỗi của một field kèm sentinel tương ứng, để errors.Is với
// errs.ErrRequired, errs.ErrInvalid dùng được. Không dùng errs.FieldRequired vì
// câu của nó khác câu client đang nhận.
func fieldError(name, detail string) *errs.FieldError {
	if detail == detailRequired {
		return errs.NewField(name, detail, errs.WithFieldCause(errs.ErrRequired))
	}
	return errs.FieldInvalid(name, detail)
}

func validationError(fields []*errs.FieldError) *errs.Error {
	noun := "fields"
	if len(fields) == 1 {
		noun = "field"
	}
	return ErrValidation.With(
		errs.WithFields(fields...),
		errs.WithDetail(fmt.Sprintf("The request body has %d invalid %s.", len(fields), noun)))
}

// fieldPath trả đường dẫn đầy đủ của field lỗi, vd address.city
// (Namespace bắt đầu bằng tên struct gốc, bỏ đi)
func fieldPath(fe validator.FieldError) string {
	if _, path, ok := strings.Cut(fe.Namespace(), "."); ok {
		return path
	}
	return fe.Field()
}

// ruleMessage đổi tag của validator thành câu dễ đọc
func ruleMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return detailRequired
	case "email":
		return "must be a valid email address"
	case "url":
		return "must be a valid URL"
	case "uuid":
		return "must be a valid UUID"
	case "min":
		if fe.Kind() == reflect.String {
			return fmt.Sprintf("must be at least %s characters", fe.Param())
		}
		return fmt.Sprintf("must be at least %s", fe.Param())
	case "max":
		if fe.Kind() == reflect.String {
			return fmt.Sprintf("must be at most %s characters", fe.Param())
		}
		return fmt.Sprintf("must be at most %s", fe.Param())
	case "len":
		return fmt.Sprintf("must be exactly %s characters", fe.Param())
	case "oneof":
		return fmt.Sprintf("must be one of: %s", strings.ReplaceAll(fe.Param(), " ", ", "))
	case "eqfield":
		return fmt.Sprintf("must match %s", fe.Param())
	default:
		return fmt.Sprintf("failed the %q rule", fe.Tag())
	}
}

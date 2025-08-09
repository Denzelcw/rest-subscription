package validator

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

func ValidateDates(startDateStr, endDateStr string) error {
	const layout = "01-2006"

	startDate, err := time.Parse(layout, startDateStr)
	if err != nil {
		return fmt.Errorf("invalid start_date format: %w", err)
	}

	if endDateStr == "" {
		return nil
	}
	fmt.Println("endDateStr", endDateStr)
	endDate, err := time.Parse(layout, endDateStr)
	if err != nil {
		return fmt.Errorf("invalid end_date format: %w", err)
	}

	if !endDate.After(startDate) {
		return fmt.Errorf("end_date must be after start_date")
	}

	return nil
}

func ValidationError(errs validator.ValidationErrors, req interface{}) string {
	var errMsgs []string

	fieldToJSON := make(map[string]string)
	t := reflect.TypeOf(req)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		jsonTag := field.Tag.Get("json")
		if jsonTag == "" {
			jsonTag = field.Name
		}
		fieldToJSON[field.Name] = jsonTag
	}

	for _, err := range errs {
		fieldName := err.Field()
		jsonName, ok := fieldToJSON[fieldName]
		if !ok {
			jsonName = fieldName
		}
		switch err.Tag() {
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is required", jsonName))
		case "min":
			switch fieldName {
			case "ServiceName":
				errMsgs = append(errMsgs, fmt.Sprintf("field %s must be at least %s characters long", jsonName, err.Param()))
			case "Price":
				errMsgs = append(errMsgs, fmt.Sprintf("field %s must be at least %s", jsonName, err.Param()))
			default:
				errMsgs = append(errMsgs, fmt.Sprintf("field %s has a minimum value requirement", jsonName))
			}

		case "max":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s must be no more than %s characters long", jsonName, err.Param()))
		case "uuid4":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s must be a valid UUID", jsonName))
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not valid (unknown reason)", jsonName))
		}
	}
	return strings.Join(errMsgs, ", ")
}

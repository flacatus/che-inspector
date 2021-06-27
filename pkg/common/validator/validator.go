package validator

import (
	"fmt"

	"github.com/flacatus/che-inspector/pkg/api"
	"github.com/go-playground/validator"
)

// The CheInspectorValidator validate all fields from structure read it from yaml file.
// Validate all basic fields like name, version etc and if exists tests in yaml file
func CheInspectorValidator(inspector *api.CheInspector) (err error) {
	validate := validator.New()
	// register validation for 'CheInspector'
	// NOTE: only have to register a non-pointer type for 'CheInspector', validator
	// internally dereferences during it's type checks.
	validate.RegisterStructValidation(validateCheInspectorStruct, api.CheInspector{})
	// returns InvalidValidationError for bad validation input, nil or ValidationErrors ( []FieldError )
	err = validate.Struct(inspector)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return
		}

		for _, err := range err.(validator.ValidationErrors) {
			if err.Tag() == "Missing" {
				return fmt.Errorf("Failed to validate che-inspector config file.  Field '%s' is missing from config file. Structure field: '%s'.", err.Field(), err.StructNamespace())
			}
			if err.Tag() == "EmptyTests" {
				return fmt.Errorf("Failed to validate che-inspector config file. No test specified.  Field '%s' is empty. Structure field: '%s'.", err.Field(), err.StructNamespace())
			}
			if err.Tag() == "testSuite" {
				return fmt.Errorf("Failed to validate che-inspector config file. Test field  '%s' is empty.", err.Field())
			}
		}
	}
	return nil
}

// Basic validator for cheInspector fields
func validateCheInspectorStruct(sl validator.StructLevel) {
	inspector := sl.Current().Interface().(api.CheInspector)

	if inspector.Name == "" {
		sl.ReportError(inspector.Name, "name", "Name", "Missing", "")
	}

	if inspector.Spec.Deployment.Cli.Namespace == "" {
		sl.ReportError(inspector.Spec.Deployment.Cli.Namespace, "namespace", "Spec.Deployment.Cli.Namespace", "Missing", "")
	}

	if inspector.Version == "" {
		sl.ReportError(inspector.Version, "name", "Name", "Missing", "")
	}

	if len(inspector.Spec.Tests) == 0 {
		sl.ReportError(inspector.Version, "tests", "Tests", "EmptyTests", "")
	}

	for _, testSuite := range inspector.Spec.Tests {
		validateTestsStruct(testSuite, sl)
	}
}

// Validate if yaml tests contain a name and a namespace
func validateTestsStruct(testSuite api.CheTestsSpec, sl validator.StructLevel) {
	if testSuite.Name == "" {
		sl.ReportError(testSuite.Name, "name", "Name", "testSuite", "")
	}
	if testSuite.Namespace == "" {
		sl.ReportError(testSuite.Name, "namespace", "Namespace", "testSuite", "")
	}
}

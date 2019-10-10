package schema

import (
	"testing"

	"github.com/elimity-com/scim/errors"
	"github.com/elimity-com/scim/optional"
)

func TestInvalidAttributeName(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("did not panic")
		}
	}()

	_ = Schema{
		ID:          "urn:ietf:params:scim:schemas:core:2.0:User",
		Name:        "User",
		Description: optional.NewString("User Account"),
		Attributes: []CoreAttribute{
			SimpleCoreAttribute(SimpleStringParams(StringParams{Name: "_Invalid"})),
		},
	}
}

var testSchema = Schema{
	ID:          "empty",
	Name:        "test",
	Description: optional.String{},
	Attributes: []CoreAttribute{
		SimpleCoreAttribute(SimpleStringParams(StringParams{
			Name:     "required",
			Required: true,
		})),
		SimpleCoreAttribute(SimpleBooleanParams(BooleanParams{
			MultiValued: true,
			Name:        "booleans",
			Required:    true,
		})),
		ComplexCoreAttribute(ComplexParams{
			MultiValued: true,
			Name:        "complex",
			SubAttributes: []SimpleParams{
				SimpleStringParams(StringParams{Name: "sub"}),
			},
		}),

		SimpleCoreAttribute(SimpleBinaryParams(BinaryParams{
			Name: "binary",
		})),
		SimpleCoreAttribute(SimpleDateTimeParams(DateTimeParams{
			Name: "dateTime",
		})),
		SimpleCoreAttribute(SimpleReferenceParams(ReferenceParams{
			Name: "reference",
		})),
		SimpleCoreAttribute(SimpleNumberParams(NumberParams{
			Name: "integer",
			Type: AttributeTypeInteger(),
		})),
		SimpleCoreAttribute(SimpleNumberParams(NumberParams{
			Name: "decimal",
			Type: AttributeTypeDecimal(),
		})),
	},
}

func TestResourceInvalid(t *testing.T) {
	var resource interface{}
	if _, scimErr := testSchema.Validate(resource); scimErr == errors.ValidationErrorNil {
		t.Error("invalid resource expected")
	}
}

func TestValidationInvalid(t *testing.T) {
	for _, test := range []map[string]interface{}{
		{ // missing required field
			"field": "present",
			"booleans": []interface{}{
				true,
			},
		},
		{ // missing required multivalued field
			"required": "present",
			"booleans": []interface{}{},
		},
		{ // wrong type element of slice
			"required": "present",
			"booleans": []interface{}{
				"present",
			},
		},
		{ // duplicate names
			"required": "present",
			"Required": "present",
			"booleans": []interface{}{
				true,
			},
		},
		{ // wrong string type
			"required": true,
			"booleans": []interface{}{
				true,
			},
		},
		{ // wrong complex type
			"required": "present",
			"complex":  "present",
			"booleans": []interface{}{
				true,
			},
		},
		{ // wrong complex element type
			"required": "present",
			"booleans": []interface{}{
				true,
			},
			"complex": []interface{}{
				"present",
			},
		},
		{ // duplicate complex element names
			"required": "present",
			"booleans": []interface{}{
				true,
			},
			"complex": []interface{}{
				map[string]interface{}{
					"sub": "present",
					"Sub": "present",
				},
			},
		},
		{ // wrong type complex element
			"required": "present",
			"booleans": []interface{}{
				true,
			},
			"complex": []interface{}{
				map[string]interface{}{
					"sub": true,
				},
			},
		},
		{ // invalid type binary
			"required": "present",
			"booleans": []interface{}{
				true,
			},
			"binary": true,
		},
		{ // invalid type dateTime
			"required": "present",
			"booleans": []interface{}{
				true,
			},
			"dateTime": "04:56:22Z2008-01-23T",
		},
		{ // invalid type integer
			"required": "present",
			"booleans": []interface{}{
				true,
			},
			"integer": 1.1,
		},
		{ // invalid type decimal
			"required": "present",
			"booleans": []interface{}{
				true,
			},
			"decimal": "1.1",
		},
	} {
		if _, scimErr := testSchema.Validate(test); scimErr == errors.ValidationErrorNil {
			t.Errorf("invalid resource expected")
		}
	}
}

func TestValidValidation(t *testing.T) {
	for _, test := range []map[string]interface{}{
		{
			"required": "present",
			"booleans": []interface{}{
				true,
			},
			"complex": []interface{}{
				map[string]interface{}{
					"sub": "present",
				},
			},
			"binary":   "ZXhhbXBsZQ==",
			"dateTime": "2008-01-23T04:56:22Z",
			"integer":  11,
			"decimal":  -2.1e5,
		},
	} {
		if _, scimErr := testSchema.Validate(test); scimErr != errors.ValidationErrorNil {
			t.Errorf("valid resource expected")
		}
	}
}
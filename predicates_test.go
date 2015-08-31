package pbc

import "testing"

func TestIsStringTrivial(t *testing.T) {
	if warn := IsString("Hello"); len(warn) != 0 {
		t.Errorf("IsString('Hello') returns an error")
	}
}

func TestIsStringInt(t *testing.T) {
	warn := IsString(1)

	if len(warn) == 0 {
		t.Error("IsString(1) returns no error")
	}

	if warn[0] != "expected string" {
		t.Error("IsString(1) does not return the right error")
	}
}

func TestArrayTrivial(t *testing.T) {
	warn := ArrayOf(IsString)([]interface{}{
		"Hello", "World",
	})

	if len(warn) != 0 {
		t.Errorf("IsArray(IsString)(['hello', 'world']) returns warnings %v", warn)
	}
}

func TestArrayOneError(t *testing.T) {
	warn := ArrayOf(IsString)([]interface{}{
		"Hello", 1,
	})

	switch len(warn) {
	case 0:
		t.Errorf("IsArray with on failing case has zero warnings")
	case 1:
		if warn[0] != "in array expected string" {
			t.Errorf("IsArray with one failing case: expected 'in array expected string' but got '%v'", warn[0])
		}
	default:
		t.Errorf("IsArray with one failing case has more then one warning: '%v'", warn)

	}
}

func TestArrayTwoErrors(t *testing.T) {
	warn := ArrayOf(IsString)([]interface{}{
		1, 2,
	})

	switch len(warn) {
	case 2:
		if warn[0] != "in array expected string" {
			t.Errorf("IsArray with two failing cases: expected warn[0] to be 'in array expected string' but got '%v'", warn[0])
		}

		if warn[1] != "in array expected string" {
			t.Errorf("IsArray with two failing cases: expected warn[1] to be 'in array expected string' but got '%v'", warn[1])
		}

	default:
		t.Errorf("IsArray with two failing cases: mismatch cases '%v'", warn)
	}
}

func TestObjectTrivial(t *testing.T) {
	warn := Object(map[string]TypeCheck{})(map[string]interface{}{})

	if len(warn) != 0 {
		t.Errorf("expected no warnings: '%v'", warn)
	}
}

func TestObjectOneError(t *testing.T) {
	warn := Object(map[string]TypeCheck{
		"uuid": IsNumber,
	})(map[string]interface{}{
		"uuid": "1",
	})

	if len(warn) == 1 {
		if warn[0] != "key:'uuid' expected number" {
			t.Errorf("Expected 'key:'uuid' expected int' but got '%v'", warn[0])
		}
	} else {
		t.Errorf("Expected 1 warning: '%v", warn)
	}
}

func TestObjectMissing(t *testing.T) {
	check := Object(map[string]TypeCheck{
		"uuid": IsNumber,
	})

	warn := check(map[string]interface{}{})

	if len(warn) != 0 {
		t.Errorf("Expected no warning: %v", warn)
	}
}

func TestRequiredTrivial(t *testing.T) {
	warn := Required("uuid")(map[string]interface{}{
		"uuid": 1,
	})

	if len(warn) != 0 {
		t.Errorf("Expected no warnings: %v", warn)
	}
}

func TestRequiredMissingKey(t *testing.T) {
	warn := Required("uuid")(map[string]interface{}{})

	if len(warn) == 1 {
		expected := "missing key 'uuid'"
		if warn[0] != expected {
			t.Errorf("Expected '%v' but got '%v'", expected, warn[0])
		}
	} else {
		t.Errorf("Expected 1 warning: %v", warn)
	}
}

func TestRequiredMissingKeys(t *testing.T) {
	warn := Required("uuid", "description")(map[string]interface{}{})
	if len(warn) == 2 {

		if warn[0] != "missing key 'uuid'" {
			t.Errorf("Expected 'missing key 'uuid'' but got '%v'", warn[0])
		}

		if warn[1] != "missing key 'description'" {
			t.Errorf("Expected 'missing key 'description'' but got '%v'", warn[1])
		}
	} else {
		t.Errorf("Wrong number of warnings: %v", warn)
	}
}

func TestEitherTrivial(t *testing.T) {
	if warn := Either(IsString, IsNumber)(1.0); len(warn) != 0 {
		t.Errorf("Expected no warnings: %v", warn)
	}

	if warn := Either(IsString, IsNumber)("Hello"); len(warn) != 0 {
		t.Errorf("Expected no warnings: %v", warn)
	}
}

func TestEitherError(t *testing.T) {
	warn := Either(IsString, ArrayOf(IsString))(map[string]interface{}{})

	if len(warn) == 1 {
		expected := "(expected string,expected array)"
		if warn[0] != expected {
			t.Errorf("Expected '%v', but got '%v'", expected, warn[0])
		}
	} else {
		t.Errorf("Wrong number of warnings: %v", warn)
	}

}

func TestAndTrivial(t *testing.T) {
	check := And(Required("uuid"), Object(map[string]TypeCheck{"uuid": IsNumber}))

	warn := check(map[string]interface{}{"uuid": 1.0})

	if len(warn) != 0 {
		t.Errorf("Expected no warnings: %v", warn)
	}
}

func TestAndMissingKey(t *testing.T) {
	check := And(Required("uuid"), Object(map[string]TypeCheck{"uuid": IsNumber}))

	warn := check(map[string]interface{}{})

	if len(warn) == 1 {
		expected := "missing key 'uuid'"
		if warn[0] != expected {
			t.Errorf("Expected '%v', but got '%v'", expected, warn[0])
		}
	} else {
		t.Errorf("Expected one warning: %v", warn)
	}
}

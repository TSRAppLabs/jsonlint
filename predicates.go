package pbc

import (
	"fmt"
	"strings"
)

/*
  A TypeCheck will return many warnings about what it's given
*/
type TypeCheck func(interface{}) Warning

/*
  Returns warning if value is not a string
*/
func IsString(value interface{}) Warning {
	_, isStr := value.(string)
	return ifError(!isStr, "expected string")
}

/*
  Returns warning if value is not a int
*/
func IsNumber(value interface{}) Warning {
	_, isNumber := value.(float64)
	return ifError(!isNumber, "expected number")
}

/*
  Returns warning if value is a bool
*/
func IsBool(value interface{}) Warning {
	_, isBool := value.(bool)
	return ifError(!isBool, "expected bool")
}

/*
  Returns warning if value is a float
*/
func IsDouble(value interface{}) Warning {
	_, isDouble := value.(float64)
	return ifError(!isDouble, "expected double")
}

/*
  Returns a check for a value to be a string and is contained in the input values
*/
func StringEnum(values ...string) TypeCheck {
	return func(val interface{}) Warning {
		str, isStr := val.(string)

		if !isStr {
			return NewWarning("expected string")
		}

		for _, v := range values {
			if v == str {
				return []string{}
			}
		}

		return NewWarning("expected one of %v", values)
	}
}

/*
  Constructs a type check which will return a Warning if value is not an array
  Or have values specified in checks TypeCheck
*/
func ArrayOf(check TypeCheck) TypeCheck {
	return func(value interface{}) Warning {
		arr, valid := value.([]interface{})

		if !valid {
			return NewWarning("expected array")
		}

		results := []string{}

		for _, e := range arr {
			warn := check(e)

			if warn != nil {
				for _, msg := range warn {
					results = append(results, fmt.Sprintf("in array %v", msg))
				}

			}
		}

		return results
	}
}

/*
  Checks the keys of an object with the corresponding functions in the passed map.
  if the key exists in the passed map
*/
func Object(kv map[string]TypeCheck) TypeCheck {
	return func(val interface{}) Warning {
		obj, isObj := val.(map[string]interface{})

		if !isObj {
			return []string{"expected obj"}
		}

		result := []string{}

		for key, check := range kv {
			val, exists := obj[key]

			if !exists {
				continue
			}

			for _, msg := range check(val) {
				result = append(result, fmt.Sprintf("key:'%v' %v", key, msg))
			}
		}

		return result
	}
}

/*
  Returns that the checked value is a object, and that it has the specified keys
*/
func Required(keys ...string) TypeCheck {
	return func(val interface{}) Warning {
		obj, isObj := val.(map[string]interface{})

		if !isObj {
			return NewWarning("expected obj")
		}

		results := []string{}

		for _, key := range keys {
			_, exists := obj[key]

			if !exists {
				results = append(results, fmt.Sprintf("missing key '%v'", key))
			}
		}

		return results

	}
}

/*
  Returns that the checked value has at least one of the input keys
*/
func WhiteList(keys ...string) TypeCheck {
	return func(val interface{}) Warning {
		obj, isObj := val.(map[string]interface{})

		if !isObj {
			return NewWarning("expected obj")
		}

		warn := []string{}
		for k, _ := range obj {
			found := false

			for _, key := range keys {
				if key == k {
					found = true
					break
				}
			}

			if !found {
				warn = append(warn, NewWarning("unexpected key '%v'", k)...)
			}

		}

		return warn
	}
}

/*
  Returns that the checked value is an object, and that has one or none of the input keys
*/
func Mutex(keys ...string) TypeCheck {
	return func(val interface{}) Warning {
		obj, isObj := val.(map[string]interface{})

		if !isObj {
			return NewWarning("expected obj")
		}

		keysFound := []string{}

		for _, key := range keys {
			if _, found := obj[key]; found {
				keysFound = append(keysFound, key)
			}
		}

		if len(keysFound) > 1 {
			return NewWarning("mutually exclusive keys found %v", strings.Join(keysFound, ","))
		}

		return []string{}
	}
}

/*
  Returns a checker, which only needs to satisfy one the the input checks
*/
func Either(checks ...TypeCheck) TypeCheck {
	return func(val interface{}) Warning {
		results := []string{}

		for _, check := range checks {
			warn := check(val)

			if len(warn) == 0 {
				return []string{}
			}

			for _, msg := range warn {
				results = append(results, msg)
			}
		}
		return NewWarning("(%v)", strings.Join(results, ","))
	}
}

/*
  Returns a checker, which must satisfy all input checkers
*/
func And(checks ...TypeCheck) TypeCheck {
	return func(val interface{}) Warning {
		results := []string{}
		for _, check := range checks {
			results = append(results, check(val)...)
		}

		return results
	}
}

func ifError(cond bool, msg string) Warning {
	if cond {
		return NewWarning(msg)
	} else {
		return []string{}
	}
}

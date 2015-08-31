package pbc

import (
	"fmt"
)

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
func IsInt(value interface{}) Warning {
	_, isInt := value.(int)
	return ifError(!isInt, "expected int")
}

/*
  Constructs a type check which will return a Warning if value is not an array
  Or have values specified in checks TypeCheck
*/
func ArrayOf(check TypeCheck) TypeCheck {
	return func(value interface{}) Warning {
		arr, valid := value.([]interface{})

		if !valid {
			return Warning("expected array")
		}

		for _, e := range arr {
			warn := check(e)

			if warn != NilWarning {
				return warn
			}
		}

		return NilWarning
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
			return Warning("expected obj")
		}

		for key, check := range kv {
			val, exists := obj[key]

			if !exists {
				continue
			}

			if warn := check(val); warn != NilWarning {
				return warn
			}

		}

		return NilWarning
	}
}

func Required(keys ...string) TypeCheck {
	return func(val interface{}) Warning {
		obj, isObj := val.(map[string]interface{})

		if isObj {
			for _, key := range keys {
				_, exists := obj[key]

				if !exists {
					return NewWarning("missing key '%v'", key)
				}
			}

			return NilWarning
		}

		return Warning("expected obj")
	}
}

func Either(checks ...TypeCheck) TypeCheck {
	return func(val interface{}) Warning {
		for _, check := range checks {
			warn := check(val)
		}

		return NilWarning
	}
}

func ifError(cond bool, msg string) Warning {
	if cond {
		return Warning(msg)
	} else {
		return ""
	}
}

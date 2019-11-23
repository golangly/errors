package errors

func Types(err error) []string {
	if err == nil {
		return nil
	}
	switch t := err.(type) {
	case *fundamental:
		types := make([]string, len(t.types))
		for i := 0; i < len(t.types); i++ {
			types[i] = t.types[i]
		}
		return types
	case *withMessage:
		types := Types(t.cause)
		return append(types, t.types...)
	case *withStack:
		types := Types(t.error)
		return append(types, t.types...)
	default:
		return []string{}
	}
}

// HasType is a helper function that will recurse up from the root error and check that the provided type
// is present using an equality check
func HasType(err error, typ string) bool {
	if err == nil {
		return false
	}
	switch t := err.(type) {
	case *fundamental:
		for _, tt := range t.types {
			if tt == typ {
				return true
			}
		}
		return false
	case *withMessage:
		for _, tt := range t.types {
			if tt == typ {
				return true
			}
		}
		return HasType(t.cause, typ)
	case *withStack:
		for _, tt := range t.types {
			if tt == typ {
				return true
			}
		}
		return HasType(t.error, typ)
	default:
		return false
	}
}

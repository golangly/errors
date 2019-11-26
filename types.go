package errors

func Types(err error) []string {
	if err == nil {
		return nil
	} else if w, ok := err.(*wrapper); ok {
		types := Types(w.cause)
		return append(types, w.types...)
	} else {
		return []string{}
	}
}

// HasType is a helper function that will recurse up from the root error and check that the provided type
// is present using an equality check
func HasType(err error, typ string) bool {
	if err == nil {
		return false
	} else if w, ok := err.(*wrapper); ok {
		for _, tt := range w.types {
			if tt == typ {
				return true
			}
		}
		return HasType(w.cause, typ)
	} else {
		return false
	}
}

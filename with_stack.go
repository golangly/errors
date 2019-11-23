package errors

import (
	"fmt"
	"io"
	"strings"
)

type withStack struct {
	error
	*stack
	types []string
	tags  []Tag
}

// Return the causing (wrapped) error of this error.
func (w *withStack) Cause() error { return w.error }

// Unwrap provides compatibility for Go 1.13 error chains.
func (w *withStack) Unwrap() error { return w.error }

// Formats this error for printf strings.
func (w *withStack) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			_, _ = fmt.Fprintf(s, "%+v", w.Cause())
			for _, tag := range w.tags {
				_, _ = fmt.Fprintf(s, " %s=%v", tag.Key, tag.Value)
			}
			if len(w.types) > 0 {
				_, _ = fmt.Fprintf(s, " types=%s", strings.Join(w.types, ","))
			}
			w.stack.Format(s, verb)
			return
		}
		fallthrough
	case 's':
		_, _ = io.WriteString(s, w.Error())
	case 'q':
		_, _ = fmt.Fprintf(s, "%q", w.Error())
	}
}

// Adds the given tag.
func (w *withStack) AddTag(key string, value interface{}) ErrorExt {
	return w.AddTags(T(key, value))
}

// Adds the given tags.
func (w *withStack) AddTags(tags ...Tag) ErrorExt {
	w.tags = append(w.tags, tags...)
	return w
}

// Adds the given types.
func (w *withStack) AddTypes(types ...string) ErrorExt {
	w.types = append(w.types, types...)
	return w
}

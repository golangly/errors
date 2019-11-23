package errors

// T is a shortcut to make a Tag
func T(key string, value interface{}) Tag {
	return Tag{Key: key, Value: value}
}

// Tag contains a single key value combination
// to be attached to your error
type Tag struct {
	Key   string
	Value interface{}
}

// Returns all tags for the given error.
func Tags(err error) map[string]interface{} {
	if err == nil {
		return nil
	}
	tags := make(map[string]interface{})
	collectTags(err, tags)
	return tags
}

// LookupTag recursively searches for the provided tag and returns it's value or nil
func LookupTag(err error, key string) interface{} {
	if err == nil {
		return nil
	}
	switch t := err.(type) {
	case *fundamental:
		for _, tag := range t.tags {
			if tag.Key == key {
				return tag.Value
			}
		}
		return nil
	case *withMessage:
		for _, tag := range t.tags {
			if tag.Key == key {
				return tag.Value
			}
		}
		return LookupTag(t.cause, key)
	case *withStack:
		for _, tag := range t.tags {
			if tag.Key == key {
				return tag.Value
			}
		}
		return LookupTag(t.error, key)
	default:
		return nil
	}
}

func collectTags(err error, tags map[string]interface{}) {
	if err == nil {
		return
	}
	switch t := err.(type) {
	case *fundamental:
		for _, tag := range t.tags {
			tags[tag.Key] = tag.Value
		}
	case *withMessage:
		collectTags(t.cause, tags)
		for _, tag := range t.tags {
			tags[tag.Key] = tag.Value
		}
	case *withStack:
		collectTags(t.error, tags)
		for _, tag := range t.tags {
			tags[tag.Key] = tag.Value
		}
	}
}

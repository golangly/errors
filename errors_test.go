package errors

import (
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	err := New("test")
	assert.Nil(t, errors.Unwrap(err))
	assert.Nil(t, err.(*wrapper).Cause())
	assert.Equal(t, "test", err.Error())
	assert.Equal(t, 0, len(Tags(err)))
	assert.Equal(t, 0, len(Types(err)))
}

func TestNewf(t *testing.T) {
	err := Newf("test %s %d", "s1", 99)
	assert.Nil(t, errors.Unwrap(err))
	assert.Nil(t, err.(*wrapper).Cause())
	assert.Equal(t, "test s1 99", err.Error())
	assert.Equal(t, 0, len(Tags(err)))
	assert.Equal(t, 0, len(Types(err)))
}

func TestWrap(t *testing.T) {
	cause := errors.New("cause")
	err := Wrap(cause, "test")
	assert.Same(t, cause, errors.Unwrap(err))
	assert.Same(t, cause, err.(*wrapper).Cause())
	assert.Equal(t, "test: cause", err.Error())
	assert.Equal(t, 0, len(Tags(err)))
	assert.Equal(t, 0, len(Types(err)))
}

func TestWrapf(t *testing.T) {
	cause := errors.New("cause")
	err := Wrapf(cause, "test %s %d", "s1", 99)
	assert.Same(t, cause, errors.Unwrap(err))
	assert.Same(t, cause, err.(*wrapper).Cause())
	assert.Equal(t, "test s1 99: cause", err.Error())
	assert.Equal(t, 0, len(Tags(err)))
	assert.Equal(t, 0, len(Types(err)))
}

func TestUnwrap(t *testing.T) {
	rootCause := errors.New("root")
	interimCause := Wrap(rootCause, "interim")
	err := Wrap(interimCause, "test")
	assert.Same(t, interimCause, Unwrap(err))
	assert.Same(t, rootCause, Unwrap(Unwrap(err)))
}

func TestRootCause(t *testing.T) {
	rootCause := errors.New("root")
	interimCause := Wrap(rootCause, "interim")
	err := Wrap(interimCause, "test")
	assert.Same(t, rootCause, RootCause(err))
}

func TestFormat(t *testing.T) {
	stackTrace := fmt.Sprintf("%+v", callers(2))
	stackTraceRe := regexp.MustCompile("errors_test.go:\\d+").ReplaceAllString(stackTrace, "errors_test.go:\\d+")
	rootCause := errors.New("root")
	interimCause := Wrap(rootCause, "interim")
	err := Wrap(interimCause, "test")
	assert.Equal(t, "test: interim: root", fmt.Sprintf("%s", err))
	assert.Equal(t, "test: interim: root", fmt.Sprintf("%s", err))
	assert.Equal(t, "\"test: interim: root\"", fmt.Sprintf("%q", err))
	assert.Equal(t, "test: interim: root", fmt.Sprintf("%v", err))
	var expected string

	expected = ""
	expected = expected + "test: interim: root"
	expected = expected + stackTraceRe + "\n"
	expected = expected + "interim: root"
	expected = expected + stackTraceRe + "\n"
	expected = expected + "root"
	assert.Regexp(t, expected, fmt.Sprintf("%+v", err))

	expected = ""
	expected = expected + "interim: root"
	expected = expected + stackTraceRe + "\n"
	expected = expected + "root"
	assert.Regexp(t, expected, fmt.Sprintf("%+v", interimCause))
}

func TestTags(t *testing.T) {
	rootCause := errors.New("root")
	interimCause1 := Wrap(rootCause, "interim1").AddTag("tag1", "value1")
	interimCause2 := Wrap(interimCause1, "interim2").AddTag("tag2", "value2")
	err := Wrap(interimCause2, "error")
	assert.Equal(t, map[string]interface{}{"tag1": "value1", "tag2": "value2"}, Tags(err))
	assert.Equal(t, map[string]interface{}{"tag1": "value1", "tag2": "value2"}, Tags(interimCause2))
	assert.Equal(t, map[string]interface{}{"tag1": "value1"}, Tags(interimCause1))
	assert.Equal(t, map[string]interface{}{}, Tags(rootCause))
	assert.Equal(t, nil, LookupTag(rootCause, "tag1"))
	assert.Equal(t, "value1", LookupTag(interimCause1, "tag1"))
	assert.Equal(t, "value1", LookupTag(interimCause2, "tag1"))
	assert.Equal(t, "value1", LookupTag(err, "tag1"))
	assert.Equal(t, nil, LookupTag(rootCause, "tag2"))
	assert.Equal(t, nil, LookupTag(interimCause1, "tag2"))
	assert.Equal(t, "value2", LookupTag(interimCause2, "tag2"))
	assert.Equal(t, "value2", LookupTag(err, "tag2"))
}

func TestTypes(t *testing.T) {
	rootCause := errors.New("root")
	interimCause1 := Wrap(rootCause, "interim1").AddTypes("t1", "t11")
	interimCause2 := Wrap(interimCause1, "interim2").AddTypes("t2", "t22")
	err := Wrap(interimCause2, "error")
	assert.Equal(t, []string{"t1", "t11", "t2", "t22"}, Types(err))
	assert.Equal(t, []string{"t1", "t11", "t2", "t22"}, Types(interimCause2))
	assert.Equal(t, []string{"t1", "t11"}, Types(interimCause1))
	assert.Equal(t, []string{}, Types(rootCause))
	assert.True(t, HasType(err, "t1"))
	assert.True(t, HasType(err, "t11"))
	assert.True(t, HasType(err, "t2"))
	assert.True(t, HasType(err, "t22"))
	assert.True(t, HasType(interimCause2, "t1"))
	assert.True(t, HasType(interimCause2, "t11"))
	assert.True(t, HasType(interimCause2, "t2"))
	assert.True(t, HasType(interimCause2, "t22"))
	assert.True(t, HasType(interimCause1, "t1"))
	assert.True(t, HasType(interimCause1, "t11"))
	assert.False(t, HasType(interimCause1, "t2"))
	assert.False(t, HasType(interimCause1, "t22"))
	assert.False(t, HasType(rootCause, "t1"))
	assert.False(t, HasType(rootCause, "t11"))
	assert.False(t, HasType(rootCause, "t2"))
	assert.False(t, HasType(rootCause, "t22"))
}

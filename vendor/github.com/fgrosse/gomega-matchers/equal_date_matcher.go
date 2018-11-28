package matchers

import (
	"fmt"
	"time"

	"github.com/onsi/gomega/format"
	"github.com/onsi/gomega/types"
)

// EqualTime succeeds if actual is a time.Time or a RFC3339 formatted string
// and equals expected.
// expected can either be a time.Time or a string in RFC3339.
func EqualTime(expected interface{}) types.GomegaMatcher {
	m := &equalTimeMatcher{}

	var err error
	m.expected, err = toTime(expected)
	if err != nil {
		panic("EqualTime: invalid expectation: " + err.Error())
	}

	return m
}

type equalTimeMatcher struct{ expected time.Time }

func (m *equalTimeMatcher) Match(actual interface{}) (success bool, err error) {
	actualTime, err := toTime(actual)

	if actualTime.IsZero() && m.expected.IsZero() {
		// in this case we do not care about the location
		return true, nil
	}

	return actualTime == m.expected, err
}

func (m *equalTimeMatcher) FailureMessage(actual interface{}) (message string) {
	return m.format(actual, "to equal", m.expected)
}

func (m *equalTimeMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return m.format(actual, "not to equal", m.expected)
}

func (m *equalTimeMatcher) format(actual interface{}, message string, expected interface{}) string {
	return fmt.Sprintf("Expected\n%s\n%s\n%s", m.formatObject(actual), message, m.formatObject(expected))
}

func (m *equalTimeMatcher) formatObject(i interface{}) string {
	t, err := toTime(i)
	if err != nil {
		panic(err)
	}

	value := t.Format(time.RFC3339)
	if t.IsZero() {
		value = "zero time"
	}

	return format.Indent + fmt.Sprintf("<%T> %s", i, value)
}

func toTime(i interface{}) (time.Time, error) {
	switch t := i.(type) {
	case time.Time:
		return t, nil
	case string:
		ti, err := time.Parse(time.RFC3339, t)
		if err != nil {
			err = fmt.Errorf("could not parse value as RFC3339 formatted time string: %s", err)
		}
		return ti, err
	default:
		return time.Time{}, fmt.Errorf("value must either be a time.Time or a RFC3339 formatted string but is as %T", i)
	}
}

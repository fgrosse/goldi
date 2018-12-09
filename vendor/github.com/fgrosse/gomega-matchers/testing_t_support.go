package matchers

import (
	"regexp"
	"runtime/debug"
	"strings"

	"github.com/mgutz/ansi"
	"github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
)

var StacktraceColor = ansi.ColorCode("gray")

// RegisterExtendedTestingT is like gomega.RegisterTestingT but it also
// colors the output and omits any gomega frames from the printed stacktrace.
func RegisterExtendedTestingT(t types.GomegaTestingT) {
	gomega.RegisterFailHandler(buildTestingTGomegaFailHandler(t))
}

func buildTestingTGomegaFailHandler(t types.GomegaTestingT) types.GomegaFailHandler {
	return func(message string, callerSkip ...int) {
		skip := 2 // initial runtime/debug.Stack frame + gomega-matchers.buildTestingTGomegaFailHandler frame
		if len(callerSkip) > 0 {
			skip += callerSkip[0]
		}
		stackTrace := pruneStack(string(debug.Stack()), skip)
		stackTrace = strings.TrimSpace(stackTrace)
		stackTrace = ansi.Color(stackTrace, StacktraceColor)

		t.Fatalf("\n%s\n%s", stackTrace, ansi.Color(message, "red"))
	}
}

func pruneStack(fullStackTrace string, skip int) string {
	stack := strings.Split(fullStackTrace, "\n")
	if len(stack) > 1+2*skip {
		stack = stack[1+2*skip:]
	}

	srcBlacklist := []string{
		"testing.tRunner",
		"created by testing.RunTests",
	}
	srcBlacklistRE := regexp.MustCompile(strings.Join(srcBlacklist, "|"))

	suffix := regexp.MustCompile(` \+0x[0-9a-f]+$`)
	trim := func(s string) string {
		return suffix.ReplaceAllString(s, "")
	}

	prunedStack := []string{}
	for i := 0; i < len(stack)/2; i++ {
		if srcBlacklistRE.Match([]byte(stack[i*2])) {
			continue
		}

		prunedStack = append(prunedStack, trim(stack[i*2]))
		prunedStack = append(prunedStack, trim(stack[i*2+1]))
	}

	return strings.Join(prunedStack, "\n")
}

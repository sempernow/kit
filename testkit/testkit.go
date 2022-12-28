// Package testkit provides testing-helper functions that are usable at any layer.
package testkit

import (
	"context"
	"encoding/json"
	"os/exec"
	"testing"
	"time"

	"github.com/sempernow/kit/web"

	"github.com/gofrs/uuid" // FORKED from github.com/satori
	"github.com/google/go-cmp/cmp"
)

// ============================================================================
//  HELPERs for testing at any layer

// Context returns an app level context for testing.
func Context() context.Context {
	values := web.Values{
		TraceID: uuid.Must(uuid.NewV4()).String(),
		Now:     time.Now(),
	}

	return context.WithValue(context.Background(), web.Key1, &values)
}

// StringPointer is a helper to get a *string from a string. It is in the tests
// package because we normally don't want to deal with pointers to basic types
// but it's useful in some tests; pointer types are useful for optional/NULLable keys.
func StringPointer(s string) *string {
	return &s
}

// IntPointer is a helper to get a *int from a int. It is in the tests package
// because we normally don't want to deal with pointers to basic types but it's
// useful in some tests; pointer types are useful for optional/NULLable keys.
func IntPointer(i int) *int {
	return &i
}

// `Success` (✓) and `Failure` (✗) are discriptive Unicode markers.
const (
	Success = "\u2713" // ✓
	Failure = "\u2717" // ✗

	// HEAVY CHECK MARK \u2714 // ✔
	// HEAVY BALLOT X \u2718 // ✘
)

// Off to report when test is turned off.
func Off(t *testing.T) {
	t.Helper()
	t.Log("Off : NO TEST.")
}

// LogPrettyStruct pretty-prints any struct
func LogPrettyStruct(t *testing.T, x interface{}) {
	t.Helper()
	s, _ := json.MarshalIndent(x, "", "\t")
	t.Logf(string(s))
}

// LogJq ... pretty-prints json string per jq
func LogJq(t *testing.T, msg, json string) {
	t.Helper()
	cmd := "printf \"%s\" '" + json + "' | jq . "
	out, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		t.Fatalf("\t%s\t%s : %v", Failure, msg, err)
	}
	t.Logf("\t%s\t%s: %s", Success, msg, out)
}

// Jq pipes `json` string input to `jq . ` and returns result.
func Jq(t *testing.T, json string) string {
	cmd := "printf \"%s\" '" + json + "' | jq . "
	out, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		return ""
	}
	return string(out)
}

// Log is a testing-helper function.
func Log(t *testing.T, msg string, err interface{}) {
	t.Helper()
	if err != nil {
		t.Fatalf("\t%s\t%s : %v", Failure, msg, err)
	}
	t.Logf("\t%s\t%s.", Success, msg)
}

// LogCmp is a testing-helper function @ `got` v. `exp` cases, per `cmp.Diff(...)`.
func LogCmp(t *testing.T, msg string, got, exp interface{}) {
	t.Helper()
	if diff := cmp.Diff(exp, got); diff != "" {
		t.Fatalf("\t%s\t%s : %v", Failure, msg, diff)
	}
	t.Logf("\t%s\t%s.", Success, msg)
}

// LogDiff is a testing-helper function @ `got` v. `exp` cases.
func LogDiff(t *testing.T, msg string, got, exp interface{}) {
	t.Helper()
	if got != exp {
		t.Logf("\t%s\t%s:", Failure, msg)
		t.Logf("\t\t  Got: %v", got)
		t.Logf("\t\t  Exp: %v", exp)
		t.Fatal()
	}
	t.Logf("\t%s\t%s.", Success, msg)
}

// LogDiffPerField is a testing-helper function @ `got` v. `exp` cases, per `field`.
func LogDiffPerField(t *testing.T, field, msg, got, exp interface{}) {
	t.Helper()
	if got != exp {
		t.Logf("\t%s\t%s:", Failure, msg)
		t.Logf("\t\t  Got @ %s: %v", field, got)
		t.Logf("\t\t  Exp @ %s: %v", field, exp)
		t.Fatal()
	}
	t.Logf("\t%s\t%s @ `%s`.", Success, msg, field)
}

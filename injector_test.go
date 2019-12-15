package main

import (
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestUpdateHTML(t *testing.T) {
	html := `<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <meta http-equiv="X-UA-Compatible" content="ie=edge" />
    <title>Document</title>
  </head>
  <body>
    <!-- INJECT_ENV_START -->
    <!-- INJECT_ENV_END -->
  </body>
</html>`
	var sb strings.Builder

	value := "test string"

	updateHTML(value, strings.NewReader(html), &sb)
	got := sb.String()
	if strings.Index(got, value) == -1 {
		t.Errorf("could not find %q in %s", value, got)
	}
}

func TestInjectRaisesErrorIfNoMarkerFound(t *testing.T) {
	html := `<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <meta http-equiv="X-UA-Compatible" content="ie=edge" />
    <title>Document</title>
  </head>
  <body>
  </body>
</html>`
	var sb strings.Builder

	value := "test string"

	err := updateHTML(value, strings.NewReader(html), &sb)
	if err == nil {
		t.Errorf("expected an error")
	}
	switch err.(type) {
	case *NoMarkerFound:
		// good
	default:
		t.Errorf("got %s instead of *NoMarkerFound", err)
	}
}

func TestGetEnvVars(t *testing.T) {
	os.Setenv("TEST_VAR", "first")
	os.Setenv("SECOND_URL", "hi")
	defer os.Unsetenv("TEST_VAR")
	defer os.Unsetenv("SECOND_URL")

	whitelist := []string{"TEST_VAR", "SECOND_URL"}

	want := make(map[string]string)
	want["TEST_VAR"] = "first"
	want["SECOND_URL"] = "hi"

	got, err := getEnvVars(whitelist)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("expected %v, got %v", want, got)
	}
}

func TestGetEnvVarsThrowsErrorForMissingVar(t *testing.T) {
	os.Setenv("TEST_VAR", "first")
	defer os.Unsetenv("TEST_VAR")

	whitelist := []string{"TEST_VAR", "SECOND_URL"}

	want := "Env var \"SECOND_URL\" not found"
	_, err := getEnvVars(whitelist)
	switch err.(type) {
	case *ErrEnvVarNotFound:
		if err.Error() != want {
			t.Errorf("got %q, want %q", err.Error(), want)
		}
	default:
		t.Error("Expected ErrEnvVarNotFound")
	}
}

func TestWhitelistParser(t *testing.T) {
	var whitelistTests = []struct {
		in  string
		out []string
	}{
		{"ENV1,ENV2", []string{"ENV1", "ENV2"}},
		{"ENV1,ENV2,", []string{"ENV1", "ENV2"}},
		{"ENV1,,,ENV2", []string{"ENV1", "ENV2"}},
	}

	for _, tt := range whitelistTests {
		t.Run(tt.in, func(t *testing.T) {
			got := parseWhitelist(tt.in)
			if !reflect.DeepEqual(got, tt.out) {
				t.Errorf("got %v, want %v", got, tt.out)
			}
		})
	}
}

// func TestInject(t *testing.T) {
// 	html := `<!DOCTYPE html>
// 	<html lang="en">
// 	  <head>
// 		<meta charset="UTF-8" />
// 		<meta name="viewport" content="width=device-width, initial-scale=1.0" />
// 		<meta http-equiv="X-UA-Compatible" content="ie=edge" />
// 		<title>Document</title>
// 	  </head>
// 	  <body>
// 	  </body>
// 	</html>`

// 	var buf bytes.Buffer
// 	out := bufio.NewWriter(&buf)

// 	in := strings.NewReader(html)
// 	// updateHTML(value, strings.NewReader(html), &sb)
// 	inject("TESTVAR=hi", in, out)
// 	got := buf.String()
// 	want := `"TESTVAR": "hi"`
// 	fmt.Println(got)
// 	if strings.Index(got, want) == -1 {
// 		t.Errorf("could not find %q in %s", want, got)
// 	}
// }

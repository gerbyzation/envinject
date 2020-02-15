package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

type NoMarkerFound struct {
	message string
}

func NewNoMarkerFoundError(message string) *NoMarkerFound {
	return &NoMarkerFound{
		message: message,
	}
}

func (e *NoMarkerFound) Error() string {
	return e.message
}

type ErrEnvVarNotFound struct {
	message string
}

func NewErrEnvVarNotFound(message string) *ErrEnvVarNotFound {
	return &ErrEnvVarNotFound{
		message: message,
	}
}

func (e *ErrEnvVarNotFound) Error() string {
	return e.message
}

// Get all env vars mentioned in `whitelist` from the environment
func getEnvVars(whitelist []string) (map[string]string, error) {
	vars := make(map[string]string)
	for index := range whitelist {
		envKey := whitelist[index]
		val, ok := os.LookupEnv(envKey)
		if !ok {
			return nil, NewErrEnvVarNotFound(fmt.Sprintf("Env var %q not found", envKey))
		}
		vars[envKey] = val
	}
	return vars, nil
}

func updateHTML(jsonString string, htmlFile io.Reader, out io.Writer) error {
	buf := new(bytes.Buffer)
	buf.ReadFrom(htmlFile)
	html := buf.String()

	markers := strings.Count(html, "__ENV_INJECT__")
	if markers == 0 {
		return NewNoMarkerFoundError("Unable to find \"__ENV_INJECT__\" marker")
	} else if markers > 1 {
		return NewNoMarkerFoundError("Multiple markers \"__ENV_INJECT__\" found.")
	}
	html = strings.Replace(html, "__ENV_INJECT__", jsonString, 1)

	fmt.Fprint(out, html)

	return nil
}

// https://github.com/golang/go/wiki/SliceTricks#filter-in-place
func filter(a []string, keep func(string) bool) []string {
	n := 0
	for _, x := range a {
		if keep(x) {
			a[n] = x
			n++
		}
	}
	a = a[:n]
	return a
}

func parseWhitelist(whitelistString string) []string {
	names := strings.Split(whitelistString, ",")
	keep := func(name string) bool {
		if name != "" {
			return true
		}
		return false
	}
	filteredNames := filter(names, keep)
	return filteredNames
}

func inject(whitelistString string, html io.Reader, out io.Writer) error {
	whitelist := parseWhitelist(whitelistString)
	envVars, err := getEnvVars(whitelist)
	if err != nil {
		return err
	}
	jsonString, err := json.Marshal(envVars)
	if err != nil {
		return err
	}
	err = updateHTML(string(jsonString), html, out)
	if err != nil {
		return err
	}
	return nil
}

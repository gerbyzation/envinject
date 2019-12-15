package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/net/html"
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
	doc, err := html.Parse(htmlFile)
	if err != nil {
		return err
	}
	markerFound := false
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.CommentNode {
			text := strings.TrimSpace(n.Data)
			if text == "INJECT_ENV_END" {
				markerFound = true
				script := html.Node{
					Type: html.ElementNode,
					Data: "script",
				}
				scriptContent := html.Node{
					Type: html.TextNode,
					Data: fmt.Sprintf("window.ENVVARS = %s;", jsonString),
				}
				script.AppendChild(&scriptContent)
				n.Parent.InsertBefore(&script, n)
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	if markerFound == false {
		return NewNoMarkerFoundError("Unable to find \"INECT_ENV_END\" marker")
	}
	html.Render(out, doc)
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

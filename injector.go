package main

import (
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

// func buildJsonString(envVars map[string]string) (string, error) {
// 	jsonString, err := json.Marshal(envVars)
// 	return string(jsonString), err
// }

func Inject(jsonString string, htmlFile io.Reader, out io.Writer) error {
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

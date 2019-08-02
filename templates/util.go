package templates

import (
	"bytes"
	"html"
	"html/template"
)

func ParseTemplate(tpl *template.Template, data interface{}) (string, error) {
	var buf bytes.Buffer
	err := tpl.Execute(&buf, data)
	if err != nil {
		return "", err
	}

	return html.UnescapeString(buf.String()), nil
}

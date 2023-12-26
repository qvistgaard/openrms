package utils

import (
	"bytes"
	"text/template"
)

func ProcessTemplate(tpl string, vars any) (string, error) {
	tmpl, err := template.New("tmpl").Parse(tpl)

	if err != nil {
		return "", err
	}
	var tmplBytes bytes.Buffer

	err = tmpl.Execute(&tmplBytes, vars)
	if err != nil {
		return "", err
	}
	return tmplBytes.String(), err
}

package application

import (
	_ "embed"
	"html/template"
)

//go:embed index.html
var UIHtml string
var UITemplate *template.Template = template.Must(template.New("").Parse(UIHtml))

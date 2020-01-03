package main

import (
	"html/template"
	"os"
)

func main() {
	t1 := template.Must(template.New("login").Parse(`{{ define "body" }} Login Body {{ end }}`))
	t1 = template.Must(t1.Parse(`{{ define "layout" }}Layout: {{ block "body" . }}{{ end }}{{ end }}`))
	t1.ExecuteTemplate(os.Stderr, "layout", nil)
}

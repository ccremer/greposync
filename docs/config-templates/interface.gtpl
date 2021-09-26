=== {{ .Interface.Name }}
[source, go]
----
{{.Interface.Decl}} {
{{- range .Interface.Methods}}{{if or .Exported $.Config.Private }}
	{{.Decl}}{{end}}
{{- end}}
}
----

{{.Interface.Doc}}
{{range .Interface.Methods}}{{if or .Exported $.Config.Private }}
.{{.Name}}
[source, go]
----
func {{ .Decl }}
----
{{.Doc}}
{{end}}{{end}}
'''

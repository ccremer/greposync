{{- define "variable" -}}
{{ .Metadata.Repository.Name | upper }}
{{- end }}

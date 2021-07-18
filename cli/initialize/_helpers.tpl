{{- define "comment" -}}
{{- if .Values.comment.text }}
{{- with .Values.comment.open -}}
{{ . }}
{{ end -}}
{{- with .Values.comment.text -}}
{{ $.Values.comment.prefix }}{{ . | splitList "\n" | join (cat "\n" $.Values.comment.prefix | replace "\n " "\n") }}
{{ end }}
{{- with .Values.comment.closed -}}
{{ . }}
{{ end -}}
{{- end -}}
{{- end -}}

{{- define "comment" -}}
{{ .Values.comment | replace "# " (default "# " .Values.commentPrefix) }}
{{- end }}

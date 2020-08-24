{{- define "gke-app.envvarValue" -}}
{{- if or (kindIs "int" .) (kindIs "string" .) (kindIs "bool" .) (kindIs "float64" .) }}
value: {{ . | quote }}
{{- else }}
{{ tpl . $ }}
{{- end }}
{{- end -}}
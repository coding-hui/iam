{{/*
console fullname
*/}}
{{- define "iam.consoleFullname" -}}
{{ include "iam.fullname" . }}-console
{{- end }}

{{/*
console common labels
*/}}
{{- define "iam.consoleLabels" -}}
{{ include "iam.labels" . }}
app.kubernetes.io/component: console
{{- end }}

{{/*
console selector labels
*/}}
{{- define "iam.consoleSelectorLabels" -}}
{{ include "iam.selectorLabels" . }}
app.kubernetes.io/component: console
{{- end }}

{{/*
authzserver fullname
*/}}
{{- define "iam.authzServerFullname" -}}
{{ include "iam.fullname" . }}-authzserver
{{- end }}

{{/*
authzserver common labels
*/}}
{{- define "iam.authzServerLabels" -}}
{{ include "iam.labels" . }}
app.kubernetes.io/component: authzserver
{{- end }}

{{/*
authzserver selector labels
*/}}
{{- define "iam.authzServerSelectorLabels" -}}
{{ include "iam.selectorLabels" . }}
app.kubernetes.io/component: authzserver
{{- end }}

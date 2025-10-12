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

{{/*
Return the proper image name (for the init container image)
*/}}
{{- define "iam.consoleImage" -}}
{{- include "iam.images.image" (dict "imageRoot" .Values.console.image "global" .Values.global) }}
{{- end -}}

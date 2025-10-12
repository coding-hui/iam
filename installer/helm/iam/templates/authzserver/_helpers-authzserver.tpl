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

{{/*
Return the proper image name (for the init container image)
*/}}
{{- define "iam.authzServerImage" -}}
{{- include "iam.images.image" (dict "imageRoot" .Values.authzServer.image "global" .Values.global) }}
{{- end -}}

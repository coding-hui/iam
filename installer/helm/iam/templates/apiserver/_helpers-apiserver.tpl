{{/*
apiserver fullname
*/}}
{{- define "iam.apiServerFullname" -}}
{{ include "iam.fullname" . }}-apiserver
{{- end }}

{{/*
apiserver common labels
*/}}
{{- define "iam.apiServerLabels" -}}
{{ include "iam.labels" . }}
app.kubernetes.io/component: apiserver
{{- end }}

{{/*
apiserver selector labels
*/}}
{{- define "iam.apiServerSelectorLabels" -}}
{{ include "iam.selectorLabels" . }}
app.kubernetes.io/component: apiserver
{{- end }}

{{/*
Return the proper image name (for the init container apiserver image)
*/}}
{{- define "iam.apiServerImage" -}}
{{- include "iam.images.image" (dict "imageRoot" .Values.apiServer.image "global" .Values.global) }}
{{- end -}}
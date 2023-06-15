{{/*
Expand the name of the chart.
*/}}
{{- define "piper.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Return secret name to be used based on provided values.
*/}}
{{- define "piper.argoWorkflows.tokenSecretName" -}}
{{- $fullName := printf "%s-argo-token" .Release.Name -}}
{{- default $fullName .Values.piper.argoWorkflows.server.tokenExistingSecret | quote -}}
{{- end -}}

{{/*
Return secret name to be used based on provided values.
*/}}
{{- define "piper.gitProvider.tokenSecretName" -}}
{{- $fullName := printf "%s-git-token" .Release.Name -}}
{{- default $fullName .Values.piper.gitProvider.tokenExistingSecret | quote -}}
{{- end -}}

{{/*
Return secret name to be used based on provided values.
*/}}
{{- define "piper.gitProvider.webhook.secretName" -}}
{{- $fullName := printf "%s-webhook-secret" .Release.Name -}}
{{- default $fullName .Values.piper.gitProvider.webhook.existingSecret | quote -}}
{{- end -}}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "piper.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "piper.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "piper.labels" -}}
helm.sh/chart: {{ include "piper.chart" . }}
{{ include "piper.selectorLabels" . }}

{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "piper.selectorLabels" -}}
app.kubernetes.io/name: {{ include "piper.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "piper.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "piper.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

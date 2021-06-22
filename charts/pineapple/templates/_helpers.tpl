{{/*
Expand the name of the chart.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
*/}}
{{- define "pineapple.name" -}}
{{- default "pineapple" .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
*/}}
{{- define "pineapple.fullname" -}}
{{- $name := default "pineapple" .Values.nameOverride -}}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/* Create chart name and version as used by the chart label. */}}
{{- define "pineapple.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/* Common labels */}}
{{- define "pineapple.labels" -}}
helm.sh/chart: {{ include "pineapple.chart" . }}
{{ include "pineapple.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/* Selector Labels */}}
{{- define "pineapple.selectorLabels" -}}
app.kubernetes.io/name: {{ include "pineapple.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end -}}

{{/* Core component */}}
{{- define "pineapple.core.fullname" -}}
  {{- printf "%s-core" (include "pineapple.fullname" .) -}}
{{- end -}}

{{- define "pineapple.core.labels" }}
{{ include "pineapple.labels" . }}
app.kubernetes.io/component: core
{{- end -}}

{{- define "pineapple.core.selectorLabels" }}
{{ include "pineapple.selectorLabels" . }}
app.kubernetes.io/component: core
{{- end -}}


{{/* Portal component */}}
{{- define "pineapple.portal.fullname" -}}
  {{- printf "%s-portal" (include "pineapple.fullname" .) -}}
{{- end -}}

{{- define "pineapple.portal.labels" }}
{{ include "pineapple.labels" . }}
app.kubernetes.io/component: portal
{{- end -}}

{{- define "pineapple.portal.selectorLabels" }}
{{ include "pineapple.selectorLabels" . }}
app.kubernetes.io/component: portal
{{- end -}}


{{/* Webterminal component */}}
{{- define "pineapple.webterminal.fullname" -}}
  {{- printf "%s-webterminal" (include "pineapple.fullname" .) -}}
{{- end -}}

{{- define "pineapple.webterminal.labels" }}
{{ include "pineapple.labels" . }}
app.kubernetes.io/component: webterminal
{{- end -}}

{{- define "pineapple.webterminal.selectorLabels" }}
{{ include "pineapple.selectorLabels" . }}
app.kubernetes.io/component: webterminal
{{- end -}}


{{/* Billing component */}}
{{- define "pineapple.billing.fullname" -}}
  {{- printf "%s-billing" (include "pineapple.fullname" .) -}}
{{- end -}}

{{- define "pineapple.billing.labels" }}
{{ include "pineapple.labels" . }}
app.kubernetes.io/component: billing
{{- end -}}

{{- define "pineapple.billing.selectorLabels" }}
{{ include "pineapple.selectorLabels" . }}
app.kubernetes.io/component: billing
{{- end -}}

{{/* Webhook component */}}
{{- define "pineapple.webhook.fullname" }}
  {{- printf "%s-webhook" (include "pineapple.fullname" .) -}}
{{- end -}}

{{- define "pineapple.webhook.labels" }}
{{ include "pineapple.labels" . }}
app.kubernetes.io/component: webhook
{{- end -}}

{{- define "pineapple.webhook.selectorLabels" }}
{{ include "pineapple.selectorLabels" . }}
app.kubernetes.io/component: webhook
{{- end -}}

{{/*
Expand the name of the chart.
*/}}
{{- define "yousoon-storage.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
*/}}
{{- define "yousoon-storage.fullname" -}}
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
{{- define "yousoon-storage.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "yousoon-storage.labels" -}}
helm.sh/chart: {{ include "yousoon-storage.chart" . }}
{{ include "yousoon-storage.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
app.kubernetes.io/component: storage
{{- end }}

{{/*
Selector labels
*/}}
{{- define "yousoon-storage.selectorLabels" -}}
app.kubernetes.io/name: {{ include "yousoon-storage.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "yousoon-storage.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "yousoon-storage.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
Get the infrastructure mode (sidecar or classic)
*/}}
{{- define "yousoon-storage.mode" -}}
{{- .Values.global.infra | default "sidecar" }}
{{- end }}

{{/*
Check if we're in sidecar mode
*/}}
{{- define "yousoon-storage.isSidecar" -}}
{{- if eq (include "yousoon-storage.mode" .) "sidecar" }}true{{- else }}false{{- end }}
{{- end }}

{{/*
MongoDB service name
*/}}
{{- define "yousoon-storage.mongodb.serviceName" -}}
{{- printf "%s-mongodb" (include "yousoon-storage.fullname" .) }}
{{- end }}

{{/*
Redis service name
*/}}
{{- define "yousoon-storage.redis.serviceName" -}}
{{- printf "%s-redis" (include "yousoon-storage.fullname" .) }}
{{- end }}

{{/*
NATS service name
*/}}
{{- define "yousoon-storage.nats.serviceName" -}}
{{- printf "%s-nats" (include "yousoon-storage.fullname" .) }}
{{- end }}

{{/*
Elasticsearch service name
*/}}
{{- define "yousoon-storage.elasticsearch.serviceName" -}}
{{- printf "%s-elasticsearch" (include "yousoon-storage.fullname" .) }}
{{- end }}

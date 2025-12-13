{{/*
Expand the name of the chart.
*/}}
{{- define "yousoon-services.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
*/}}
{{- define "yousoon-services.fullname" -}}
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
{{- define "yousoon-services.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "yousoon-services.labels" -}}
helm.sh/chart: {{ include "yousoon-services.chart" . }}
{{ include "yousoon-services.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
app.kubernetes.io/component: services
{{- end }}

{{/*
Selector labels
*/}}
{{- define "yousoon-services.selectorLabels" -}}
app.kubernetes.io/name: {{ include "yousoon-services.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Full image path for a service
*/}}
{{- define "yousoon-services.image" -}}
{{- $registry := .Values.global.imageRegistry -}}
{{- $repository := .repository -}}
{{- $tag := .tag -}}
{{- printf "%s/%s:%s" $registry $repository $tag -}}
{{- end }}

{{/*
Common environment variables for all services
*/}}
{{- define "yousoon-services.commonEnv" -}}
- name: MONGODB_HOST
  value: {{ .Values.storageConnection.mongodbHost | quote }}
- name: MONGODB_PORT
  value: {{ .Values.storageConnection.mongodbPort | quote }}
- name: MONGODB_USERNAME
  value: "admin"
- name: MONGODB_PASSWORD
  valueFrom:
    secretKeyRef:
      name: yousoon-storage-secrets
      key: mongodb-root-password
- name: MONGODB_URI
  value: "mongodb://admin:$(MONGODB_PASSWORD)@{{ .Values.storageConnection.mongodbHost }}:{{ .Values.storageConnection.mongodbPort }}"
- name: REDIS_HOST
  value: {{ .Values.storageConnection.redisHost | quote }}
- name: REDIS_PORT
  value: {{ .Values.storageConnection.redisPort | quote }}
- name: REDIS_PASSWORD
  valueFrom:
    secretKeyRef:
      name: yousoon-storage-secrets
      key: redis-password
- name: NATS_URL
  value: "nats://{{ .Values.storageConnection.natsHost }}:{{ .Values.storageConnection.natsPort }}"
- name: ELASTICSEARCH_URL
  value: "http://{{ .Values.storageConnection.elasticsearchHost }}:{{ .Values.storageConnection.elasticsearchPort }}"
{{- range .Values.common.env }}
- name: {{ .name }}
  value: {{ .value | quote }}
{{- end }}
{{- end }}


# ip-tracker-chart/templates/_helpers.tpl
{{/*
Return the fully qualified name.
*/}}
{{- define "ip-tracker.fullname" -}}
{{ include "ip-tracker.name" . }}-{{ .Release.Name }}
{{- end }}

{{- define "ip-tracker.labels" -}}
app/name: "{{ include "ip-tracker.name" . }}"
{{- end }}



{{/*
Chart name
*/}}
{{- define "ip-tracker.name" -}}
{{ .Chart.Name }}
{{- end }}


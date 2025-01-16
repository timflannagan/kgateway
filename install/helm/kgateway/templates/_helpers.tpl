{{- define "kgateway.name" -}}
kgateway
{{- end }}

{{- define "kgateway.labels" -}}
app: "kgateway-controller"
chart: "{{ include "kgateway.name" . }}-{{ .Chart.Version }}"
release: "{{ .Release.Name }}"
{{- end }}
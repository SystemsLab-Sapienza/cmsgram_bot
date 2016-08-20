{{- range .}}
{{- if .Nome}}
<b>Nome:</b> {{.Nome}}
{{- end}}

{{- if .Email}}
<b>Email:</b> {{.Email}}
{{- end}}

{{- if .Telefono}}
<b>Telefono:</b> {{.Telefono}}
{{- end}}

{{- if .URL}}
<b>Sito:</b> {{.URL}}
{{- end}}
{{end}}
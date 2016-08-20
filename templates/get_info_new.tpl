{{- if .Nome}}
<b>Nome:</b> {{.Nome}}
{{- end}}
<b>Cognome:</b> {{.Cognome}}
<b>Email:</b> {{.Email}}

{{- if .EmailAltre}}
<b>Altri indirizzi email:</b>
{{- range .EmailAltre}}
{{.}}
{{- end}}
{{- end}}

{{- if .Indirizzo}}
<b>Indirizzo:</b> {{.Indirizzo}}
{{- end}}

{{- if .Telefono}}
<b>Telefono:</b> {{.Telefono}}
{{- end}}

{{- if .Sito}}
<b>Sito:</b> {{.Sito}}
{{- end}}

{{- if .SitoAltri}}
<b>Altre pagine web:</b>
{{- range .SitoAltri}}
{{.}}
{{- end}}
{{- end}}

<i>Iscriviti:</i> /s_d{{.ID}}

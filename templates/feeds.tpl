Sei iscritto ai seguenti feed:

{{- range .Feeds}}
<b>{{.Name}}</b>
<i>Cancella iscrizione:</i> /u_{{.ID}}
{{end}}
{{- $repositoryURL := .Info.RepositoryURL -}}
{{- if .Versions -}}
{{- if .Unreleased.CommitGroups -}}
<a name="unreleased"></a>
## [Unreleased]({{ $repositoryURL }}/compare/{{ $latest := index .Versions 0 }}{{ $latest.Tag.Name }}...HEAD)
{{ range .Unreleased.CommitGroups -}}
{{- $commits := list -}}
{{- $title := .Title -}}
{{- range .Commits -}}
{{- if not (and (hasPrefix $title "Chore") (hasPrefix .Scope "release")) -}}
{{- $commits = append $commits . -}}
{{- end -}}
{{- end -}}
{{- if $commits }}
### {{ .Title }}
{{ range $commits -}}
- {{ if .Scope }}**{{ .Scope }}:** {{ end }}[{{ .Hash.Short }}]({{ $repositoryURL }}/commit/{{ .Hash.Long }}) {{ .Subject }}
{{ end -}}
{{- end -}}
{{- end -}}
{{- range .Unreleased.NoteGroups -}}
### {{ .Title }}
{{- range .Notes }}
{{ .Body }}
{{ end -}}
{{- end -}}
{{- end -}}
{{- range .Versions }}
<a name="{{ .Tag.Name }}"></a>
## {{ if .Tag.Previous }}[{{ .Tag.Name }}]({{ $repositoryURL }}/compare/{{ .Tag.Previous.Name }}...{{ .Tag.Name }}){{ else }}{{ .Tag.Name }}{{ end }} ({{ datetime "2006-01-02" .Tag.Date }})
{{ range .CommitGroups }}
{{- $commits := list -}}
{{- $title := .Title -}}
{{- range .Commits -}}
{{- if not (and (hasPrefix $title "Chore") (hasPrefix .Scope "release")) -}}
{{- $commits = append $commits . -}}
{{- end -}}
{{- end -}}
{{- if $commits }}
### {{ .Title }}
{{ range $commits -}}
- {{ if .Scope }}**{{ .Scope }}:** {{ end }}[{{ .Hash.Short }}]({{ $repositoryURL }}/commit/{{ .Hash.Long }}) {{ .Subject }}
{{ end -}}
{{- end -}}
{{- end -}}
{{- range .NoteGroups }}
### {{ .Title }}
{{ range .Notes }}
{{ .Body }}
{{ end -}}
{{- end -}}
{{- end -}}
{{- end -}}
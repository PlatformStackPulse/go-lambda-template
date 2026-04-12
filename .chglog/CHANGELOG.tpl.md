{{ if .Versions -}}
<a name="unreleased"></a>
## [Unreleased]
{{ if .Unreleased.CommitGroups -}}
{{ range .Unreleased.CommitGroups -}}
### {{ .Title }}
{{ range .Commits -}}
- {{ if .Scope }}**{{ .Scope }}:** {{ end }}{{ .Subject }} ([{{ .Hash.Short }}](https://github.com/PlatformStackPulse/go-template/commit/{{ .Hash.Long }}))
{{ end }}
{{ end }}
{{ else -}}
- No notable changes.
{{ end }}

{{ range .Versions -}}
<a name="{{ .Tag.Name }}"></a>
## [{{ if .Tag.Previous }}{{ .Tag.Name }}{{ else }}{{ .Tag.Name }}{{ end }}] - {{ datetime "2006-01-02" .Tag.Date }}
{{ range .CommitGroups -}}
### {{ .Title }}
{{ range .Commits -}}
- {{ if .Scope }}**{{ .Scope }}:** {{ end }}{{ .Subject }} ([{{ .Hash.Short }}](https://github.com/PlatformStackPulse/go-template/commit/{{ .Hash.Long }}))
{{ end }}
{{ end }}

{{ end -}}
{{ else -}}
<a name="unreleased"></a>
## [Unreleased]
- Initial scaffolding.
{{ end -}}

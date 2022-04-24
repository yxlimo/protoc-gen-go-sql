package gen

import (
	"text/template"

	pgs "github.com/lyft/protoc-gen-star"
	pgsgo "github.com/lyft/protoc-gen-star/lang/go"
	"google.golang.org/protobuf/encoding/protojson"
)

type Option struct {
	Marshaler   protojson.MarshalOptions
	Unmarshaler protojson.UnmarshalOptions
}

var _ pgs.Module = &Sql{}

type Sql struct {
	option     Option
	base       *pgs.ModuleBase
	ctx        pgsgo.Context
	tpl        *template.Template
	messageTpl *template.Template
	enumTpl    *template.Template
}

type FileWrapper struct {
	File   pgs.File
	Option Option
}

func New(opt Option) *Sql {
	return &Sql{
		option: opt,
		base:   &pgs.ModuleBase{},
	}
}

func (s *Sql) Name() string {
	return "sql"
}

func (s *Sql) InitContext(c pgs.BuildContext) {
	s.base.InitContext(c)
	s.ctx = pgsgo.InitContext(c.Parameters())
	s.tpl = template.Must(template.New("sql").Funcs(map[string]interface{}{
		"package": s.ctx.PackageName,
		"name":    s.ctx.Name,
		"debug": func(v ...interface{}) string {
			s.base.Debug(v...)
			return ""
		},
	}).Parse(sqlTpl))
}

func (s *Sql) Execute(targets map[string]pgs.File, packages map[string]pgs.Package) []pgs.Artifact {
	for _, t := range targets {
		s.generate(t)
	}
	return s.base.Artifacts()
}

func (s *Sql) generate(file pgs.File) {
	wrapper := FileWrapper{
		File:   file,
		Option: s.option,
	}
	name := s.ctx.OutputPath(file).SetExt(".sql.go")
	s.base.AddGeneratorTemplateFile(name.String(), s.tpl, wrapper)
}

var sqlTpl = `package {{ package .File }}

import (
	"database/sql"
	"database/sql/driver"
	"fmt"

	"google.golang.org/protobuf/encoding/protojson"
	"github.com/spf13/cast"
)

var _ sql.Scanner
var _ driver.Valuer
var _ fmt.Stringer
var _ = cast.ToString
var _ = protojson.Marshal

{{ range .File.AllMessages }}
	{{- $fieldConflict := false}}
	{{- range .Fields}}
		{{- if or (eq (name .) "Value") (eq (name .) "Scan") }}
			{{- $fieldConflict = true }}{{break}}
		{{- end }}
	{{- end }}
{{- if $fieldConflict }}{{continue}}{{end}}
{{ debug "start generate message" (name .) }}
// Scan implements sql.Scanner
func (msg *{{ name . }}) Scan(src interface{}) error {
	if msg == nil {
		return fmt.Errorf("scan into nil {{ name . }}")
	}
	var value []byte
	switch v := src.(type) {
	case []byte:
		value = v
	case string:
		value = []byte(v)
	default:
		return fmt.Errorf("can't convert %v to {{ name . }}, unsupported type %T", src, src)
	}

	err := protojson.UnmarshalOptions{
		AllowPartial:   {{$.Option.Unmarshaler.AllowPartial}},
		DiscardUnknown: {{$.Option.Unmarshaler.DiscardUnknown}},
	}.Unmarshal(value, msg)
	if err != nil {
		return fmt.Errorf("can't unmarshal {{ name . }}: %w", err)
	}
	return nil
}

// Value implements driver.Valuer
func (msg {{ name . }}) Value() (driver.Value, error) {
	value, err := protojson.MarshalOptions{
		Multiline:       {{$.Option.Marshaler.Multiline}},
		UseEnumNumbers:  {{$.Option.Marshaler.UseEnumNumbers}},
		EmitUnpopulated: {{$.Option.Marshaler.EmitUnpopulated}},
		UseProtoNames:   {{$.Option.Marshaler.UseProtoNames}},
		AllowPartial:    {{$.Option.Marshaler.AllowPartial}},
	}.Marshal(&msg)
	if err != nil {
		return nil, fmt.Errorf("can't marshal {{ name . }}: %w", err)
	}
	return value, nil
}

{{ debug "generate" (name .) "ok" }}
{{ end }}

{{ debug "finish generate message" }}

{{ range .File.AllEnums }}
func (x *{{ name . }}) Scan(src any) error {

	switch v := src.(type) {
	case []byte, string:
		*x = {{ name . }}({{ name . }}_value[cast.ToString(v)])
	case int, int8, int16, int32, int64:
		*x = {{ name . }}(cast.ToInt32(v))
	default:
		return fmt.Errorf("cannot scan type %T into {{ name . }}: %v", src, src)
	}
	return nil
}

func (x {{ name . }}) Value() (driver.Value, error) {
	{{- if $.Option.Marshaler.UseEnumNumbers }}
	return x, nil
	{{ else }}
	return x.String(), nil
	{{- end }}
}


{{ end }}

`

package sql

import (
	"bytes"
	"database/sql"
	"sort"
	"text/template"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

const tpl = `
BEGIN;
{{ range $tblIdx, $table := .Tables -}}
DROP TABLE IF EXISTS {{$table.Name}};
CREATE TABLE {{$table.Name}} (
{{ range $index, $col := $table.Columns -}}
{{$col.Name}} {{$col.DataType}}{{if (or (notLast $index $table.Columns) (ne $table.PrimaryKey ""))}},{{end}}
{{ end }}
{{ if ne .PrimaryKey "" }}PRIMARY KEY ({{.PrimaryKey}}){{ end }}
);
{{ end }}
COMMIT;
`

type column struct {
	Name     string
	DataType string
}

type table struct {
	Name       string
	PrimaryKey string
	Columns    []column
}

type schema struct {
	Tables []table
}

func parse(yml string) (schema, error) {
	m := make(map[string]map[string]interface{})

	var s schema

	err := yaml.Unmarshal([]byte(yml), &m)

	if err != nil {
		return schema{}, err
	}

	for k, val := range m {

		var columns []column

		var pk string

		for colName, colDataType := range val {
			if colName == "_pk" {
				pk = colDataType.(string)
				continue
			}
			columns = append(columns, column{
				DataType: colDataType.(string),
				Name:     colName,
			})
		}

		s.Tables = append(s.Tables, table{
			Name:       k,
			Columns:    columns,
			PrimaryKey: pk,
		})

	}

	sort.SliceStable(s.Tables, func(i, j int) bool {
		return s.Tables[i].Name < s.Tables[j].Name
	})

	for ix, _ := range s.Tables {
		sort.SliceStable(s.Tables[ix].Columns, func(i, j int) bool {
			return s.Tables[ix].Columns[i].Name < s.Tables[ix].Columns[j].Name
		})
	}

	return s, nil
}

// BuildSetupSchema builds a sql query string to tear down and setup a schema
func BuildSetupSchema(yml string) (string, error) {

	data, err := parse(yml)

	var buf bytes.Buffer

	if err != nil {
		return "", errors.Errorf("Parsed fail: %v", err)
	}

	t, err := template.New("sql").Funcs(funcMap).Parse(tpl)

	if err != nil {
		return "", errors.Errorf("TPL Parsed fail: %v", err)
	}

	err = t.Execute(&buf, data)

	if err != nil {
		return "", errors.Errorf("Template failed: %v", err)
	}

	return buf.String(), nil
}

func Exec(db *sql.DB, yml string) error {
	q, err := BuildSetupSchema(yml)

	if err != nil {
		return err
	}

	_, err = db.Exec(q)

	return err
}
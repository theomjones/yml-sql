package sql

import (
	require "github.com/stretchr/testify/require"
	"strings"
	"testing"
)

const fixture1 = `
people:
  id: "char(36) NOT NULL"
  name: "varchar(255) NOT NULL"
  _pk: id
dogs:
  id: "char(36)"
  name: "varchar(255) NOT NULL"
  breed: "varchar(255)"
cats:
  id: "char(36)"
  name: "varchar(255) NOT NULL"
  breed: "varchar(255)"
  _pk: id
`

const result1 = `
BEGIN;
DROP TABLE IF EXISTS cats;
CREATE TABLE cats (
 breed VARCHAR(255),
 id CHAR(36),
 name VARCHAR(255) NOT NULL,
 PRIMARY KEY (id)
);
DROP TABLE IF EXISTS dogs;
CREATE TABLE dogs (
 breed VARCHAR(255),
 id CHAR(36),
 name VARCHAR(255) NOT NULL
);
DROP TABLE IF EXISTS people;
CREATE TABLE people (
 id CHAR(36) NOT NULL,
 name VARCHAR(255) NOT NULL,
 PRIMARY KEY (id)
);
COMMIT;
`

func TestBuild(t *testing.T) {

	req := require.New(t)

	tests := []struct {
		want string
		give string
	}{
		{result1, fixture1},
	}

	for _, row := range tests {
		res, err := BuildSetupSchema(row.give)
		//t.Log("\nRESULT:\n", res)
		if err != nil {
			t.Error(err)
		}
		req.Equal(normaliseString(row.want), normaliseString(res))
	}

}

func normaliseString(s string) string {
	return strings.ToLower(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(s, "\n", " "), "\t", " "), " ", ""))
}

func TestParse(t *testing.T) {
	req := require.New(t)
	s, err := parse(fixture1)

	if err != nil {
		t.Error(err)
	}

	people := getTable(s.Tables, "people")
	cats := getTable(s.Tables, "cats")
	dogs := getTable(s.Tables, "dogs")

	req.Contains(people.Columns, column{Name: "id", DataType: "char(36) NOT NULL"})
	req.Contains(people.Columns, column{Name: "name", DataType: "varchar(255) NOT NULL"})
	req.Equal(people.PrimaryKey, "id")

	req.Contains(cats.Columns, column{Name: "breed", DataType: "varchar(255)"})
	req.Contains(cats.Columns, column{Name: "id", DataType: "char(36)"})
	req.Contains(cats.Columns, column{Name: "name", DataType: "varchar(255) NOT NULL"})
	req.Equal(cats.PrimaryKey, "id")

	req.Contains(dogs.Columns, column{Name: "breed", DataType: "varchar(255)"})
	req.Contains(dogs.Columns, column{Name: "id", DataType: "char(36)"})
	req.Contains(dogs.Columns, column{Name: "name", DataType: "varchar(255) NOT NULL"})
	req.Zero(dogs.PrimaryKey)

}

func getTable(ts []table, k string) table {
	var t table
	for _, v := range ts {
		if v.Name == k {
			t = v
			break
		}
	}
	return t
}

package main

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	sq "github.com/theomjones/yml-sql/sql"
	"html/template"
	"io"
	"io/ioutil"

	_ "github.com/go-sql-driver/mysql"
	uuid "github.com/google/uuid"
	"github.com/romanyx/polluter"
)

type SeedData struct {
	Name string
}

func newUUIDString() string {
	return uuid.NewString()
}

var funcMap = template.FuncMap{
	"UUID": newUUIDString,
}

// TODO: big clean up
// then create tables of first objects

func getSeedQuery(filename string) ([]byte, error) {
	tpl, err := template.New("seed.yml").Funcs(funcMap).ParseFiles(filename)

	if err != nil {
		return nil, errors.New("No seed.yml")
	}

	var buff bytes.Buffer

	sd := SeedData{
		Name: "Theoz",
	}

	err = tpl.Execute(&buff, sd)

	if err != nil {
		return nil, err
	}

	return buff.Bytes(), nil
}

func getSetupYaml(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}

func main() {

	db, err := sql.Open("mysql", "root:password@/test")
	defer db.Close()

	if err != nil {
		panic(err)
	}

	seedQuery, err := getSeedQuery("./db/seed.yml")

	if err != nil {
		panic(err)
	}


	defer seed(db, bytes.NewReader(seedQuery))

	setupYml, err := getSetupYaml("./db/schema.yml")

	if err != nil {
		panic(err)
	}

	setup, err := sq.BuildSetupSchema(string(setupYml))

	if err != nil {
		panic(err)
	}


	fmt.Println("sql", setup)

	_, err = db.Exec(`
CREATE TABLE IF NOT EXISTS test.dogs
(
    breed varchar(255),
    id    char,
    name  varchar(255) NOT NULL,

    PRIMARY KEY (id)
);
`)

	if err != nil {
		panic(err)
	}

	//if err = seed(db, bytes.NewReader(seedQuery)); err != nil {
	//	panic(err)
	//}
}

func seed(db *sql.DB, r io.Reader) error {

	eng := polluter.MySQLEngine(db)
	p := polluter.New(eng)

	if err := p.Pollute(r); err != nil {
		return err
	}

	return nil
}

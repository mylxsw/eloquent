package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/mylxsw/eloquent/examples/models"
	"github.com/mylxsw/eloquent/query"
	"gopkg.in/guregu/null.v3"
)

func main() {

	connURI := "root:@tcp(127.0.0.1:3306)/wizard?parseTime=true"
	db, err := sql.Open("mysql", connURI)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	projectModel := models.NewProjectModel(db)

	id, err := projectModel.Save(models.Project{
		Name: "test",
	})
	assertError(err)

	fmt.Println("insert id=", id)

	_, err = projectModel.UpdateById(id, models.Project{
		Name:      "test2",
		CatalogId: null.IntFrom(100),
	})
	assertError(err)

	_, err = projectModel.DeleteById(id)
	assertError(err)

	projects, err := projectModel.Get(query.Builder().OrderBy("id", "desc").Limit(10))
	assertError(err)

	for _, p := range projects {
		fmt.Printf("%v\n", p)
	}
}

func assertError(err error) {
	if err != nil {
		panic(err)
	}
}

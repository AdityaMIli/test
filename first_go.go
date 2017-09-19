package main

import (
	"net/http"
	"fmt"
	"html/template"
	"database/sql"
	_ "github.com/lib/pq"
	"time"
)

type userinfo struct {
	uid    int
	Name  string
	Department string
	Created time.Time
}


func connect() (*sql.DB, error) {
	db, err := sql.Open("postgres", "user=postgres password='cari cari remote tv' dbname=postgres sslmode=disable")
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return db, nil
}


func userInfo(w http.ResponseWriter, r *http.Request) {
	db, err := connect()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer db.Close()

	if r.Method == "GET" {

		rows, err := db.Query("SELECT uid, username, departname, created FROM public.userinfo")
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		defer rows.Close()

		funcMap := template.FuncMap{
			"inc": func(i int) int {
				return i + 1
			},
		}

		tempt, err := template.New("template.html").Funcs(funcMap).ParseFiles("template.html")
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		//fmt.Println("Success")

		var data []userinfo

		for rows.Next() {
			var each = userinfo{}
			var err = rows.Scan(&each.uid,&each.Name, &each.Department,&each.Created)

			if err != nil {
				fmt.Println(err.Error())
				return
			}

			data = append(data, each)
		}

		if err = rows.Err(); err != nil {
			fmt.Println(err.Error())
			return
		}

		tempt.Execute(w, map[string]interface{}{
			"row": data,
		})
	}
		return
}

func main() {
	http.HandleFunc("/userinfo", userInfo)
	fmt.Println("starting web server at http://localhost:8088/")
	http.ListenAndServe(":8088", nil)
}

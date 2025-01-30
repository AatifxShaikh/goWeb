package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"

	"aatif.net/snippetbox/pkg/models/mysql"
	_ "github.com/go-sql-driver/mysql"
)

type application struct {
	errLog   *log.Logger
	infoLog  *log.Logger
	snippets *mysql.SnippetModel
}

// testing
func main() {
	addr := flag.String("addr", ":4000", "Http network address")
	dsn := flag.String("dsn", "web:tafia934@/snippetbox?parseTime=true", "MySql databse")
	flag.Parse()
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		errLog.Fatal(err)
	}

	defer db.Close()
	app := &application{
		errLog:   errLog,
		infoLog:  infoLog,
		snippets: &mysql.SnippetModel{DB: db},
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServe()
	errLog.Fatal(err)
}
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err

	}
	return db, nil
}

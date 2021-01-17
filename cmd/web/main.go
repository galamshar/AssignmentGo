package main

import (
	"AssignmentGo/pkg/models/postgres"
	"context"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	snippets *postgres.SnippetModel
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	pool, err := pgxpool.Connect(context.Background(), "user=postgres password=Zawer021 host=localhost port=5432 dbname=snippetbox sslmode=disable pool_max_conns=10")
	if err != nil {
		log.Fatalf("Unable to connection to database : %v\n", err)
	}
	snippets := postgres.SnippetModel{Pool: pool}
	defer pool.Close()

	var app = &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		snippets: &snippets,
	}
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(), // Call the new app.routes() method
	}
	infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

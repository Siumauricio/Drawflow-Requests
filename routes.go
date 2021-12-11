package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/dgraph-io/dgo/v210/protos/api"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)
type Code struct {
    Code []string `json:"Code"`
	Uid string `json:"uid"`
}

const queryCode string = `
{
	getAll(func: has(Code)) {
		uid
		Code
	}
}
`
const getFileWithId string = `
{
	node(func: uid(%s)) {
	  uid
	  Code
	  expand(_all_) {
		uid
		expand(_all_)
	  }
	}
  }
  `


func Route() *chi.Mux {
	mux := chi.NewMux()
	cors := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, 
	})

	//Globals middlewares
	mux.Use(
		middleware.Logger,
		middleware.Recoverer,
		cors.Handler,
	)

	mux.Get("/getAll", getAllFiles)
	mux.Post("/insertCode", insertCode)
	mux.Get("/getFile", getFile);
	return mux
}

func getAllFiles(w http.ResponseWriter, r *http.Request) {
	 w.Header().Set("Content-Type", "application/json")
	 dgClient := newClient()
	 txn := dgClient.NewTxn()
	 resp , err := txn.Query(context.Background(), queryCode)

	 if err != nil {
	 	log.Fatal(err)
	 }
	w.Write(resp.Json)
}
func insertCode(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var rawCode Code
	_ = json.NewDecoder(r.Body).Decode(&rawCode)
	 p := Code { Code: rawCode.Code }
	pb, err := json.Marshal(p)
	if err != nil {
		log.Fatal(err)
	}
	dgClient := newClient()
	txn := dgClient.NewTxn()

	mutBuyers := &api.Mutation{
		CommitNow: true,
		SetJson: pb,
	}
	resp , err := txn.Mutate(context.Background(), mutBuyers)

	if err != nil {
		log.Fatal(err)
	}
	 w.Write(resp.Json)
}
func getFile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var rawCode Code
	_ = json.NewDecoder(r.Body).Decode(&rawCode)
	query := getQuery(rawCode.Uid)
	dgClient := newClient()
	txn := dgClient.NewTxn()
	resp , err := txn.Query(context.Background(), query)

	if err != nil {
		log.Fatal(err)
	}
	w.Write(resp.Json)
}

func getQuery( uid string )string{
	return fmt.Sprintf(getFileWithId,uid )
}


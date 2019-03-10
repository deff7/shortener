package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/boltdb/bolt"
)

var bucketName = []byte("urls")

const baseOffset = 100000000

type rootHandler struct {
	db   *bolt.DB
	tmpl *template.Template
	gen  *Generator
}

func (h *rootHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		if req.URL.Path == "/" {
			h.tmpl.Execute(w, nil)
			return
		}

		h.db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket(bucketName)
			v := b.Get([]byte(req.URL.Path[1:]))
			http.Redirect(w, req, string(v), 301)
			return nil
		})
	}

	if req.Method == http.MethodPost {
		h.db.Update(func(tx *bolt.Tx) error {
			b, err := tx.CreateBucketIfNotExists(bucketName)
			if err != nil {
				return err
			}
			id, err := b.NextSequence()
			if err != nil {
				return err
			}
			id += baseOffset

			var (
				key = h.gen.EncodeID(id)
				val = req.FormValue("url")
			)

			err = b.Put(key, []byte(val))
			if err == nil {
				w.Write([]byte("localhost:8080/" + string(key)))
			}

			return err
		})
	}
}

func main() {
	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	http.Handle("/", &rootHandler{
		db:   db,
		tmpl: template.Must(template.ParseFiles("./web/template/index.html")),
		gen:  NewGenerator(),
	})
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}

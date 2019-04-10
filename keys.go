package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"

	"github.com/belljustin/hancock/models"
)

func newGetKey(keys models.Keys) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		w.Header().Set("Content-Type", "application/json")

		sid := ps.ByName("id")
		if sid == "" {
			log.Println("Router matched path without `id` param")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		id, err := uuid.Parse(sid)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		k, err := keys.Get(id)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		enc := json.NewEncoder(w)
		if err := enc.Encode(&k); err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

type CreateKeyRequest struct {
	Algorithm string `json:"alg"`
}

func newCreateKey(keys models.Keys) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		w.Header().Set("Content-Type", "application/json")

		decoder := json.NewDecoder(req.Body)
		var ck CreateKeyRequest
		if err := decoder.Decode(&ck); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		jsonCk, err := json.Marshal(ck)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		w.Write(jsonCk)
	}
}

func RegisterKeyHandlers(r *httprouter.Router, keys models.Keys) {
	root := "/keys"

	r.GET(root+"/:id", newGetKey(keys))
	r.POST(root, newCreateKey(keys))
}

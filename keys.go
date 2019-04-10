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

func newCreateKey(keys models.Keys, algs map[string]Alg) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		w.Header().Set("Content-Type", "application/json")

		decoder := json.NewDecoder(req.Body)
		var ck CreateKeyRequest
		if err := decoder.Decode(&ck); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		alg := algs[ck.Algorithm]
		if alg == nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		k, err := alg.NewKey("belljust.in/justin")
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err = keys.Create(k); err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
		}

		jsonKey, err := json.Marshal(k)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		w.Write(jsonKey)
	}
}

func RegisterKeyHandlers(r *httprouter.Router, keys models.Keys, algs map[string]Alg) {
	root := "/keys"

	r.GET(root+"/:id", newGetKey(keys))
	r.POST(root, newCreateKey(keys, algs))
}

package main

import (
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"

	"github.com/belljustin/hancock/models"
)

type httpError struct {
	statusCode int
}

func (err *httpError) Error() string {
	return http.StatusText(err.statusCode)
}

func getKey(keys models.Keys, sid string) (*models.Key, error) {
	if sid == "" {
		log.Println("Router matched path without `id` param")
		return nil, &httpError{http.StatusInternalServerError}
	}
	id, err := uuid.Parse(sid)
	if err != nil {
		return nil, &httpError{http.StatusBadRequest}
	}

	k, err := keys.Get(id)
	if err != nil {
		log.Println(err.Error())
		return nil, &httpError{http.StatusInternalServerError}
	} else if k == nil {
		return nil, &httpError{http.StatusNotFound}
	}

	return k, nil
}

func newGetKey(keys models.Keys) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		w.Header().Set("Content-Type", "application/json")

		k, err := getKey(keys, ps.ByName("id"))
		if err != nil {
			if err, ok := err.(*httpError); ok {
				w.WriteHeader(err.statusCode)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
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

type CreateSignatureRequest struct {
	Digest string `json:"digest"`
	Hash   string `json:"hash"`
}

type CreateSignatureResponse struct {
	Signature []byte `json:"signature"`
}

func newCreateSignature(keys models.Keys, algs map[string]Alg) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		w.Header().Set("Content-Type", "application/json")

		k, err := getKey(keys, ps.ByName("id"))
		if err != nil {
			if err, ok := err.(*httpError); ok {
				w.WriteHeader(err.statusCode)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}

		alg := algs[k.Algorithm]
		if alg == nil {
			log.Println("Tried to sign key with unsupported algorithm")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		decoder := json.NewDecoder(req.Body)
		var cs CreateSignatureRequest
		if err := decoder.Decode(&cs); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		bDigest, err := hex.DecodeString(cs.Digest)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		sig, err := alg.Sign(k.Priv, bDigest, cs.Hash)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		jsonSig, err := json.Marshal(sig)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		w.Write(jsonSig)
	}
}

func RegisterKeyHandlers(r *httprouter.Router, keys models.Keys, algs map[string]Alg) {
	root := "/keys"

	r.POST(root, newCreateKey(keys, algs))
	r.GET(root+"/:id", newGetKey(keys))
	r.POST(root+"/:id/signature", newCreateSignature(keys, algs))
}

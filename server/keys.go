package server

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"

	"github.com/belljustin/hancock/models"
)

func getKey(keys models.Keys, sid string) (*models.Key, error) {
	if sid == "" {
		return nil, &httpError{
			http.StatusInternalServerError,
			"Router matched path without id param",
		}
	}
	id, err := uuid.Parse(sid)
	if err != nil {
		return nil, &httpError{
			http.StatusBadRequest,
			fmt.Sprintf("Could not parse key id '%s'", sid),
		}
	}

	k, err := keys.Get(id)
	if err != nil {
		return nil, newInternalServerError(err)
	} else if k == nil {
		return nil, &httpError{
			http.StatusNotFound,
			fmt.Sprintf("Could not find key id '%s'", id.String()),
		}
	}

	return k, nil
}

func newGetKey(keys models.Keys) httprouter.Handle {
	f := func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) error {
		w.Header().Set("Content-Type", "application/json")

		k, err := getKey(keys, ps.ByName("id"))
		if err != nil {
			return err
		}

		enc := json.NewEncoder(w)
		if err := enc.Encode(&k); err != nil {
			return newInternalServerError(err)
		}

		return nil
	}

	return appHandler(f).Handle
}

type CreateKeyRequest struct {
	Algorithm string `json:"alg"`
}

func newCreateKey(keys models.Keys, algs map[string]Alg) httprouter.Handle {
	f := func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) error {
		w.Header().Set("Content-Type", "application/json")

		decoder := json.NewDecoder(req.Body)
		var ck CreateKeyRequest
		if err := decoder.Decode(&ck); err != nil {
			return &httpError{
				http.StatusBadRequest,
				"Malformed request",
			}
		}

		alg := algs[ck.Algorithm]
		if alg == nil {
			return &httpError{
				http.StatusBadRequest,
				fmt.Sprintf("Unsupported algorithm '%s'", ck.Algorithm),
			}
		}

		k, err := alg.NewKey("belljust.in/justin")
		if err != nil {
			return newInternalServerError(err)
		}

		if err = keys.Create(k); err != nil {
			return newInternalServerError(err)
		}

		jsonKey, err := json.Marshal(k)
		if err != nil {
			return newInternalServerError(err)
		}

		w.WriteHeader(http.StatusCreated)
		w.Write(jsonKey)
		return nil
	}

	return appHandler(f).Handle
}

type CreateSignatureRequest struct {
	Digest string `json:"digest"`
	Hash   string `json:"hash"`
}

type CreateSignatureResponse struct {
	Signature []byte `json:"signature"`
}

func newCreateSignature(keys models.Keys, algs map[string]Alg) httprouter.Handle {
	f := func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) error {
		w.Header().Set("Content-Type", "application/json")

		k, err := getKey(keys, ps.ByName("id"))
		if err != nil {
			return err
		}

		alg := algs[k.Algorithm]
		if alg == nil {
			return &httpError{
				http.StatusInternalServerError,
				fmt.Sprintf("Tried to sign key with unsupported algorithm '%s'", k.Algorithm),
			}
		}

		decoder := json.NewDecoder(req.Body)
		var cs CreateSignatureRequest
		if err := decoder.Decode(&cs); err != nil {
			return &httpError{
				http.StatusBadRequest,
				"Malformed Request",
			}
		}

		bDigest, err := hex.DecodeString(cs.Digest)
		if err != nil {
			return &httpError{
				http.StatusBadRequest,
				"Digest must be hex encoded",
			}
		}

		sig, err := alg.Sign(k.Priv, bDigest, cs.Hash)
		if err != nil {
			return &httpError{
				http.StatusInternalServerError,
				err.Error(),
			}
		}

		res := &CreateSignatureResponse{Signature: sig}
		jsonSig, err := json.Marshal(res)
		if err != nil {
			return newInternalServerError(err)
		}

		w.WriteHeader(http.StatusCreated)
		w.Write(jsonSig)
		return nil
	}

	return appHandler(f).Handle
}

func RegisterKeyHandlers(r *httprouter.Router, keys models.Keys, algs map[string]Alg) {
	root := "/keys"

	r.POST(root, newCreateKey(keys, algs))
	r.GET(root+"/:id", newGetKey(keys))
	r.POST(root+"/:id/signature", newCreateSignature(keys, algs))
}

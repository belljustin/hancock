package server

import (
	"encoding/hex"
	_ "encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"

	"github.com/belljustin/hancock/models"
)

func getKey(keys models.Keys, sid string) (*models.Key, error) {
	id, err := uuid.Parse(sid)
	if err != nil {
		return nil, &httpError{
			http.StatusBadRequest,
			fmt.Sprintf("Could not parse key id '%s'", sid),
		}
	}

	k, err := keys.Get(id)
	if err != nil {
		return nil, err
	} else if k == nil {
		return nil, &httpError{
			http.StatusNotFound,
			fmt.Sprintf("Could not find key id '%s'", id.String()),
		}
	}

	return k, nil
}

func newGetKey(keys models.Keys) gin.HandlerFunc {
	return func(c *gin.Context) {
		k, err := getKey(keys, c.Param("id"))
		if err != nil {
			handleError(c, err)
			return
		}
		c.JSON(200, &k)
	}
}

type CreateKeyRequest struct {
	Algorithm string `json:"alg" binding:"required"`
}

func newCreateKey(keys models.Keys, algs map[string]Alg) gin.HandlerFunc {
	return func(c *gin.Context) {
		var ck CreateKeyRequest
		if err := c.ShouldBind(&ck); err != nil {
			handleError(c, &httpError{
				http.StatusBadRequest,
				"Malformed request",
			})
			return
		}

		alg := algs[ck.Algorithm]
		if alg == nil {
			handleError(c, &httpError{
				http.StatusBadRequest,
				fmt.Sprintf("Unsupported algorithm '%s'", ck.Algorithm),
			})
			return
		}

		k, err := alg.NewKey("belljust.in/justin")
		if err != nil {
			panic(err)
		}

		if err = keys.Create(k); err != nil {
			panic(err)
		}

		c.JSON(http.StatusCreated, &k)
	}
}

type CreateSignatureRequest struct {
	Digest string `json:"digest" binding:"required"`
	Hash   string `json:"hash" binding:"required"`
}

type CreateSignatureResponse struct {
	Signature []byte `json:"signature"`
}

func newCreateSignature(keys models.Keys, algs map[string]Alg) gin.HandlerFunc {
	return func(c *gin.Context) {
		k, err := getKey(keys, c.Param("id"))
		if err != nil {
			handleError(c, err)
			return
		}

		alg := algs[k.Algorithm]
		if alg == nil {
			err = fmt.Errorf("Tried to sign key with unsupported algorithm '%s'", k.Algorithm)
			panic(err)
		}

		var cs CreateSignatureRequest
		if err := c.ShouldBind(&cs); err != nil {
			handleError(c, &httpError{
				http.StatusBadRequest,
				"Malformed Request",
			})
			return
		}
		fmt.Printf("%+v\n", cs)

		bDigest, err := hex.DecodeString(cs.Digest)
		if err != nil {
			handleError(c, &httpError{
				http.StatusBadRequest,
				"Digest must be hex encoded",
			})
			return
		}

		sig, err := alg.Sign(k.Priv, bDigest, cs.Hash)
		if err != nil {
			panic(err)
		}

		res := &CreateSignatureResponse{Signature: sig}
		c.JSON(http.StatusCreated, &res)
	}
}

func RegisterKeyHandlers(r *gin.Engine, keys models.Keys, algs map[string]Alg) {
	root := "/keys"

	r.POST(root, newCreateKey(keys, algs))
	r.GET(root+"/:id", newGetKey(keys))
	r.POST(root+"/:id/signature", newCreateSignature(keys, algs))
}

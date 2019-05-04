package server

import (
	"crypto"
	"crypto/rand"
	"encoding/hex"
	_ "encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/gin-gonic/gin/binding"

	"github.com/belljustin/hancock/key"
)

var hashes = map[string]crypto.Hash{
	"sha256": crypto.SHA256,
}

type KeysHandler struct {
	keys key.Storage
}

func (h *KeysHandler) getKeyById(id string) (*key.Key, error) {
	k, err := h.keys.Get(id)
	if err != nil {
		return nil, err
	} else if k == nil {
		return nil, &httpError{
			http.StatusNotFound,
			fmt.Sprintf("Could not find key id '%s'", id),
		}
	}

	return k, nil
}

type GetKeyResponse struct {
	Id        string           `json:"id"`
	Algorithm string           `json:"alg"`
	Owner     string           `json:"owner"`
	PublicKey crypto.PublicKey `json:"public_key"`
}

func (h *KeysHandler) getKey(c *gin.Context) {
	k, err := h.getKeyById(c.Param("id"))
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(200, &GetKeyResponse{
		Id:        k.Id,
		Algorithm: k.Algorithm,
		Owner:     k.Owner,
		PublicKey: k.Signer.Public(),
	})
}

type CreateKeyRequest struct {
	Algorithm string   `json:"alg" binding:"required"`
	Opts      key.Opts `json:"opts"`
}

type CreateKeyResponse struct {
	Id string `json:"id" binding:"required"`
}

func (h *KeysHandler) createKey(c *gin.Context) {
	var ck CreateKeyRequest
	if err := c.ShouldBind(&ck); err != nil {
		handleError(c, &httpError{
			http.StatusBadRequest,
			"Malformed request",
		})
		return
	}

	// TODO: cleanup opts and owner
	k, err := h.keys.Create("belljust.in/justin", ck.Algorithm, ck.Opts)
	if err != nil {
		// TODO: figure out how to handle errors
		handleError(c, &httpError{
			http.StatusBadRequest,
			err.Error(),
		})
		return
	}

	res := &CreateKeyResponse{k.Id}
	c.JSON(http.StatusCreated, res)
}

type CreateSignatureRequest struct {
	Digest string `json:"digest" binding:"required"`
	Hash   string `json:"hash" binding:"required"`
}

type CreateSignatureResponse struct {
	Signature []byte `json:"signature"`
}

func (h *KeysHandler) createSignature(c *gin.Context) {
	k, err := h.getKeyById(c.Param("id"))
	if err != nil {
		handleError(c, err)
		return
	}

	var cs CreateSignatureRequest
	if err := c.ShouldBind(&cs); err != nil {
		handleError(c, &httpError{
			http.StatusBadRequest,
			err.Error(),
		})
		return
	}

	hash, ok := hashes[cs.Hash]
	if !ok {
		handleError(c, &httpError{
			http.StatusBadRequest,
			fmt.Sprintf("Hash '%s' is not supported", cs.Hash),
		})
		return
	}

	bDigest, err := hex.DecodeString(cs.Digest)
	if err != nil {
		handleError(c, &httpError{
			http.StatusBadRequest,
			"Digest must be hex encoded",
		})
		return
	}

	sig, err := k.Signer.Sign(rand.Reader, bDigest, hash)
	if err != nil {
		panic(err)
	}

	res := &CreateSignatureResponse{Signature: sig}
	c.JSON(http.StatusCreated, &res)
}

func registerKeyHandlers(r *gin.Engine, s key.Storage) {
	h := &KeysHandler{s}

	kr := r.Group("/keys")

	kr.POST("/", h.createKey)
	kr.GET("/:id", h.getKey)
	kr.POST("/:id/signature", h.createSignature)
}

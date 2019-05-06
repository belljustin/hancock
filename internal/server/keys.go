package server

import (
	"crypto"
	"crypto/rand"
	"encoding/hex"
	_ "encoding/json" // for tagging structs
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/gin-gonic/gin/binding" // for gin bindings

	"github.com/belljustin/hancock/key"
)

var hashes = map[string]crypto.Hash{
	"sha256": crypto.SHA256,
}

type keysHandler struct {
	keys key.Storage
}

func (h *keysHandler) getKeyByID(id string) (*key.Key, error) {
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

type getKeyResponse struct {
	ID        string           `json:"id"`
	Algorithm string           `json:"alg"`
	PublicKey crypto.PublicKey `json:"public_key"`
}

func (h *keysHandler) getKey(c *gin.Context) {
	k, err := h.getKeyByID(c.Param("id"))
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(200, &getKeyResponse{
		ID:        k.ID,
		Algorithm: k.Algorithm,
		PublicKey: k.Signer.Public(),
	})
}

type createKeyRequest struct {
	Algorithm string   `json:"alg" binding:"required"`
	Opts      key.Opts `json:"opts"`
}

type createKeyResponse struct {
	ID string `json:"id" binding:"required"`
}

func (h *keysHandler) createKey(c *gin.Context) {
	var ck createKeyRequest
	if err := c.ShouldBind(&ck); err != nil {
		handleError(c, &httpError{
			http.StatusBadRequest,
			"Malformed request",
		})
		return
	}

	// TODO: cleanup opts
	k, err := h.keys.Create(ck.Algorithm, ck.Opts)
	if err != nil {
		// TODO: figure out how to handle errors
		handleError(c, &httpError{
			http.StatusBadRequest,
			err.Error(),
		})
		return
	}

	res := &createKeyResponse{k.ID}
	c.JSON(http.StatusCreated, res)
}

type createSignatureRequest struct {
	Digest string `json:"digest" binding:"required"`
	Hash   string `json:"hash" binding:"required"`
}

type createSignatureResponse struct {
	Signature []byte `json:"signature"`
}

func (h *keysHandler) createSignature(c *gin.Context) {
	k, err := h.getKeyByID(c.Param("id"))
	if err != nil {
		handleError(c, err)
		return
	}

	var cs createSignatureRequest
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

	res := &createSignatureResponse{Signature: sig}
	c.JSON(http.StatusCreated, &res)
}

func registerKeyHandlers(r *gin.Engine, s key.Storage) {
	h := &keysHandler{s}

	kr := r.Group("/keys")

	kr.POST("/", h.createKey)
	kr.GET("/:id", h.getKey)
	kr.POST("/:id/signature", h.createSignature)
}

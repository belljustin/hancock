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

type KeysHandler struct {
	keys models.Keys
	algs map[string]Alg
}

func (h *KeysHandler) getKeyById(sid string) (*models.Key, error) {
	id, err := uuid.Parse(sid)
	if err != nil {
		return nil, &httpError{
			http.StatusBadRequest,
			fmt.Sprintf("Could not parse key id '%s'", sid),
		}
	}

	k, err := h.keys.Get(id)
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

func (h *KeysHandler) getKey(c *gin.Context) {
	k, err := h.getKeyById(c.Param("id"))
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(200, &k)
}

type CreateKeyRequest struct {
	Algorithm string `json:"alg" binding:"required"`
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

	alg := h.algs[ck.Algorithm]
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

	if err = h.keys.Create(k); err != nil {
		panic(err)
	}

	c.JSON(http.StatusCreated, &k)
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

	alg := h.algs[k.Algorithm]
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

func registerKeyHandlers(r *gin.Engine, keys models.Keys, algs map[string]Alg) {
	h := &KeysHandler{
		keys,
		algs,
	}

	kr := r.Group("/keys")

	kr.POST("/", h.createKey)
	kr.GET("/:id", h.getKey)
	kr.POST("/:id/signature", h.createSignature)
}

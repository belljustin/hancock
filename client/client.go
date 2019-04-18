package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/google/uuid"

	"github.com/belljustin/hancock/models"
	"github.com/belljustin/hancock/server"
)

type Hancock struct {
	url string
	*http.Client
}

func NewHancockClient(url string) *Hancock {
	c := &http.Client{}
	return &Hancock{url, c}
}

func (c *Hancock) NewKey(alg string) (*models.Key, error) {
	ck := server.CreateKeyRequest{Algorithm: alg}

	data := new(bytes.Buffer)
	enc := json.NewEncoder(data)
	if err := enc.Encode(ck); err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", c.url+"/keys", data)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusCreated {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, errors.New(resp.Status)
		}
		return nil, errors.New(string(body))
	}

	var k models.Key
	dec := json.NewDecoder(resp.Body)
	if err = dec.Decode(&k); err != nil {
		return nil, err
	}

	return &k, nil
}

func (c *Hancock) GetKey(id uuid.UUID) (*models.Key, error) {
	url := fmt.Sprintf("%s/keys/%s", c.url, id.String())

	data := new(bytes.Buffer)
	req, err := http.NewRequest("GET", url, data)
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, errors.New(resp.Status)
		}
		return nil, errors.New(string(body))
	}

	var k models.Key
	dec := json.NewDecoder(resp.Body)
	if err = dec.Decode(&k); err != nil {
		return nil, err
	}

	return &k, nil
}

func (c *Hancock) CreateSignature(id uuid.UUID, hash, digest string) ([]byte, error) {
	url := fmt.Sprintf("%s/keys/%s/signature", c.url, id.String())

	cs := server.CreateSignatureRequest{Digest: digest, Hash: hash}
	data := new(bytes.Buffer)
	enc := json.NewEncoder(data)
	if err := enc.Encode(cs); err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, data)
	req.Header.Add("Content-Type", "application/json")
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusCreated {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, errors.New(resp.Status)
		}
		return nil, errors.New(string(body))
	}

	var s server.CreateSignatureResponse
	dec := json.NewDecoder(resp.Body)
	if err = dec.Decode(&s); err != nil {
		return nil, err
	}

	return s.Signature, nil
}

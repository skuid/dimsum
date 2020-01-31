package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/heroku/docker-registry-client/registry"
)

//
type QuayResponse struct {
	SchemaVersion int    `json:"schemaVersion"`
	Tag           string `json:"tag"`
	Name          string `json:"name"`
	Architecture  string `json:"architecture"`
	FsLayers      []struct {
		BlobSum string `json:"blobSum"`
	} `json:"fsLayers"`
	History []struct {
		V1Compatibility string `json:"v1Compatibility"`
	} `json:"history"`
}

// manifestQ is a termorary replacement for the Regsitry.Manifest() method from the heroku/docker-registry-client package
// Quay.io recently brought their API more in line with the Docker Registry API spec as part of their efforts to move to supporting
// the V2 format. The original method requested an unsigned manifest, but was actually expected a signed manifest.
// Now that Quay.io complies more with the spec, this means that the method was unable to unmarshal the response from Quay.
// This function replicates the original Registry.Manifest() method, but does not set a header.

// UPDATE: Quay has moved to docker v2 registry format by default, and while V1 compatibility info can be retrieved, the format is different.
// The above QuayResponse struct allows us to unmarshal the response from Quay and access the information we need withing
// the History []struct, returning that to the client.
func manifestQ(registry *registry.Registry, repository string, reference string) (*QuayResponse, error) {
	endpoint := fmt.Sprintf("/v2/%s/manifests/%s", repository, reference)
	url := registry.URL + endpoint
	registry.Logf("registry.manifest.get url=%s repository=%s reference=%s", url, repository, reference)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := registry.Client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var manifest QuayResponse
	err = json.Unmarshal(body, &manifest)
	if err != nil {
		return nil, err
	}

	return &manifest, nil
}

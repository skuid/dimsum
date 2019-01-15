package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/docker/distribution/manifest/schema1"
	"github.com/heroku/docker-registry-client/registry"
)

// manifestQ is a termorary replacement for the Regsitry.Manifest() method from the heroku/docker-registry-client package
// Quay.io recently brought their API more in line with the Docker Registry API spec as part of their efforts to move to supporting
// the V2 format. The original method requested an unsigned manifest, but was actually expected a signed manifest.
// Now that Quay.io complies more with the spec, this means that the method was unable to unmarshal the response from Quay.
// This function replicates the original Registry.Manifest() method, but does not set a header.
func manifestQ(registry *registry.Registry, repository string, reference string) (*schema1.SignedManifest, error) {
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

	signedManifest := &schema1.SignedManifest{}
	err = signedManifest.UnmarshalJSON(body)
	if err != nil {
		return nil, err
	}

	return signedManifest, nil
}

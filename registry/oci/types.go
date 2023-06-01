package oci

import "encoding/json"

const MediaTypeManifest = "application/vnd.oci.image.index.v1+json"

type Manifest struct {
	MediaType     string `json:"mediaType"`
	SchemaVersion int    `json:"schemaVersion"`
	Manifests     []struct {
		MediaType string `json:"mediaType"`
		Digest    string `json:"digest"`
		Size      int    `json:"size"`
		Platform  struct {
			Architecture string `json:"architecture"`
			Os           string `json:"os"`
		} `json:"platform"`
		Annotations struct {
			VndDockerReferenceDigest string `json:"vnd.docker.reference.digest"`
			VndDockerReferenceType   string `json:"vnd.docker.reference.type"`
		} `json:"annotations,omitempty"`
	} `json:"manifests"`
}

// DeserializedManifest wraps Manifest with a copy of the original JSON.
// It satisfies the distribution.Manifest interface.
type DeserializedManifest struct {
	Manifest

	// canonical is the canonical byte representation of the Manifest.
	canonical []byte
}

// UnmarshalJSON populates a new Manifest struct from JSON data.
func (m *DeserializedManifest) UnmarshalJSON(b []byte) error {
	m.canonical = make([]byte, len(b), len(b))
	// store manifest in canonical
	copy(m.canonical, b)

	// Unmarshal canonical JSON into Manifest object
	var manifest Manifest
	if err := json.Unmarshal(m.canonical, &manifest); err != nil {
		return err
	}

	m.Manifest = manifest

	return nil
}

// Payload returns the raw content of the manifest. The contents can be used to
// calculate the content identifier.
func (m DeserializedManifest) Payload() (string, []byte, error) {
	return m.MediaType, m.canonical, nil
}

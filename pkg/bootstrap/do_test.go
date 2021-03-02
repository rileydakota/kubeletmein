package bootstrap

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/4armed/kubeletmein/pkg/mocks"
	"github.com/stretchr/testify/assert"
	yaml "gopkg.in/yaml.v2"
)

var (
	exampleUserData = `k8saas_ca_cert: aWFtLi4uLmlycmVsZXZhbnQ=
k8saas_bootstrap_token: aWFtLi4uLmlycmVsZXZhbnQ=
k8saas_master_domain_name: 1.1.1.1`
)

func TestMetadataFromDOService(t *testing.T) {
	metadataClient := mocks.NewTestClient(func(req *http.Request) *http.Response {
		assert.Equal(t, "http://169.254.169.254/metadata/v1/user-data", req.URL.String(), "should be equal")
		return &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewBufferString(exampleUserData)),
			Header:     make(http.Header),
		}
	})

	kubeenv, err := fetchMetadataFromDOService(metadataClient)
	if err != nil {
		t.Errorf("want user-data, got %q", err)
	}

	m := Metadata{}
	err = yaml.Unmarshal(kubeenv, &m)
	if err != nil {
		t.Errorf("unable to parse YAML from kube-env: %v", err)
	}

	assert.Equal(t, "1.1.1.1", m.KubeMaster, "they should be equal")

}

func TestMetadataFromDOFile(t *testing.T) {
	tempFile, err := ioutil.TempFile("", "")
	if err != nil {
		t.Errorf("couldn't create temp file for user-data: %v", err)
	}
	_, err = tempFile.WriteString(exampleKubeEnv)
	if err != nil {
		t.Errorf("couldn't write user-data to temp file: %v", err)
	}

	kubeenv, err := fetchMetadataFromFile(tempFile.Name())
	if err != nil {
		t.Errorf("want user-data, got %q", err)
	}

	// Clean up
	err = os.Remove(tempFile.Name())
	if err != nil {
		t.Errorf("couldn't remove tempFile: %v", err)
	}

	k := Kubeenv{}
	err = yaml.Unmarshal(kubeenv, &k)
	if err != nil {
		t.Errorf("unable to parse YAML from kube-env: %v", err)
	}

	assert.Equal(t, "1.1.1.1", k.KubeMasterName, "they should be equal")
}

func TestBootstrapDOCmd(t *testing.T) {
	// TODO: Write test for end-to-end
}

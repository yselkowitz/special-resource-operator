package registry

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"io"
	"log"
	"strings"

	"github.com/google/go-containerregistry/pkg/crane"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/openshift-psap/special-resource-operator/pkg/exit"
	"github.com/openshift-psap/special-resource-operator/pkg/warn"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func LastLayer(entry string) v1.Layer {

	var repo string

	if hash := strings.Split(entry, "@"); len(hash) > 1 {
		repo = hash[0]
	} else if tag := strings.Split(entry, ":"); len(tag) > 1 {
		repo = tag[0]
	}

	options := crane.StrictValidation

	manifest, err := crane.Manifest(entry, options)
	if err != nil {
		warn.OnError(err)
		return nil
	}

	release := unstructured.Unstructured{}
	err = json.Unmarshal(manifest, &release.Object)
	exit.OnError(err)

	layers, _, err := unstructured.NestedSlice(release.Object, "layers")
	exit.OnError(err)

	last := layers[len(layers)-1]

	digest := last.(map[string]interface{})["digest"].(string)

	layer, err := crane.PullLayer(repo+"@"+digest, options)
	exit.OnError(err)

	return layer
}

func ReleaseManifests(name string, layer v1.Layer) (key string, value string) {

	targz, err := layer.Compressed()
	defer dclose(targz)
	exit.OnError(err)

	gr, err := gzip.NewReader(targz)
	defer dclose(gr)
	exit.OnError(err)

	tr := tar.NewReader(gr)

	version := ""
	imageURL := ""

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}

		if header.Name == "release-manifests/image-references" {

			buff, err := io.ReadAll(tr)
			exit.OnError(err)

			obj := unstructured.Unstructured{}

			err = json.Unmarshal(buff, &obj.Object)
			exit.OnError(err)

			tags, _, err := unstructured.NestedSlice(obj.Object, "spec", "tags")
			exit.OnError(err)

			for _, tag := range tags {
				if tag.(map[string]interface{})["name"] == "driver-toolkit" {
					from := tag.(map[string]interface{})["from"]
					imageURL = from.(map[string]interface{})["name"].(string)
				}
			}

		}

		if header.Name == "release-manifests/release-metadata" {

			buff, err := io.ReadAll(tr)
			exit.OnError(err)

			obj := unstructured.Unstructured{}

			err = json.Unmarshal(buff, &obj.Object)
			exit.OnError(err)

			version, _, err = unstructured.NestedString(obj.Object, "version")
			exit.OnError(err)
		}

		if version != "" && imageURL != "" {
			break
		}

	}
	if version == "" || imageURL == "" {
		return "", ""
	}

	return version, imageURL
}

func dclose(c io.Closer) {
	if err := c.Close(); err != nil {
		log.Fatal(err)
	}
}

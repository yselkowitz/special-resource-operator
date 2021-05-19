package crane

import (
	"github.com/google/go-containerregistry/pkg/name"
)

func StrictValidation(o *options) {
	o.name = append(o.name, name.StrictValidation)
}

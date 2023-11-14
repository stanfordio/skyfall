package hydrator

import (
	"github.com/dgraph-io/ristretto"
)

type Hydrator struct {
	Cache *ristretto.Cache
}

func (h *Hydrator) Hydrate(val interface{}) (result interface{}, err error) {
	err = nil
	result = val

	return
}

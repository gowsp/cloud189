package drive

import (
	"sync"

	"github.com/gowsp/cloud189/pkg"
)

type Cache struct {
	id   sync.Map
	name sync.Map
}

func (c *Cache) Store(file pkg.File) {
	c.id.Store(file.Id, file)
	c.name.Store(file.Name, file)
}

func (c *Cache) Delete(file pkg.File) {
	c.id.Delete(file.Id)
	c.name.Delete(file.Name)
}

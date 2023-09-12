package drive

import (
	"sync"
	"time"

	"github.com/gowsp/cloud189/pkg"
	"github.com/gowsp/cloud189/pkg/file"
)

var nodes sync.Map

func load(id string) *node {
	if val, ok := nodes.Load(id); ok {
		return val.(*node)
	}
	if id == file.Root.Id() {
		return newNode(file.Root)
	}
	return nil
}

func newNode(file pkg.File) *node {
	node := &node{info: file}
	nodes.Store(file.Id(), node)
	return node
}

type node struct {
	info   pkg.File
	node   sync.Map
	exp    time.Time
	loaded bool
}

func (n *node) invalid() {
	n.loaded = false
}
func (n *node) valid() bool {
	return n.loaded && n.exp.After(time.Now())
}
func (n *node) enable() {
	n.exp = time.Now().Add(time.Minute * 1)
	n.loaded = true
}
func (n *node) add(children ...pkg.File) {
	for _, child := range children {
		node := newNode(child)
		n.node.Store(child.Name(), node)
	}
}
func (n *node) search(name string, loader func() (pkg.File, error)) (pkg.File, error) {
	if val, ok := n.node.Load(name); ok {
		return val.(*node).info, nil
	}
	result, err := loader()
	if err != nil {
		return nil, err
	}
	n.add(result)
	return result, nil
}

func (n *node) list(loader func() ([]pkg.File, error)) ([]pkg.File, error) {
	if n.valid() {
		result := make([]pkg.File, 0)
		n.node.Range(func(key, value any) bool {
			result = append(result, value.(*node).info)
			return true
		})
		return result, nil
	}
	result, err := loader()
	if err != nil {
		return nil, err
	}
	n.add(result...)
	n.enable()
	return result, nil
}

func (n *node) delete(child pkg.File) {
	p := load(child.Id())
	if child.IsDir() && p != nil {
		p.node.Range(func(key, value any) bool {
			p.delete(value.(*node).info)
			return true
		})
	}
	nodes.Delete(child.Id())
	n.node.Delete(child.Name())
	n.loaded = false
}
func invalid(files ...pkg.File) {
	for _, file := range files {
		if file == nil {
			continue
		}
		node := load(file.PId())
		if node == nil {
			continue
		}
		node.delete(file)
	}
}

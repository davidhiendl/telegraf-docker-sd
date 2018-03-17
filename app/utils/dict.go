package utils

import "strings"

// a kv-tree to hold values for the template
type Dict struct {
	values   map[string]string
	children map[string]*Dict
}

// Create new Dict
func NewDict() (*Dict) {
	d := Dict{}
	return &d
}

func (d *Dict) Get(key string) (string, bool) {
	val, ok := d.values[key];
	if !ok {
		return "", false
	}

	return val, true
}

func (d *Dict) GetOrDefault(key string, def string) (string) {
	val, ok := d.values[key];
	if !ok {
		return def
	}

	return val
}

func (d *Dict) Set(key string, value string) {
	if d.HasChild(key) {
		panic("Cannot access value if a child with the same key exists")
	}

	d.values[key] = value
}

// check if value exists
func (d *Dict) HasValue(key string) bool {
	_, ok := d.values[key]
	return ok
}

// check if child exists
func (d *Dict) HasChild(key string) bool {
	_, ok := d.children[key]
	return ok
}

// retrieve or create child at key
func (d *Dict) Child(key string) *Dict {
	if val, ok := d.children[key]; ok {
		return val
	}

	if d.HasValue(key) {
		panic("Cannot access child if a value with the same key exists")
	}

	child := NewDict()
	d.children[key] = child

	return child
}

// delete value or child at key
func (d *Dict) Del(key string) {
	// delete value
	if _, ok := d.values[key]; ok {
		delete(d.values, key);
	}

	// delete child
	if _, ok := d.children[key]; ok {
		delete(d.children, key);
	}
}

// traverse the dictionary and look for a specific value at the path end
func (d *Dict) PathGet(path string) (string, bool) {
	parts := strings.Split(path, ".")

	cur := d
	// traverse tree except for last element
	for i := 0; i < len(parts)-1; i++ {
		if !cur.HasChild(parts[i]) {
			return "", false
		}
		cur = cur.Child(parts[i])
	}

	// retrieve last element
	return cur.Get(parts[len(parts)-1])

}

// traverse the dictionary and look for a specific value at the path end
func (d *Dict) PathGetOrDefault(path string, def string) (string) {
	val, ok := d.PathGet(path)
	if !ok {
		return def
	}

	return val
}

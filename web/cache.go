package web

import (
	cmap "github.com/orcaman/concurrent-map"
)

// ----------------------------------------------------------------------------
// ConcurrentMap Natives

// NewCMap returns an embedded k-v store (`Cache`).
func NewCMap() cmap.ConcurrentMap {
	return cmap.New()
}

// ************************************
//  Golang forbids non-local receivers
// ************************************

// SetCMap any type per `key`
func SetCMap(c *cmap.ConcurrentMap, key string, val interface{}) {
	c.Set(key, val)
}

// GetResourceCMap (`*Resource`) per `key`
func GetResourceCMap(c *cmap.ConcurrentMap, key string) *Resource {
	var rs *Resource
	v, ok := c.Get(key)
	if ok {
		rs = v.(*Resource)
	}
	return rs
}

// GetStrCMap (`*string`) per `key`
func GetStrCMap(c *cmap.ConcurrentMap, key string) *string {
	var s *string
	v, ok := c.Get(key)
	if ok {
		s = v.(*string)
	}
	return s
}

// GetBoolCMap (`*bool`) per `key`
func GetBoolCMap(c *cmap.ConcurrentMap, key string) *bool {
	var b *bool
	v, ok := c.Get(key)
	if ok {
		b = v.(*bool)
	}
	return b
}

// ----------------------------------------------------------------------------
// Wrappers for `cmap.ConcurrentMap` to decouple caller from the native pkg.

// Cache ...
type Cache cmap.ConcurrentMap

// NewCache returns an embedded k-v store (`Cache`).
func NewCache() Cache {
	return Cache(cmap.New())
}

// ConcurrentMap Adapters

// Set any type `val` per `key`
func (c *Cache) Set(key string, val interface{}) {
	cc := cmap.ConcurrentMap(*c)
	cc.Set(key, val)
}

// GetResource (`*Resource`) per `key`. Return zero-value if `!ok`.
func (c *Cache) GetResource(key string) *Resource {
	var (
		cc = cmap.ConcurrentMap(*c)
		rs *Resource
	)
	v, ok := cc.Get(key)
	if ok {
		rs = v.(*Resource)
	}
	return rs
}

// GetStr (`*string`) per `key`. Return zero-value if `!ok`.
func (c *Cache) GetStr(key string) *string {
	var (
		cc = cmap.ConcurrentMap(*c)
		s  *string
	)
	v, ok := cc.Get(key)
	if ok {
		s = v.(*string)
	}
	return s
}

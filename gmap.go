// Copyright 2023 james dotter. All rights reserved.typVal
// Use of this source code is governed by a
// license that can be found in the gotype LICENSE file.

package gotype

import "sort"

// Gmap is in essences an ordered map
type Gmap []GmapEl

type GmapEl struct {
	Key   string
	Value VALUE
}

// Map returns gotype VALUE as map[string]VALUE
func (v VALUE) Gmap() Gmap {
	switch v.KIND() {
	case Map:
		return (MAP)(v).Gmap()
	case Struct:
		return (STRUCT)(v).Gmap()
	}
	panic("cannot convert value to Gmap")
}

func (g *Gmap) Set(key string, a any) {
	value := ValueOfV(a)
	for i, el := range *g {
		if el.Key == key {
			(*g)[i].Value = value
			return
		}
	}
	*g = append(*g, GmapEl{key, value})
}

func (g *Gmap) Get(key string) (VALUE, bool) {
	for _, el := range *g {
		if el.Key == key {
			return el.Value, true
		}
	}
	return VALUE{}, false
}

func (g *Gmap) Del(key string) {
	for i, el := range *g {
		if el.Key == key {
			*g = append((*g)[:i], (*g)[i+1:]...)
			return
		}
	}
}

func (g *Gmap) Keys() []string {
	keys := make([]string, len(*g))
	for i, el := range *g {
		keys[i] = el.Key
	}
	return keys
}

func (g *Gmap) Values() []VALUE {
	values := make([]VALUE, len(*g))
	for i, el := range *g {
		values[i] = el.Value
	}
	return values
}

func (g *Gmap) SortByKeys() *Gmap {
	sort.Slice(*g, func(i, j int) bool {
		return (*g)[i].Key < (*g)[j].Key
	})
	return g
}

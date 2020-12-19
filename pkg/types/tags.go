package types

import (
	"reflect"
	"strings"
	"unicode"
	"unicode/utf8"
)

const (
	tagId         = "watson"
	attrOmitEmpty = "omitempty"
)

type tag struct {
	name      string
	f         *reflect.StructField
	omitempty bool
}

func parseTag(f *reflect.StructField) *tag {
	tag := &tag{f: f}
	name := f.Tag.Get(tagId)
	if name == "" {
		return tag
	}
	attrs := strings.Split(name, ",")
	first := true
	for _, attr := range attrs {
		if first {
			tag.name = attr
			first = false
			continue
		}
		switch attr {
		case attrOmitEmpty:
			tag.omitempty = true
		}
	}
	return tag
}

func (t *tag) Key() string {
	if t.name == "" {
		return strings.ToLower(t.f.Name)
	}
	return t.name
}

func (t *tag) ShouldAlwaysOmit() bool {
	r, _ := utf8.DecodeRuneInString(t.f.Name)
	return unicode.IsLower(r)
}

func (t *tag) OmitEmpty() bool {
	return t.omitempty
}
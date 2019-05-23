package ktt

import (
	"path/filepath"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

func TestKTT(t *testing.T) {
	tpl, err := LoadTemplate(filepath.Join(".", "testdata", "ktt"))
	require.NoError(t, err, "unexpected error")
	spew.Dump(tpl)
}

func loadValues(b []byte) (Values, error) {
	r := make(map[interface{}]interface{}, 64)
	if err := yaml.Unmarshal(b, &r); err != nil {
		return nil, err
	}
	return r, nil
}

func TestValues(t *testing.T) {
	v1, err := loadValues([]byte(`
abc: 123
xyz:
  x: 1
  "y": 2
  z: 3
  inner:
    inside:
      val: 123
`))
	require.NoError(t, err, "unexpected error")
	v2, err := loadValues([]byte(`
abc: 42
hello: world
xyz:
  x: 42
  new: 111
  inner:
    val: 42
    inside: 42
  new: field
`))
	require.NoError(t, err, "unexpected error")
	exp, err := loadValues([]byte(`
hello: world
abc: 42
xyz:
  x: 42
  "y": 2
  z: 3
  new: 111
  inner:
    val: 42
    inside: 42
  new: field
`))
	require.NoError(t, err, "unexpected error")
	MergeValues(v1, v2)
	require.Equal(t, exp, v1, "mismatch")
}

func TestSetValues(t *testing.T) {
	v1, err := loadValues([]byte(`
abc: 123
xyz:
  z: 3
  inner:
    inside:
      val: 123
`))
	require.NoError(t, err, "unexpected error")
	SetValuesKey(v1, "new", "value")
	SetValuesKey(v1, "xyz.z", 42)
	SetValuesKey(v1, "xyz.inner.inside.val.new", 42)
	exp, err := loadValues([]byte(`
abc: 123
new: value
xyz:
  z: 42
  inner:
    inside:
      val:
        new: 42
`))
	require.NoError(t, err, "unexpected error")
	require.Equal(t, exp, v1, "mismatch")
}

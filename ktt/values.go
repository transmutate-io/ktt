package ktt

import (
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

type Values map[interface{}]interface{}

func LoadValuesFile(dir, file string) (map[interface{}]interface{}, error) {
	r := make(map[interface{}]interface{}, 64)
	if err := loadYAMLFile(filepath.Join(dir, file), r); err != nil {
		return nil, err
	}
	return r, nil
}

func MergeValues(a, b map[interface{}]interface{}) map[interface{}]interface{} {
	for vk, vv := range b {
		// setting a map
		if vm, ok := vv.(map[interface{}]interface{}); ok {
			// key doesn't exist
			vval, ok := a[vk]
			if !ok {
				a[vk] = vv
				continue
			}
			// key exists
			if tv, ok := vval.(map[interface{}]interface{}); ok {
				// it's a map
				MergeValues(tv, vm)
			} else {
				// it's a value
				a[vk] = vv
			}
		} else {
			a[vk] = vv
		}
	}
	return a
}

func SetValuesKey(vals map[interface{}]interface{}, keyPath string, val interface{}) {
	kp := strings.Split(keyPath, ".")
	for len(kp) > 1 {
		var createNew bool
		newVals, ok := vals[kp[0]]
		if !ok {
			createNew = true
		} else {
			if _, ok = newVals.(map[interface{}]interface{}); !ok {
				createNew = true
			}
		}
		if createNew {
			newVals = make(map[interface{}]interface{}, 4)
			vals[kp[0]] = newVals
		}
		vals = newVals.(map[interface{}]interface{})
		kp = kp[1:]
	}
	vals[kp[0]] = val
}

func ValuesYAML(v map[interface{}]interface{}) string {
	if v == nil {
		return ""
	}
	b, err := yaml.Marshal(v)
	if err != nil {
		panic(err)
	}
	return string(b)
}

package yamlutil

import (
	"bytes"
	"fmt"
	"io"
	"reflect"

	yaml "gopkg.in/yaml.v2"
)

func Query(in interface{}, path string, docIdx int) (interface{}, error) {
	var mappedDoc interface{}
	var dataBucket interface{}
	var currentIndex int
	var err error

	if docIdx < 0 {
		return nil, fmt.Errorf("document index %d out of range", docIdx)
	}

	decoder, err := newDecoder(in)
	for {
		err = decoder.Decode(&dataBucket)
		if err != nil {
			if err == io.EOF {
				if currentIndex <= docIdx {
					err = fmt.Errorf("document index %d out of range", docIdx)
				} else {
					err = nil
				}
			}
			break
		}
		if currentIndex == docIdx {
			mappedDoc = dataBucket
			if path != "" {
				parts := parsePath(path)
				mappedDoc, err = recurse(dataBucket, parts[0], parts[1:])
				if err != nil {
					break
				}
			}
		}
		currentIndex++
	}

	if err != nil {
		return nil, err
	}

	return mappedDoc, nil
}

func QueryAll(in interface{}, path string) (interface{}, error) {
	var mappedDocs []interface{}
	var dataBucket yaml.MapSlice
	var err error

	decoder, err := newDecoder(in)
	for {
		err = decoder.Decode(&dataBucket)
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			break
		}
		var mappedDoc interface{} = dataBucket
		if path != "" {
			parts := parsePath(path)
			mappedDoc, err = recurse(dataBucket, parts[0], parts[1:])
			if err != nil {
				break
			}
		}
		mappedDocs = append(mappedDocs, mappedDoc)
	}

	if err != nil {
		return nil, err
	}

	return mappedDocs, nil
}

func QueryInto(in interface{}, out interface{}, path string, docIdx int) error {
	if reflect.ValueOf(out).Kind() != reflect.Ptr {
		return fmt.Errorf("`out` is not a pointer")
	}
	res, err := Query(in, path, docIdx)
	if err != nil {
		return err
	}

	bytes, err := yaml.Marshal(res)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(bytes, out)
}

func QueryAllInto(in interface{}, out interface{}, path string) error {
	if reflect.ValueOf(out).Kind() != reflect.Ptr {
		return fmt.Errorf("`out` is not a pointer")
	}
	res, err := QueryAll(in, path)
	if err != nil {
		return err
	}

	bytes, err := yaml.Marshal(res)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(bytes, out)
}

func newDecoder(in interface{}) (*yaml.Decoder, error) {
	switch tin := in.(type) {
	case io.Reader:
		return yaml.NewDecoder(tin), nil
	case string:
		return yaml.NewDecoder(bytes.NewBufferString(tin)), nil
	case []byte:
		return yaml.NewDecoder(bytes.NewReader(tin)), nil
	default:
		encoded, err := yaml.Marshal(in)
		if err != nil {
			return nil, err
		}
		return yaml.NewDecoder(bytes.NewReader(encoded)), nil
	}
}

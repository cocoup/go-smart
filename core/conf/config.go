package conf

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
)

func MustLoad(path string, v interface{}) {
	if err := Load(path, v); err != nil {
		log.Fatalf("error: config file %s, %s", path, err.Error())
	}
}

func Load(path string, v interface{}) error {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(content, v)
}

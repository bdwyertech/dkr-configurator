package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"os"

	log "github.com/sirupsen/logrus"
)

var CONFIG_PATH = "/config"

func main() {
	var cfg string
	if b64cfg := os.Getenv("B64_CONFIG"); b64cfg != "" {
		decoded, err := base64.StdEncoding.DecodeString(b64cfg)
		if err != nil {
			log.Fatal(err)
		}
		if cfg, err = Render(bufio.NewScanner(bytes.NewReader(decoded))); err != nil {
			log.Fatal(err)
		}
	} else if ssmPath := os.Getenv("SSM_PATH"); ssmPath != "" {
		cfg = string(GetParametersByPathYAML(ssmPath))
	} else {
		log.Fatal("No configuration specified...")
	}
	if cfgPath := os.Getenv("CONFIG_PATH"); cfgPath != "" {
		CONFIG_PATH = cfgPath
	}
	if err := os.WriteFile(CONFIG_PATH, []byte(cfg), os.ModePerm); err != nil {
		log.Fatal(err)
	}
}

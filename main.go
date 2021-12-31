package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"
)

var CONFIG_PATH = "/config"

func main() {
	if cfgPath := os.Getenv("CONFIGURATOR_PATH"); cfgPath != "" {
		CONFIG_PATH = cfgPath
	}
	var cfg string
	if o, err := os.Stat(CONFIG_PATH); err == nil {
		f, err := os.Open(o.Name())
		if err != nil {
			log.Fatal(err)
		}
		if cfg, err = Render(bufio.NewScanner(f)); err != nil {
			log.Fatal(err)
		}
	} else if b64cfg := os.Getenv("CONFIGURATOR_B64"); b64cfg != "" {
		decoded, err := base64.StdEncoding.DecodeString(b64cfg)
		if err != nil {
			log.Fatal(err)
		}
		if cfg, err = Render(bufio.NewScanner(bytes.NewReader(decoded))); err != nil {
			log.Fatal(err)
		}
	} else if ssmPath := os.Getenv("CONFIGURATOR_SSM_PATH"); ssmPath != "" {
		if os.Getenv("CONFIGURATOR_FORMAT") == "json" {
			cfg = string(GetParametersByPathJSON(ssmPath))
		} else {
			cfg = string(GetParametersByPathYAML(ssmPath))
		}
	} else {
		log.Fatal("No configuration specified...")
	}
	if err := os.WriteFile(CONFIG_PATH, []byte(cfg), os.ModePerm); err != nil {
		log.Fatal(err)
	}

	uid, gid := os.Getenv("CONFIGURATOR_UID"), os.Getenv("CONFIGURATOR_GID")
	if uid != "" || gid != "" {
		uid, err := strconv.Atoi(uid)
		if err != nil {
			uid = -1
		}
		gid, err := strconv.Atoi(gid)
		if err != nil {
			gid = -1
		}
		if err = os.Chown(CONFIG_PATH, uid, gid); err != nil {
			log.Fatal(err)
		}
	}
}

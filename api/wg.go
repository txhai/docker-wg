package api

import (
	"io/ioutil"
	"os"
	"sync"
)

const tmpDir = "/tmp"

type wireguard struct {
	sync.Mutex
	api *Api
}

func (api *Api) initWgApi() {
	wg := new(wireguard)
	wg.api = api
	api.wg = wg
}

func (wg *wireguard) saveConf(itf string) {
	if err := execCommand("wg-quick save %s", itf); err != nil {
		wg.api.log.Errorf("save conf exec %v", err)
	}
	wg.api.log.Printf("save conf success")
}

func (wg *wireguard) addConf(itf string, conf string) {
	api := wg.api
	var err error
	file, err := ioutil.TempFile(tmpDir, "wgc.*.conf")
	if err != nil {
		api.log.Errorf("add conf %v", err)
		return
	}
	defer os.Remove(file.Name())
	if _, err = file.Write([]byte(conf)); err != nil {
		api.log.Errorf("add conf write %v", err)
		return
	}
	if err = execCommand("wireguard addconf %s %s", itf, file.Name()); err != nil {
		api.log.Errorf("add conf (%s) exec error %v", conf, err)
		return
	}
	api.log.Printf("add conf (%s) success", conf)
}

func (wg *wireguard) removeConf(itf string, publicKey string) {
	if err := execCommand("wg set %s peer %s remove", itf, publicKey); err != nil {
		wg.api.log.Errorf("remove conf (%s, %s) exec %v", itf, publicKey, err)
	}
	wg.api.log.Printf("remove conf success")
}

package api

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strings"
	"sync"
)

const tmpDir = "/tmp"

type wireguard struct {
	sync.Mutex // sync all wireguard-tools call
	api        *Api
}

func (api *Api) initWgApi() {
	wg := new(wireguard)
	wg.api = api
	api.wg = wg
}

func (wg *wireguard) saveConf(itf string) {
	wg.Lock()
	defer wg.Unlock()
	if err := execCmd("wg-quick save %s", itf); err != nil {
		wg.api.log.Errorf("save conf exec %v", err)
	}
	wg.api.log.Printf("save conf success")
}

func (wg *wireguard) addConf(itf string, conf string) {
	wg.Lock()
	defer wg.Unlock()
	api := wg.api
	var err error
	file, err := ioutil.TempFile(tmpDir, "wgc.*.conf")
	if err != nil {
		api.log.Errorf("add conf tempfile %v", err)
		return
	}
	defer os.Remove(file.Name())
	if _, err = file.Write([]byte(conf)); err != nil {
		api.log.Errorf("add conf write %v", err)
		return
	}
	if err = execCmd("wireguard addconf %s %s", itf, file.Name()); err != nil {
		api.log.Errorf("add conf (%s) exec error %v", conf, err)
		return
	}
	api.log.Printf("add conf (%s) success", conf)
}

func (wg *wireguard) removeConf(itf string, publicKey string) {
	wg.Lock()
	defer wg.Unlock()
	if err := execCmd("wg set %s peer %s remove", itf, publicKey); err != nil {
		wg.api.log.Errorf("remove conf (%s, %s) exec %v", itf, publicKey, err)
	}
	wg.api.log.Printf("remove conf success")
}

func (wg *wireguard) getPublicKey(itf string) string {
	wg.Lock()
	defer wg.Unlock()
	output, err := execCmdGetOutput("wg show %s public-key", itf)
	if err != nil {
		wg.api.log.Errorf("wg show %s public-key - exec %v", itf, err)
		return ""
	}
	for _, s := range strings.Split(output, "\n") {
		if strings.TrimSpace(s) != "" {
			return s
		}
	}
	wg.api.log.Errorf("wg show %s public-key - empty output", itf)
	return ""
}

func (wg *wireguard) getPeerIPs(itf string) map[string]string {
	wg.Lock()
	defer wg.Unlock()
	output, err := execCmdGetOutput("wg show %s allowed-ips", itf)
	if err != nil {
		wg.api.log.Errorf("wg.getPeerIPs(%s) - exec %v", itf, err)
		return nil
	}
	var dict map[string]string = nil
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		parts := strings.Split(line, "\t")
		if len(parts) != 2 {
			continue
		}
		if dict == nil {
			dict = make(map[string]string)
		}
		_, ip, err := net.ParseCIDR(parts[1])
		if err != nil {
			wg.api.log.Errorf("wg.getPeerIPs(%s) - line (%s) - ParseCIDR %v", itf, line, err)
			continue
		}
		dict[parts[0]] = ip.IP.String()
	}
	return dict
}
func (wg *wireguard) dumps(itf string) []peerUsage {
	wg.Lock()
	defer wg.Unlock()
	api := wg.api
	output, err := execCmdGetOutput("wg show %s dump", itf)
	if err != nil {
		api.log.Errorf("wg.dumps(%s) - exec %v", itf, err)
		return nil
	}
	var peers []peerUsage
	for _, line := range strings.Split(output, "\n") {
		parts := strings.Split(line, "\t")
		if len(parts) != 8 {
			api.log.Errorf("wg.dumps(%s) - line(%s) - expect 8 found %d", itf, line, len(parts))
			continue
		}
		peer := peerUsage{
			parseStr(parts[0]),
			parseStr(parts[2]),
			parseStr(parts[3]),
			parseInt64(parts[4]),
			parseInt64(parts[5]),
			parseInt64(parts[6]),
			parseStr(parts[7]),
		}
		peers = append(peers, peer)
	}
	return peers
}

func (wg *wireguard) isLinkUp(itf string) bool {
	f, err := os.Open(fmt.Sprintf("/sys/class/net/%s/carrier", itf))
	if err != nil {
		return false
	}
	defer f.Close()
	return true
}

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

func (wg *wireguard) saveConf(itf string) error {
	wg.Lock()
	defer wg.Unlock()
	if err := execCmd("wg-quick save %s", itf); err != nil {
		return fmt.Errorf("save conf exec %v", err)
	}
	wg.api.log.Printf("save conf success")
	return nil
}

func (wg *wireguard) addConf(itf string, conf string) error {
	wg.Lock()
	defer wg.Unlock()
	api := wg.api
	var err error
	file, err := ioutil.TempFile(tmpDir, "wgc.*.conf")
	if err != nil {
		return fmt.Errorf("add conf tempfile %v", err)
	}
	defer os.Remove(file.Name())
	if _, err = file.Write([]byte(conf)); err != nil {
		return fmt.Errorf("add conf write tempfile %v", err)
	}
	if err = execCmd("wg addconf %s %s", itf, file.Name()); err != nil {
		return fmt.Errorf("add conf (%s) exec error %v", conf, err)
	}
	api.log.Printf("add conf (%s) success", conf)
	return nil
}

func (wg *wireguard) removeConf(itf string, publicKey string) {
	wg.Lock()
	defer wg.Unlock()
	if err := execCmd("wg set %s peer %s remove", itf, publicKey); err != nil {
		wg.api.log.Errorf("remove conf (%s, %s) exec %v", itf, publicKey, err)
	}
	wg.api.log.Printf("remove conf success")
}

func (wg *wireguard) getPublicKey(itf string) (string, error) {
	wg.Lock()
	defer wg.Unlock()
	output, err := execCmdGetOutput("wg show %s public-key", itf)
	if err != nil {
		return "", fmt.Errorf("wg show %s public-key - exec %v", itf, err)
	}
	for _, s := range strings.Split(output, "\n") {
		if strings.TrimSpace(s) != "" {
			return s, nil
		}
	}
	return "", fmt.Errorf("wg show %s public-key - empty output", itf)
}

func (wg *wireguard) getPeerIPs(itf string) ([]net.IP, error) {
	wg.Lock()
	defer wg.Unlock()
	output, err := execCmdGetOutput("wg show %s allowed-ips", itf)
	if err != nil {
		return nil, fmt.Errorf("wg.getPeerIPs(%s) - exec %v", itf, err)
	}
	var ips = make([]net.IP, 0)
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		parts := strings.Split(line, "\t")
		if len(parts) != 2 {
			continue
		}
		_, ip, err := net.ParseCIDR(parts[1])
		if err != nil {
			wg.api.log.Errorf("wg.getPeerIPs(%s) - line (%s) - ParseCIDR %v", itf, line, err)
			continue
		}
		ips = append(ips, ip.IP)
	}
	return ips, nil
}
func (wg *wireguard) dumps(itf string) ([]peerUsage, error) {
	wg.Lock()
	defer wg.Unlock()
	output, err := execCmdGetOutput("wg show %s dump", itf)
	if err != nil {
		return nil, fmt.Errorf("wg.dumps(%s) - exec %v", itf, err)
	}
	var peers = make([]peerUsage, 0)
	for _, line := range strings.Split(output, "\n") {
		parts := strings.Split(line, "\t")
		if len(parts) != 8 {
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
	return peers, nil
}

func (wg *wireguard) isLinkUp(itf string) bool {
	f, err := os.Open(fmt.Sprintf("/sys/class/net/%s/carrier", itf))
	if err != nil {
		return false
	}
	defer f.Close()
	return true
}

package api

import (
	"docker_wg/loggin"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"golang.org/x/exp/slices"
	"net"
	"net/http"
	"strings"
	"sync"
)

type Api struct {
	wg              *wireguard
	log             *loggin.Logger
	peerMutex       sync.Mutex
	findAvailableIp func(ips []net.IP) string
	keepAlive       string
}

func NewApi(subnet string, keepAlive string, logger *loggin.Logger) (*Api, error) {
	_, ipNet, err := net.ParseCIDR(subnet)
	if err != nil {
		return nil, fmt.Errorf("net.ParseCIDR %v", err)
	}
	mask := binary.BigEndian.Uint32(ipNet.Mask)
	start := binary.BigEndian.Uint32(ipNet.IP)
	finish := (start & mask) | (mask ^ 0xffffffff)
	findAvailableIpFunc := func(ips []net.IP) string {
		j := -1
		for i := start; i <= finish; i++ {
			j++
			if j < 3 {
				continue
			}
			ip := make(net.IP, 4)
			binary.BigEndian.PutUint32(ip, i)
			idx := slices.IndexFunc(ips, func(pip net.IP) bool {
				return ip.Equal(pip)
			})
			if idx == -1 {
				return ip.String()
			}
		}
		return ""
	}
	api := new(Api)
	api.keepAlive = keepAlive
	api.log = logger
	api.findAvailableIp = findAvailableIpFunc
	api.initWgApi()
	return api, nil
}

func (api *Api) ServerKeyHandler(w http.ResponseWriter, r *http.Request) {
	itf := getInterface(r)
	key, err := api.wg.getPublicKey(itf)
	if err != nil {
		api.log.Errorf("ServerKeyHandler %v", err)
		responseError(w, http.StatusInternalServerError, err)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"public_key": key})
}

func (api *Api) AddPeerHandler(w http.ResponseWriter, r *http.Request) {
	wg := api.wg

	// get request data
	itf := getInterface(r)
	decoder := json.NewDecoder(r.Body)
	var payload map[string]string
	err := decoder.Decode(&payload)
	if err != nil {
		responseError(w, http.StatusBadRequest, err)
		return
	}
	peerPublicKey := payload["key"]
	if peerPublicKey == "" {
		responseError(w, http.StatusBadRequest, fmt.Errorf("empty peer public key"))
		return
	}

	api.peerMutex.Lock()
	defer api.peerMutex.Unlock()

	// find available ip for peer
	peerIps, err := wg.getPeerIPs(itf)
	if err != nil {
		responseError(w, http.StatusInternalServerError, fmt.Errorf("AddPeerHandler - %v", err))
		return
	}
	availableIp := api.findAvailableIp(peerIps)
	if availableIp == "" {
		responseError(w, http.StatusInternalServerError, fmt.Errorf("no IP available for peer"))
		return
	}

	// build & add peer conf
	conf := fmt.Sprintf("[Peer]\nPublicKey = %s\nAllowedIPs = %s/32\nPersistentKeepalive = %s\n", peerPublicKey, availableIp, api.keepAlive)
	err = wg.addConf(itf, conf)
	if err != nil {
		responseError(w, http.StatusInternalServerError, fmt.Errorf("addConf %v", err))
		return
	}
	err = wg.saveConf(itf)
	if err != nil {
		responseError(w, http.StatusInternalServerError, fmt.Errorf("saveConf %v", err))
		return
	}

	// add route
	err = ipRouteAdd(itf, availableIp)
	if err != nil {
		// maybe route already exists
		api.log.Errorf("AddPeerHandler.ipRouteAdd - %v", err)
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{"success": true, "ip": availableIp})
}

func (api *Api) RemovePeerHandler(w http.ResponseWriter, r *http.Request) {
	wg := api.wg
	itf := getInterface(r)
	peerPublicKey := r.FormValue("key")
	peerPublicKey = strings.TrimSpace(peerPublicKey)
	if peerPublicKey == "" {
		responseError(w, http.StatusBadRequest, fmt.Errorf("empty peer public key"))
		return
	}
	api.peerMutex.Lock()
	defer api.peerMutex.Unlock()
	err := wg.removeConf(itf, peerPublicKey)
	if err != nil {
		responseError(w, http.StatusInternalServerError, fmt.Errorf("RemovePeerHandler %v", err))
		return
	}
	err = wg.saveConf(itf)
	if err != nil {
		responseError(w, http.StatusInternalServerError, fmt.Errorf("RemovePeerHandler.saveConf %v", err))
		return
	}
	json.NewEncoder(w).Encode(map[string]interface{}{"success": true})
}

func (api *Api) ListPeerHandler(w http.ResponseWriter, r *http.Request) {
	itf := getInterface(r)
	peers, err := api.wg.dumps(itf)
	if err != nil {
		api.log.Errorf("ListPeerHandler %v", err)
		responseError(w, http.StatusInternalServerError, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string][]peerUsage{"peers": peers})
}

func (api *Api) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	itf := getInterface(r)
	isUp := api.wg.isLinkUp(itf)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]bool{itf: isUp})
}

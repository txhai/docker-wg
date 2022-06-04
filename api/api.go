package api

import (
	"docker_wg/loggin"
	"sync"
)

type Api struct {
	wg      *wireguard
	log     *loggin.Logger
	ipPools struct {
		sync.Mutex
	}
}

func NewApi(subnet string, logger *loggin.Logger) *Api {
	api := new(Api)
	api.log = logger
	api.initWgApi()
	return api
}

func (api *Api) addPeer(publicKey string) {
	api.wg.Lock()
	defer api.wg.Unlock()

}

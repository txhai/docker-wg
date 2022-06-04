package api

type peerUsage struct {
	PublicKey           string `json:"public_key"`
	Endpoint            string `json:"endpoint"`
	AllowedIps          string `json:"allowed_ips"`
	LatestHandshake     int64  `json:"latest_handshake"`
	Rx                  int64  `json:"rx"`
	Tx                  int64  `json:"tx"`
	PersistentKeepAlive string `json:"persistent_keepalive"`
}

type peersUsage struct {
	Peers []peerUsage `json:"peers"`
}

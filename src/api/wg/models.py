from typing import Optional, NamedTuple


class KeyPair(NamedTuple):
    public: str
    private: str


class PeerUsage(NamedTuple):
    public_key: str
    endpoint: Optional[str]
    allowed_ips: Optional[str]
    latest_handshake: Optional[int]
    rx: Optional[int]
    tx: Optional[int]
    persistent_keepalive: Optional[str]

    def to_dict(self):
        return {
            'public_key': self.public_key,
            'endpoint': self.endpoint,
            'allowed_ips': self.allowed_ips,
            'latest_handshake': self.latest_handshake,
            'rx': self.rx,
            'tx': self.tx,
            'persistent_keepalive': self.persistent_keepalive
        }

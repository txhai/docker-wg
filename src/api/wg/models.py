import json
from typing import Optional


class Model:
    def to_dict(self):
        props = [p for p in dir(self) if not p.startswith('_') and not callable(getattr(self, p))]
        data = {}
        for prop in props:
            val = getattr(self, prop, None)
            data.update({prop: val})
        return data

    def to_json(self):
        return json.dumps(self.to_dict())


class PeerIdentity(Model):
    def __init__(self, public_key: str = None, private_key: Optional[str] = None):
        super().__init__()
        self.private_key = private_key
        self.public_key = public_key


class Peer(PeerIdentity):
    def __init__(self,
                 public_key: str,
                 private_key: Optional[str],
                 allowed_ips: Optional[str],
                 config: Optional[str]):
        super().__init__(public_key, private_key)
        self.allowed_ips = allowed_ips
        self.config = config


class PeerUsage(Peer):
    def __init__(self,
                 public_key: str,
                 private_key: Optional[str],
                 endpoint: Optional[str],
                 allowed_ips: Optional[str],
                 latest_handshake: Optional[int],
                 rx: Optional[int],
                 tx: Optional[int],
                 persistent_keepalive: Optional[str]):
        super().__init__(public_key, private_key, allowed_ips, None)
        self.rx = rx
        self.tx = tx
        self.endpoint = endpoint
        self.latest_handshake = latest_handshake
        self.persistent_keepalive = persistent_keepalive

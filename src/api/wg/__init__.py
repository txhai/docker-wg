from .ip import get_available_ip
from .models import Peer
from .wg import *

from ..logger import Logger, DEBUG

logger = Logger('wg', DEBUG)


def register_peer(interface: str) -> Peer:
    # create key pair
    key = create_keypair()
    assert key.public_key and len(key.public_key) > 0
    assert key.private_key and len(key.private_key) > 0

    # get available ip
    ips = get_peer_ips(interface)
    peer_ip = get_available_ip(list(ips.values()))

    # init peer config
    config = '[Peer]\n' \
             f'PublicKey = {key.public_key}\n' \
             f'AllowedIPs = {peer_ip}/32\n'

    # set new peer config to wg
    add_conf(interface, config)

    # save config
    save_conf(interface)

    return Peer(
        private_key=key.private_key,
        public_key=key.public_key,
        allowed_ips=peer_ip,
        config=config)


def unregister_peer(interface: str, key: str):
    # Remove peer from current config
    remove_peer(interface, key)
    # Save config
    save_conf(interface)


def list_peers(interface) -> List[PeerUsage]:
    return dumps(interface)

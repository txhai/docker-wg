import json

from .wg import *


from .logger import Logger, DEBUG

logger = Logger('wg', DEBUG)


class Action:
    CREATE_PEER = 0
    ADD_PEER = 1
    REMOVE_PEER = 2

    @staticmethod
    def get_name(action):
        if action == Action.CREATE_PEER:
            return "CREATE_PEER"
        if action == Action.ADD_PEER:
            return "ADD_PEER"
        if action == Action.REMOVE_PEER:
            return "REMOVE_PEER"


def add_peer(interface, public_key) -> str:
    ips = get_peer_ips(interface)
    # logger.info(ips)
    # logger.info(json.dumps({'str': public_key}))
    # logger.info(json.dumps({'str': public_key.strip()}))
    _ip = ips[public_key] if public_key in ips else get_available_ip(list(ips.values()))

    # create peer config
    config = '[Peer]\n' \
             f'PublicKey = {public_key}\n' \
             f'AllowedIPs = {_ip}/32\n'

    # add to wireguard
    add_conf(interface, config)
    save_conf(interface)
    ip_route_add(interface, f'{_ip}/32')

    return _ip


def create_peer(interface):
    keypair = create_keypair()
    # add peer
    _ip = add_peer(interface, keypair.public_key)
    return {
        "private_key": keypair.private_key,
        "public_key": keypair.public_key,
        "ip": _ip
    }


def remove_peer(interface, public_key):
    remove_conf(interface, public_key)
    save_conf(interface)

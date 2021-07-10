from .wg import *


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
    _ip = get_available_ip(ips)

    # create peer config
    config = '[Peer]\n' \
             f'PublicKey = {public_key}\n' \
             f'AllowedIPs = {_ip}/32\n'

    # add to wireguard
    add_conf(interface, config)
    save_conf(interface)

    return _ip


def create_peer(interface):
    keypair = create_keypair()
    # add peer
    _ip = add_peer(interface, keypair.public)
    return {
        "private_key": keypair.private,
        "public_key": keypair.public,
        "ip": _ip
    }


def remove_peer(interface, public_key):
    remove_conf(interface, public_key)
    save_conf(interface)

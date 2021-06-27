import os
import tempfile
import uuid
from typing import Dict, List

from .executor import execute
from .models import PeerIdentity, PeerUsage
from .utils import parse_val


def create_keypair() -> PeerIdentity:
    (private_fd, private_key_path) = tempfile.mkstemp(prefix='wg', suffix=f'{uuid.uuid4().hex}')
    (public_fd, public_key_path) = tempfile.mkstemp(prefix='wg', suffix=f'{uuid.uuid4().hex}')
    try:
        execute(f'wg genkey > {private_key_path} && '
                f'wg pubkey < {private_key_path} > {public_key_path}', shell=True)
        with open(private_fd, 'r') as rd:
            private_key = rd.read()
        with open(public_fd, 'r') as rd:
            public_key = rd.read()
    finally:
        os.remove(private_key_path)
        os.remove(public_key_path)
    return PeerIdentity(private_key, public_key)


def get_peer_ips(interface: str) -> Dict[str, str]:
    output, _ = execute(f'wg show {interface} allowed-ips')
    data = {}
    if output and len(output) > 0:
        for line in output.split('\n'):
            l = line.split('\t')
            if len(l) != 2:
                continue
            peer, ip = l[0], l[1]
            data[peer] = ip[:ip.index('/')]
    return data


def add_conf(interface: str, config: str):
    (fd, config_path) = tempfile.mkstemp(prefix='wg.config', suffix=f'{uuid.uuid4().hex}')
    try:
        with open(fd, 'w') as wd:
            wd.write(config)
        execute(f'wg addconf {interface} {config_path}')
    finally:
        os.remove(config_path)


def save_conf(interface):
    execute(f'wg-quick save {interface}')


def remove_peer(interface, key):
    execute(f'wg set {interface} peer {key} remove')


def dumps(interface) -> List[PeerUsage]:
    output, _ = execute(f'wg show {interface} dump')
    peers: List[PeerUsage] = []
    if not output or len(output) == 0:
        return peers
    lines = output.split('\n')[1:]
    for line in lines:
        l = line.split('\t')
        if len(l) != 8:
            continue
        peer = PeerUsage(
            public_key=parse_val(l[0]),
            private_key=None,
            endpoint=parse_val(l[2]),
            allowed_ips=parse_val(l[3]),
            latest_handshake=parse_val(l[4], int),
            rx=parse_val(l[5], int),
            tx=parse_val(l[6], int),
            persistent_keepalive=parse_val(l[7]))
        peers.append(peer)
    return peers

import os
from typing import List

from netaddr import IPNetwork

subnet = os.environ.get('INTERNAL_SUBNET', '10.13.13.0/24')

ip_pool = IPNetwork(subnet)


def get_available_ip(exists: List[str]):
    for idx, ip in enumerate(ip_pool):
        if idx < 3:
            continue
        if str(ip) not in exists:
            return str(ip)
    raise EOFError('No IP available')

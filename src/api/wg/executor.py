import subprocess as sp
from typing import Tuple

from ..logger import Logger, DEBUG

logger = Logger('wg', DEBUG)


def execute(command: str, shell=False) -> Tuple[str, str]:
    if shell:
        sp.run(command, shell=True)
    else:
        result = sp.run([*command.split(' ')], capture_output=True, text=True)
        stdout = result.stdout if result.stdout else ''
        stderr = result.stderr if result.stderr else ''
        return stdout, stderr

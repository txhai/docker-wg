import json
import queue
import time
import uuid
from multiprocessing import Process, Manager
from typing import Dict, TYPE_CHECKING, Optional

from .logger import Logger, DEBUG
from .wg import register_peer, unregister_peer

if TYPE_CHECKING:
    from multiprocessing import Queue as Q

logger = Logger('task', DEBUG)


class Task:
    NEW_PEER = 0
    REMOVE_PEER = 1

    task_id: str

    def __init__(self, cmd: int, interface: str, peer: Optional[str] = None, task_id: Optional[str] = None):
        self.command = cmd
        self.interface = interface
        self.peer = peer
        self.task_id = f'{uuid.uuid4()}' if not task_id else task_id

    def to_json(self):
        return json.dumps({
            'task_id': self.task_id,
            'command': self.command,
            'interface': self.interface,
            'peer': self.peer
        })

    @staticmethod
    def from_json(data):
        d: Dict = json.loads(data)
        return Task(d.get('command', None), d.get('interface', None), d.get('peer', None), d['task_id'])


class Queue(Process):
    def __init__(self):
        super().__init__()
        self.manager = Manager()
        self.tasks: Dict = self.manager.dict()
        self.queue: Q = self.manager.Queue()

    def exec_task(self, task: Task, timeout=0):
        task_id = task.task_id

        self.tasks[task_id] = None

        try:
            self.queue.put(task.to_json(), timeout=5)
            logger.info("Put to queue")
        except Exception as e:
            logger.error(e)
            raise e

        start_ts = time.time()
        result = None
        while True:
            logger.info(f"Still wait [{task_id}]")
            if timeout != 0 and time.time() - start_ts > timeout:
                raise TimeoutError()
            if self.tasks[task_id] is None:
                time.sleep(0.1)
                continue
            if not self.tasks[task_id]:
                # Task failed
                raise RuntimeError('Execute task failed')
            result = self.tasks[task_id]
            del self.tasks[task_id]
            break
        return result

    def process(self, task: Task):
        task_id = task.task_id
        logger.info(f"[.] Got task {task_id}")
        try:
            if task.command == Task.NEW_PEER:
                peer = register_peer(task.interface)
                self.tasks[task_id] = peer.to_json()
            elif task.command == task.REMOVE_PEER:
                assert task.peer is not None
                unregister_peer(task.interface, task.peer)
                self.tasks[task_id] = True
            else:
                self.tasks[task_id] = True
        except Exception as e:
            logger.error(e)
            self.tasks[task_id] = False
        logger.info(f"[!] Done task {task_id}")

    def reset_queue(self):
        self.queue.close()
        for k in self.tasks.keys():
            self.tasks[k] = False
        self.queue: Q = self.manager.Queue()

    def run(self) -> None:
        logger.info("[!] Started task queue")
        while True:
            try:
                task = self.queue.get(block=False)
                self.process(Task.from_json(task))
            except queue.Empty:
                time.sleep(0.1)
                continue
            except Exception as e:
                logger.error(f'Queue error: {e}')
                self.reset_queue()
                time.sleep(1)

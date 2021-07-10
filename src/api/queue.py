import queue
import time
from multiprocessing import Process, Manager
from typing import TYPE_CHECKING

from .logger import Logger, DEBUG
from .task import *

if TYPE_CHECKING:
    from multiprocessing import Queue as Q

logger = Logger('Queue', DEBUG)


class Queue(Process):
    def __init__(self):
        super().__init__()
        self.manager = Manager()
        self.tasks: Dict[str] = self.manager.dict()
        self.queue: Q = self.manager.Queue()
        self.counter = self.manager.Value('i', 0)

    def wait_for(self, task: Task, timeout=0) -> TaskResult:
        task_id = task.task_id
        self.tasks[task_id] = None
        try:
            self.queue.put(serialize_task(task), timeout=5)
            logger.info(f"[T][{task_id}] {task.name}")
        except Exception as e:
            logger.error(e)
            raise e
        # Wait until task done
        start_ts = time.time()
        result = None
        while True:
            # logger.info(f"[T][{task_id}] WAIT")
            if timeout != 0 and time.time() - start_ts > timeout:
                raise TimeoutError()
            if self.tasks[task_id] is None:
                # Not done yet, keep trying fetch
                time.sleep(0.05)
                continue
            else:
                result = deserialize_task_result(self.tasks[task_id])
            del self.tasks[task_id]
            break
        logger.info(f"[T][{task_id}] DONE")
        return result

    def process(self, encoded_task: str):
        task = deserialize_task(encoded_task)
        logger.info(f"[C][{task.task_id}] Consuming task {self.counter.get()}")
        result = task.execute()
        self.tasks[task.task_id] = serialize_task_result(result)
        logger.info(f"[C][{task.task_id}] Consumed [{'OK' if result.success else 'FAIL'}]")
        self.counter.set(self.counter.get() + 1)

    def reset_queue(self):
        self.queue.close()
        for k in self.tasks.keys():
            self.tasks[k] = False
        self.queue: Q = self.manager.Queue()

    def run(self) -> None:
        logger.info("[!] Started task queue")
        while True:
            try:
                encoded_task = self.queue.get(block=False)
                self.process(encoded_task)
            except queue.Empty:
                time.sleep(0.05)
                continue
            except Exception as e:
                logger.error(f'Queue error: {e}')
                self.reset_queue()
                time.sleep(1)

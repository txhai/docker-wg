import json
from typing import Dict, NamedTuple, Optional, Any
from uuid import uuid4
from .action import *


class TaskResult(NamedTuple):
    success: bool
    meta: Optional[Any]


class Task:
    def __init__(self, task_id: str, action: int, interface: str, meta: Optional[str]):
        self.task_id = task_id
        self.action = action
        self.interface = interface
        self.meta = meta

    @property
    def name(self) -> str:
        return Action.get_name(self.action)

    def execute(self) -> TaskResult:
        try:
            result = self._execute()
        except Exception as e:
            result = TaskResult(False, str(e))
        return result

    def _execute(self):
        if self.action == Action.CREATE_PEER:
            meta = create_peer(self.interface)
            return TaskResult(True, meta)

        if self.action == Action.ADD_PEER:
            meta = add_peer(self.interface, self.meta)
            return TaskResult(True, meta)

        if self.action == Action.REMOVE_PEER:
            remove_peer(self.interface, self.meta)
            return TaskResult(True, None)

        return TaskResult(False, None)

    @staticmethod
    def create(action: int, interface: str, meta: Optional[str] = None) -> "Task":
        task_id = f'{uuid4()}'
        return Task(task_id, action, interface, meta)


def serialize_task(task: Task) -> str:
    payload = {
        'task_id': task.task_id,
        'action': task.action,
        'interface': task.interface,
        'meta': task.meta
    }
    return json.dumps(payload, separators=(',', ':'))


def deserialize_task(payload: str) -> Task:
    data: Dict = json.loads(payload)
    return Task(data['task_id'], data['action'], data['interface'], data['meta'])


def serialize_task_result(result: TaskResult) -> str:
    payload = {
        'success': result.success,
        'meta': result.meta
    }
    return json.dumps(payload, separators=(',', ':'))


def deserialize_task_result(payload: str) -> TaskResult:
    data: Dict = json.loads(payload)
    return TaskResult(data['success'], data['meta'])

from typing import Dict

from flask import Flask, jsonify, request

from .action import Action
from .queue import Queue
from .task import Task
from .wg import dumps

app = Flask(__name__)
queue = Queue()


# We put registering peer, and removing peer in a queue to keep WireGuard configs consistently

@app.route('/<interface>/new_peer', methods=["POST"])
def new_peer(interface):
    t = Task.create(Action.CREATE_PEER, interface)
    result = queue.wait_for(t, timeout=10)
    if result.success:
        return jsonify({
            'success': True,
            **result.meta
        })
    return jsonify({'success': False}), 500


@app.route('/<interface>/add_peer', methods=["POST"])
def add_peer(interface):
    data: Dict = request.get_json(force=True, silent=True, cache=False)
    key = data.get('key', None)
    if not key:
        return jsonify({'success': False}), 400
    t = Task.create(Action.ADD_PEER, interface, key)
    result = queue.wait_for(t, timeout=10)
    if result.success:
        return jsonify({
            'success': True,
            'ip': result.meta
        })
    return jsonify({'success': False}), 500


@app.route('/<interface>/remove_peer', methods=["DELETE"])
def remove_peer(interface):
    key = request.args.get('key')
    if not key:
        return jsonify({'success': False}), 400
    t = Task.create(Action.REMOVE_PEER, interface, key)
    result = queue.wait_for(t, timeout=10)
    return jsonify({
        'success': result.success
    })


@app.route('/<interface>/peer/list', methods=["GET"])
def list_peer(interface):
    peers = dumps(interface)
    return jsonify({
        'peers': [p.to_dict() for p in peers]
    })

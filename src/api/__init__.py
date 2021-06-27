from flask import Flask, jsonify, request

from .task_queue import Queue, Task
from .wg import list_peers

app = Flask(__name__)
queue = Queue()


# We put registering peer, and removing peer in a queue to keep WireGuard configs consistently

@app.route('/<interface>/peer/register', methods=["POST"])
def register_peer(interface):
    peer_json = queue.exec_task(Task(Task.NEW_PEER, interface), timeout=10)
    return app.response_class(response=peer_json, status=200, mimetype='application/json')


@app.route('/<interface>/peer', methods=["DELETE"])
def unregister_peer(interface):
    peer = request.args.get('key')
    success = queue.exec_task(Task(Task.REMOVE_PEER, interface, peer=peer), timeout=10)
    return jsonify({'success': success}), 200


@app.route('/<interface>/peer/list', methods=["GET"])
def list_peer(interface):
    peers = list_peers(interface)
    return jsonify({
        'peers': [p.to_dict() for p in peers]
    })

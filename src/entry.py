import os

from waitress import serve

from api import app, queue

HOST = os.environ.get('HOST', '0.0.0.0')
PORT = os.environ.get('PORT', 8000)
NUM_THREAD = os.environ.get('NUM_THREAD', 4)

if __name__ == '__main__':
    queue.start()
    serve(app, host=HOST, port=PORT, threads=NUM_THREAD)

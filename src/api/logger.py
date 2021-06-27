import logging

DEBUG = logging.DEBUG


class Logger:
    def __init__(self, module_name, logging_level, enable_telegram=False):
        self.__logger = logging.getLogger('wg')
        self.__logging_level = logging_level
        self.__module_name = module_name
        self.__telegram = enable_telegram
        if len(self.__logger.handlers) == 0:
            # Stream handler
            handler = logging.StreamHandler()
            handler.setLevel(logging.DEBUG)
            handler.setFormatter(logging.Formatter(f'%(message)s'))
            self.__logger.addHandler(handler)

            self.__logger.propagate = False
            self.__logger.setLevel(logging.DEBUG)

    # std flush
    def flush(self):
        pass

    # std log
    def write(self, buf):
        self.__logger.log(self.__logging_level, f'[{self.__module_name}] {buf.rstrip()}')

    def info(self, msg):
        self.__logger.info(f'[{self.__module_name}] {msg}')

    def error(self, msg):
        self.__logger.error(f'[{self.__module_name}] {msg}')

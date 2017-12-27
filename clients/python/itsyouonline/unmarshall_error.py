class UnmarshallError(Exception):
    def __init__(self, resp, message=''):
        self.response = resp
        self.message = message

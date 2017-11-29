"""
Auto-generated class for KeyData
"""
from datetime import datetime
from six import string_types

from . import client_support


class KeyData(object):
    """
    auto-generated. don't touch.
    """

    @staticmethod
    def create(**kwargs):
        """
        :type algorithm: str
        :type comment: str
        :type timestamp: datetime
        :rtype: KeyData
        """

        return KeyData(**kwargs)

    def __init__(self, json=None, **kwargs):
        if json is None and not kwargs:
            raise ValueError('No data or kwargs present')

        class_name = 'KeyData'
        data = json or kwargs

        # set attributes
        data_types = [string_types]
        self.algorithm = client_support.set_property('algorithm', data, data_types, False, [], False, True, class_name)
        data_types = [string_types]
        self.comment = client_support.set_property('comment', data, data_types, False, [], False, False, class_name)
        data_types = [datetime]
        self.timestamp = client_support.set_property('timestamp', data, data_types, False, [], False, False, class_name)

    def __str__(self):
        return self.as_json(indent=4)

    def as_json(self, indent=0):
        return client_support.to_json(self, indent=indent)

    def as_dict(self):
        return client_support.to_dict(self)

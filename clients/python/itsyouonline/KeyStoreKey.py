"""
Auto-generated class for KeyStoreKey
"""
from .KeyData import KeyData
from .Label import Label
from six import string_types

from . import client_support


class KeyStoreKey(object):
    """
    auto-generated. don't touch.
    """

    @staticmethod
    def create(**kwargs):
        """
        :type globalid: str
        :type key: str
        :type keydata: KeyData
        :type label: Label
        :type username: str
        :rtype: KeyStoreKey
        """

        return KeyStoreKey(**kwargs)

    def __init__(self, json=None, **kwargs):
        if json is None and not kwargs:
            raise ValueError('No data or kwargs present')

        class_name = 'KeyStoreKey'
        data = json or kwargs

        # set attributes
        data_types = [string_types]
        self.globalid = client_support.set_property('globalid', data, data_types, False, [], False, False, class_name)
        data_types = [string_types]
        self.key = client_support.set_property('key', data, data_types, False, [], False, True, class_name)
        data_types = [KeyData]
        self.keydata = client_support.set_property('keydata', data, data_types, False, [], False, True, class_name)
        data_types = [Label]
        self.label = client_support.set_property('label', data, data_types, False, [], False, True, class_name)
        data_types = [string_types]
        self.username = client_support.set_property('username', data, data_types, False, [], False, False, class_name)

    def __str__(self):
        return self.as_json(indent=4)

    def as_json(self, indent=0):
        return client_support.to_json(self, indent=indent)

    def as_dict(self):
        return client_support.to_dict(self)

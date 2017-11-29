"""
Auto-generated class for Contract
"""
from .Party import Party
from .Signature import Signature
from datetime import datetime
from six import string_types

from . import client_support


class Contract(object):
    """
    auto-generated. don't touch.
    """

    @staticmethod
    def create(**kwargs):
        """
        :type content: str
        :type contractId: str
        :type contractType: str
        :type expires: datetime
        :type extends: list[str]
        :type invalidates: list[str]
        :type parties: list[Party]
        :type signatures: list[Signature]
        :rtype: Contract
        """

        return Contract(**kwargs)

    def __init__(self, json=None, **kwargs):
        if json is None and not kwargs:
            raise ValueError('No data or kwargs present')

        class_name = 'Contract'
        data = json or kwargs

        # set attributes
        data_types = [string_types]
        self.content = client_support.set_property('content', data, data_types, False, [], False, True, class_name)
        data_types = [string_types]
        self.contractId = client_support.set_property('contractId', data, data_types, False, [], False, True, class_name)
        data_types = [string_types]
        self.contractType = client_support.set_property('contractType', data, data_types, False, [], False, True, class_name)
        data_types = [datetime]
        self.expires = client_support.set_property('expires', data, data_types, False, [], False, True, class_name)
        data_types = [string_types]
        self.extends = client_support.set_property('extends', data, data_types, False, [], True, False, class_name)
        data_types = [string_types]
        self.invalidates = client_support.set_property('invalidates', data, data_types, False, [], True, False, class_name)
        data_types = [Party]
        self.parties = client_support.set_property('parties', data, data_types, False, [], True, True, class_name)
        data_types = [Signature]
        self.signatures = client_support.set_property('signatures', data, data_types, False, [], True, True, class_name)

    def __str__(self):
        return self.as_json(indent=4)

    def as_json(self, indent=0):
        return client_support.to_json(self, indent=indent)

    def as_dict(self):
        return client_support.to_dict(self)

"""
Auto-generated class for AddApiKeyReqBody
"""
from .Label import Label

from . import client_support


class AddApiKeyReqBody(object):
    """
    auto-generated. don't touch.
    """

    @staticmethod
    def create(**kwargs):
        """
        :type label: Label
        :rtype: AddApiKeyReqBody
        """

        return AddApiKeyReqBody(**kwargs)

    def __init__(self, json=None, **kwargs):
        if json is None and not kwargs:
            raise ValueError('No data or kwargs present')

        class_name = 'AddApiKeyReqBody'
        data = json or kwargs

        # set attributes
        data_types = [Label]
        self.label = client_support.set_property('label', data, data_types, False, [], False, True, class_name)

    def __str__(self):
        return self.as_json(indent=4)

    def as_json(self, indent=0):
        return client_support.to_json(self, indent=indent)

    def as_dict(self):
        return client_support.to_dict(self)

"""
Auto-generated class for RequiredScope
"""
from six import string_types

from . import client_support


class RequiredScope(object):
    """
    auto-generated. don't touch.
    """

    @staticmethod
    def create(**kwargs):
        """
        :type accessscopes: list[str]
        :type scope: str
        :rtype: RequiredScope
        """

        return RequiredScope(**kwargs)

    def __init__(self, json=None, **kwargs):
        if json is None and not kwargs:
            raise ValueError('No data or kwargs present')

        class_name = 'RequiredScope'
        data = json or kwargs

        # set attributes
        data_types = [string_types]
        self.accessscopes = client_support.set_property('accessscopes', data, data_types, False, [], True, True, class_name)
        data_types = [string_types]
        self.scope = client_support.set_property('scope', data, data_types, False, [], False, True, class_name)

    def __str__(self):
        return self.as_json(indent=4)

    def as_json(self, indent=0):
        return client_support.to_json(self, indent=indent)

    def as_dict(self):
        return client_support.to_dict(self)

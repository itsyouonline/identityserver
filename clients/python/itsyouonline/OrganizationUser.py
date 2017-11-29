"""
Auto-generated class for OrganizationUser
"""
from six import string_types

from . import client_support


class OrganizationUser(object):
    """
    auto-generated. don't touch.
    """

    @staticmethod
    def create(**kwargs):
        """
        :type missingscopes: list[str]
        :type role: str
        :type username: str
        :rtype: OrganizationUser
        """

        return OrganizationUser(**kwargs)

    def __init__(self, json=None, **kwargs):
        if json is None and not kwargs:
            raise ValueError('No data or kwargs present')

        class_name = 'OrganizationUser'
        data = json or kwargs

        # set attributes
        data_types = [string_types]
        self.missingscopes = client_support.set_property('missingscopes', data, data_types, False, [], True, True, class_name)
        data_types = [string_types]
        self.role = client_support.set_property('role', data, data_types, False, [], False, True, class_name)
        data_types = [string_types]
        self.username = client_support.set_property('username', data, data_types, False, [], False, True, class_name)

    def __str__(self):
        return self.as_json(indent=4)

    def as_json(self, indent=0):
        return client_support.to_json(self, indent=indent)

    def as_dict(self):
        return client_support.to_dict(self)

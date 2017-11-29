"""
Auto-generated class for GithubAccount
"""
from six import string_types

from . import client_support


class GithubAccount(object):
    """
    auto-generated. don't touch.
    """

    @staticmethod
    def create(**kwargs):
        """
        :type avatar_url: str
        :type html_url: str
        :type id: int
        :type login: str
        :type name: str
        :rtype: GithubAccount
        """

        return GithubAccount(**kwargs)

    def __init__(self, json=None, **kwargs):
        if json is None and not kwargs:
            raise ValueError('No data or kwargs present')

        class_name = 'GithubAccount'
        data = json or kwargs

        # set attributes
        data_types = [string_types]
        self.avatar_url = client_support.set_property('avatar_url', data, data_types, False, [], False, True, class_name)
        data_types = [string_types]
        self.html_url = client_support.set_property('html_url', data, data_types, False, [], False, True, class_name)
        data_types = [int]
        self.id = client_support.set_property('id', data, data_types, False, [], False, True, class_name)
        data_types = [string_types]
        self.login = client_support.set_property('login', data, data_types, False, [], False, True, class_name)
        data_types = [string_types]
        self.name = client_support.set_property('name', data, data_types, False, [], False, True, class_name)

    def __str__(self):
        return self.as_json(indent=4)

    def as_json(self, indent=0):
        return client_support.to_json(self, indent=indent)

    def as_dict(self):
        return client_support.to_dict(self)

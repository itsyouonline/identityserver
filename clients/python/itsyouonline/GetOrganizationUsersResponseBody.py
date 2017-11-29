"""
Auto-generated class for GetOrganizationUsersResponseBody
"""
from .OrganizationUser import OrganizationUser

from . import client_support


class GetOrganizationUsersResponseBody(object):
    """
    auto-generated. don't touch.
    """

    @staticmethod
    def create(**kwargs):
        """
        :type haseditpermissions: bool
        :type users: list[OrganizationUser]
        :rtype: GetOrganizationUsersResponseBody
        """

        return GetOrganizationUsersResponseBody(**kwargs)

    def __init__(self, json=None, **kwargs):
        if json is None and not kwargs:
            raise ValueError('No data or kwargs present')

        class_name = 'GetOrganizationUsersResponseBody'
        data = json or kwargs

        # set attributes
        data_types = [bool]
        self.haseditpermissions = client_support.set_property('haseditpermissions', data, data_types, False, [], False, True, class_name)
        data_types = [OrganizationUser]
        self.users = client_support.set_property('users', data, data_types, False, [], True, True, class_name)

    def __str__(self):
        return self.as_json(indent=4)

    def as_json(self, indent=0):
        return client_support.to_json(self, indent=indent)

    def as_dict(self):
        return client_support.to_dict(self)

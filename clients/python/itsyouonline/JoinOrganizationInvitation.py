"""
Auto-generated class for JoinOrganizationInvitation
"""
from .EnumJoinOrganizationInvitationMethod import EnumJoinOrganizationInvitationMethod
from .EnumJoinOrganizationInvitationRole import EnumJoinOrganizationInvitationRole
from .EnumJoinOrganizationInvitationStatus import EnumJoinOrganizationInvitationStatus
from datetime import datetime
from six import string_types

from . import client_support


class JoinOrganizationInvitation(object):
    """
    auto-generated. don't touch.
    """

    @staticmethod
    def create(**kwargs):
        """
        :type created: datetime
        :type emailaddress: str
        :type isorganization: bool
        :type method: EnumJoinOrganizationInvitationMethod
        :type organization: str
        :type phonenumber: str
        :type role: EnumJoinOrganizationInvitationRole
        :type status: EnumJoinOrganizationInvitationStatus
        :type user: str
        :rtype: JoinOrganizationInvitation
        """

        return JoinOrganizationInvitation(**kwargs)

    def __init__(self, json=None, **kwargs):
        if json is None and not kwargs:
            raise ValueError('No data or kwargs present')

        class_name = 'JoinOrganizationInvitation'
        data = json or kwargs

        # set attributes
        data_types = [datetime]
        self.created = client_support.set_property('created', data, data_types, False, [], False, False, class_name)
        data_types = [string_types]
        self.emailaddress = client_support.set_property('emailaddress', data, data_types, False, [], False, True, class_name)
        data_types = [bool]
        self.isorganization = client_support.set_property('isorganization', data, data_types, False, [], False, True, class_name)
        data_types = [EnumJoinOrganizationInvitationMethod]
        self.method = client_support.set_property('method', data, data_types, False, [], False, True, class_name)
        data_types = [string_types]
        self.organization = client_support.set_property('organization', data, data_types, False, [], False, True, class_name)
        data_types = [string_types]
        self.phonenumber = client_support.set_property('phonenumber', data, data_types, False, [], False, True, class_name)
        data_types = [EnumJoinOrganizationInvitationRole]
        self.role = client_support.set_property('role', data, data_types, False, [], False, True, class_name)
        data_types = [EnumJoinOrganizationInvitationStatus]
        self.status = client_support.set_property('status', data, data_types, False, [], False, True, class_name)
        data_types = [string_types]
        self.user = client_support.set_property('user', data, data_types, False, [], False, True, class_name)

    def __str__(self):
        return self.as_json(indent=4)

    def as_json(self, indent=0):
        return client_support.to_json(self, indent=indent)

    def as_dict(self):
        return client_support.to_dict(self)

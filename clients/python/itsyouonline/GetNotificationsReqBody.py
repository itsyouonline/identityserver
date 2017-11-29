"""
Auto-generated class for GetNotificationsReqBody
"""
from .ContractSigningRequest import ContractSigningRequest
from .JoinOrganizationInvitation import JoinOrganizationInvitation
from .MissingScopes import MissingScopes

from . import client_support


class GetNotificationsReqBody(object):
    """
    auto-generated. don't touch.
    """

    @staticmethod
    def create(**kwargs):
        """
        :type approvals: list[JoinOrganizationInvitation]
        :type contractRequests: list[ContractSigningRequest]
        :type invitations: list[JoinOrganizationInvitation]
        :type missingscopes: list[MissingScopes]
        :rtype: GetNotificationsReqBody
        """

        return GetNotificationsReqBody(**kwargs)

    def __init__(self, json=None, **kwargs):
        if json is None and not kwargs:
            raise ValueError('No data or kwargs present')

        class_name = 'GetNotificationsReqBody'
        data = json or kwargs

        # set attributes
        data_types = [JoinOrganizationInvitation]
        self.approvals = client_support.set_property('approvals', data, data_types, False, [], True, True, class_name)
        data_types = [ContractSigningRequest]
        self.contractRequests = client_support.set_property('contractRequests', data, data_types, False, [], True, True, class_name)
        data_types = [JoinOrganizationInvitation]
        self.invitations = client_support.set_property('invitations', data, data_types, False, [], True, True, class_name)
        data_types = [MissingScopes]
        self.missingscopes = client_support.set_property('missingscopes', data, data_types, False, [], True, True, class_name)

    def __str__(self):
        return self.as_json(indent=4)

    def as_json(self, indent=0):
        return client_support.to_json(self, indent=indent)

    def as_dict(self):
        return client_support.to_dict(self)

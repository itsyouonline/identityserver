"""
Auto-generated class for SetOrgOwnerReqBody
"""
from six import string_types

from . import client_support


class SetOrgOwnerReqBody(object):
    """
    auto-generated. don't touch.
    """

    @staticmethod
    def create(**kwargs):
        """
        :type orgowner: str
        :rtype: SetOrgOwnerReqBody
        """

        return SetOrgOwnerReqBody(**kwargs)

    def __init__(self, json=None, **kwargs):
        if json is None and not kwargs:
            raise ValueError('No data or kwargs present')

        class_name = 'SetOrgOwnerReqBody'
        data = json or kwargs

        # set attributes
        data_types = [string_types]
        self.orgowner = client_support.set_property('orgowner', data, data_types, False, [], False, True, class_name)

    def __str__(self):
        return self.as_json(indent=4)

    def as_json(self, indent=0):
        return client_support.to_json(self, indent=indent)

    def as_dict(self):
        return client_support.to_dict(self)

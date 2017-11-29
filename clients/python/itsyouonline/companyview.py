"""
Auto-generated class for companyview
"""
from datetime import datetime
from six import string_types

from . import client_support


class companyview(object):
    """
    auto-generated. don't touch.
    """

    @staticmethod
    def create(**kwargs):
        """
        :type expire: datetime
        :type globalid: str
        :type info: list[str]
        :type organizations: list[str]
        :type publicKeys: list[str]
        :type taxnr: str
        :rtype: companyview
        """

        return companyview(**kwargs)

    def __init__(self, json=None, **kwargs):
        if json is None and not kwargs:
            raise ValueError('No data or kwargs present')

        class_name = 'companyview'
        data = json or kwargs

        # set attributes
        data_types = [datetime]
        self.expire = client_support.set_property('expire', data, data_types, False, [], False, False, class_name)
        data_types = [string_types]
        self.globalid = client_support.set_property('globalid', data, data_types, False, [], False, True, class_name)
        data_types = [string_types]
        self.info = client_support.set_property('info', data, data_types, False, [], True, False, class_name)
        data_types = [string_types]
        self.organizations = client_support.set_property('organizations', data, data_types, False, [], True, False, class_name)
        data_types = [string_types]
        self.publicKeys = client_support.set_property('publicKeys', data, data_types, False, [], True, True, class_name)
        data_types = [string_types]
        self.taxnr = client_support.set_property('taxnr', data, data_types, False, [], False, False, class_name)

    def __str__(self):
        return self.as_json(indent=4)

    def as_json(self, indent=0):
        return client_support.to_json(self, indent=indent)

    def as_dict(self):
        return client_support.to_dict(self)

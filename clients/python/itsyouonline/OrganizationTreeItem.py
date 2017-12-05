"""
Auto-generated class for OrganizationTreeItem
"""
from six import string_types

from . import client_support


class OrganizationTreeItem(object):
    """
    auto-generated. don't touch.
    """

    @staticmethod
    def create(**kwargs):
        """
        :type children: list[OrganizationTreeItem]
        :type globalid: str
        :rtype: OrganizationTreeItem
        """

        return OrganizationTreeItem(**kwargs)

    def __init__(self, json=None, **kwargs):
        if json is None and not kwargs:
            raise ValueError('No data or kwargs present')

        class_name = 'OrganizationTreeItem'
        data = json or kwargs

        # set attributes
        data_types = [OrganizationTreeItem]
        self.children = client_support.set_property('children', data, data_types, False, [], True, True, class_name)
        data_types = [string_types]
        self.globalid = client_support.set_property('globalid', data, data_types, False, [], False, True, class_name)

    def __str__(self):
        return self.as_json(indent=4)

    def as_json(self, indent=0):
        return client_support.to_json(self, indent=indent)

    def as_dict(self):
        return client_support.to_dict(self)

"""
Auto-generated class for SeeVersion
"""
from six import string_types

from . import client_support


class SeeVersion(object):
    """
    auto-generated. don't touch.
    """

    @staticmethod
    def create(**kwargs):
        """
        :type category: str
        :type content_type: str
        :type creation_date: str
        :type end_date: str
        :type keystore_label: str
        :type link: str
        :type markdown_full_description: str
        :type markdown_short_description: str
        :type signature: str
        :type start_date: str
        :type version: int
        :rtype: SeeVersion
        """

        return SeeVersion(**kwargs)

    def __init__(self, json=None, **kwargs):
        if json is None and not kwargs:
            raise ValueError('No data or kwargs present')

        class_name = 'SeeVersion'
        data = json or kwargs

        # set attributes
        data_types = [string_types]
        self.category = client_support.set_property('category', data, data_types, False, [], False, True, class_name)
        data_types = [string_types]
        self.content_type = client_support.set_property('content_type', data, data_types, False, [], False, True, class_name)
        data_types = [string_types]
        self.creation_date = client_support.set_property('creation_date', data, data_types, False, [], False, True, class_name)
        data_types = [string_types]
        self.end_date = client_support.set_property('end_date', data, data_types, False, [], False, True, class_name)
        data_types = [string_types]
        self.keystore_label = client_support.set_property('keystore_label', data, data_types, False, [], False, True, class_name)
        data_types = [string_types]
        self.link = client_support.set_property('link', data, data_types, False, [], False, True, class_name)
        data_types = [string_types]
        self.markdown_full_description = client_support.set_property('markdown_full_description', data, data_types, False, [], False, True, class_name)
        data_types = [string_types]
        self.markdown_short_description = client_support.set_property('markdown_short_description', data, data_types, False, [], False, True, class_name)
        data_types = [string_types]
        self.signature = client_support.set_property('signature', data, data_types, False, [], False, True, class_name)
        data_types = [string_types]
        self.start_date = client_support.set_property('start_date', data, data_types, False, [], False, True, class_name)
        data_types = [int]
        self.version = client_support.set_property('version', data, data_types, False, [], False, True, class_name)

    def __str__(self):
        return self.as_json(indent=4)

    def as_json(self, indent=0):
        return client_support.to_json(self, indent=indent)

    def as_dict(self):
        return client_support.to_dict(self)

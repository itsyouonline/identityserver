"""
Auto-generated class for Ownerof
"""
from .EmailAddress import EmailAddress

from . import client_support


class Ownerof(object):
    """
    auto-generated. don't touch.
    """

    @staticmethod
    def create(**kwargs):
        """
        :type emailaddresses: list[EmailAddress]
        :rtype: Ownerof
        """

        return Ownerof(**kwargs)

    def __init__(self, json=None, **kwargs):
        if json is None and not kwargs:
            raise ValueError('No data or kwargs present')

        class_name = 'Ownerof'
        data = json or kwargs

        # set attributes
        data_types = [EmailAddress]
        self.emailaddresses = client_support.set_property('emailaddresses', data, data_types, False, [], True, True, class_name)

    def __str__(self):
        return self.as_json(indent=4)

    def as_json(self, indent=0):
        return client_support.to_json(self, indent=indent)

    def as_dict(self):
        return client_support.to_dict(self)

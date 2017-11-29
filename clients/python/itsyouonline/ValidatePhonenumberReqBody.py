"""
Auto-generated class for ValidatePhonenumberReqBody
"""
from six import string_types

from . import client_support


class ValidatePhonenumberReqBody(object):
    """
    auto-generated. don't touch.
    """

    @staticmethod
    def create(**kwargs):
        """
        :type validationkey: str
        :rtype: ValidatePhonenumberReqBody
        """

        return ValidatePhonenumberReqBody(**kwargs)

    def __init__(self, json=None, **kwargs):
        if json is None and not kwargs:
            raise ValueError('No data or kwargs present')

        class_name = 'ValidatePhonenumberReqBody'
        data = json or kwargs

        # set attributes
        data_types = [string_types]
        self.validationkey = client_support.set_property('validationkey', data, data_types, False, [], False, True, class_name)

    def __str__(self):
        return self.as_json(indent=4)

    def as_json(self, indent=0):
        return client_support.to_json(self, indent=indent)

    def as_dict(self):
        return client_support.to_dict(self)

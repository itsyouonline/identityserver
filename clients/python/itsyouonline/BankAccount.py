"""
Auto-generated class for BankAccount
"""
from .Label import Label
from six import string_types

from . import client_support


class BankAccount(object):
    """
    auto-generated. don't touch.
    """

    @staticmethod
    def create(**kwargs):
        """
        :type bic: str
        :type country: str
        :type iban: str
        :type label: Label
        :rtype: BankAccount
        """

        return BankAccount(**kwargs)

    def __init__(self, json=None, **kwargs):
        if json is None and not kwargs:
            raise ValueError('No data or kwargs present')

        class_name = 'BankAccount'
        data = json or kwargs

        # set attributes
        data_types = [string_types]
        self.bic = client_support.set_property('bic', data, data_types, False, [], False, True, class_name)
        data_types = [string_types]
        self.country = client_support.set_property('country', data, data_types, False, [], False, True, class_name)
        data_types = [string_types]
        self.iban = client_support.set_property('iban', data, data_types, False, [], False, True, class_name)
        data_types = [Label]
        self.label = client_support.set_property('label', data, data_types, False, [], False, True, class_name)

    def __str__(self):
        return self.as_json(indent=4)

    def as_json(self, indent=0):
        return client_support.to_json(self, indent=indent)

    def as_dict(self):
        return client_support.to_dict(self)

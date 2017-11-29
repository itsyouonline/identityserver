"""
Auto-generated class for CreateUserBankAccountReqBody
"""
from .BankAccount import BankAccount

from . import client_support


class CreateUserBankAccountReqBody(object):
    """
    auto-generated. don't touch.
    """

    @staticmethod
    def create(**kwargs):
        """
        :type type: BankAccount
        :rtype: CreateUserBankAccountReqBody
        """

        return CreateUserBankAccountReqBody(**kwargs)

    def __init__(self, json=None, **kwargs):
        if json is None and not kwargs:
            raise ValueError('No data or kwargs present')

        class_name = 'CreateUserBankAccountReqBody'
        data = json or kwargs

        # set attributes
        data_types = [BankAccount]
        self.type = client_support.set_property('type', data, data_types, False, [], False, True, class_name)

    def __str__(self):
        return self.as_json(indent=4)

    def as_json(self, indent=0):
        return client_support.to_json(self, indent=indent)

    def as_dict(self):
        return client_support.to_dict(self)

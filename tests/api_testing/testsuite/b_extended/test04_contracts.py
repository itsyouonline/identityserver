from api_testing.utils import BaseTest
import types
import unittest
from random import randint
import json
import time


class ContractsTestsB(BaseTest):

    def setUp(self):
        super(ContractsTestsB, self).setUp()

    def tearDown(self):
        super(ContractsTestsB, self).tearDown()

    @unittest.skip('bug: #469')
    def test000_get_post_contract(self):
        """
            ITSYOU-059
            - Create a new contract, should succeed with 201
            - Get contract by contractid, should succeed with 200
            - Get invalid contract, should fail with 404
            - Add signature, should succeed with 201
            - Add signature with invalid inputs, should fail with 400
        """

        self.lg('Create a new contract, should succeed with 201')
        contractid = self.random_value()
        expire = '2030-10-02T22:00:00Z'
        data = {'content':'contract_1', 'contractId':contractid, 'contractType':'partnership','expires':expire}
        response = self.client_1.api.CreateUserContract(data, self.user_1)
        self.assertEqual(response.status_code, 201)

        self.lg('Get contract by contractid, should succeed with 200')
        response = self.client_1.api.GetContract(contractid)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(contractid, response.json()['contractId'])

        #bug #469
        self.lg('Get invalid contract, should fail with 404')
        response = self.client_1.api.GetContract('fake_contract_654564')
        self.assertEqual(response.status_code, 404)

        self.lg('Add signature, should succeed with 201')
        signature = self.random_value()
        signedby = self.random_value()
        date = self.time_rfc3339_format()
        invalid_date = self.random_value()
        publicKey = self.random_value()

        data = {"date":date,"publicKey":publicKey,"signature":signature,"signedBy":signedby}
        response = self.client_1.api.SignContract(data, contractid)
        self.assertEqual(response.status_code, 201)
        response = self.client_1.api.GetContract(contractid)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(data, response.json()['signatures'][-1])

        self.lg('Add signature with invalid inputs, should fail with 400')
        data = {"date":invalid_date,"publicKey":publicKey,"signature":signature,"signedBy":signedby}
        response = self.client_1.api.SignContract(data, contractid)
        self.assertEqual(response.status_code, 400)

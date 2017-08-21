from api_testing.utils import BaseTest
import types
import unittest


class ContractsTests(BaseTest):

    def setUp(self):
        super(ContractsTests, self).setUp()
        self.lg('Create a new contract')
        self.contractid_1 = self.random_value()
        expire = '2030-10-02T22:00:00Z'
        data = {'content':'contract_1', 'contractId':self.contractid_1, 'contractType':'partnership','expires':expire}
        response = self.client_1.api.CreateUserContract(data, self.user_1)
        self.assertEqual(response.status_code, 201)

    def test001_get_contracts(self):
        """ ITSYOU-023
        *Test case for check get a contract GET /contracts/{contractId}.*

        **Test Scenario:**

        #. check get a contract, should succeed
        #. validate all expected keys in the returned response
        """
        self.lg('%s STARTED' % self._testID)
        response = self.client_1.api.GetContract(self.contractid_1)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(type(response.json()), types.DictType)
        self.lg('%s ENDED' % self._testID)

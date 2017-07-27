from api_testing.utils import BaseTest
import types
import unittest
from random import randint
import json
import time

class CompaniesTestsB(BaseTest):

    def setUp(self):
        super(CompaniesTestsB, self).setUp()

    @unittest.skip('#462 #463 #464 #465')
    def test001_get_post_put_company(self):
        """
            ITSYOU-057
            - Create company (1), should succeed with 201
            - Create company with same globalid of company (1), should fail with 409
            - Get companies, should succeed with 200
            - Get company by globalid, should succeed with 200
            - Get company info, should succeed with 200
            - Update company globalid, should fail with 403
            - Update company info, should succeed with 201
            - Get validate, should succeed with 200
        """

        self.lg('%s STARTED' % self._testID)

        self.lg('Create company (1), should succeed with 201')
        globalid = self.random_value()
        data = {"globalid":globalid}
        response = self.client_1.api.CreateCompany(data)
        self.assertEqual(response.status_code, 201)

        self.lg('Create company with same globalid of company (1), should fail with 409')
        response = self.client_1.api.CreateCompany(data)
        self.assertEqual(response.status_code, 409)

        #bug #463
        self.lg('Get companies, should succeed with 200')
        response = self.client_1.api.GetCompanyList()
        self.assertEqual(response.status_code, 200)
        self.assertEqual(globalid, response.json()[-1]['globalid'])

        #bug #465
        self.lg('Get company by globalid, should succeed with 200')
        response = self.client_1.api.GetCompany(globalid)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(globalid, response.json()['globalid'])

        #bug #464
        self.lg('Get company info, should succeed with 200')
        response = self.client_1.api.GetCompanyInfo(globalid)
        self.assertEqual(response.status_code, 200)

        self.lg('Update company globalid, should fail with 403')
        data = {"globalid": self.random_value()}
        response = self.client_1.api.UpdateCompany(data, globalid)
        self.assertEqual(response.status_code, 403)

        #bug #466
        self.lg('Update company info, should succeed with 201')
        data = {"globalid": globalid, "taxnr":self.random_value()}
        response = self.client_1.api.UpdateCompany(data, globalid)
        self.assertEqual(response.status_code, 201)

        #bug #462
        self.lg('Get validate, should succeed with 200')
        response = self.client_1.api.companies_byGlobalId_validate_get(globalid)
        self.assertEqual(response.status_code, 200)

        self.lg('%s ENDED' % self._testID)

    @unittest.skip('bug: #468')
    def test002_get_post_contract(self):
        """
            ITSYOU-058
            - Create company, should succeed with 201
            - Create a new contract (1), should succeed with 201
            - Create a new contract (2), should succeed with 201
            - Create a new expired contract (3), should succeed with 201
            - Get company\'s contracts, should succeed with 200
            - Get company\'s contracts & include the expired contracts, should succeed with 200
            - Get company\'s contracts with page size 1, should succeed with 200
            - Get company\'s contracts with start page 2, should succeed with 200
        """
        self.lg('%s STARTED' % self._testID)

        self.lg('Create company, should succeed with 201')
        globalid = self.random_value()
        data = {"globalid":globalid}
        response = self.client_1.api.CreateCompany(data)
        self.assertEqual(response.status_code, 201)

        self.lg('Create a new contract (1), should succeed with 201')
        contractid_1 = self.random_value()
        expire = '2030-10-02T22:00:00Z'
        data = {'content':'contract_1', 'contractId':contractid_1, 'contractType':'partnership','expires':expire}
        response = self.client_1.api.CreateCompanyContract(data, globalid)
        self.assertEqual(response.status_code, 201)

        self.lg('Create a new contract (2), should succeed with 201')
        contractid_2 = self.random_value()
        expire = '2030-10-02T22:00:00Z'
        data = {'content':'contract_2', 'contractId':contractid_2, 'contractType':'partnership','expires':expire}
        response = self.client_1.api.CreateCompanyContract(data, globalid)
        self.assertEqual(response.status_code, 201)

        self.lg('Create a new expired contract (3), should succeed with 201')
        contractid_3 = self.random_value()
        expire = '2010-10-02T22:00:00Z'
        data = {'content':'contract_3', 'contractId':contractid_3, 'contractType':'partnership','expires':expire}
        response = self.client_1.api.CreateCompanyContract(data, globalid)
        self.assertEqual(response.status_code, 201)

        self.lg('[GET] Get company contracts, should succeed with 200')
        response = self.client_1.api.GetOrganizationContracts(globalid)
        self.assertEqual(response.status_code, 200)
        self.assertNotEqual(contractid_3, response.json()[-1]['contractId'])
        self.assertEqual(contractid_1, response.json()[-2]['contractId'])
        self.assertEqual(contractid_2, response.json()[-1]['contractId'])

        response = self.client_1.api.GetCompanyContracts(globalid, query_params={"max":1000,"includeExpired":True})
        self.assertEqual(response.status_code, 200)
        number_of_contracts = len(response.json())-1

        self.lg('[GET] Get company contracts & include the expired contracts, should succeed with 200')
        response = self.client_1.api.GetCompanyContracts(globalid, query_params={"max":1000,"includeExpired":True})
        self.assertEqual(response.status_code, 200)
        self.assertEqual(contractid_3, response.json()[-1]['contractId'])

        self.lg('[GET] Get company contracts with page size 1, should succeed with 200')
        response = self.client_1.api.GetCompanyContracts(globalid, query_params={"max":1, "start":number_of_contracts})
        self.assertEqual(contractid_2, response.json()[0]['contractId'])
        self.assertEqual(response.status_code, 200)

        self.lg('[GET] Get company contracts with page size 2, should succeed with 200')
        response = self.client_1.api.GetCompanyContracts(globalid, query_params={"max":1, "start":number_of_contracts-1})
        self.assertEqual(response.status_code, 200)
        self.assertEqual(contractid_1, response.json()[0]['contractId'])

        self.lg('%s ENDED' % self._testID)

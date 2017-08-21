from api_testing.utils import BaseTest
import types
import unittest


class CompaniesTests(BaseTest):

    def setUp(self):
        super(CompaniesTests, self).setUp()
        self.response = self.client_1.api.GetCompanyList()
        self.assertEqual(self.response.status_code, 200)
        self.company = self.response.json()[0]
        self.lg('GetCompanyList [%s] response [%s]' % (self.user_1, self.response.json()))

    @unittest.skip("bug: #463")
    def test001_get_companies(self):
        """ ITSYOU-018
        *Test case for check get companies GET /companies.*

        **Test Scenario:**

        #. check get companies, should succeed
        #. validate all expected keys in the returned response
        """
        self.lg('%s STARTED' % self._testID)
        response = self.response
        self.assertEqual(response.status_code, 200)
        self.assertEqual(type(response.json()), types.ListType)
        self.lg('%s ENDED' % self._testID)

    @unittest.skip("bug: #465")
    def test002_get_company(self):
        """ ITSYOU-019
        *Test case for check get company GET /companies/{globalId}.*

        **Test Scenario:**

        #. check get company, should succeed
        #. validate all expected keys in the returned response
        """
        self.lg('%s STARTED' % self._testID)
        response = self.client_1.api.GetCompany(self.company['globalId'])
        self.assertEqual(response.status_code, 200)
        self.assertEqual(type(response.json()), types.DictType)
        self.assertEqual(response.json()['globalId'], self.company['globalId'])
        self.lg('%s ENDED' % self._testID)

    @unittest.skip("bug: #468")
    def test003_get_company_contracts(self):
        """ ITSYOU-020
        *Test case for check get company contracts GET /companies/{globalId}/contracts.*

        **Test Scenario:**

        #. check get company contracts, should succeed
        #. validate all expected keys in the returned response
        """
        self.lg('%s STARTED' % self._testID)
        response = self.client_1.api.GetCompanyContracts(self.company['globalId'])
        self.assertEqual(response.status_code, 200)
        self.assertEqual(type(response.json()), types.ListType)
        self.lg('%s ENDED' % self._testID)

    @unittest.skip("bug: #464")
    def test004_get_company_info(self):
        """ ITSYOU-021
        *Test case for check get company info GET /companies/{globalId}/info.*

        **Test Scenario:**

        #. check get company info, should succeed
        #. validate all expected keys in the returned response
        """
        self.lg('%s STARTED' % self._testID)
        response = self.client_1.api.GetCompanyInfo(self.company['globalId'])
        self.assertEqual(response.status_code, 200)
        self.assertEqual(type(response.json()), types.DictType)
        self.assertEqual(response.json(), self.company['info'])
        self.lg('%s ENDED' % self._testID)

    @unittest.skip("bug: #462")
    def test005_get_company_validate(self):
        """ ITSYOU-022
        *Test case for check get company validate GET /companies/{globalId}/validate.*

        **Test Scenario:**

        #. check get company validate, should succeed
        #. validate all expected keys in the returned response
        """
        self.lg('%s STARTED' % self._testID)
        response = self.client_1.api.companies_byGlobalId_validate_get(self.company['globalId'])
        self.assertEqual(response.status_code, 200)
        self.assertEqual(type(response.json()), types.DictType)
        self.lg('%s ENDED' % self._testID)

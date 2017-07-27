from api_testing.utils import BaseTest
import types
import unittest


class OrganizationsTests(BaseTest):


    def setUp(self):
        super(OrganizationsTests, self).setUp()
        org_1_data = {"globalid":self.organization_1}
        response = self.client_1.api.CreateNewOrganization(org_1_data)
        self.assertEqual(response.status_code, 201)

    def tearDown(self):
        response = self.client_1.api.GetUserOrganizations(self.user_1)
        self.assertEqual(response.status_code, 200)
        for org in response.json()['owner']:
            response = self.client_1.api.DeleteOrganization(org)

    def test001_get_organization(self):
        """ ITSYOU-013
        *Test case for check get organization GET /organizations/{globalid}.*

        **Test Scenario:**

        #. check get organizations, should succeed
        #. validate all expected keys in the returned response
        """
        self.lg('%s STARTED' % self._testID)
        response = self.client_1.api.GetOrganization(self.organization_1)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(type(response.json()), types.DictType)
        self.lg('%s ENDED' % self._testID)

    def test002_get_organization_tree(self):
        """ ITSYOU-014
        *Test case for check get organization tree GET /organizations/{globalid}/tree.*

        **Test Scenario:**

        #. check get organizations tree, should succeed
        #. validate all expected keys in the returned response
        """
        self.lg('%s STARTED' % self._testID)
        response = self.client_1.api.GetOrganizationTree(self.organization_1)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(type(response.json()), types.DictType)
        self.lg('%s ENDED' % self._testID)


    def test003_get_organization_contracts(self):
        """ ITSYOU-015
        *Test case for check get organization contracts GET /organizations/{globalid}/contracts.*

        **Test Scenario:**

        #. check get organizations contracts, should succeed
        #. validate all expected keys in the returned response
        """
        self.lg('%s STARTED' % self._testID)
        response = self.client_1.api.GetOrganizationContracts(self.organization_1)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(type(response.json()), types.ListType)
        self.lg('%s ENDED' % self._testID)

    def test004_get_organization_invitations(self):
        """ ITSYOU-016
        *Test case for check get organization invitations GET /organizations/{globalid}/invitations.*

        **Test Scenario:**

        #. check get organizations invitations, should succeed
        #. validate all expected keys in the returned response
        """
        self.lg('%s STARTED' % self._testID)
        response = self.client_1.api.GetPendingOrganizationInvitations(self.organization_1)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(type(response.json()), types.ListType)
        self.lg('%s ENDED' % self._testID)

    def test005_get_organization_apikeys(self):
        """ ITSYOU-017
        *Test case for check get organization apikeys GET /organizations/{globalid}/apikeys.*

        **Test Scenario:**

        #. check get organizations apikeys, should succeed
        #. validate all expected keys in the returned response
        """
        self.lg('%s STARTED' % self._testID)
        response = self.client_1.api.GetOrganizationAPIKeyLabels(self.organization_1)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(type(response.json()), types.ListType)
        self.lg('%s ENDED' % self._testID)

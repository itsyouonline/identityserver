from api_testing.utils import BaseTest
import types
import unittest

class UsersTestsA(BaseTest):


    def setUp(self):
        super(UsersTestsA, self).setUp()
        self.response = self.client_1.api.GetUser(self.user_1)
        self.lg('GetUser [%s] response [%s]' % (self.user_1, self.response.json()))


    @unittest.skip("bug: #467")
    def test001_get_user(self):
        """ ITSYOU-001
        *Test case for check get user /users/{username}.*

        **Test Scenario:**

        #. check get user, should succeed
        #. validate all expected keys in the returned response
        """
        self.lg('%s STARTED' % self._testID)
        response = self.response
        self.assertEqual(response.status_code, 200)
        self.assertEqual(type(response.json()), types.DictType)
        self.assertIn('username', response.json().keys())
        self.assertEqual(type(response.json()['username']), type(u''))
        self.assertIn('firstname', response.json().keys())
        self.assertEqual(type(response.json()['firstname']), type(u''))
        self.assertIn('lastname', response.json().keys())
        self.assertEqual(type(response.json()['lastname']), type(u''))
        self.assertIn('expire', response.json().keys())
        self.assertEqual(type(response.json()['expire']), type(u''))
        self.assertIn('addresses', response.json().keys())
        self.assertEqual(type(response.json()['addresses']), types.ListType)
        self.assertIn('bankaccounts', response.json().keys())
        self.assertEqual(type(response.json()['bankaccounts']), types.ListType)
        self.assertIn('emailaddresses', response.json().keys())
        self.assertEqual(type(response.json()['emailaddresses']), types.ListType)
        self.assertIn('phonenumbers', response.json().keys())
        self.assertEqual(type(response.json()['phonenumbers']), types.ListType)
        self.assertIn('publicKeys', response.json().keys())
        self.assertEqual(type(response.json()['publicKeys']), types.ListType)
        self.assertIn('digitalwallet', response.json().keys())
        self.assertEqual(type(response.json()['digitalwallet']), types.ListType)
        self.assertIn('facebook', response.json().keys())
        self.assertEqual(type(response.json()['facebook']), types.DictType)
        self.assertIn('github', response.json().keys())
        self.assertEqual(type(response.json()['github']), types.DictType)
        self.lg('%s ENDED' % self._testID)

    def test002_get_user_addresses(self):
        """ ITSYOU-002
        *Test case for check get user addressess /users/{username}/addresses.*

        **Test Scenario:**

        #. check get user addressess, should succeed
        #. validate all expected user addressess in the returned response
        """
        self.lg('%s STARTED' % self._testID)
        response = self.client_1.api.GetUserAddresses(self.user_1)
        self.lg('GetUserAddresses [%s] response [%s]' % (self.user_1, response.json()))
        self.assertEqual(response.status_code, 200)
        self.assertEqual(type(response.json()), types.ListType)
        self.assertEqual(response.json(), self.response.json()['addresses'])
        self.lg('%s ENDED' % self._testID)

    def test003_get_user_emailaddresses(self):
        """ ITSYOU-003
        *Test case for check get user email addressess /users/{username}/emailaddresses.*

        **Test Scenario:**

        #. check get user email addressess, should succeed
        #. validate all expected email addressess in the returned response
        """
        self.lg('%s STARTED' % self._testID)
        response = self.client_1.api.GetEmailAddresses(self.user_1)
        self.lg('GetEmailAddresses [%s] response [%s]' % (self.user_1, response.json()))
        self.assertEqual(response.status_code, 200)
        self.assertEqual(type(response.json()), types.ListType)
        self.assertEqual(response.json(), self.response.json()['emailaddresses'])
        self.lg('%s ENDED' % self._testID)

    def test004_get_user_bankaccounts(self):
        """ ITSYOU-004
        *Test case for check get user bankaccounts /users/{username}/bankaccounts.*

        **Test Scenario:**

        #. check get user bankaccounts, should succeed
        #. validate all expected bankaccounts in the returned response
        """
        self.lg('%s STARTED' % self._testID)
        response = self.client_1.api.GetUserBankAccounts(self.user_1)
        self.lg('GetUserBankAccounts [%s] response [%s]' % (self.user_1, response.json()))
        self.assertEqual(response.status_code, 200)
        self.assertEqual(type(response.json()), types.ListType)
        self.assertEqual(response.json(), self.response.json()['bankaccounts'])
        self.lg('%s ENDED' % self._testID)

    def test005_get_user_bankaccounts(self):
        """ ITSYOU-005
        *Test case for check get user bankaccounts /users/{username}/bankaccounts.*

        **Test Scenario:**

        #. check get user bankaccounts, should succeed
        #. validate all expected bankaccounts in the returned response
        """
        self.lg('%s STARTED' % self._testID)
        response = self.client_1.api.GetUserBankAccounts(self.user_1)
        self.lg('GetUserBankAccounts [%s] response [%s]' % (self.user_1, response.json()))
        self.assertEqual(response.status_code, 200)
        self.assertEqual(type(response.json()), types.ListType)
        self.assertEqual(response.json(), self.response.json()['bankaccounts'])
        self.lg('%s ENDED' % self._testID)


class UsersTestsB(BaseTest):


    def test006_get_user_username(self):
        """ ITSYOU-006
        *Test case for check get user username /users/{username}/username.*

        **Test Scenario:**

        #. check get user username, should succeed
        #. validate all expected username in the returned response
        """
        self.lg('%s STARTED' % self._testID)
        response = self.client_1.api.GetUser(self.user_1)
        self.lg('GetUser [%s] response [%s]' % (self.user_1, response.json()))
        self.assertEqual(response.status_code, 200)
        self.assertEqual(type(response.json()['username']), type(u''))
        self.assertEqual(response.json()['username'], self.user_1)
        self.lg('%s ENDED' % self._testID)

    def test007_get_user_notifications(self):
        """ ITSYOU-007
        *Test case for check get user notifications /users/{username}/notifications.*

        **Test Scenario:**

        #. check get user notifications, should succeed
        #. validate all expected notifications in the returned response
        """
        self.lg('%s STARTED' % self._testID)
        response = self.client_1.api.GetNotifications(self.user_1)
        self.lg('GetNotifications [%s] response [%s]' % (self.user_1, response.json()))
        self.assertEqual(response.status_code, 200)
        self.assertEqual(type(response.json()), types.DictType)
        self.lg('%s ENDED' % self._testID)

    def test008_get_user_organizations(self):
        """ ITSYOU-008
        *Test case for check get user organizations /users/{username}/organizations.*

        **Test Scenario:**

        #. check get user organizations, should succeed
        #. validate all expected organizations in the returned response
        """
        self.lg('%s STARTED' % self._testID)
        response = self.client_1.api.GetUserOrganizations(self.user_1)
        self.lg('GetUserOrganizations [%s] response [%s]' % (self.user_1, response.json()))
        self.assertEqual(response.status_code, 200)
        self.assertEqual(type(response.json()), types.DictType)
        self.lg('%s ENDED' % self._testID)

    def test009_get_user_apikeys(self):
        """ ITSYOU-009
        *Test case for check get user apikeys /users/{username}/apikeys.*

        **Test Scenario:**

        #. check get user apikeys, should succeed
        #. validate all expected apikeys in the returned response
        """
        self.lg('%s STARTED' % self._testID)
        response = self.client_1.api.ListAPIKeys(self.user_1)
        self.lg('ListAPIKeys [%s] response [%s]' % (self.user_1, response.json()))
        self.assertEqual(response.status_code, 200)
        self.assertEqual(type(response.json()), types.ListType)
        self.lg('%s ENDED' % self._testID)

    def test010_get_user_info(self):
        """ ITSYOU-010
        *Test case for check get user info /users/{username}/info.*

        **Test Scenario:**

        #. check get user apikeys, should succeed
        #. validate all expected apikeys in the returned response
        """
        self.lg('%s STARTED' % self._testID)
        response = self.client_1.api.GetUserInformation(self.user_1)
        self.lg('GetUserInformation [%s] response [%s]' % (self.user_1, response.json()))
        self.assertEqual(response.status_code, 200)
        self.assertEqual(type(response.json()), types.DictType)
        self.assertIn('username', response.json().keys())
        self.assertEqual(type(response.json()['username']), type(u''))
        self.assertIn('firstname', response.json().keys())
        self.assertEqual(type(response.json()['firstname']), type(u''))
        self.assertIn('lastname', response.json().keys())
        self.assertEqual(type(response.json()['lastname']), type(u''))
        self.assertIn('addresses', response.json().keys())
        self.assertEqual(type(response.json()['addresses']), types.ListType)
        self.assertIn('emailaddresses', response.json().keys())
        self.assertEqual(type(response.json()['emailaddresses']), types.ListType)
        self.assertIn('phonenumbers', response.json().keys())
        self.assertEqual(type(response.json()['phonenumbers']), types.ListType)
        self.assertIn('publicKeys', response.json().keys())
        self.assertEqual(type(response.json()['publicKeys']), types.ListType)
        self.assertIn('facebook', response.json().keys())
        self.assertEqual(type(response.json()['facebook']), types.DictType)
        self.assertIn('github', response.json().keys())
        self.assertEqual(type(response.json()['github']), types.DictType)
        self.lg('%s ENDED' % self._testID)

    def test011_get_user_contracts(self):
        """ ITSYOU-011
        *Test case for check get user info /users/{username}/contracts.*

        **Test Scenario:**

        #. check get user contracts, should succeed
        #. validate all expected contracts in the returned response
        """
        self.lg('%s STARTED' % self._testID)
        response = self.client_1.api.GetUserContracts(self.user_1)
        self.lg('GetUserContracts [%s] response [%s]' % (self.user_1, response.json()))
        self.assertEqual(response.status_code, 200)
        self.assertEqual(type(response.json()), types.ListType)
        self.lg('%s ENDED' % self._testID)

    def test012_get_user_authorizations(self):
        """ ITSYOU-012
        *Test case for check get user authorizations /users/{username}/authorizations.*

        **Test Scenario:**

        #. check get user authorizations, should succeed
        #. validate all expected authorizations in the returned response
        """
        self.lg('%s STARTED' % self._testID)
        response = self.client_1.api.GetAllAuthorizations(self.user_1)
        self.lg('GetAllAuthorizations [%s] response [%s]' % (self.user_1, response.json()))
        self.assertEqual(response.status_code, 200)
        self.assertEqual(type(response.json()), types.ListType)
        self.lg('%s ENDED' % self._testID)

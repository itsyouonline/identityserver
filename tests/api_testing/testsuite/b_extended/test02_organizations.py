from api_testing.utils import BaseTest
import types
import unittest
from random import randint
import json
import time

class OrganizationsTestsB(BaseTest):

    def setUp(self):
        super(OrganizationsTestsB, self).setUp()

        self.org_11_globalid = self.random_value()
        self.org_12_globalid = self.random_value()
        self.org_13_globalid = self.random_value()
        self.org_21_globalid = self.random_value()
        self.org_22_globalid = self.random_value()

        self.org_13_suborg_globalid = self.org_13_globalid+'.'+self.random_value()

        #### org-11 user_1 owner
        self.lg('Create org %s for user %s' %(self.org_11_globalid, self.user_1))
        org_11_data = {"globalid":self.org_11_globalid}
        response = self.client_1.api.CreateNewOrganization(org_11_data)
        self.assertEqual(response.status_code, 201)

        #### org-12 user_1 owner & user_2 member
        self.lg('Create org %s for user %s' %(self.org_12_globalid, self.user_1))
        org_12_data = {"globalid":self.org_12_globalid}
        response = self.client_1.api.CreateNewOrganization(org_12_data)
        self.assertEqual(response.status_code, 201)

        #### org-13 user_1 owner
        self.lg('Create org %s for user %s' %(self.org_13_globalid, self.user_1))
        org_13_data = {"globalid":self.org_13_globalid}
        response = self.client_1.api.CreateNewOrganization(org_13_data)
        self.assertEqual(response.status_code, 201)
        ### create suborganization in org-13 and invite user_2 to be owner
        self.lg('Create suborganization %s in organization %s' %(self.org_13_suborg_globalid, self.org_13_globalid))
        sub_org_1_data = {"globalid":self.org_13_suborg_globalid}
        response = self.client_1.api.CreateNewSubOrganization(sub_org_1_data, self.org_13_globalid)
        self.assertEqual(response.status_code, 201)
        ## user_1 invite user_2 to join suborg org_13_suborg
        self.lg('user %s invite %s to join suborg %s' %(self.user_1, self.user_2, self.org_13_suborg_globalid))
        data = {'searchstring': self.user_2}
        response = self.client_1.api.AddOrganizationOwner(data, self.org_13_suborg_globalid)
        self.assertEqual(response.status_code, 201)
        ## user_2 accept invitation
        self.lg('user %s accept invitation to join %s' %(self.user_2, self.org_13_suborg_globalid))
        response = self.client_2.api.AcceptMembership(self.org_13_suborg_globalid , 'owner', self.user_2)
        self.assertEqual(response.status_code, 201)

        ## user_1 invite user_2 to join org-12
        self.lg('user %s invite %s to join %s' %(self.user_1, self.user_2, self.org_12_globalid))
        data = {'searchstring': self.user_2}
        response = self.client_1.api.AddOrganizationMember(data, self.org_12_globalid)
        self.assertEqual(response.status_code, 201)
        ## user_2 accept invitation
        self.lg('user %s accept invitation to join %s' %(self.user_2, self.org_12_globalid))
        response = self.client_2.api.AcceptMembership(self.org_12_globalid , 'member', self.user_2)
        self.assertEqual(response.status_code, 201)
        #### org-21 user_2 owner
        self.lg('Create org %s for user %s' %(self.org_21_globalid, self.user_2))
        org_21_data = {"globalid":self.org_21_globalid}
        response = self.client_2.api.CreateNewOrganization(org_21_data)
        self.assertEqual(response.status_code, 201)
        #### org-22 user_2 owner & user_1 member
        self.lg('Create org %s for user %s' %(self.org_22_globalid, self.user_2))
        org_22_data = {"globalid":self.org_22_globalid}
        response = self.client_2.api.CreateNewOrganization(org_22_data)
        self.assertEqual(response.status_code, 201)
        ## user_2 invite user_1
        self.lg('user %s invite %s to join %s' %(self.user_2, self.user_1, self.org_22_globalid))
        data = {'searchstring': self.user_1}
        response = self.client_2.api.AddOrganizationMember(data, self.org_22_globalid)
        self.assertEqual(response.status_code, 201)
        ## user_1 accept invitation
        self.lg('user %s accept invitation to join %s' %(self.user_1, self.org_22_globalid))
        response = self.client_1.api.AcceptMembership(self.org_22_globalid , 'member', self.user_1)
        self.assertEqual(response.status_code, 201)

    def tearDown(self):
        # user_1 orgs
        response = self.client_1.api.GetUserOrganizations(self.user_1)
        self.assertEqual(response.status_code, 200)
        for org in response.json()['owner']:
            response = self.client_1.api.DeleteOrganization(org)
        # user_2 orgs
        response = self.client_2.api.GetUserOrganizations(self.user_2)
        self.assertEqual(response.status_code, 200)
        for org in response.json()['owner']:
            response = self.client_2.api.DeleteOrganization(org)
        super(OrganizationsTestsB, self).tearDown()

    @unittest.skip('bug: #442')
    def test001_post_get_organization(self):
        """
            ITSYOU-040
            - Create new organization, should succeed with 201
            - Create new organization with globalid already exists, should fail with 409
            - Get organization by globalid, should succeed with 200
            - Get invalid organization, should fail with 404
            - Delete organization, should succeed with 204
        """
        self.lg('%s STARTED' % self._testID)

        self.lg('[POST] Create new organization, succeed with 201')

        globalid = self.random_value()
        data = {"globalid":globalid}
        response = self.client_1.api.CreateNewOrganization(data)
        self.assertEqual(response.status_code, 201)

        self.lg('[POST] Create new organization with globalid already exists,  succeed with 409')
        response = self.client_1.api.CreateNewOrganization(data)
        self.assertEqual(response.status_code, 409)

        self.lg('[GET] Get organization by globalid, should succeed with 200')
        response = self.client_1.api.GetOrganization(globalid)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(globalid, response.json()['globalid'])
        self.assertEqual([self.user_1], response.json()['owners'])

        #bug #442
        self.lg('[GET] Get invalid organization - should fail with 404')
        response = self.client_1.api.GetOrganization('fake_organization')
        self.assertEqual(response.status_code, 404)

        self.lg('[DEL] Delete organization by globalid, should succeed with 204')
        response = self.client_1.api.DeleteOrganization(globalid)
        self.assertEqual(response.status_code, 204)

        self.lg('%s ENDED' % self._testID)

    @unittest.skip('bug: #448 #450')
    def test002_get_post_put_delete_organization(self):
        """
            ITSYOU-041
            - Create suborganization (1) inside org-11, should succeed with 201
            - Get organization org-11 tree, should succeed with 200
            - Create suborganization with globalid already exists (globalid of suborganization(1)), should succeed with 409
            - Update organization org-11 globalid, should fail with 403
            - Update organization org-11 info, should succeed with 201
            - Delete suborganization (1), should succeed with 204
            - Delete nonexisting organization, should fail with 404
            - Delete org-11, should succeed with 204
        """

        self.lg('%s STARTED' % self._testID)

        globalid = self.org_11_globalid

        self.lg('[POST] Create suborganization, should succeed with 201')
        sub_org_1_globalid = globalid+'.'+self.random_value()
        sub_org_1_data = {"globalid":sub_org_1_globalid}
        response = self.client_1.api.CreateNewSubOrganization(sub_org_1_data, globalid)
        self.assertEqual(response.status_code, 201)

        self.lg('[GET] Get organization org-11 tree, should succeed with 200')
        response = self.client_1.api.GetOrganizationTree(globalid)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(globalid, response.json()['globalid'])
        self.assertIn(sub_org_1_globalid, [x['globalid'] for x in response.json()['children']])

        self.lg('[POST] Create suborganization (1) with globalid already exists (globalid of suborganization(1)), should succeed with 409')
        response = self.client_1.api.CreateNewSubOrganization(sub_org_1_data, globalid)
        self.assertEqual(response.status_code, 409)

        #bug #450
        self.lg('[PUT] Update suborganization (1) globalid, should fail with 403')
        sub_org_1_data = {"globalid":self.random_value()}
        response = self.client_1.api.UpdateOrganization(sub_org_1_data, sub_org_1_globalid)
        self.assertEqual(response.status_code, 403)

        self.lg('[PUT] Update organization info, should succeed with 201')
        data_new = {"globalid":globalid, "dns":["www.a.com"]}
        response = self.client_1.api.UpdateOrganization(data_new, globalid)
        self.assertEqual(response.status_code, 201)

        self.lg('[DEL] Delete suborganization (1), should succeed with 204')
        response = self.client_1.api.DeleteOrganization(sub_org_1_globalid)
        self.assertEqual(response.status_code, 204)

        #bug #448
        self.lg('[DEL] Delete nonexisting organization, should fail with 404')
        response = self.client_1.api.DeleteOrganization('fake_organization')
        self.assertEqual(response.status_code, 404)

        self.lg('[DEL] Delete organization org-11, should succeed with 204')
        response = self.client_1.api.DeleteOrganization(self.org_11_globalid)
        self.assertEqual(response.status_code, 204)

        self.lg('%s ENDED' % self._testID)

    @unittest.skip('#453')
    def test003_get_post_put_delete_apikey(self):
        """
            ITSYOU-042
            - Register a new apikey (1), should succeed with 201
            - Register a new apikey (2), should succeed with 201
            - Get organization\'s apikeys, should succeed with 200
            - Get apikey (2) by label, should succeed with 200
            - Get nonexisting apikey, should fail with 404
            - Register new apikey with label already exists (label of apikey (2)), should fail with 409
            - Update the apikey (2), should succeed with 201
            - Update the apikey (2) with label already exists (label of apikey (1)), should fail with 409
            - Delete apikey (1), should succeed with 204
            - Delete apikey (2), should succeed with 204
        """

        self.lg('%s STARTED' % self._testID)

        globalid = self.org_11_globalid

        self.lg('[POST] Register a new apikey (1), should succeed with 201')
        label = self.random_value()
        callbackURL = self.random_value()
        data = {'label' : label, 'callbackURL':callbackURL}
        response = self.client_1.api.CreateNewOrganizationAPIKey(data, globalid)
        self.assertEqual(response.status_code, 201)

        self.lg('[POST] Register a new apikey (2), should succeed with 201')
        new_label = self.random_value()
        new_callbackURL = self.random_value()
        new_data = {'label' : new_label, 'callbackURL':new_callbackURL}
        response = self.client_1.api.CreateNewOrganizationAPIKey(new_data, globalid)
        self.assertEqual(response.status_code, 201)

        self.lg('[GET] Get organization\'s apikeys, should succeed with 200')
        response = self.client_1.api.GetOrganizationAPIKeyLabels(globalid)
        self.assertEqual(response.status_code, 200)
        self.assertIn(label, response.json())
        self.assertIn(new_label, response.json())

        self.lg('[GET] Get apikey (2) by label, should succeed with 200')
        response = self.client_1.api.GetOrganizationAPIKey(new_label, globalid)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(new_label, response.json()['label'])
        self.assertEqual(new_callbackURL, response.json()['callbackURL'])

        self.lg('[GET] - Get nonexisting apikey, should fail with 404')
        response = self.client_1.api.GetOrganizationAPIKey('fake_label', globalid)
        self.assertEqual(response.status_code, 404)

        self.lg('[POST] Register new apikey with label already exists (label of apikey (2)), should fail with 409')
        response = self.client_1.api.GetOrganizationAPIKeyLabels(globalid)
        self.assertEqual(response.status_code, 200)
        new_label = response.json()[-1]
        new_callbackURL = self.random_value()
        new_data = {'label' : new_label, 'callbackURL':new_callbackURL}
        response = self.client_1.api.CreateNewOrganizationAPIKey(new_data, globalid)
        self.assertEqual(response.status_code, 409)

        #bug 453 bad request require secret
        self.lg('[PUT] Update the apikey (2), should succeed with 201')
        response = self.client_1.api.GetOrganizationAPIKeyLabels(globalid)
        self.assertEqual(response.status_code, 200)
        label = response.json()[-1]
        print label
        new_label = self.random_value()
        new_callbackURL = self.random_value()
        new_data = {'label' : new_label, 'callbackURL':new_callbackURL}
        response = self.client_1.api.UpdateOrganizationAPIKey(new_data, label, globalid)
        self.assertEqual(response.status_code, 201)

        self.lg('[PUT] Update the apikey (2) with label already exists (label of apikey (1)), should fail with 409')
        response = self.client_1.api.GetOrganizationAPIKeyLabels(globalid)
        self.assertEqual(response.status_code, 200)
        label = response.json()[-1]
        new_label = response.json()[-2]
        new_callbackURL = self.random_value()
        new_data = {'label' : new_label, 'callbackURL':new_callbackURL}
        response = self.client_1.api.UpdateOrganizationAPIKey(new_data, label, globalid)
        self.assertEqual(response.status_code, 409)

        self.lg('[DELETE] Delete apikey (1), should succeed with 204')
        response = self.client_1.api.GetOrganizationAPIKeyLabels(globalid)
        self.assertEqual(response.status_code, 200)
        label = response.json()[-1]
        response = self.client_1.api.DeleteOrganizationAPIKey(label, globalid)
        self.assertEqual(response.status_code, 204)

        self.lg('[DELETE] Delete apikey (2), should succeed with 204')
        response = self.client_1.api.GetOrganizationAPIKeyLabels(globalid)
        self.assertEqual(response.status_code, 200)
        label = response.json()[-1]
        response = self.client_1.api.DeleteOrganizationAPIKey(label, globalid)
        self.assertEqual(response.status_code, 204)

        response = self.client_1.api.GetOrganizationAPIKeyLabels(globalid)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(response.json(), [])

        self.lg('%s ENDED' % self._testID)

    @unittest.skip('bug: #443')
    def test004_get_post_put_delete_registry(self):

        """
            ITSYOU-043
            - Register a new registry (1), should succeed with 201
            - Register a new registry (2), should succeed with 201.
            - Get user registries, should succeed with 200
            - Get registry (2) by key, should succeed with 200
            - Get nonexisting registry, should fail with 404
            - Register a new registry with key already exists (key of registry (2)), should succeed with 201
            - Register a new registry with invalid inputs, should fail with 400
            - Delete registry, should succeed with 204
            - Delete registry again, should succeed with 204
            - Delete nonexisting registry, should fail with 404
        """

        self.lg('%s STARTED' % self._testID)

        globalid = self.org_11_globalid

        self.lg('[POST] Register a new registry (1) - should succeed with 201')
        key = self.random_value()
        value = self.random_value()
        data = {"Key":key,"Value":value}
        response = self.client_1.api.CreateOrganizationRegistry(data, globalid)
        self.assertEqual(response.status_code, 201)

        self.lg('[POST] Register a new registry (2) - should succeed with 201')
        key = self.random_value()
        value = self.random_value()
        new_data = {"Key": key,"Value": value}
        response = self.client_1.api.CreateOrganizationRegistry(new_data, globalid)
        self.assertEqual(response.status_code, 201)

        self.lg('[GET] Get user registries - should succeed with 200')
        response = self.client_1.api.GetOrganizationRegistries(globalid)
        self.assertEqual(response.status_code, 200)
        self.assertIn(data, response.json())
        self.assertIn(new_data, response.json())

        self.lg('[GET] Get registry (2) by key - should succeed with 200')
        key = response.json()[-1]['Key']
        response = self.client_1.api.GetOrganizationRegistry(key, globalid)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(new_data, response.json())

        self.lg('[GET] Get nonexisting registry - should fail with 404')
        response = self.client_1.api.GetOrganizationRegistry('fake_key', globalid)
        self.assertEqual(response.status_code, 404)

        self.lg('[POST] Register a new registry with key already exists (key of registry (2)), should succeed with 201')
        response = self.client_1.api.GetOrganizationRegistries(globalid)
        self.assertEqual(response.status_code, 200)
        key = response.json()[-1]['Key']
        value = self.random_value()
        new_data = {"Key": key,"Value": value}
        response = self.client_1.api.CreateOrganizationRegistry(new_data, globalid)
        self.assertEqual(response.status_code, 201)


        self.lg('[POST] Register a new registry with invalid inputs - should fail with 400')
        key = ''
        value = ''
        new_data = {"Key": key,"Value": value}
        response = self.client_1.api.CreateOrganizationRegistry(new_data, globalid)
        self.assertEqual(response.status_code, 400)

        self.lg('[DELETE] Delete registry (2) - should succeed with 204')
        response = self.client_1.api.GetOrganizationRegistries(globalid)
        self.assertEqual(response.status_code, 200)
        key = response.json()[-1]['Key']
        response = self.client_1.api.DeleteOrganizaitonRegistry(key, globalid)
        self.assertEqual(response.status_code, 204)

        self.lg('[DELETE] Delete registry (1) - should succeed with 204')
        response = self.client_1.api.GetOrganizationRegistries(globalid)
        self.assertEqual(response.status_code, 200)
        key = response.json()[-1]['Key']
        response = self.client_1.api.DeleteOrganizaitonRegistry(key, globalid)
        self.assertEqual(response.status_code, 204)

        self.lg('[DELETE] Delete nonexisting registry - should fail with 404')
        response = self.client_1.api.DeleteOrganizaitonRegistry('fake_key', globalid)
        self.assertEqual(response.status_code, 404)

        self.lg('%s ENDED' % self._testID)

    @unittest.skip('bug: #441 #444')
    def test005_get_post_put_delete_dns(self):

        """
            ITSYOU-044
            - Register a new dns name (1), should succeed with 201
            - Register a new dns name (2), should succeed with 201
            - Register a new dns name with name already exists (name of dns (1)), should fail with 409
            - Get organization dns names, should succeed with 200
            - Register a new dns name with invalid name, should fail with 400
            - Update organization dns name (2), should succeed with 201
            - Update organization dns name (2) with name already exists (name of dns (1)), should fail with 409
            - Update organization dns name with invalid name, should fail with 400
            - Update nonexisting organization dns name, should fail with 404
            - Delete dns name (2), should succeed with 204
            - Delete dns name (1), should succeed with 204
            - Delete invalid dns name, should fail with 404
        """

        self.lg('%s STARTED' % self._testID)

        globalid = self.org_11_globalid

        self.lg('[POST] Register a new dns name (1), should succeed with 201')
        dnsname_1 = "www.abc.com"
        data_1 = {"name":dnsname_1}
        response = self.client_1.api.CreateOrganizationDNS(data_1, globalid)
        self.assertEqual(response.status_code, 201)

        self.lg('[POST] Register a new dns name (2), should succeed with 201')
        dnsname_2 = "www.xyz.com"
        data_2 = {"name": dnsname_2}
        response = self.client_1.api.CreateOrganizationDNS(data_2, globalid)
        self.assertEqual(response.status_code, 201)

        #bug #441
        self.lg('[POST] Register a new dns name with name already exists (name of dns (1)) , should fail with 409')
        response = self.client_1.api.CreateOrganizationDNS(data_1, globalid)
        self.assertEqual(response.status_code, 409)

        self.lg('[GET] Get organization dns names, should succeed with 200')
        response = self.client_1.api.GetOrganization(globalid)
        self.assertEqual(response.status_code, 200)
        self.assertIn(dnsname_1, response.json()['dns'])
        self.assertIn(dnsname_2, response.json()['dns'])

        self.lg('[POST] Register a new dns name with invalid name , should fail with 400')
        data_invalid = {"name": self.random_value()}
        response = self.client_1.api.CreateOrganizationDNS(data_invalid, globalid)
        self.assertEqual(response.status_code, 400)

        self.lg('[PUT] Update organization dns name (2), should succeed with 201')
        new_dnsname = "www.mno.com"
        new_data = {"name":new_dnsname}
        response = self.client_1.api.UpdateOrganizationDNS(new_data, dnsname_2, globalid)
        self.assertEqual(response.status_code, 201)

        response = self.client_1.api.GetOrganization(globalid)
        self.assertEqual(response.status_code, 200)
        self.assertIn(new_dnsname, response.json()['dns'])
        self.assertNotIn(dnsname_2, response.json()['dns'])
        dnsname_2 = new_dnsname

        #bug #441
        self.lg('[PUT] Update organization dns name (2) with name already exists (name of dns (1)), should fail with 409')
        new_data = {"name":dnsname_2}
        response = self.client_1.api.UpdateOrganizationDNS(new_data, dnsname_1, globalid)
        self.assertEqual(response.status_code, 409)

        self.lg('[PUT] Update organization dns name with invalid name, should fail with 400')
        new_dnsname = self.random_value()
        new_data = {"name":new_dnsname}
        response = self.client_1.api.UpdateOrganizationDNS(new_data, dnsname_2, globalid)
        self.assertEqual(response.status_code, 400)

        # bug #444
        self.lg('[PUT] Update nonexisting organization dns name, should fail with 404')
        new_data = {"name":"www.qwe.com"}
        response = self.client_1.api.UpdateOrganizationDNS(new_data, 'fakse_dnsname', globalid)
        self.assertEqual(response.status_code, 404)

        self.lg('[DELETE] Delete dns name(1), should succeed with 204')
        response = self.client_1.api.DeleteOrganizaitonDNS(dnsname_1, globalid)
        self.assertEqual(response.status_code, 204)

        self.lg('[DELETE] Delete dns name again, should succeed with 204')
        response = self.client_1.api.DeleteOrganizaitonDNS(dnsname_2, globalid)
        self.assertEqual(response.status_code, 204)

        response = self.client_1.api.GetOrganization(globalid)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(response.json()['dns'], [])

        self.lg('[DELETE] Delete invalid dns name, should fail with 404')
        response = self.client_1.api.DeleteOrganizaitonDNS('fake_dns', globalid)
        self.assertEqual(response.status_code, 404)

        self.lg('%s ENDED' % self._testID)

    @unittest.skip('bug: #440 #452')
    def test006_get_post_twofa(self):
        """
            ITSYOU-045
            - Update 2fa with owner user, should succeed with 200
            - Update 2fa with invalid value, should fail with 400
            - Update 2fa with member user, should fail with 403
            - Get 2fa, should succeed with 200
            - Get 2fa of nonexisting organization , should fail with 404
            - Update 2fa of nonexisting organization, should fail with 404
        """
        self.lg('%s STARTED' % self._testID)

        globalid = self.org_12_globalid

        self.lg('[POST] Update 2fa with owner user, should succeed with 200')
        secondsvalidity = 2000
        data = {"secondsvalidity":secondsvalidity}
        response = self.client_1.api.UpdateTwoFA(data, globalid)
        self.assertEqual(response.status_code, 200)

        self.lg('[POST] Update 2fa with invalid value, should fail with 400')
        secondsvalidity = self.random_value()
        data = {"secondsvalidity":secondsvalidity}
        response = self.client_1.api.UpdateTwoFA(data, globalid)
        self.assertEqual(response.status_code, 400)

        self.lg('[POST] Update 2fa with member user, should fail with 403')
        secondsvalidity = 2000
        data = {"secondsvalidity":secondsvalidity}
        response = self.client_2.api.UpdateTwoFA(data, globalid)
        self.assertEqual(response.status_code, 403)

        self.lg('[GET] Get 2fa, should succeed with 200')
        response = response = self.client_1.api.GetTwoFA(globalid)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(response.json()['secondsvalidity'], secondsvalidity)

        #bug #440 (internal server error)
        self.lg('[GET] Get 2fa of nonexisting organization , should fail with 404')
        response = self.client_1.api.GetTwoFA('fake_globalid')
        self.assertEqual(response.status_code, 404)

        ##bug #452 (403 instead of 404 put not post)
        self.lg('[POST] Update 2fa of nonexisting organization, should fail with 404')
        data = {"secondsvalidity":secondsvalidity}
        response = self.client_2.api.UpdateTwoFA(data, 'fake_globalid')
        self.assertEqual(response.status_code, 404)

        self.lg('%s ENDED' % self._testID)

    @unittest.skip('bug: #439')
    def test007_get_put_delete_logo(self):
            """
                ITSYOU-046
                - Get organization logo, should succeed with 200
                - Update organization logo, should succeed with 200
                - Update organization logo with large file, should fail with 413
                - Remove organization logo, should succeed with 204
                - Remove nonexisting organization logo, should fail with 404
            """
            self.lg('%s STARTED' % self._testID)

            globalid = self.org_11_globalid

            self.lg('[GET] Get organization logo, should succeed with 200')
            response = self.client_1.api.GetOrganizationLogo(globalid)
            self.assertEqual(response.status_code, 200)
            self.assertEqual(response.json()['logo'], "")

            self.lg('[PUT] Update organization logo, should succeed with 200')
            logo = 'data:image/png;base64,' + self.random_value(1024)
            data = {"logo":logo}
            response = self.client_1.api.ganizationLogo(data, globalid)
            self.assertEqual(response.status_code, 200)

            response = self.client_1.api.GetOrganizationLogo(globalid)
            self.assertEqual(response.status_code, 200)
            self.assertEqual(response.json()['logo'], logo)

            self.lg('[PUT] Update organization logo with large file, should fail with 413')
            large_logo = 'data:image/png;base64,' + self.random_value(1024*1024*6)
            data = {"logo":large_logo}
            response = self.client_1.api.SetOrganizationLogo(data, globalid)
            self.assertEqual(response.status_code, 413)

            self.lg('[DEL] Remove organization logo, should succeed with 204')
            response = self.client_1.api.DeleteOrganizationLogo(globalid)
            self.assertEqual(response.status_code, 204)

            #bug #439 (403 instead of 404)
            self.lg('[DEL] Remove nonexisting organization logo, should fail with 404')
            response = self.client_1.api.DeleteOrganizationLogo('fake_globalid')
            self.assertEqual(response.status_code, 404)

            self.lg('%s ENDED' % self._testID)

    @unittest.skip('bug: #445')
    def test008_get_post_contracts(self):
            """
                ITSYOU-047
                - Create a new contract (1), should succeed with 201
                - Create a new contract (2), should succeed with 201
                - Create a new expired contract (3), should succeed with 201
                - Get organization\'s contracts, should succeed with 200
                - Get organization\'s contracts & include the expired contracts, should succeed with 200
                - Get organization\'s contracts with page size 1, should succeed with 200
                - Get organization\'s contracts with start page 2, should succeed with 200
            """
            self.lg('%s STARTED' % self._testID)

            globalid = self.org_11_globalid

            self.lg('Create a new contract (1), should succeed with 201')
            contractid_1 = self.random_value()
            expire = '2030-10-02T22:00:00Z'
            data = {'content':'contract_1', 'contractId':contractid_1, 'contractType':'partnership','expires':expire}
            response = self.client_1.api.CreateOrganizationContracts(data, globalid)
            self.assertEqual(response.status_code, 201)

            self.lg('Create a new contract (2), should succeed with 201')
            contractid_2 = self.random_value()
            expire = '2030-10-02T22:00:00Z'
            data = {'content':'contract_2', 'contractId':contractid_2, 'contractType':'partnership','expires':expire}
            response = self.client_1.api.CreateOrganizationContracts(data, globalid)
            self.assertEqual(response.status_code, 201)

            self.lg('Create a new expired contract (3), should succeed with 201')
            contractid_3 = self.random_value()
            expire = '2010-10-02T22:00:00Z'
            data = {'content':'contract_3', 'contractId':contractid_3, 'contractType':'partnership','expires':expire}
            response = self.client_1.api.CreateOrganizationContracts(data, globalid)
            self.assertEqual(response.status_code, 201)

            self.lg('[GET] Get organization\'s contracts, should succeed with 200')
            response = self.client_1.api.GetOrganizationContracts(globalid)
            self.assertEqual(response.status_code, 200)
            self.assertNotEqual(contractid_3, response.json()[-1]['contractId'])
            self.assertEqual(contractid_1, response.json()[-2]['contractId'])
            self.assertEqual(contractid_2, response.json()[-1]['contractId'])

            response = self.client_1.api.GetOrganizationContracts(globalid, query_params={"max":1000,"includeExpired":True})
            self.assertEqual(response.status_code, 200)
            number_of_contracts = len(response.json())-1

            self.lg('[GET] Get organization\'s contracts & include the expired contracts, should succeed with 200')
            response = self.client_1.api.GetOrganizationContracts(globalid, query_params={"max":1000,"includeExpired":True})
            self.assertEqual(response.status_code, 200)
            self.assertEqual(contractid_3, response.json()[-1]['contractId'])

            self.lg('[GET] Get organization\'s contracts with page size 1, should succeed with 200')
            response = self.client_1.api.GetOrganizationContracts(globalid, query_params={"max":1, "start":number_of_contracts})
            self.assertEqual(contractid_2, response.json()[0]['contractId'])
            self.assertEqual(response.status_code, 200)

            self.lg('[GET] Get organization\'s contracts with page size 2, should succeed with 200')
            response = self.client_1.api.GetOrganizationContracts(globalid, query_params={"max":1, "start":number_of_contracts-1})
            self.assertEqual(response.status_code, 200)
            self.assertEqual(contractid_1, response.json()[0]['contractId'])

            self.lg('%s ENDED' % self._testID)

    def test009_get_post_put_delete_description(self):
        """
            ITSYOU-048
            - Add new description with langkey "en", should succeed with 201
            - Get description by langkey, should succeed with 200
            - Add new description with random langkey, should succeed with 201
            - Update description, should succeed with 200
            - Delete description, should succeed with 204
            - Get description withfallback, should succeed with 200
        """
        self.lg('%s STARTED' % self._testID)

        globalid = self.org_11_globalid

        self.lg('Add new description with langkey "en", should succeed with 201')
        text = self.random_value()
        langkey = 'en'
        data = {"langkey": langkey, "text": text}
        response = self.client_1.api.AddOrganizationDescription(data, globalid)
        self.assertEqual(response.status_code, 201)
        self.assertEqual(response.json(), data)

        self.lg('Get description by langkey, should succeed with 200')
        response = self.client_1.api.GetOrganizationDescription(langkey, globalid)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(response.json()['langkey'], langkey)

        self.lg('Add new description with random langkey, should succeed with 201')
        text = self.random_value()
        new_langkey = self.random_value(2)
        new_data = {"langkey": new_langkey, "text": text}
        response = self.client_1.api.AddOrganizationDescription(new_data, globalid)
        self.assertEqual(response.status_code, 201)
        self.assertEqual(response.json(), new_data)

        self.lg('Update description, should succeed with 200')
        text = self.random_value()
        new_data = {"langkey": new_langkey, "text": text}
        response = self.client_1.api.UpdateOrganizationDescription(new_data, globalid)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(response.json(), new_data)

        self.lg('Delete description, should succeed with 204')
        response = self.client_1.api.DeleteOrganizationDescription(new_langkey, globalid)
        self.assertEqual(response.status_code, 204)

        self.lg('Get description withfallback, should succeed with 200')
        response = self.client_1.api.GetOrganizationDescriptionWithfallback(new_langkey, globalid)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(response.json()['langkey'], 'en')

        self.lg('%s ENDED' % self._testID)

    def test010_get_delete_users_invitation(self):
        """
            ITSYOU-049
            * Same steps for member and owner roles
            - User_1 invite nonexisting user to join org_11, should fail with 404
            - User_1 invite user_2 to join org_11, should succeed with 201
            - Get org-11 pending invitations, should succeed with 200
            - Check user_2 recevied the invitation, should succeed with 200
            - Cancel org-11 pending invitation for user_2 to join org_11, should succeed with 204
            - Check org-11 pending invitations are empty, should succeed with 200
            - Check user_2 invitations are empty, should succeed with 200
            - Cancel org-11 pending invitation for nonexisting user to join org_11, should fail with 404
        """

        self.lg('%s STARTED' % self._testID)

        globalid = self.org_11_globalid

        for role in ['member', 'owner']:

            self.lg('User_1 invite nonexisting user to join org_11 role %s , should succeed with 201' % role)
            data = {'searchstring': 'fake_user'}
            if role =="member":
                response = self.client_1.api.AddOrganizationMember(data, globalid)
            elif role == "owner":
                response = self.client_1.api.AddOrganizationOwner(data, globalid)
            self.assertEqual(response.status_code, 404)

            self.lg('User_1 invite user_2 to join org_11 role %s , should succeed with 201' % role)
            data = {'searchstring': self.user_2}
            if role =="member":
                response = self.client_1.api.AddOrganizationMember(data, globalid)
            elif role == "owner":
                response = self.client_1.api.AddOrganizationOwner(data, globalid)
            self.assertEqual(response.status_code, 201)

            self.lg('Get org-11 pending invitations, should succeed with 200')
            response = self.client_1.api.GetPendingOrganizationInvitations(globalid)
            self.assertEqual(response.status_code, 200)
            self.assertEqual(response.json()[-1]['status'], 'pending')
            self.assertEqual(response.json()[-1]['organization'], globalid)
            self.assertEqual(response.json()[-1]['user'], self.user_2)
            self.assertEqual(response.json()[-1]['role'], role)

            self.lg('Check user_2 invitations, should succeed with 200')
            response = self.client_2.api.GetNotifications(self.user_2)
            self.assertEqual(response.status_code, 200)
            self.assertEqual(self.org_11_globalid, response.json()['invitations'][-1]['organization'])
            self.assertEqual(role, response.json()['invitations'][-1]['role'])
            self.assertEqual('pending', response.json()['invitations'][-1]['status'])

            self.lg('Cancel org-11 pending invitation for user_2 to join org_11 role %s , should succeed with 204' % role)
            response = self.client_1.api.RemovePendingOrganizationInvitation(self.user_2, globalid)
            self.assertEqual(response.status_code, 204)

            self.lg('Check org-11 pending invitations are empty, should succeed with 200')
            response = self.client_1.api.GetPendingOrganizationInvitations(globalid)
            self.assertEqual(response.status_code, 200)
            self.assertEqual(response.json(), [])

            self.lg('Check user_2 invitations are empty, should succeed with 200')
            response = self.client_2.api.GetNotifications(self.user_2)
            self.assertEqual(response.status_code, 200)
            self.assertEqual(response.json()['invitations'], [])

            self.lg('Cancel org-11 pending invitation for nonexisting user to join org_11 role %s , should fail with 404' % role)
            response = self.client_1.api.RemovePendingOrganizationInvitation('fake_user', globalid)
            self.assertEqual(response.status_code, 404)

        self.lg('%s ENDED' % self._testID)

    def test011_get_delete_organizations_invitation(self):
        """
            ITSYOU-050
            * Same steps for orgmember and orgowner roles
            - org-11 invite nonexisting organization to join org_11 as orgmember, should fail with 404
            - org-11 invite org-21 to join org_11 as orgmember, should succeed with 201
            - Get org-11 pending invitations, should succeed with 200
            - Check org-21 recevied the invitation, should succeed with 200
            - Cancel org-11 pending invitation for org-21 to join org_11, should succeed with 201
            - Check org-11 pending invitations are empty, should succeed with 200
            - Check org-21 invitations are empty, should succeed with 200
            - Cancel org-11 pending invitation for nonexisting organization to join org_11, should fail with 404
        """

        self.lg('%s STARTED' % self._testID)

        globalid = self.org_11_globalid

        for role in ['orgmember', 'orgowner']:

            self.lg('org-11 invite nonexisting organization to join org_11 role %s , should fail with 404' % role)
            data = {role: 'fake_organization'}
            if role =="orgmember":
                response = self.client_1.api.SetOrgMember(data, globalid)
            elif role == "orgowner":
                response = self.client_1.api.SetOrgOwner(data, globalid)
            self.assertEqual(response.status_code, 404)

            self.lg('org-11 invite org-21 to join org_11 role %s , should succeed with 201' % role)
            data = {role: self.org_21_globalid}
            if role =="orgmember":
                response = self.client_1.api.SetOrgMember(data, globalid)
            elif role == "orgowner":
                response = self.client_1.api.SetOrgOwner(data, globalid)
            self.assertEqual(response.status_code, 201)

            self.lg('Get org-11 pending invitations, should succeed with 200')
            response = self.client_1.api.GetPendingOrganizationInvitations(globalid)
            self.assertEqual(response.status_code, 200)
            self.assertEqual(response.json()[-1]['status'], 'pending')
            self.assertEqual(response.json()[-1]['organization'], self.org_11_globalid)
            self.assertEqual(response.json()[-1]['user'], self.org_21_globalid)
            self.assertEqual(response.json()[-1]['role'], role)
            self.assertTrue(response.json()[-1]['isorganization'])

            self.lg('Check user_2 invitations, should succeed with 200')
            response = self.client_2.api.GetNotifications(self.user_2)
            self.assertEqual(response.status_code, 200)
            self.assertEqual(self.org_11_globalid, response.json()['organizationinvitations'][-1]['organization'])
            self.assertEqual(role, response.json()['organizationinvitations'][-1]['role'])
            self.assertEqual('pending', response.json()['organizationinvitations'][-1]['status'])

            self.lg('Cancel org-11 pending invitation for user_2 to join org_11 role %s , should succeed with 204' % role)
            response = self.client_1.api.RemovePendingOrganizationInvitation(self.org_21_globalid, globalid)
            self.assertEqual(response.status_code, 204)

            self.lg('Check org-11 pending invitations are empty, should succeed with 200')
            response = self.client_1.api.GetPendingOrganizationInvitations(globalid)
            self.assertEqual(response.status_code, 200)
            self.assertEqual(response.json(), [])

            self.lg('Check user_2 invitations are empty, should succeed with 200')
            response = self.client_2.api.GetNotifications(self.user_2)
            self.assertEqual(response.status_code, 200)
            self.assertEqual(response.json()['organizationinvitations'], [])

            self.lg('Cancel org-11 pending invitation for nonexisting organization to join org_11 role %s , should fail with 404' % role)
            response = self.client_1.api.RemovePendingOrganizationInvitation('fake_organization', globalid)
            self.assertEqual(response.status_code, 404)

        self.lg('%s ENDED' % self._testID)

    @unittest.skip('bug: #454')
    def test012_post_put_delete_members(self):
        """
            ITSYOU-051
            - Create org-11 organization and make user_1 an owner
            - User_1 invite user_2 to be a member of org-11, should succeed with 201
            - User_2 accept the invitations, should succeed with 201
            - User_2 remove user_1 as owner of org-11, should fail with 403
            - User_1 invite user_2 to be a member of org-11 again, should fail with 409
            - User_1 update user_2 role to be an owner of org-11, should succeed with 200
            - User_1 update user_2 role to be a member, should succeed with 200
            - User_1 update nonexisting user\'s role , should fail with 404
            - User_1 Delete user_2 from org-11, should succeed with 204
            - User_1 Delete nonexisting user from org-11, should fail with 404
        """
        self.lg('%s STARTED' % self._testID)

        self.lg('[POST] User_1 invite user_2 to be a member of org-11, should succeed with 201')
        data = {"searchstring":self.user_2}
        response = self.client_1.api.AddOrganizationMember(data, self.org_11_globalid)
        self.assertEqual(response.status_code, 201)

        self.lg('[POST] User_2 accept the invitations, should succeed with 201')
        response = self.client_2.api.AcceptMembership(self.org_11_globalid, 'member', self.user_2)
        self.assertEqual(response.status_code, 201)

        self.lg('[DEL] User_2 remove user_1 as owner of org-11, should fail with 403')
        response = self.client_2.api.RemoveOrganizationMember(self.user_1, self.org_11_globalid)
        self.assertEqual(response.status_code, 403)

        response = self.client_2.api.GetOrganization(self.org_11_globalid)
        self.assertEqual(response.status_code, 200)
        self.assertIn(self.user_2, response.json()['members'])

        self.lg('[POST] User_1 invite user_2 to be a member of org-11 again, should succeed with 201')
        data = {"searchstring":self.user_2}
        response = self.client_1.api.AddOrganizationMember(data, self.org_11_globalid)
        self.assertEqual(response.status_code, 409)

        self.lg('[PUT] User_1 update user_2 role to be an owner of org-11, should succeed with 200')
        data = {"username":self.user_2, "role":"owners"}
        response = self.client_1.api.UpdateOrganizationMemberShip(data, self.org_11_globalid)
        self.assertEqual(response.status_code, 200)

        response = self.client_1.api.GetOrganization(self.org_11_globalid)
        self.assertEqual(response.status_code, 200)
        self.assertIn(self.user_2, response.json()['owners'])

        self.lg('[PUT] User_1 update user_2 role to be a member, should succeed with 200')
        data = {"username":self.user_2, "role":"members"}
        response = self.client_1.api.UpdateOrganizationMemberShip(data, self.org_11_globalid)
        self.assertEqual(response.status_code, 200)

        response = self.client_1.api.GetOrganization(self.org_11_globalid)
        self.assertEqual(response.status_code, 200)
        self.assertIn(self.user_2, response.json()['members'])

        #bug #454 internal server error
        self.lg('[PUT] User_1 update nonexisting user\'s role , should fail with 404')
        data = {"username":"fake_user", "role":"owners"}
        response = self.client_1.api.UpdateOrganizationMemberShip(data, self.org_11_globalid)
        self.assertEqual(response.status_code, 404)

        self.lg('[DEL] User_1 Delete user_2 from org-11, should succeed with 204')
        response = self.client_1.api.RemoveOrganizationMember(self.user_2, self.org_11_globalid)
        self.assertEqual(response.status_code, 204)

        self.lg('[DEL] User_1 Delete nonexisting user from org-11, should fail with 404')
        response = self.client_1.api.RemoveOrganizationMember(self.user_1, self.org_11_globalid)
        self.assertEqual(response.status_code, 204)

        self.lg('%s ENDED' % self._testID)

    def test013_post_delete_owners(self):
        """
            ITSYOU-052
            - Create organization org-11 and make user_1 an owner
            - User_1 invite user_2 to be an owner of org-11, should succeed with 201
            - User_2 accept the invitations, should succeed with 201
            - User_1 invite user_2 to be an owner of org-11 again, should fail with 409
            - User_1 Delete user_2 from org-11, should succeed with 204
            - User_1 Delete nonexisting user from org-11, should fail with 404
        """
        self.lg('%s STARTED' % self._testID)

        self.lg('[POST] User_1 invite user_2 to be an owner of org-11, should succeed with 201')
        data = {"searchstring":self.user_2}
        response = self.client_1.api.AddOrganizationOwner(data, self.org_11_globalid)
        self.assertEqual(response.status_code, 201)

        self.lg('[POST] User_2 accept the invitations, should succeed with 201')
        response = self.client_2.api.AcceptMembership(self.org_11_globalid, 'owner', self.user_2)
        self.assertEqual(response.status_code, 201)

        response = self.client_2.api.GetOrganization(self.org_11_globalid)
        self.assertEqual(response.status_code, 200)
        self.assertIn(self.user_2, response.json()['owners'])

        self.lg('[POST] User_1 invite user_2 to be an owner of org-11 again, should succeed with 201')
        data = {"searchstring":self.user_2}
        response = self.client_1.api.AddOrganizationOwner(data, self.org_11_globalid)
        self.assertEqual(response.status_code, 409)

        self.lg('[DEL] User_1 Delete user_2 from org-11, should succeed with 204')
        response = self.client_1.api.RemoveOrganizationMember(self.user_2, self.org_11_globalid)
        self.assertEqual(response.status_code, 204)

        self.lg('[DEL] User_1 Delete nonexisting user from org-11, should fail with 404')
        response = self.client_1.api.RemoveOrganizationMember(self.user_1, self.org_11_globalid)
        self.assertEqual(response.status_code, 204)

        self.lg('%s ENDED' % self._testID)

    @unittest.skip('bug: #461')
    def test014_post_put_delete_orgmembers(self):
        """
            ITSYOU-053
            - Create organization org-11 and make user_1 an owner
            - Create organization org-12 and make user_1 an owner & user_2 a member
            - Create organization org-21 and make user_2 an owner
            - Create organization org-13 and make user_1
            - Create suborganization inside org-13 and make user_2 an owner
            - Add org-12 to org-11 orgmembers, should succeed with 201
            - Add org-12 to org-11 orgmembers again, should fail with 409
            - User_2 remove user_1 as owner of org-11, should fail with 403
            - invite org-21 to be orgmember in org-11, should succeed with 201
            - check org-11 pending invitations, should succeed with 200
            - User_2 accept User_1 invitation
            - Add nonexisting organization to org-11 orgmembers, should fail with 404
            - Update org-12 from orgmember to orgowner of org-11, should succeed with 200
            - Update org-12 from orgowner to orgmember of org-11, should succeed with 200
            - Update org-21 organizations membership from orgmember to orgowner, should succeed with 200    ## BUG: 461
            - Update org-21 organizations membership from orgowner to orgmember, should fail with 200       ## BUG: 461
            - Update nonexisting organizations membership, should fail with 404
            - Remove org-12 as orgmember of org-11, should succeed with 204
            - Remove org-21 as orgmember of org-11, should succeed with 204
            - Remove nonexisting organization as orgmember of org-11, should fail with 404
            - Add org-13 to org-11 orgmembers, should succeed with 201
            - User_2 get org-11 info , should fail with 403
            - Include suborganizations of org-13 to be orgmembers of org-11, should succeed with 201
            - User_2 get org-11 info , should succeed with 200
            - Remove suborganizations of org-13 from orgmembers of org-11, should succeed with 204
            - User_2 get org-11 info , should fail with 403
            - Include suborganization of nonexisting organization to be orgmembers of org-11, should fail with 404
            - Remove suborganizations of nonexisting organization from orgmembers of org-11, should fail with 404
        """
        self.lg('%s STARTED' % self._testID)

        self.lg('[POST] Add org-12 to org-11 orgmembers, should succeed with 201')
        data = {"orgmember":self.org_12_globalid}
        response = self.client_1.api.SetOrgMember(data, self.org_11_globalid)
        self.assertEqual(response.status_code, 201)

        response = self.client_1.api.GetOrganization(self.org_11_globalid)
        self.assertEqual(response.status_code, 200)
        self.assertIn(self.org_12_globalid, response.json()['orgmembers'])

        self.lg('[POST] Add org-12 to org-11 orgmembers again, should fail with 409')
        data = {"orgmember":self.org_12_globalid}
        response = self.client_1.api.SetOrgMember(data, self.org_11_globalid)
        self.assertEqual(response.status_code, 409)

        self.lg('[DEL] User_2 remove user_1 as owner of org-11, should fail with 403')
        response = self.client_2.api.RemoveOrganizationMember(self.user_1, self.org_11_globalid)
        self.assertEqual(response.status_code, 403)

        self.lg('[POST] invite org-21 to be orgmember in  org-11 , should succeed with 201')
        data = {"orgmember":self.org_21_globalid}
        response = self.client_1.api.SetOrgMember(data, self.org_11_globalid)
        self.assertEqual(response.status_code, 201)

        self.lg('Get org-11 pending invitations, should succeed with 200')
        response = self.client_1.api.GetPendingOrganizationInvitations(self.org_11_globalid)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(response.json()[-1]['status'], 'pending')
        self.assertEqual(response.json()[-1]['organization'], self.org_11_globalid)
        self.assertEqual(response.json()[-1]['user'], self.org_21_globalid)
        self.assertEqual(response.json()[-1]['role'], 'orgmember')
        self.assertTrue(response.json()[-1]['isorganization'])

        self.lg('Check user_2 invitations, should succeed with 200')
        response = self.client_2.api.GetNotifications(self.user_2)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(self.org_11_globalid, response.json()['organizationinvitations'][-1]['organization'])
        self.assertEqual('orgmember', response.json()['organizationinvitations'][-1]['role'])
        self.assertEqual('pending', response.json()['organizationinvitations'][-1]['status'])

        self.lg('[POST] User_2 accept User_1 invitation')
        response = self.client_2.api.AcceptOrgMembership(self.org_21_globalid, 'orgmember', self.org_11_globalid)
        self.assertEqual(response.status_code, 201)

        self.lg('[POST] Add nonexisting to org-11 orgmembers, should fail with 404')
        data = {"orgmember":'fake_org'}
        response = self.client_1.api.SetOrgMember(data, self.org_11_globalid)
        self.assertEqual(response.status_code, 404)

        self.lg('[PUT] Update org-12 from orgmember to orgowner of org-11, should succeed with 200')
        data = {"org":self.org_12_globalid, "role":"owners"}
        response = self.client_1.api.UpdateOrganizationOrgMemberShip(data, self.org_11_globalid)
        self.assertEqual(response.status_code, 200)
        self.assertIn(self.org_12_globalid, response.json()['orgowners'])

        self.lg('[PUT] Update org-12 from orgowner to orgmember of org-11, should succeed with 200')
        data = {"org":self.org_12_globalid, "role":"members"}
        response = self.client_1.api.UpdateOrganizationOrgMemberShip(data, self.org_11_globalid)
        self.assertEqual(response.status_code, 200)
        self.assertIn(self.org_12_globalid, response.json()['orgmembers'])

        # self.lg('[PUT] Update org-21 organizations membership from orgmember to orgowner, should succeed with 200')
        # data = {"org":self.org_21_globalid, "role":"members"}
        # response = self.client_1.api.UpdateOrganizationOrgMemberShip(data, self.org_11_globalid)
        # self.assertEqual(response.status_code, 200)
        #
        # self.lg('[PUT] Update org-21 organizations membership from orgowner to orgmember, should fail with 200')
        # data = {"org":self.org_21_globalid, "role":"members"}
        # response = self.client_1.api.UpdateOrganizationOrgMemberShip(data, self.org_11_globalid)
        # self.assertEqual(response.status_code, 200)

        self.lg('[PUT] Update nonexisting organizations membership, should fail with 404')
        data = {"org":"fake_org", "role":"members"}
        response = self.client_1.api.UpdateOrganizationOrgMemberShip(data, self.org_11_globalid)
        self.assertEqual(response.status_code, 404)

        self.lg('[DEL] Remove org-12 as orgmember of org-11, should succeed with 204')
        response = self.client_1.api.DeleteOrgMember(self.org_12_globalid, self.org_11_globalid)
        self.assertEqual(response.status_code, 204)

        response = self.client_1.api.GetOrganization(self.org_11_globalid)
        self.assertEqual(response.status_code, 200)
        self.assertNotIn(self.org_12_globalid, response.json()['orgmembers'])

        self.lg('[DEL] Remove org-21 as orgmember of org-11, should succeed with 204')
        response = self.client_1.api.DeleteOrgMember(self.org_21_globalid, self.org_11_globalid)
        self.assertEqual(response.status_code, 204)

        response = self.client_1.api.GetOrganization(self.org_11_globalid)
        self.assertEqual(response.status_code, 200)
        self.assertNotIn(self.org_21_globalid, response.json()['orgmembers'])

        self.lg('[DEL] Remove nonexisting organization as orgmember of org-11, should fail with 404')
        response = self.client_1.api.DeleteOrgMember('fake_org', self.org_11_globalid)
        self.assertEqual(response.status_code, 404)

        self.lg('[POST]  Add org-13 to org-11 orgmembers, should succeed with 201')
        data = {"orgmember":self.org_13_globalid}
        response = self.client_1.api.SetOrgMember(data, self.org_11_globalid)
        self.assertEqual(response.status_code, 201)

        self.lg('[GET] User_2 get org-11 info , should fail with 403')
        response = self.client_2.api.GetOrganization(self.org_11_globalid)
        self.assertEqual(response.status_code, 403)

        self.lg('[POST] Include suborganizations of org-13 to be orgmembers of org-11, should succedd with 201')
        data = {"globalid":self.org_13_globalid}
        response = self.client_1.api.IncludeSuborgsof(data, self.org_11_globalid)
        self.assertEqual(response.status_code, 201)

        self.lg('[GET] User_2 get org-11 info , should succeed with 200')
        response = self.client_2.api.GetOrganization(self.org_11_globalid)
        self.assertEqual(response.status_code, 200)

        self.lg('[DEL] Remove suborganizations of org-13 from orgmembers of org-11, should succedd with 204')
        response = self.client_1.api.RemoveIncludeSuborgsof(self.org_11_globalid, self.org_13_globalid)
        self.assertEqual(response.status_code, 204)

        self.lg('[GET] User_2 get org-11 info , should fail with 403')
        response = self.client_2.api.GetOrganization(self.org_11_globalid)
        self.assertEqual(response.status_code, 403)

        self.lg('[POST] Include suborganizations of nonexisting organization to be orgmembers of org-11, should fail with 404')
        data = {"globalid":'fake_organization'}
        response = self.client_1.api.IncludeSuborgsof(data, self.org_11_globalid)
        self.assertEqual(response.status_code, 404)

        self.lg('[DEL] Remove suborganizations of nonexisting organization from orgmembers of org-11, should fail with 404')
        response = self.client_1.api.RemoveIncludeSuborgsof(self.org_11_globalid, 'fake_organization')
        self.assertEqual(response.status_code, 404)

        self.lg('%s ENDED' % self._testID)

    def test015_post_delete_orgowners(self):
        """
            ITSYOU-054
            - Create organization org-11 and make user_1 an owner
            - Create organization org-12 and make user_1 an owner & user_2 a member
            - Create organization org-21 and make user_2 an owner & user_1 a member
            - Create organization org-21 and make user_2 an owner & user_1 a member
            - Create organization org-13 and make user_1
            - Create suborganization inside org-13 and make user_2 an owner
            - Add org-12 to org-11 orgowners, should succeed with 201
            - Add org-12 to org-11 orgowners again, should fail with 409
            - invite org-21 to be org-11 orgmembers, should succeed with 201
            - Add nonexisting organization to org-11 orgowners, should fail with 404
            - Remove org-12 as orgowners of org-11, should succeed with 204
            - Remove org-21 as orgowners of org-11, should succeed with 204
            - Remove nonexisting organization as orgowners of org-11, should fail with 404
            - Add org-13 to org-11 orgowners, should succeed with 201
            - User_2 get org-11 info , should fail with 403
            - Include suborganization of org-13 to be orgowners of org-11, should succeed with 201
            - User_2 get org-11 info , should succeed with 200
            - Remove suborganization of org-13 from orgowners of org-11, should succeed with 204
            - User_2 get org-11 info , should fail with 403
            - Include suborganization of nonexisting organization to be orgowners of org-11, should fail with 404
            - Remove suborganization of nonexisting organization from orgowners of org-11, should fail with 404
        """
        self.lg('%s STARTED' % self._testID)

        self.lg('[POST] Add org-12 to org-11 orgowner, should succeed with 201')
        data = {"orgowner":self.org_12_globalid}
        response = self.client_1.api.SetOrgOwner(data, self.org_11_globalid)
        self.assertEqual(response.status_code, 201)

        response = self.client_1.api.GetOrganization(self.org_11_globalid)
        self.assertEqual(response.status_code, 200)
        self.assertIn(self.org_12_globalid, response.json()['orgowners'])

        self.lg('[POST] Add org-12 to org-11 orgowner again, should fail with 409')
        data = {"orgowner":self.org_12_globalid}
        response = self.client_1.api.SetOrgOwner(data, self.org_11_globalid)
        self.assertEqual(response.status_code, 409)

        self.lg('[POST] Add org-21 org-11 orgmembers, should succeed with 201')
        data = {"orgowner":self.org_21_globalid}
        response = self.client_1.api.SetOrgOwner(data, self.org_11_globalid)
        self.assertEqual(response.status_code, 201)

        self.lg('Get org-11 pending invitations, should succeed with 200')
        response = self.client_1.api.GetPendingOrganizationInvitations(self.org_11_globalid)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(response.json()[-1]['status'], 'pending')
        self.assertEqual(response.json()[-1]['organization'], self.org_11_globalid)
        self.assertEqual(response.json()[-1]['user'], self.org_21_globalid)
        self.assertEqual(response.json()[-1]['role'], 'orgowner')
        self.assertTrue(response.json()[-1]['isorganization'])

        self.lg('Check user_2 invitations, should succeed with 200')
        response = self.client_2.api.GetNotifications(self.user_2)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(self.org_11_globalid, response.json()['organizationinvitations'][-1]['organization'])
        self.assertEqual('orgowner', response.json()['organizationinvitations'][-1]['role'])
        self.assertEqual('pending', response.json()['organizationinvitations'][-1]['status'])

        self.lg('[POST] User_2 accept User_1 invitation')
        response = self.client_2.api.AcceptOrgMembership(self.org_21_globalid, 'orgowner', self.org_11_globalid)
        self.assertEqual(response.status_code, 201)

        self.lg('[POST] Add nonexisting organization to org-11 orgowner, should fail with 404')
        data = {"orgowner":'fake_org'}
        response = self.client_1.api.SetOrgOwner(data, self.org_11_globalid)
        self.assertEqual(response.status_code, 404)

        self.lg('[DEL] Remove org-12 as orgowner of org-11, should succeed with 204')
        response = self.client_1.api.DeleteOrgOwner(self.org_12_globalid, self.org_11_globalid)
        self.assertEqual(response.status_code, 204)
        response = self.client_1.api.GetOrganization(self.org_11_globalid)
        self.assertEqual(response.status_code, 200)
        self.assertNotIn(self.org_12_globalid, response.json()['orgowners'])

        self.lg('[DEL] Remove org-21 as orgowner of org-11, should succeed with 204')
        response = self.client_1.api.DeleteOrgOwner(self.org_21_globalid, self.org_11_globalid)
        self.assertEqual(response.status_code, 204)
        response = self.client_1.api.GetOrganization(self.org_11_globalid)
        self.assertEqual(response.status_code, 200)
        self.assertNotIn(self.org_21_globalid, response.json()['orgowners'])

        self.lg('[DEL] Remove nonexisting organization as orgowner of org-11, should fail with 404')
        response = self.client_1.api.DeleteOrgOwner('fake_org', self.org_11_globalid)
        self.assertEqual(response.status_code, 404)

        self.lg('[POST]  Add org-13 to org-11 orgowners, should succeed with 201')
        data = {"orgowner":self.org_13_globalid}
        response = self.client_1.api.SetOrgOwner(data, self.org_11_globalid)
        self.assertEqual(response.status_code, 201)

        self.lg('[GET] User_2 get org-11 info , should fail with 403')
        response = self.client_2.api.GetOrganization(self.org_11_globalid)
        self.assertEqual(response.status_code, 403)

        self.lg('[POST] Include suborganization of org-13 to be orgowners of org-11, should succedd with 201')
        data = {"globalid":self.org_13_globalid}
        response = self.client_1.api.IncludeSuborgsof(data, self.org_11_globalid)
        self.assertEqual(response.status_code, 201)

        self.lg('[GET] User_2 get org-11 info , should succeed with 200')
        response = self.client_2.api.GetOrganization(self.org_11_globalid)
        self.assertEqual(response.status_code, 200)

        self.lg('[DEL] Remove suborganization of org-13 from orgowners of org-11, should succedd with 204')
        response = self.client_1.api.RemoveIncludeSuborgsof(self.org_11_globalid, self.org_13_globalid)
        self.assertEqual(response.status_code, 204)

        self.lg('[GET] User_2 get org-11 info , should fail with 403')
        response = self.client_2.api.GetOrganization(self.org_11_globalid)
        self.assertEqual(response.status_code, 403)

        self.lg('[POST] Include suborganization of nonexisting organization to be orgowners of org-11, should fail with 404')
        data = {"globalid":'fake_organization'}
        response = self.client_1.api.IncludeSuborgsof(data, self.org_11_globalid)
        self.assertEqual(response.status_code, 404)

        self.lg('[DEL] Remove suborganization of nonexisting organization from orgowners of org-11, should fail with 404')
        response = self.client_1.api.RemoveIncludeSuborgsof(self.org_11_globalid, 'fake_organization')
        self.assertEqual(response.status_code, 404)

        self.lg('%s ENDED' % self._testID)

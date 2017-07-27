from api_testing.utils import BaseTest
import types
import unittest
import json
import time

class UsersTestsB(BaseTest):

    def setUp(self):
        super(UsersTestsB, self).setUp()
        self.SetTotp(self.user_1)

        org_1_data = {"globalid":self.organization_1}
        response = self.client_1.api.CreateNewOrganization(org_1_data)
        self.assertEqual(response.status_code, 201)
        org_2_data = {"globalid":self.organization_2}
        response = self.client_2.api.CreateNewOrganization(org_2_data)
        self.assertEqual(response.status_code, 201)
        ## user_2 invite user_1 to be owner
        data = {'searchstring': self.user_1}
        response = self.client_2.api.AddOrganizationOwner(data, self.organization_2)
        self.assertEqual(response.status_code, 201)
        ## user_1 accept invitation
        response = self.client_1.api.AcceptMembership(self.organization_2 , 'owner', self.user_1)
        self.assertEqual(response.status_code, 201)

    def tearDown(self):

        self.SetTotp(self.user_1)
        self.lg('Delete all user\'s email addresses, should succeed with 204')
        self.DeleteAllUserEmails(self.user_1)
        self.lg('Delete all user\'s addresses, should succeed with 204')
        self.DeleteAllUserAddresses(self.user_1)
        self.lg('Delete all user\'s publicKeys, should succeed with 204')
        self.DeleteAllUserPublicKeys(self.user_1)
        self.lg('Delete all user\'s phonenumbers, should succeed with 204')
        self.DeleteAllUserPhonenumbers(self.user_1)
        self.lg('Delete all user\'s registries, should succeed with 204')
        self.DeleteAllUserRegistries(self.user_1)
        self.lg('Delete all user\'s bank accounts, should succeed with 204')
        self.DeleteAllUserBankAccounts(self.user_1)
        self.lg('Delete all user\'s digital wallets, should succeed with 204')
        self.DeleteAllUserDigitalWallet(self.user_1)
        #bug 405
        # self.lg('Delete all user\'s apikeys, should succeed with 204')
        #self.DeleteAllUserApiKeys(self.user_1)

        response = self.client_1.api.DeleteOrganization(self.organization_1)
        self.assertEqual(response.status_code, 204)
        response = self.client_1.api.DeleteOrganization(self.organization_2)
        self.assertEqual(response.status_code, 204)

        super(UsersTestsB, self).tearDown()

    def test001_get_user_info(self):

        """
            ITSYOU-024
            - Get username information using /{username}, should succeed with 200
            - Get username information using /{username}/info, should succeed with 200
        """
        self.lg('%s STARTED' % self._testID)

        self.lg('[GET] Get username information using /{username}, should succeed with 200')
        response_1 = self.client_1.api.GetUser(self.user_1)
        self.assertEqual(response_1.status_code, 200)
        self.assertEqual(response_1.json()['username'], self.user_1)

        #bug 412
        self.lg('[GET] Get username information using /{username}/info, should succeed with 200')
        response_2 = self.client_1.api.GetUser(self.user_1)
        self.assertEqual(response_2.status_code, 200)
        self.assertEqual(response_2.json()['username'], self.user_1)

        self.assertEqual(response_1.json()['addresses'], response_2.json()['addresses'])
        self.assertEqual(response_1.json()['bankaccounts'], response_2.json()['bankaccounts'])
        self.assertEqual(response_1.json()['emailaddresses'], response_2.json()['emailaddresses'])
        self.assertEqual(response_1.json()['phonenumbers'], response_2.json()['phonenumbers'])
        self.assertEqual(response_1.json()['publicKeys'], response_2.json()['publicKeys'])


        self.lg('%s ENDED' % self._testID)

    @unittest.skip('bug:')
    def test002_put_name(self):

        """
            ITSYOU-025
            - Change firstname & lastname with valid user, should succeed with 204
            - Change firstname & lastname with invalid user, should fail with 404
        """

        self.lg('%s STARTED' % self._testID)

        self.lg('[PUT] Change firstname & lastname to a valid user, should succeed with 204')
        firstname = self.random_value()
        lastname  = self.random_value()
        data = {"firstname": firstname, "lastname": lastname}
        response = self.client_1.api.UpdateUserName(data, self.user_1)
        self.assertEqual(response.status_code, 204)
        response = self.client_1.api.GetUser(self.user_1)
        self.assertEqual(response.json()['firstname'], firstname)
        self.assertEqual(response.json()['lastname'], lastname)

        # bug 403 instead of 404
        self.lg('[PUT] Change firstname & lastname to invalid user, should fail with 404')
        response = self.client_1.api.UpdateUserName(data, 'fake_user')
        self.assertEqual(response.status_code, 404)

        self.lg('%s ENDED' % self._testID)

    def test003_put_password(self):
        """
            ITSYOU-026
            - Change password with valid current password, should succeed with 204
            - Change password with valid current password again, should succeed with 204
            - Change password with wrong current password, should fail with 422
            - Change password with invalid new password, should fail with 422
        """

        self.lg('%s STARTED' % self._testID)

        self.lg('[PUT] Change password with valid current password, should succeed with 204')
        current_password = self.user_1_password
        new_password = 'TheNewPassword123456'
        data = {'currentpassword':current_password, 'newpassword':new_password}
        response = self.client_1.api.UpdatePassword(data, self.user_1)
        self.assertEqual(response.status_code, 204)

        self.lg('[PUT] Change password with valid current password again, should succeed with 204')
        current_password = new_password
        new_password = self.user_1_password
        data = {'currentpassword':current_password, 'newpassword':new_password}
        response = self.client_1.api.UpdatePassword(data, self.user_1)
        self.assertEqual(response.status_code, 204)

        self.lg('[PUT] Change password with wrong current password, should fail with 422')
        current_password = self.random_value()
        new_password = self.random_value()
        data = {'currentpassword':current_password, 'newpassword':new_password}
        response = self.client_1.api.UpdatePassword(data, self.user_1)
        self.assertEqual(response.status_code, 422)
        self.assertEqual(response.json()['error'], 'incorrect_password')

        self.lg('Change password with invalid new password, should fail with 422')
        current_password = self.user_1_password
        new_password = self.random_value()[0:3]
        data = {'currentpassword':current_password, 'newpassword':new_password}
        response = self.client_1.api.UpdatePassword(data, self.user_1)
        self.assertEqual(response.status_code, 422)
        self.assertEqual(response.json()['error'], 'invalid_password')

        self.lg('%s ENDED' % self._testID)

    def test004_get_post_put_delete_email_address(self):

        """
            ITSYOU-027
            - Create new email address, should succeed with 201
            - Get user\'s email addresses, should succeed with 200
            - Create new email address with label already exists, should fail with 409
            - Update email address & label, should succeed with 201
            - Update email label with label already exists, should fail with 409
            - Sends validation email to email address, should succeed with 204
            - Sends validation email to invalid email address, should fail with 404
            - Delete email address, should succeed with 204
            - Delete invalid email address, should fail with 404
            - Delete last email address, should fail with 409
        """

        self.lg('%s STARTED' % self._testID)

        self.lg('[POST] Create new email address, should succeed with 201')
        label = self.random_value()
        new_email_address = self.random_value() + "@gig.com"
        data = {"emailaddress":new_email_address, "label":label}
        response = self.client_1.api.RegisterNewEmailAddress(data, self.user_1)
        self.assertEqual(response.status_code, 201)

        self.lg('[GET] Get email addresses, should succeed with 200')
        response = self.client_1.api.GetEmailAddresses(self.user_1)
        self.assertEqual(response.status_code, 200)
        self.assertIn(data, response.json())

        self.lg('[POST] Create new email address with label already exist, should fail with 409')
        response = self.client_1.api.GetEmailAddresses(self.user_1)
        self.assertEqual(response.status_code, 200)
        label = response.json()[-1]['label']
        new_email_address = self.random_value() + "@gig.com"
        data = {'emailaddress':new_email_address, 'label':label}
        response = self.client_1.api.RegisterNewEmailAddress(data, self.user_1)
        self.assertEqual(response.status_code, 409)

        self.lg('[PUT] Update email address & label, should succeed with 201')
        response = self.client_1.api.GetEmailAddresses(self.user_1)
        self.assertEqual(response.status_code, 200)
        label = response.json()[-1]['label']
        new_label = self.random_value()
        new_email_address = self.random_value() + "@gig.com"
        data = {'emailaddress':new_email_address, 'label':new_label}
        response = self.client_1.api.UpdateEmailAddress(data,label, self.user_1)
        self.assertEqual(response.status_code, 201)
        response = self.client_1.api.GetEmailAddresses(self.user_1)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(new_email_address, response.json()[-1]['emailaddress'])
        self.assertEqual(new_label, response.json()[-1]['label'])

        self.lg('[PUT] Update email label with label already exists, should fail with 409')
        response = self.client_1.api.GetEmailAddresses(self.user_1)
        self.assertEqual(response.status_code, 200)
        label = response.json()[-1]['label']
        new_label = response.json()[-2]['label']
        new_email_address = self.random_value() + "@gig.com"
        data = {'emailaddress':new_email_address, 'label':new_label}
        response = self.client_1.api.UpdateEmailAddress(data, label, self.user_1)
        self.assertEqual(response.status_code, 409)

        self.lg('[POST] Create new email address & validate it, should succeed with 201')
        label = "validation email"
        #check if validation email already exist then modify it
        response = self.client_1.api.GetEmailAddresses(self.user_1)
        self.assertEqual(response.status_code, 200)
        if label in [x['label'] for x in response.json()]:
            data = {'emailaddress':self.random_value()+"@gig.com", 'label':self.random_value()}
            response = self.client_1.api.UpdateEmailAddress(data,label, self.user_1)

        new_email_address = self.validation_email
        data = {"emailaddress":new_email_address, "label":label}
        response = self.client_1.api.RegisterNewEmailAddress(data, self.user_1)
        self.assertEqual(response.status_code, 201)
        time.sleep(20)
        response = self.client_1.api.ValidateEmailAddress(label, self.user_1)
        self.assertEqual(response.status_code, 204)
        time.sleep(20)
        self.lg('Check the validation message & validate, should succeed with 200')
        response = self.UserValidateEmail()
        self.assertEqual(response.status_code, 200)

        response = self.client_1.api.GetEmailAddresses(self.user_1, query_params={'validated':True})
        self.assertEqual(response.status_code, 200)
        self.assertIn(label, [x['label'] for x in response.json()])

        self.lg('[POST] Sends validation email to invalid email address, should fail with 404')
        response = self.client_1.api.ValidateEmailAddress('fake_label', self.user_1)
        self.assertEqual(response.status_code, 404)

        self.lg('[DELETE] Delete email address, should succeed with 204')
        response = self.client_1.api.GetEmailAddresses(self.user_1)
        self.assertEqual(response.status_code, 200)
        label = response.json()[-1]['label']
        response = self.client_1.api.DeleteEmailAddress(label, self.user_1)
        self.assertEqual(response.status_code, 204)
        response = self.client_1.api.GetEmailAddresses(self.user_1)
        self.assertEqual(response.status_code, 200)
        self.assertNotEqual(label, response.json()[-1]['label'])

        self.lg('[DELETE] Delete invalid email address, should fail with 404')
        response = self.client_1.api.DeleteEmailAddress("fake_label", self.user_1)
        self.assertEqual(response.status_code, 404)

        self.lg('[DELETE] Delete last email address, should fail with 409')
        response = self.client_1.api.GetEmailAddresses(self.user_1)
        self.assertEqual(response.status_code, 200)
        labels = [x['label'] for x in response.json()]
        for i in range(0, len(labels)-1):
            response = self.client_1.api.DeleteEmailAddress(labels[i], self.user_1)
            self.assertEqual(response.status_code, 204)
        else:
            response = self.client_1.api.DeleteEmailAddress(labels[len(labels)-1], self.user_1)
            self.assertEqual(response.status_code, 409)

        self.lg('%s ENDED' % self._testID)

    @unittest.skip(' ')
    def test005_get_post_put_delete_phonenumber(self):
        """
            ITSYOU-028
            - Register a new phonenumber (1), should succeed with 201
            - Register a new phonenumber (2), should succeed with 201
            - Get user\'s phonenumbers, should succeed with 200
            - Get phonenumber (2) by label, should succeed with 200
            - Register a new phonenumber with invalid number, should fail with 400
            - Register a new phonenumber with existing label, should fail with 409
            - Get phonenumber by nonexisting label, should fail with 404
            - Update phonenumber (2), should succeed with 201
            - Update phonenumber with invalid number, should fail with 400
            - Update phonenumber (2) label with label already exists (label of phonenumber (1)), should fail with 409
            - Register a new phonenumber (3), send validation sms & verify it, should succeed with 201 & 200
            - Delete totp method, should succeed with 204
            - Delete verified phonenumber (3), should fail with 409
            - Force Delete verified phonenumber (3), should fail with 409
            - Set totp method, should succeed with 201
            - Delete verified phonenumber (3), should fail with 409
            - Force Delete verified phonenumber (3), should succeed with 204
            - Validate & verify invalid phonenumber, should fail with 404
            - Delete phonenumber (2), should succeed with 204
            - Delete phonenumber (1), should succeed with 204
            - Delete invalid phone number, should fail with 404
        """

        self.lg('%s STARTED' % self._testID)

        self.lg('[POST] Register a new phonenumber (1), should succeed with 201')
        label = self.random_value()
        phonenumber = "+0123456789"
        data = {"label":label, "phonenumber":phonenumber}
        response = self.client_1.api.RegisterNewUserPhonenumber(data, self.user_1)
        self.assertEqual(response.status_code, 201)
        self.assertEqual(response.json(), data)

        self.lg('[POST] Register a new phonenumber (2), should succeed with 201')
        new_label = self.random_value()
        new_phonenumber = "+01987654321"
        new_data = {"label":new_label, "phonenumber":new_phonenumber}
        response = self.client_1.api.RegisterNewUserPhonenumber(new_data, self.user_1)
        self.assertEqual(response.status_code, 201)
        self.assertEqual(response.json(), new_data)

        self.lg('[GET] Get user phonenumbers, should succeed with 200')
        response = self.client_1.api.GetUserPhoneNumbers(self.user_1)
        self.assertEqual(response.status_code, 200)
        self.assertIn(data, response.json())
        self.assertIn(new_data, response.json())

        self.lg('[GET] Get phonenumber (2) by label, should succeed with 200')
        label = response.json()[-1]['label']
        response = self.client_1.api.GetUserPhonenumberByLabel(label, self.user_1)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(new_data, response.json())

        self.lg('[POST] Register a new phonenumber with invalid number, should fail with 400')
        new_label = self.random_value()
        new_phonenumber = self.random_value()
        new_data = {"label":new_label, "phonenumber":new_phonenumber}
        response = self.client_1.api.RegisterNewUserPhonenumber(new_data, self.user_1)
        self.assertEqual(response.status_code, 400)

        self.lg('[POST] Register a new phonenumber with existing label, should fail with 409')
        response = self.client_1.api.GetUserPhoneNumbers(self.user_1)
        self.assertEqual(response.status_code, 200)
        new_label = response.json()[-1]['label']
        new_phonenumber = "+123456789"
        new_data = {"label":new_label, "phonenumber":new_phonenumber}
        response = self.client_1.api.RegisterNewUserPhonenumber(new_data, self.user_1)
        self.assertEqual(response.status_code, 409)

        self.lg('[GET] Get phonenumber by nonexisting label, should fail with 404')
        response = self.client_1.api.GetUserPhonenumberByLabel('fake_label', self.user_1)
        self.assertEqual(response.status_code, 404)

        self.lg('[PUT] Update phonenumber (2), should succeed with 201')
        response = self.client_1.api.GetUserPhoneNumbers(self.user_1)
        self.assertEqual(response.status_code, 200)
        label = response.json()[-1]['label']
        new_label = self.random_value()
        new_phonenumber = "+987654321"
        new_data = {"label":new_label, "phonenumber":new_phonenumber}
        response = self.client_1.api.UpdateUserPhonenumber(new_data, label, self.user_1)
        self.assertEqual(response.status_code, 201)

        self.lg('[PUT] Update phonenumber with invalid number, should fail with 400')
        response = self.client_1.api.GetUserPhoneNumbers(self.user_1)
        self.assertEqual(response.status_code, 200)
        label = response.json()[-1]['label']
        new_label = self.random_value()
        new_phonenumber = self.random_value()
        new_data = {"label":new_label, "phonenumber":new_phonenumber}
        response = self.client_1.api.UpdateUserPhonenumber(new_data, label, self.user_1)
        self.assertEqual(response.status_code, 400)

        self.lg('[PUT] Update phonenumber (2) label with label already exists (label of phonenumber (1)), should fail with 409')
        response = self.client_1.api.GetUserPhoneNumbers(self.user_1)
        self.assertEqual(response.status_code, 200)
        label = response.json()[-1]['label']
        new_label = response.json()[-2]['label']
        new_phonenumber = "+123456789"
        new_data = {"label":new_label, "phonenumber":new_phonenumber}
        response = self.client_1.api.UpdateUserPhonenumber(new_data, label, self.user_1)
        self.assertEqual(response.status_code, 409)

        self.lg('[POST] Register a new phonenumber (3) , send validation sms & verify it, should succeed with 201 & 200')
        label = 'validation number'
        phonenumber = self.get_valid_phonenumber()
        data = {"label":label, "phonenumber":phonenumber}
        response = self.client_1.api.RegisterNewUserPhonenumber(data, self.user_1)
        self.assertEqual(response.status_code, 201)

        response = self.client_1.api.ValidatePhonenumber(label, self.user_1)
        self.assertEqual(response.status_code, 200)
        validationkey = response.json()['validationkey']

        time.sleep(25)

        smscode = self.get_mobile_verification_code()
        self.assertTrue(smscode, 'error while getting sms code, verification message not received with virtual number %s' % phonenumber)
        data = {"smscode":smscode, "validationkey":validationkey}
        response = self.client_1.api.VerifyPhoneNumber(data, label, self.user_1)
        self.assertEqual(response.status_code, 204)

        response = self.client_1.api.GetUserPhoneNumbers(self.user_1, query_params={'validated':True})
        self.assertEqual(response.status_code, 200)
        self.assertIn(label, [x['label'] for x in response.json()])

        self.lg('[DEL] Delete totp, should succeed with 204')
        response = self.client_1.api.DeleteTotp(self.user_1)
        self.assertEqual(response.status_code, 204)

        self.lg('[DELETE] Delete verified phonenumber (3), should fail with 409')
        response = self.client_1.api.DeleteUserPhonenumber(label, self.user_1)
        self.assertEqual(response.status_code, 409)
        self.assertEqual(response.json()['error'], 'warning_delete_last_verified_phone_number')

        self.lg('[DELETE] Force Delete verified phonenumber (3), should fail with 409')
        response = self.client_1.api.DeleteUserPhonenumber(label, self.user_1, query_params={'force':'true'})
        self.assertEqual(response.status_code, 409)
        self.assertEqual(response.json()['error'], 'cannot_delete_last_verified_phone_number')

        self.lg('[POST] Set totp code, should succeed with 204')
        secret = self.totp_secret
        totpcode = self.get_totp_code(secret)
        data = {"totpcode":totpcode, "totpsecret":secret}
        response = self.client_1.api.EditTotp(data, self.user_1)
        self.assertEqual(response.status_code, 204)

        self.lg('[DELETE] Delete verified phonenumber (3), should fail with 409')
        response = self.client_1.api.DeleteUserPhonenumber(label, self.user_1)
        self.assertEqual(response.status_code, 409)
        self.assertEqual(response.json()['error'], 'warning_delete_last_verified_phone_number')

        self.lg('[DELETE] Force Delete verified phonenumber (3), should succeed with 204')
        response = self.client_1.api.DeleteUserPhonenumber(label, self.user_1, query_params={'force':'true'})
        self.assertEqual(response.status_code, 204)

        self.lg('[POST] Validate & verify invalid phonenumber, should fail with 404')
        response = self.client_1.api.ValidatePhonenumber('fake_label', self.user_1)
        self.assertEqual(response.status_code, 404)

        self.lg('[DELETE] Delete phonenumber (2), should succeed with 204')
        response = self.client_1.api.GetUserPhoneNumbers(self.user_1)
        self.assertEqual(response.status_code, 200)
        label = response.json()[-1]['label']
        response = self.client_1.api.DeleteUserPhonenumber(label, self.user_1)
        self.assertEqual(response.status_code, 204)

        self.lg('[DELETE] Delete phonenumber (1), should succeed with 204')
        response = self.client_1.api.GetUserPhoneNumbers(self.user_1)
        self.assertEqual(response.status_code, 200)
        label = response.json()[-1]['label']
        response = self.client_1.api.DeleteUserPhonenumber(label, self.user_1)
        self.assertEqual(response.status_code, 204)

        self.lg('[DELETE] Delete invalid phone number, should fail with 404')
        response = self.client_1.api.DeleteUserPhonenumber(label, self.user_1)
        self.assertEqual(response.status_code, 404)

        self.lg('%s ENDED' % self._testID)

    def test006_get_post_put_delete_address(self):
        """
        ITSYOU-029
        - Register a new address (1), should succeed with 201
        - Register a new address (2), should succeed with 201
        - Get user\'s addresses, should succeed with 200
        - Get address (2) by label, should succeed with 200
        - Register a new address with label already exists (label of address (2)), should fail with 409
        - Register a new address with invalid inputs, should fail with 400
        - Update address (2), should succeed with 201
        - Update address (2) with label already exists (label of address (1)), should fail with 409
        - Update address with invalid inputs, should fail with 400
        - Delete address (2), should succeed with 204
        - Delete address (1), should succeed with 204
        - Delete invalid address, should fail with 404
        """

        elements = {"city":30, "country":40 ,"nr":10, "postalcode":20, "street":50, "other":30}

        self.lg('%s STARTED' % self._testID)

        self.lg('[POST] Register a new address (1), should succeed with 201')
        label = self.random_value()
        data = {"city": self.random_value(),
                "country": self.random_value(),
                "nr": self.random_value(),
                "other": self.random_value(),
                "postalcode": self.random_value(),
                "street": self.random_value(),
                "label": label}

        response = self.client_1.api.RegisterNewUserAddress(data, self.user_1)
        self.assertEqual(response.status_code, 201)
        self.assertEqual(response.json(), data)

        self.lg('[POST] Register a new address (2), should succeed with 201')
        label = self.random_value()
        new_data = {"city": self.random_value(),
                "country": self.random_value(),
                "nr": self.random_value(),
                "other": self.random_value(),
                "postalcode": self.random_value(),
                "street": self.random_value(),
                "label": label}

        response = self.client_1.api.RegisterNewUserAddress(new_data, self.user_1)
        self.assertEqual(response.status_code, 201)
        self.assertEqual(response.json(), new_data)

        self.lg('[GET] Get user addresses, should succeed with 200')
        response = self.client_1.api.GetUserAddresses(self.user_1)
        self.assertEqual(response.status_code, 200)
        self.assertIn(data, response.json())
        self.assertIn(new_data, response.json())

        self.lg('[GET] Get address (2) by label, should succeed with 200')
        response = self.client_1.api.GetUserAddresses(self.user_1)
        self.assertEqual(response.status_code, 200)
        label = response.json()[-1]['label']
        response = self.client_1.api.GetUserAddressByLabel(label, self.user_1)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(new_data, response.json())

        self.lg('[POST] Register a new address with label already exists (label of address (2)), should fail with 409')
        new_data['label'] = label
        response = self.client_1.api.RegisterNewUserAddress(data, self.user_1)
        self.assertEqual(response.status_code, 409)

        self.lg('[POST] Register a new address with invalid inputs, should fail with 400')
        new_data = dict(data)
        for key, value in elements.items():
            new_data['label'] = self.random_value()
            new_data[key] = self.random_value(value+1)
            response = self.client_1.api.RegisterNewUserAddress(new_data, self.user_1)
            self.assertEqual(response.status_code, 400)
            new_data = dict(data)

        self.lg('[PUT] Update address (2), should succeed with 201')
        response = self.client_1.api.GetUserAddresses(self.user_1)
        self.assertEqual(response.status_code, 200)
        label = response.json()[-1]['label']
        new_data = {"city": self.random_value(),
                "country": self.random_value(),
                "nr": self.random_value(),
                "other": self.random_value(),
                "postalcode": self.random_value(),
                "street": self.random_value(),
                "label": self.random_value()}

        response = self.client_1.api.UpdateUserAddress(new_data, label, self.user_1)
        self.assertEqual(response.status_code, 201)

        self.lg('[PUT] Update address (2) with label already exists (label of address (1)), should fail with 409')
        response = self.client_1.api.GetUserAddresses(self.user_1)
        self.assertEqual(response.status_code, 200)
        label = response.json()[-1]['label']
        new_label = response.json()[-2]['label']
        new_data = {"city": self.random_value(),
                    "country": self.random_value(),
                    "nr": self.random_value(),
                    "other": self.random_value(),
                    "postalcode": self.random_value(),
                    "street": self.random_value(),
                    "label": new_label}

        response = self.client_1.api.UpdateUserAddress(new_data, label, self.user_1)
        self.assertEqual(response.status_code, 409)

        self.lg('[PUT] Update address with invalid inputs, should fail with 400')

        new_data = dict(data)
        new_data['label'] = self.random_value()
        for key, value in elements.items():
            new_data[key] = self.random_value(value+1)
            response = self.client_1.api.UpdateUserAddress(new_data, label, self.user_1)
            self.assertEqual(response.status_code, 400)
            new_data = dict(data)

        self.lg('[DELETE] Delete address (2), should succeed with 204')
        response = self.client_1.api.GetUserAddresses(self.user_1)
        self.assertEqual(response.status_code, 200)
        label = response.json()[-1]['label']
        response = self.client_1.api.DeleteUserAddress(label, self.user_1)
        self.assertEqual(response.status_code, 204)

        self.lg('[DELETE] Delete address (1), should succeed with 204')
        response = self.client_1.api.GetUserAddresses(self.user_1)
        self.assertEqual(response.status_code, 200)
        label = response.json()[-1]['label']
        response = self.client_1.api.DeleteUserAddress(label, self.user_1)
        self.assertEqual(response.status_code, 204)

        response = self.client_1.api.GetUserAddresses(self.user_1)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(response.json(), [])

        self.lg('[DELETE] Delete invalid address, should fail with 404')
        response = self.client_1.api.DeleteUserAddress('fake_address', self.user_1)
        self.assertEqual(response.status_code, 404)

        self.lg('%s ENDED' % self._testID)

    @unittest.skip('bug #415')
    def test007_get_post_put_delete_banks(self):
        """
            ITSYOU-030
            - Register a new bank account (1), should succeed with 201
            - Register a new bank account (2), should succeed with 201
            - Get user\'s bank accounts, should succeed with 200
            - Get bank account (2) by label, should succeed with 200
            - Register a new bank account with label already exists (label of bank account (2)), should fail with 409
            - Register a new bank account with invalid inputs, should fail with 400
            - Update bank account (2), should succeed with 201
            - Update bank account (2) label with label already exists (label of bank account (1)), should fail with 409
            - Update bank account with invalid inputs, should fail with 400
            - Delete bank account (2), should succeed with 204
            - Delete bank account (1), should succeed with 204
            - Delete invalid bank account, should fail with 404
        """

        elements = {"bic":11,"country":40, "iban":30}

        self.lg('%s STARTED' % self._testID)

        self.lg('[POST] Register a new bank account (1), should succeed with 201')
        label = self.random_value()
        data = {"bic": self.random_value(8),
                "country": self.random_value(),
                "iban": self.random_value(),
                "label": label}
        response = self.client_1.api.CreateUserBankAccount(data, self.user_1)
        self.assertEqual(response.status_code, 201)
        self.assertEqual(response.json(), data)

        self.lg('[POST] Register a new bank account (2), should succeed with 201')
        new_label = self.random_value()
        new_data = {"bic": self.random_value(11),
                    "country": self.random_value(),
                    "iban": self.random_value(),
                    "label": new_label}
        response = self.client_1.api.CreateUserBankAccount(new_data, self.user_1)
        self.assertEqual(response.status_code, 201)
        self.assertEqual(response.json(), new_data)

        self.lg('[GET] Get user bank accounts, should succeed with 200')
        response = self.client_1.api.GetUserBankAccounts(self.user_1)
        self.assertEqual(response.status_code, 200)
        self.assertIn(data, response.json())
        self.assertIn(new_data, response.json())

        self.lg('[GET] Get bank account (2) by label, should succeed with 200')
        label = response.json()[-1]['label']
        response = self.client_1.api.GetUserBankAccountByLabel(label, self.user_1)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(new_data, response.json())

        self.lg('[POST] Register a new bank account with label already exists (label of bank account (2)), should fail with 409')
        response = self.client_1.api.GetUserBankAccounts(self.user_1)
        self.assertEqual(response.status_code, 200)
        label = response.json()[-1]['label']
        new_data = {"bic": self.random_value(8),
                    "country": self.random_value(),
                    "iban": self.random_value(),
                    "label": label}

        response = self.client_1.api.CreateUserBankAccount(new_data, self.user_1)
        self.assertEqual(response.status_code, 409)

        self.lg('[POST] Register a new bank account with invalid inputs, should fail with 400')
        new_data = dict(data)
        for key, value in elements.items():
            new_data['label'] = self.random_value()
            new_data[key] = self.random_value(value+1)
            response = self.client_1.api.CreateUserBankAccount(new_data, self.user_1)
            self.assertEqual(response.status_code, 400)
            new_data = dict(data)

        self.lg('[PUT] Edit bank account (2), should succeed with 201')
        response = self.client_1.api.GetUserBankAccounts(self.user_1)
        self.assertEqual(response.status_code, 200)
        label = response.json()[-1]['label']
        new_data = {"bic": self.random_value(8),
                    "country": self.random_value(),
                    "iban": self.random_value(),
                    "label": self.random_value()}
        response = self.client_1.api.UpdateUserBankAccount(new_data, label, self.user_1)
        self.assertEqual(response.status_code, 200) ### response section is empty

        self.lg('[PUT] Edit bank account (2) label with label already exists (label of bank account (1)), should fail with 409')
        response = self.client_1.api.GetUserBankAccounts(self.user_1)
        self.assertEqual(response.status_code, 200)
        label = response.json()[-1]['label']
        new_label = response.json()[-2]['label']
        new_data = {"bic": self.random_value(8),
                    "country": self.random_value(),
                    "iban": self.random_value(),
                    "label": new_label}
        response = self.client_1.api.UpdateUserBankAccount(new_data, label, self.user_1)
        self.assertEqual(response.status_code, 409)

        self.lg('[PUT] Edit bank account with invalid inputs, should fail with 400')
        new_data = dict(data)
        for key, value in elements.items():
            new_data[key] = self.random_value(value+1)
            response = self.client_1.api.UpdateUserBankAccount(new_data, label, self.user_1)
            self.assertEqual(response.status_code, 400)
            new_data = dict(data)

        self.lg('[DELETE] Delete bank account (2), should succeed with 204')
        response = self.client_1.api.GetUserBankAccounts(self.user_1)
        self.assertEqual(response.status_code, 200)
        label = response.json()[-1]['label']
        response = self.client_1.api.DeleteUserBankAccount(label, self.user_1)
        self.assertEqual(response.status_code, 204)

        self.lg('[DELETE] Delete bank account (1), should succeed with 204')
        response = self.client_1.api.GetUserBankAccounts(self.user_1)
        self.assertEqual(response.status_code, 200)
        label = response.json()[-1]['label']
        response = self.client_1.api.DeleteUserBankAccount(label, self.user_1)
        self.assertEqual(response.status_code, 204)

        response = self.client_1.api.GetUserBankAccounts(self.user_1)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(response.json(), [])

        self.lg('[DELETE] Delete invalid bank account, should fail with 404')
        response = self.client_1.api.DeleteUserBankAccount('fake_bank_account', self.user_1)
        self.assertEqual(response.status_code, 404)

        self.lg('%s ENDED' % self._testID)

    def test008_get_post_put_delete_publickey(self):
        """
            ITSYOU-031
            - Register a new publickey (1), should succeed with 201
            - Register a new publickey (2), should succeed with 201
            - Get user\'s publickeys, should succeed with 200
            - Get publickey (2) by label, should succeed with 200
            - Register a new publickey with label already exists (label of publickey (2)), should fail with 409
            - Register a new publickey with invalid inputs, should fail with 400
            - Update publickey (2), should succeed with 201
            - Update publickey (2) with label already exists (label of publickey (1)), should fail with 409
            - Update publickey with invalid inputs, should fail with 400
            - Delete publickey (2), should succeed with 204
            - Delete publickey (1), should succeed with 204
            - Delete invalid publickey, should fail with 404
        """

        self.lg('%s STARTED' % self._testID)

        self.lg('[POST] Register a new publickey (1), should succeed with 201')
        label = self.random_value()
        publickey = 'ssh-rsa AAAAB3NzaC1yc2E'+self.random_value(30)
        data = {"label": label,"publickey": publickey}
        response = self.client_1.api.RegisterUserPublicKey(data, self.user_1)
        self.assertEqual(response.status_code, 201)
        self.assertEqual(response.json(), data)

        self.lg('[POST] Register a new publickey (2), should succeed with 201')
        new_label = self.random_value()
        publickey = 'ssh-rsa AAAAB3NzaC1yc2E'+self.random_value(30)
        new_data = {"label": new_label,"publickey": publickey}
        response = self.client_1.api.RegisterUserPublicKey(new_data, self.user_1)
        self.assertEqual(response.status_code, 201)
        self.assertEqual(response.json(), new_data)

        self.lg('[GET] Get user publickeys, should succeed with 200')
        response = self.client_1.api.GetUserPublicKeys(self.user_1)
        self.assertEqual(response.status_code, 200)
        self.assertIn(data, response.json())
        self.assertIn(new_data, response.json())

        self.lg('[GET] Get publickey (2) by label, should succeed with 200')
        label = response.json()[-1]['label']
        response = self.client_1.api.GetUserPublicKeyByLabel(label, self.user_1)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(new_data, response.json())

        self.lg('[POST] Register a new publickey with label already exists (label of publickey (2)), should fail with 409')
        response = self.client_1.api.GetUserPublicKeys(self.user_1)
        self.assertEqual(response.status_code, 200)
        new_label = response.json()[-1]['label']
        new_publickey = 'ssh-rsa AAAAB3NzaC1yc2E'+self.random_value(30)
        new_data = {"label": new_label,"publickey": new_publickey}
        response = self.client_1.api.RegisterUserPublicKey(new_data, self.user_1)
        self.assertEqual(response.status_code, 409)

        self.lg('[POST] Register a new publickey with invalid inputs, should fail with 400')
        new_label = self.random_value()
        new_publickey = self.random_value(30)
        new_data = {"label": new_label,"publickey": new_publickey}
        response = self.client_1.api.RegisterUserPublicKey(new_data, self.user_1)
        self.assertEqual(response.status_code, 400)

        self.lg('[PUT] Update publickey (2), should succeed with 201')
        response = self.client_1.api.GetUserPublicKeys(self.user_1)
        self.assertEqual(response.status_code, 200)
        label = response.json()[-1]['label']
        new_label = self.random_value()
        new_publickey = 'ssh-rsa AAAAB3NzaC1yc2E'+self.random_value(30)
        new_data = {"label": new_label,"publickey": new_publickey}
        response = self.client_1.api.UpdateUserPublicKey(new_data, label, self.user_1)
        self.assertEqual(response.status_code, 201)

        self.lg('[PUT] Update publickey (2) with label already exists (label of publickey (1)), should fail with 409')
        response = self.client_1.api.GetUserPublicKeys(self.user_1)
        self.assertEqual(response.status_code, 200)
        label = response.json()[-1]['label']
        new_label = response.json()[-2]['label']
        new_publickey = 'ssh-rsa AAAAB3NzaC1yc2E'+self.random_value(30)
        new_data = {"label": new_label,"publickey": new_publickey}
        response = self.client_1.api.UpdateUserPublicKey(new_data, label, self.user_1)
        self.assertEqual(response.status_code, 409)

        self.lg('[PUT] Update publickey with invalid inputs, should fail with 400')
        response = self.client_1.api.GetUserPublicKeys(self.user_1)
        self.assertEqual(response.status_code, 200)
        label = response.json()[-1]['label']
        new_label = self.random_value()
        new_publickey = self.random_value(30)
        new_data = {"label": new_label,"publickey": new_publickey}
        response = self.client_1.api.UpdateUserPublicKey(new_data, label, self.user_1)
        self.assertEqual(response.status_code, 400)

        self.lg('[DELETE] Delete publickey (2), should succeed with 204')
        response = self.client_1.api.GetUserPublicKeys(self.user_1)
        self.assertEqual(response.status_code, 200)
        label = response.json()[-1]['label']
        response = self.client_1.api.DeleteUserPublicKey(label, self.user_1)
        self.assertEqual(response.status_code, 204)

        self.lg('[DELETE] Delete publickey (1), should succeed with 204')
        response = self.client_1.api.GetUserPublicKeys(self.user_1)
        self.assertEqual(response.status_code, 200)
        label = response.json()[-1]['label']
        response = self.client_1.api.DeleteUserPublicKey(label, self.user_1)
        self.assertEqual(response.status_code, 204)

        response = self.client_1.api.GetUserPublicKeys(self.user_1)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(response.json(), [])

        self.lg('[DELETE] Delete invalid publickey, should fail with 404')
        response = self.client_1.api.DeleteUserPublicKey('fake_publickey', self.user_1)
        self.assertEqual(response.status_code, 404)

        self.lg('%s ENDED' % self._testID)

    @unittest.skip('bug: #402 #403 #404 #405')
    def test009_get_post_put_delete_apikeys(self):

        """
        ITSYOU-032
        - Register a new apikey (1), should succeed with 201
        - Register a new apikey (2), should succeed with 201
        - Get user\'s apikeys, should succeed with 200
        - Get apikey (2) by label, should succeed with 200
        - Register a new apikey with label already exists (label of apikey (2)), should fail with 409
        - Update the apikey (2), should succeed with 201
        - Update the apikey (2) with label already exists (label of apikey (1)), should fail with 409
        - Delete apikey (2), should succeed with 204
        - Delete apikey (1), should succeed with 204
        - Delete apikey with fake label, should fail with 404
        """

        self.lg('%s STARTED' % self._testID)

        self.lg('[POST] Register a new apikey (1), should succeed with 201')
        label = self.random_value()
        data = {'label' : label}
        response = self.client_1.api.AddApiKey(data, self.user_1)
        self.assertEqual(response.status_code, 201)
        self.assertIn('apikey', response.json().keys())
        self.assertIn('applicationid', response.json().keys())
        self.assertIn('label', response.json().keys())
        self.assertIn('scopes', response.json().keys())
        self.assertIn('username', response.json().keys())
        self.assertEqual(response.json()['label'], label)


        self.lg('[POST] Register a new apikey (2), should succeed with 201')
        new_label = self.random_value()
        new_data = {'label' : new_label}
        response = self.client_1.api.AddApiKey(new_data, self.user_1)
        self.assertEqual(response.status_code, 201)
        self.assertEqual(response.json()['label'], new_label)

        apikey = response.json()['apikey']
        applicationid = response.json()['applicationid']
        scopes = response.json()['scopes']
        username = response.json()['username']

        self.lg('[GET] Get user\'s apikeys, should succeed with 200')
        response = self.client_1.api.ListAPIKeys(self.user_1)
        self.assertEqual(response.status_code, 200)
        self.assertIn(label, [x['label'] for x in response.json()])
        self.assertIn(new_label, [x['label'] for x in response.json()])

        self.lg('[GET] Get apikey (2) by label, should succeed with 200')
        label = response.json()[-1]['label']
        response = self.client_1.api.GetAPIkey(label, self.user_1)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(label, response.json()['label'])
        self.assertEqual(apikey, response.json()['apikey'])
        self.assertEqual(applicationid, response.json()['applicationid'])
        self.assertEqual(scopes, response.json()['scopes'])
        self.assertEqual(username, response.json()['username'])

        #bug #402
        self.lg('[POST] Register a new apikey with label already exists (label of apikey (2)), should fail with 409')
        response = self.client_1.api.ListAPIKeys(self.user_1)
        self.assertEqual(response.status_code, 200)
        new_label = response.json()[-1]['label']
        new_data = {'label' : new_label}
        response = self.client_1.api.AddApiKey(new_data, self.user_1)
        self.assertEqual(response.status_code, 409)

        #bug #404
        self.lg('[PUT] Update the apikey (2), should succeed with 201')
        response = self.client_1.api.ListAPIKeys(self.user_1)
        self.assertEqual(response.status_code, 200)
        label = response.json()[-1]['label']
        new_label = self.random_value()
        new_data = {'label': new_label}
        response = self.client_1.api.UpdateAPIkey(new_data, label, self.user_1)
        self.assertEqual(response.status_code, 201)

        #bug #403
        self.lg('[PUT] Update the apikey (2) with label already exists (label of apikey (1)), should fail with 409')
        response = self.client_1.api.ListAPIKeys(self.user_1)
        self.assertEqual(response.status_code, 200)
        label = response.json()[-1]['label']
        new_label = response.json()[-2]['label']
        data = {'label': new_label}
        response = self.client_1.api.UpdateAPIkey(data, label, self.user_1)
        self.assertEqual(response.status_code, 409)

        #bug #405
        self.lg('[DELETE] Delete apikey (2), should succeed with 204')
        response = self.client_1.api.ListAPIKeys(self.user_1)
        self.assertEqual(response.status_code, 200)
        label = response.json()[-1]['label']
        response = self.client_1.api.DeleteAPIkey(label, self.user_1)
        self.assertEqual(response.status_code, 204)

        self.lg('[DELETE] Delete apikey (1), should succeed with 204')
        response = self.client_1.api.ListAPIKeys(self.user_1)
        self.assertEqual(response.status_code, 200)
        label = response.json()[-1]['label']
        response = self.client_1.api.DeleteAPIkey(label, self.user_1)
        self.assertEqual(response.status_code, 204)

        self.lg('[DELETE] Delete apikey with fake label, should fail with 404')
        response = self.client_1.api.DeleteAPIkey('fake_label', self.user_1)
        self.assertEqual(response.status_code, 404)

        self.lg('%s ENDED' % self._testID)

    @unittest.skip('bug: #407 #424')
    def test010_get_post_put_delete_digitalwallet(self):
        """
            ITSYOU-033
            - Register a new digital wallet (1), should succeed with 201
            - Register a new digital wallet (2), should succeed with 201
            - Get user\'s digital wallets, should succeed with 200
            - Get digital wallet (2) by label, should succeed with 200
            - Register a new digital wallet with label already exists (label of digital wallet (2)), should fail with 409
            - Register a new digital wallet with invalid inputs, should fail with 400
            - Update digital wallet (2), should succeed with 201
            - Update digital wallet (2) with label already exists (label of digital wallet (1)), should fail with 409
            - Update digital wallet with invalid inputs, should fail with 400
            - Delete digital wallet (2), should succeed with 204
            - Delete digital wallet (1), should succeed with 204
            - Delete invalid digital wallet, should fail with 404
        """

        self.lg('%s STARTED' % self._testID)

        #bug #424
        self.lg('[POST] Register a new digital wallet (1), should succeed with 201')
        label = self.random_value()
        data = {"currencysymbol":self.random_value(),
                "address":self.random_value(),
                "label":label,
                "expire":"2018-01-16T15:35:14.507Z",
                "noexpiration":False}
        response = self.client_1.api.RegisterDigitalWallet(data, self.user_1)
        self.assertEqual(response.status_code, 201)
        self.assertEqual(response.json(), data)

        self.lg('[POST] Register a new digital wallet (2), should succeed with 201')
        new_label = self.random_value()
        new_data =  {"currencysymbol": self.random_value(),
                    "address": self.random_value(),
                    "label": new_label,
                    "noexpiration": True}
        response = self.client_1.api.RegisterDigitalWallet(new_data, self.user_1)
        self.assertEqual(response.status_code, 201)
        self.assertEqual(response.json(), new_data)

        self.lg('[GET] Get user\'s digital wallets, should succeed with 200')
        response = self.client_1.api.GetUserDigitalWallets(self.user_1)
        self.assertEqual(response.status_code, 200)
        self.assertIn(data, response.json())
        self.assertIn(new_data, response.json())

        self.lg('[GET] Get digital wallet (2) by label, should succeed with 200')
        response = self.client_1.api.GetUserDigitalWalletByLabel(new_label, self.user_1)
        self.assertEqual(response.status_code, 200)
        self.assertIn(new_data, response.json())

        self.lg('[POST] Register a new digital wallet with label already exists (label of digital wallet (2)), should fail with 409')
        data['label'] = label
        response = self.client_1.api.RegisterDigitalWallet(data, self.user_1)
        self.assertEqual(response.status_code, 409)

        self.lg('[POST] Register a new digital wallet with invalid inputs, should fail with 400')
        label = self.random_value()
        data['label'] = label
        data['expire'] = self.random_value()
        response = self.client_1.api.RegisterDigitalWallet(data, self.user_1)
        self.assertEqual(response.status_code, 400)

        self.lg('[PUT] Update digital wallet (2), should succeed with 201')
        response = self.client_1.api.GetUserDigitalWallets(self.user_1)
        self.assertEqual(response.status_code, 200)
        label = response.json()[-1]['label']
        new_label = self.random_value()
        new_data = {"currencysymbol": self.random_value(),
                    "address": self.random_value(),
                    "label": new_label,
                    "expire": "2018-01-16T15:35:14.507Z",
                    "noexpiration": False}
        response = self.client_1.api.UpdateUserDigitalWallet(new_data, label, self.user_1)
        self.assertEqual(response.status_code, 201)

        # bug #407
        self.lg('[PUT] Update digital wallet (2) with label already exists (label of digital wallet (1)), should fail with 409')
        response = self.client_1.api.GetUserDigitalWallets(self.user_1)
        self.assertEqual(response.status_code, 200)
        label = response.json()[-1]['label']
        new_label = response.json()[-2]['label']
        new_data = {"currencysymbol": self.random_value(),
                    "address": self.random_value(),
                    "label": new_label,
                    "expire": "2018-01-16T15:35:14.507Z",
                    "noexpiration": False}
        response = self.client_1.api.UpdateUserDigitalWallet(new_data, label, self.user_1)
        self.assertEqual(response.status_code, 409)

        self.lg('[PUT] Update digital wallet with invalid inputs, should fail with 400')
        response = self.client_1.api.GetUserDigitalWallets(self.user_1)
        self.assertEqual(response.status_code, 200)
        label = response.json()[randint(0, len(response.json())-1)]['label']
        new_label = self.random_value()
        new_data = {"currencysymbol": self.random_value(),
                    "address": self.random_value(),
                    "label": new_label,
                    "expire": self.random_value(),
                    "noexpiration": False}
        response = self.client_1.api.UpdateUserDigitalWallet(new_data, label, self.user_1)
        self.assertEqual(response.status_code, 400)

        self.lg('[DELETE] Delete digital wallet (2), should succeed with 204')
        response = self.client_1.api.GetUserDigitalWallets(self.user_1)
        self.assertEqual(response.status_code, 200)
        label = response.json()[-1]['label']
        response = self.client_1.api.DeleteUserDigitalWallet(label, self.user_1)
        self.assertEqual(response.status_code, 204)

        self.lg('[DELETE] Delete digital wallet (1), should succeed with 204')
        response = self.client_1.api.GetUserDigitalWallets(self.user_1)
        self.assertEqual(response.status_code, 200)
        label = response.json()[-1]['label']
        response = self.client_1.api.DeleteUserDigitalWallet(label, self.user_1)
        self.assertEqual(response.status_code, 204)

        response = self.client_1.api.GetUserDigitalWallets(self.user_1)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(response.json(), [])

        self.lg('[DELETE] Delete invalid digital wallet, should fail with 404')
        response = self.client_1.api.DeleteUserDigitalWallet('fake_digital_wallet', self.user_1)
        self.assertEqual(response.status_code, 404)

        self.lg('%s ENDED' % self._testID)

    @unittest.skip('bug: #414')
    def test011_get_post_delete_organizations_auth(self):

        """
            ITSYOU-034
            * Same steps for owner and member roles
            - Create organization org_1 and make user_1 an owner
            - Create organization org_2 and make user_1 and user_2 owners
            - User_1 send invitation to user_2 to join org_1, should succeed with 201
            - User_2 reject the invitation, should succeed with 204
            - User_1 send invitation to user_2 to join org_1 again, should succeed with 201
            - User_2 accept the invitation, should succeed with 201
            - Modify certain information for user_2 to be granted to org_1, should succeed with 201
            - Get all Authorizations of user_2, should succeed with 200
            - Get the Authorizations of user_2 for specific organization (org_1), should succeed with 200
            - User_2 remove the authorization for the org_1, should succeed with 204
            - User_2 leave organization org_1, should succeed with 204
            - Unothorized user (user_1) try to make user_2 leave org_2, should fail with 403
            - User_2 leave fake organization, should fail with 404
            - Fake user leave org_1, should fail with 404
        """
        self.lg('%s STARTED' % self._testID)

        response = self.client_1.api.GetUserOrganizations(self.user_1)
        self.assertEqual(response.status_code, 200)
        self.assertIn(self.organization_1, response.json()['owner'])

        for role in ['member', 'owner']:

            self.lg('[POST] User_1 send invitation to user_2 to join org_1 role %s, should succeed with 201' % role)
            data = {'searchstring': self.user_2}
            if role =="member":
                response = self.client_1.api.AddOrganizationMember(data, self.organization_1)
            elif role == "owner":
                response = self.client_1.api.AddOrganizationOwner(data, self.organization_1)
            self.assertEqual(response.status_code, 201)

            response = self.client_2.api.GetNotifications(self.user_2)
            self.assertEqual(response.status_code, 200)
            self.assertEqual(self.organization_1, response.json()['invitations'][-1]['organization'])
            self.assertEqual(role, response.json()['invitations'][-1]['role'])
            self.assertEqual('pending', response.json()['invitations'][-1]['status'])

            self.lg('[DELETE] User_2 reject the invitation, should succeed with 204')
            response = self.client_2.api.RejectMembership(self.organization_1, role, self.user_2)
            self.assertEqual(response.status_code, 204)
            response = self.client_2.api.GetUserOrganizations(self.user_2)
            self.assertNotIn(self.organization_1, response.json()[role])

            self.lg('[POST] User_1 send invitation to user_2 to join org_1 again role %s, should succeed with 201' % role)
            data = {'searchstring': self.user_2}
            if role =="member":
                response = self.client_1.api.AddOrganizationMember(data, self.organization_1)
            elif role == "owner":
                response = self.client_1.api.AddOrganizationOwner(data, self.organization_1)
            self.assertEqual(response.status_code, 201)

            response = self.client_2.api.GetNotifications(self.user_2)
            self.assertEqual(response.status_code, 200)
            self.assertEqual(self.organization_1, response.json()['invitations'][-1]['organization'])
            self.assertEqual(role, response.json()['invitations'][-1]['role'])
            self.assertEqual('pending', response.json()['invitations'][-1]['status'])

            self.lg('[POST] User_2 accept the invitation, should succeed with 201')
            response = self.client_2.api.AcceptMembership(self.organization_1, role, self.user_2)
            self.assertEqual(response.status_code, 201)
            response = self.client_2.api.GetUserOrganizations(self.user_2)
            self.assertIn(self.organization_1, response.json()[role])

            self.lg('[PUT] Modify certain information for user_2 to be granted to org_1, should succeed with 201')
            response = self.client_2.api.GetEmailAddresses(self.user_2)
            self.assertEqual(response.status_code, 200)
            label = response.json()[-1]['label']

            grantedto = self.organization_1
            data = {"username":self.user_2,
                    "grantedTo":grantedto,
                    "emailaddresses":[{"requestedlabel": self.random_value(), "reallabel": label}]}

            response = self.client_2.api.UpdateAuthorization(data, grantedto, self.user_2)
            self.assertEqual(response.status_code, 201)

            self.lg('[GET] Get all Authorizations of user_2, should succeed with 200')
            response = self.client_2.api.GetAllAuthorizations(self.user_2)
            self.assertEqual(response.status_code, 200)
            self.assertEqual(self.user_2, response.json()[-1]['username'])
            self.assertEqual(grantedto, response.json()[-1]['grantedTo'])
            self.assertEqual(data['emailaddresses'][-1]['reallabel'], response.json()[-1]['emailaddresses'][-1]['reallabel'])
            self.assertEqual(data['emailaddresses'][-1]['requestedlabel'], response.json()[-1]['emailaddresses'][-1]['requestedlabel'])

            self.lg('[GET] Get the Authorizations of user_2 for specific organization (org_1), should succeed with 200')
            response = self.client_2.api.GetAuthorization(grantedto, self.user_2)
            self.assertEqual(response.status_code, 200)
            self.assertEqual(self.user_2, response.json()['username'])
            self.assertEqual(grantedto, response.json()['grantedTo'])
            self.assertEqual(data['emailaddresses'][-1]['reallabel'], response.json()['emailaddresses'][-1]['reallabel'])
            self.assertEqual(data['emailaddresses'][-1]['requestedlabel'], response.json()['emailaddresses'][-1]['requestedlabel'])

            self.lg('[DELETE] User_2 remove the authorization for the org_1, should succeed with 204')
            response = self.client_2.api.DeleteAuthorization(grantedto, self.user_2)
            self.assertEqual(response.status_code, 204)

            response = self.client_2.api.GetAuthorization(grantedto, self.user_2)
            self.assertEqual(response.status_code, 404)

            response = self.client_2.api.GetAllAuthorizations(self.user_2)
            self.assertEqual(response.status_code, 200)
            self.assertEqual(response.json(), [])

            self.lg('[DELETE] User_2 leave organization org_1, should succeed with 204')
            response = self.client_2.api.LeaveOrganization(self.organization_1, self.user_2)
            self.assertEqual(response.status_code, 204)
            response = self.client_2.api.GetUserOrganizations(self.user_2)
            self.assertNotIn(self.organization_1, response.json()[role])

            self.lg('[DELETE] Unothorized user (user_1) try to make user_2 leave org_2, should fail with 403')
            response = self.client_1.api.LeaveOrganization(self.organization_2, self.user_2)
            self.assertEqual(response.status_code, 403)

            self.lg('[DELETE] User_2 leave fake organization, should fail with 404')
            response = self.client_2.api.LeaveOrganization('fake_organization', self.user_2)
            self.assertEqual(response.status_code, 404)

            # #bug 414
            self.lg('[DELETE] Fake user leave org_1, should fail with 404')
            response = self.client_2.api.LeaveOrganization(self.organization_1, 'fake_user')
            self.assertEqual(response.status_code, 404)

            self.lg('%s ENDED' % self._testID)
    
    @unittest.skip(' ')
    def test012_get_post_delete_totp_twofamethods(self):
        """
            ITSYOU-035
            - Get totp secret, should succeed with 200
            - Set totp with invalid secret, should fail with 422
            - Set totp with valid secret and invalid code, should fail with 422
            - Set totp with valid code and valid secret, should succeed with 204
            - Delete totp, should fail with 409
            - Register a new phonenumber, send validation sms & verify it, should succeed with 201 & 200
            - Delete totp, should succeed with 204
            - Set totp code again, should succeed with 204
            - Force Delete the verified phonenumber, should succeed with 204
        """

        self.lg('%s STARTED' % self._testID)

        self.lg('[GET] totp secret, should succeed with 200')
        response = self.client_1.api.GetTotp(self.user_1)
        self.assertEqual(response.status_code, 200)
        secret = response.json()['totpsecret']

        self.lg('[POST] Set totp with invalid secret, should fail with 422')
        new_secret = self.random_value()
        totpcode = self.random_value()
        data = {"totpcode":totpcode, "totpsecret":new_secret}
        response = self.client_1.api.EditTotp(data, self.user_1)
        self.assertEqual(response.status_code, 422)

        self.lg('[POST] Set totp with valid secret and invalid code, should fail with 422')
        totpcode = self.random_value()
        data = {"totpcode":totpcode, "totpsecret":secret}
        response = self.client_1.api.EditTotp(data, self.user_1)
        self.assertEqual(response.status_code, 422)

        self.lg('[POST] Set totp with valid code and valid secret, should succeed with 204')
        totpcode = self.get_totp_code(secret)
        data = {"totpcode":totpcode, "totpsecret":secret}
        response = self.client_1.api.EditTotp(data, self.user_1)
        self.assertEqual(response.status_code, 204)

        response = self.client_1.api.GetTwofamethods(self.user_1)
        self.assertEqual(response.status_code, 200)
        self.assertTrue(response.json()['totp'])

        self.lg('[DEL] Delete totp, should fail with 409')
        response = self.client_1.api.DeleteTotp(self.user_1)
        self.assertEqual(response.status_code, 409)

        self.lg('[POST] Register a new phonenumber & send validation sms & verify it, should succeed with 201 & 200')
        label = 'validation number'
        phonenumber = self.get_valid_phonenumber()
        data = {"label":label, "phonenumber":phonenumber}
        response = self.client_1.api.RegisterNewUserPhonenumber(data, self.user_1)
        self.assertEqual(response.status_code, 201)

        response = self.client_1.api.ValidatePhonenumber(label, self.user_1)
        self.assertEqual(response.status_code, 200)
        validationkey = response.json()['validationkey']

        time.sleep(25)

        smscode = self.get_mobile_verification_code()
        self.assertTrue(smscode, 'error while getting sms code, verification message not received with virtual number %s' % phonenumber)
        data = {"smscode":smscode, "validationkey":validationkey}
        response = self.client_1.api.VerifyPhoneNumber(data, label, self.user_1)
        self.assertEqual(response.status_code, 204)

        response = self.client_1.api.GetUserPhoneNumbers(self.user_1, query_params={'validated':True})
        self.assertEqual(response.status_code, 200)
        self.assertIn(label, [x['label'] for x in response.json()])

        self.lg('[DEL] Delete totp, should succeed with 204')
        response = self.client_1.api.DeleteTotp(self.user_1)
        self.assertEqual(response.status_code, 204)

        response = self.client_1.api.GetTwofamethods(self.user_1)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(response.json()['sms'][-1]['label'], 'validation number')
        self.assertEqual(response.json()['sms'][-1]['phonenumber'], phonenumber)

        self.lg('[POST] Set totp code again, should succeed with 204')
        secret = self.totp_secret
        totpcode = self.get_totp_code(secret)
        data = {"totpcode":totpcode, "totpsecret":secret}
        response = self.client_1.api.EditTotp(data, self.user_1)
        self.assertEqual(response.status_code, 204)

        response = self.client_1.api.GetTwofamethods(self.user_1)
        self.assertEqual(response.status_code, 200)
        self.assertTrue(response.json()['totp'])

        self.lg('[DELETE] Force Delete verified phonenumber, should succeed with 204')
        response = self.client_1.api.DeleteUserPhonenumber(label, self.user_1, query_params={'force':'true'})
        self.assertEqual(response.status_code, 204)

        response = self.client_1.api.GetTwofamethods(self.user_1)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(response.json()['sms'], None)

        self.lg('%s ENDED' % self._testID)

    def test013_delete_facebook_account(self):

        """
            ITSYOU-036
            - Check if facebook account exists or not, should succeed with 200
            - Delete facebook account if exists, should succeed with 204
            - Check if the facebook account is deleted, should succeed with 200
        """

        self.lg('%s STARTED' % self._testID)

        self.lg('Check if facebook account exists or not, should succeed')
        response = self.client_1.api.GetUser(self.user_1)
        self.assertEqual(response.status_code, 200)
        self.assertIn('facebook', response.json().keys())
        self.assertEqual(type(response.json()['facebook']), types.DictType)
        empty_account = {'id':'', 'link':'', 'name':'', 'picture':''}

        if response.json()['facebook'] != empty_account:
            self.lg('Delete facebook account, should succeed')
            response = self.client_1.api.DeleteFacebookAccount(self.user_1)
            self.assertEqual(response.status_code, 204)
            self.lg('Check if the facebook account is deleted, should succeed')
            response = self.client_1.api.GetUser(self.user_1)
            self.assertEqual(response.status_code, 200)
            self.assertEqual(response.json()['facebook'], empty_account)
        else:
            self.lg('facebook account is already deleted')

        self.lg('%s ENDED' % self._testID)

    def test014_delete_github_account(self):
        """
            ITSYOU-037
            - Check if github account exists or not, should succeed with 200
            - Delete github account if exists, should succeed with 204
            - Check if the github account is deleted, should succeed with 200
        """

        self.lg('%s STARTED' % self._testID)

        self.lg('Check if github account exists, should succeed')
        response = self.client_1.api.GetUser(self.user_1)
        self.assertEqual(response.status_code, 200)
        self.assertIn('github', response.json().keys())
        self.assertEqual(type(response.json()['github']), types.DictType)
        empty_account = {u'avatar_url': u'', u'html_url': u'', u'id': 0, u'login': u'', u'name': u''}

        if response.json()['github'] != empty_account:
            self.lg('Delete github account, should succeed')
            response = self.client_1.api.DeleteGithubAccount(username=self.user_1)
            self.assertEqual(response.status_code, 204)
            self.lg('Check if the github account is deleted, should succeed')
            response = self.client_1.api.GetUser(self.user_1)
            self.assertEqual(response.status_code, 200)
            self.assertIn('github', response.json().keys())
            self.assertEqual(type(response.json()['github']), types.DictType)
            self.assertEqual(response.json()['github'], empty_account)
        else:
            self.lg('github account is already deleted')

        self.lg('%s ENDED' % self._testID)

    @unittest.skip("bug: #458")
    def test015_get_post_contract(self):
        """
            ITSYOU-038
            - Create a new contract (1), should succeed with 201
            - Create a new contract (2), should succeed with 201
            - Create a new expired contract (3), should succeed with 201
            - Get user\'s contracts, should succeed with 200
            - Get user\'s contracts & include the expired contracts, should succeed with 200
            - Get user\'s contracts with page size 1, should succeed with 200
            - Get user\'s contracts with start page 2, should succeed with 200
        """

        self.lg('%s STARTED' % self._testID)

        self.lg('Create a new contract (1), should succeed with 201')
        contractid_1 = self.random_value()
        expire = '2030-10-02T22:00:00Z'
        data = {'content':'contract_1', 'contractId':contractid_1, 'contractType':'partnership','expires':expire}
        response = self.client_1.api.CreateUserContract(data, self.user_1)
        self.assertEqual(response.status_code, 201)

        self.lg('Create a new contract (2), should succeed with 201')
        contractid_2 = self.random_value()
        expire = '2030-10-02T22:00:00Z'
        data = {'content':'contract_2', 'contractId':contractid_2, 'contractType':'partnership','expires':expire}
        response = self.client_1.api.CreateUserContract(data, self.user_1)
        self.assertEqual(response.status_code, 201)
        #bug #458
        self.lg('Create a new expired contract (3), should succeed with 201')
        contractid_3 = self.random_value()
        expire = '2010-10-02T22:00:00Z'
        data = {'content':'contract_3', 'contractId':contractid_3, 'contractType':'partnership','expires':expire}
        response = self.client_1.api.CreateUserContract(data, self.user_1)
        self.assertEqual(response.status_code, 201)

        self.lg('[GET] Get user\'s contracts, should succeed with 200')
        response = self.client_1.api.GetUserContracts(self.user_1)
        self.assertEqual(response.status_code, 200)
        self.assertNotEqual(contractid_3, response.json()[-1]['contractId'])
        self.assertEqual(contractid_1, response.json()[-2]['contractId'])
        self.assertEqual(contractid_2, response.json()[-1]['contractId'])

        response = self.client_1.api.GetUserContracts(self.user_1, query_params={"max":1000,"includeExpired":True})
        self.assertEqual(response.status_code, 200)
        number_of_contracts = len(response.json())-1

        self.lg('[GET] Get user\'s contracts & include the expired contracts, should succeed with 200')
        response = self.client_1.api.GetUserContracts(self.user_1, query_params={"max":1000,"includeExpired":True})
        self.assertEqual(response.status_code, 200)
        self.assertEqual(contractid_3, response.json()[-1]['contractId'])

        self.lg('[GET] Get user\'s contracts with page size 1, should succeed with 200')
        response = self.client_1.api.GetUserContracts(self.user_1, query_params={"max":1, "start":number_of_contracts})
        self.assertEqual(contractid_2, response.json()[0]['contractId'])
        self.assertEqual(response.status_code, 200)

        self.lg('[GET] Get user\'s contracts with page size 2, should succeed with 200')
        response = self.client_1.api.GetUserContracts(self.user_1, query_params={"max":1, "start":number_of_contracts-1})
        self.assertEqual(response.status_code, 200)
        self.assertEqual(contractid_1, response.json()[0]['contractId'])


        self.lg('%s ENDED' % self._testID)

    @unittest.skip('bug: #413')
    def test016_get_post_delete_registry(self):

        """
            ITSYOU-039
            - Register a new registry (1), should succeed with 201
            - Register a new registry (2), should succeed with 201.
            - Get user'\s registries, should succeed with 200
            - Get registry (2) by key, should succeed with 200
            - Get invalid registry, should fail with 404
            - Register a new registry with key already exists (key of registry (1)), should fail with 409
            - Register a new registry with invalid inputs, should fail with 400
            - Delete registry (2), should succeed with 204
            - Delete registry (1), should succeed with 204
            - Delete invalid registry, should fail with 404
        """

        self.lg('%s STARTED' % self._testID)

        #bug #413
        self.lg('[POST] Register a new registry (1), should succeed with 201')
        key = self.random_value()
        value = self.random_value()
        data = {"Key": key,"Value": value}
        response = self.client_1.api.CreateNewRegistry(data, self.user_1)
        self.assertEqual(response.status_code, 201)
        self.assertEqual(response.json(), data)

        self.lg('[POST] Register a new registry (2), should succeed with 201')
        key = self.random_value()
        value = self.random_value()
        new_data = {"Key": key,"Value": value}
        response = self.client_1.api.CreateNewRegistry(new_data, self.user_1)
        self.assertEqual(response.status_code, 201)
        self.assertEqual(response.json(), new_data)

        self.lg('[GET] Get user\'s registries, should succeed with 200')
        response = self.client_1.api.GetRegistries(self.user_1)
        self.assertEqual(response.status_code, 200)
        self.assertIn(data, response.json())
        self.assertIn(new_data, response.json())

        self.lg('[GET] Get registry (2) by key, should succeed with 200')
        key = response.json()[-1]['Key']
        response = self.client_1.api.GetRegistry(key, self.user_1)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(new_data, response.json())

        self.lg('[GET] Get invalid registry, should fail with 404')
        response = self.client_1.api.GetRegistry('fake_key', self.user_1)
        self.assertEqual(response.status_code, 404)

        self.lg('[POST] Register a new registry with key already exists (key of registry (1)), should fail with 409')
        response = self.client_1.api.GetRegistries(self.user_1)
        self.assertEqual(response.status_code, 200)
        key = response.json()[-1]['Key']
        value = self.random_value()
        new_data = {"Key": key,"Value": value}
        response = self.client_1.api.CreateNewRegistry(new_data, self.user_1)
        self.assertEqual(response.status_code, 409)

        self.lg('[POST] Register a new registry with invalid inputs, should fail with 400')
        key = ''
        value = ''
        new_data = {"Key": key,"Value": value}
        response = self.client_1.api.CreateNewRegistry(new_data, self.user_1)
        self.assertEqual(response.status_code, 400)

        self.lg('[DELETE] Delete registry (2), should succeed with 204')
        response = self.client_1.api.GetRegistries(self.user_1)
        self.assertEqual(response.status_code, 200)
        key = response.json()[-1]['Key']
        response = self.client_1.api.DeleteRegistry(key, self.user_1)
        self.assertEqual(response.status_code, 204)

        self.lg('[DELETE] Delete registry (1), should succeed with 204')
        response = self.client_1.api.GetRegistries(self.user_1)
        self.assertEqual(response.status_code, 200)
        key = response.json()[-1]['Key']
        response = self.client_1.api.DeleteRegistry(key, self.user_1)
        self.assertEqual(response.status_code, 204)

        response = self.client_1.api.GetRegistries(self.user_1)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(response.json(), [])

        self.lg('[DELETE] Delete invalid registry, should fail with 404')
        response = self.client_1.api.DeleteRegistry('fake_key', self.user_1)
        self.assertEqual(response.status_code, 404)

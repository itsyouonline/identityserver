import requests
from client_utils import build_query_string


class Client:
    def __init__(self, env_url):
        self.url = env_url + 'api/'
        self.session = requests.Session()

    def CreateUser(self, data,  query_params=None):
        """
        Create a new user
        It is method for POST /users
        """
        uri = self.url + "users"
        uri = uri + build_query_string(query_params)
        return self.session.post(uri, json=data)

    def GetNotifications(self, username,  query_params=None):
        """
        Get the list of notifications, these are pending invitations or approvals
        It is method for GET /users/{username}/notifications
        """
        uri = self.url + "users/"+username+"/notifications"
        uri = uri + build_query_string(query_params)
        return self.session.get(uri)

    def GetUserOrganizations(self, username,  query_params=None):
        """
        Get the list organizations a user is owner or member of
        It is method for GET /users/{username}/organizations
        """
        uri = self.url + "users/"+username+"/organizations"
        uri = uri + build_query_string(query_params)
        return self.session.get(uri)

    def AcceptMembership(self, globalid, role, username,  query_params=None):
        """
        Accept membership in organization
        It is method for POST /users/{username}/organizations/{globalid}/roles/{role}
        """
        uri = self.url + "users/"+username+"/organizations/"+globalid+"/roles/"+role
        uri = uri + build_query_string(query_params)
        return self.session.post(uri, json={})

    def AcceptOrgMembership(self, globalid, role, invitingorg,  query_params=None):
        """
        Accept membership in organization
        It is method for POST /organizations/{globalid}/organizations/{invitingorg}/roles/{role}
        """
        uri = self.url + "organizations/"+globalid+"/organizations/"+invitingorg+"/roles/"+role
        uri = uri + build_query_string(query_params)
        return self.session.post(uri, json={})


    def RejectMembership(self, globalid, role, username,  query_params=None):
        """
        Reject membership invitation in an organization.
        It is method for DELETE /users/{username}/organizations/{globalid}/roles/{role}
        """
        uri = self.url + "/users/"+username+"/organizations/"+globalid+"/roles/"+role
        uri = uri + build_query_string(query_params)
        return self.session.delete(uri)

    def LeaveOrganization(self, globalid, username,  query_params=None):
        """
        Removes the user from an organization.
        It is method for DELETE /users/{username}/organizations/{globalid}/leave
        """
        uri = self.url + "users/"+username+"/organizations/"+globalid+"/leave"
        uri = uri + build_query_string(query_params)
        return self.session.delete(uri)

    def GetUser(self, username,  query_params=None):
        """
        It is method for GET /users/{username}
        """
        uri = self.url + "users/"+username
        uri = uri + build_query_string(query_params)
        return self.session.get(uri)

    def UpdateUserName(self, data, username,  query_params=None):
        """
        Update the user his firstname and lastname
        It is method for PUT /users/{username}/name
        """
        uri = self.url + "users/"+username+"/name"
        uri = uri + build_query_string(query_params)
        return self.session.put(uri, json=data)

    def UpdatePassword(self, data, username,  query_params=None):
        """
        Update the user his password
        It is method for PUT /users/{username}/password
        """
        uri = self.url + "users/"+username+"/password"
        uri = uri + build_query_string(query_params)
        return self.session.put(uri, json=data)

    def GetEmailAddresses(self, username,  query_params=None):
        """
        Get a list of the user his email addresses.
        It is method for GET /users/{username}/emailaddresses
        """
        uri = self.url + "users/"+username+"/emailaddresses"
        uri = uri + build_query_string(query_params)
        return self.session.get(uri)


    def RegisterNewEmailAddress(self, data, username,  query_params=None):
        """
        Register a new email address
        It is method for POST /users/{username}/emailaddresses
        """
        uri = self.url + "users/"+username+"/emailaddresses"
        uri = uri + build_query_string(query_params)
        return self.session.post(uri, json=data)

    def UpdateEmailAddress(self, data, label, username, query_params=None):
        """
        Updates the label and/or value of an email address
        It is method for PUT /users/{username}/emailaddresses/{label}
        """
        uri = self.url + "users/"+username+"/emailaddresses/"+label
        uri = uri + build_query_string(query_params)
        return self.session.put(uri, json=data)

    def DeleteEmailAddress(self, label, username, query_params=None):
        """
        Removes an email address
        It is method for DELETE /users/{username}/emailaddresses/{label}
        """
        uri = self.url + "users/"+username+"/emailaddresses/"+label
        uri = uri + build_query_string(query_params)
        return self.session.delete(uri)

    def ValidateEmailAddress(self, label, username, query_params=None):
        """
        Sends validation email to email address
        It is method for POST /users/{username}/emailaddresses/{label}/validate
        """
        uri = self.url + "users/"+username+"/emailaddresses/"+label+"/validate"
        uri = uri + build_query_string(query_params)
        return self.session.post(uri)
    def ListAPIKeys(self, username, query_params=None):
        """
        Lists the API keys
        It is method for GET /users/{username}/apikeys
        """
        uri = self.url + "users/"+username+"/apikeys"
        uri = uri + build_query_string(query_params)
        return self.session.get(uri)

    def AddApiKey(self, data, username, query_params=None):
        """
        Adds an APIKey to the user
        It is method for POST /users/{username}/apikeys
        """
        uri = self.url + "users/"+username+"/apikeys"
        uri = uri + build_query_string(query_params)
        return self.session.post(uri, json=data)

    def GetAPIkey(self, label, username, query_params=None):
        """
        Get an API key by label
        It is method for GET /users/{username}/apikeys/{label}
        """
        uri = self.url + "users/"+username+"/apikeys/"+label
        uri = uri + build_query_string(query_params)
        return self.session.get(uri)

    def UpdateAPIkey(self, data, label, username, query_params=None):
        """
        Updates the label for the api key
        It is method for PUT /users/{username}/apikeys/{label}
        """
        uri = self.url + "users/"+username+"/apikeys/"+label
        uri = uri + build_query_string(query_params)
        return self.session.put(uri, json=data)

    def DeleteAPIkey(self, label, username, query_params=None):
        """
        Removes an API key
        It is method for DELETE /users/{username}/apikeys/{label}
        """
        uri = self.url + "users/"+username+"/apikeys/"+label
        uri = uri + build_query_string(query_params)
        return self.session.delete(uri)

    def DeleteGithubAccount(self, username, query_params=None):
        """
        Unlink Github Account
        It is method for DELETE /users/{username}/github
        """
        uri = self.url + "users/"+username+"/github"
        uri = uri + build_query_string(query_params)
        return self.session.delete(uri)

    def DeleteFacebookAccount(self, username, query_params=None):
        """
        Delete the associated facebook account
        It is method for DELETE /users/{username}/facebook
        """
        uri = self.url + "users/"+username+"/facebook"
        uri = uri + build_query_string(query_params)
        return self.session.delete(uri)

    def GetUserInformation(self, username, query_params=None):
        """
        It is method for GET /users/{username}/info
        """
        uri = self.url + "users/"+username+"/info"
        uri = uri + build_query_string(query_params)
        return self.session.get(uri)

    def users_byUsernamevalidate_get(self, username, query_params=None):
        """
        It is method for GET /users/{username}/validate
        """
        uri = self.url + "users/"+username+"/validate"
        uri = uri + build_query_string(query_params)
        return self.session.get(uri)

    def GetUserAddresses(self, username, query_params=None):
        """
        It is method for GET /users/{username}/addresses
        """
        uri = self.url + "users/"+username+"/addresses"
        uri = uri + build_query_string(query_params)
        return self.session.get(uri)

    def RegisterNewUserAddress(self, data, username, query_params=None):
        """
        Register a new address
        It is method for POST /users/{username}/addresses
        """
        uri = self.url + "users/"+username+"/addresses"
        uri = uri + build_query_string(query_params)
        return self.session.post(uri, json=data)

    def GetUserAddressByLabel(self, label, username, query_params=None):
        """
        It is method for GET /users/{username}/addresses/{label}
        """
        uri = self.url + "users/"+username+"/addresses/"+label
        uri = uri + build_query_string(query_params)
        return self.session.get(uri)

    def UpdateUserAddress(self, data, label, username, query_params=None):
        """
        Update the label and/or value of an existing address.
        It is method for PUT /users/{username}/addresses/{label}
        """
        uri = self.url + "users/"+username+"/addresses/"+label
        uri = uri + build_query_string(query_params)
        return self.session.put(uri, json=data)

    def DeleteUserAddress(self, label, username, query_params=None):
        """
        Removes an address
        It is method for DELETE /users/{username}/addresses/{label}
        """
        uri = self.url + "users/"+username+"/addresses/"+label
        uri = uri + build_query_string(query_params)
        return self.session.delete(uri)

    def GetUserBankAccounts(self, username, query_params=None):
        """
        It is method for GET /users/{username}/banks
        """
        uri = self.url + "users/"+username+"/banks"
        uri = uri + build_query_string(query_params)
        return self.session.get(uri)

    def CreateUserBankAccount(self, data, username, query_params=None):
        """
        Create new bank account
        It is method for POST /users/{username}/banks
        """
        uri = self.url + "users/"+username+"/banks"
        uri = uri + build_query_string(query_params)
        return self.session.post(uri, json=data)

    def GetUserBankAccountByLabel(self, label, username, query_params=None):
        """
        It is method for GET /users/{username}/banks/{label}
        """
        uri = self.url + "users/"+username+"/banks/"+label
        uri = uri + build_query_string(query_params)
        return self.session.get(uri)

    def UpdateUserBankAccount(self, data, label, username, query_params=None):
        """
        Update an existing bankaccount and label.
        It is method for PUT /users/{username}/banks/{label}
        """
        uri = self.url + "users/"+username+"/banks/"+label
        uri = uri + build_query_string(query_params)
        return self.session.put(uri, json=data)

    def DeleteUserBankAccount(self, label, username, query_params=None):
        """
        Delete a BankAccount
        It is method for DELETE /users/{username}/banks/{label}
        """
        uri = self.url + "users/"+username+"/banks/"+label
        uri = uri + build_query_string(query_params)
        return self.session.delete(uri)

    def GetUserDigitalWallets(self, username, query_params=None):
        """
        It is method for GET /users/{username}/digitalwallet
        """
        uri = self.url + "users/"+username+"/digitalwallet"
        uri = uri + build_query_string(query_params)
        return self.session.get(uri)

    def GetUserDigitalWalletByLabel(self, label, username, query_params=None):
            """
            It is method for GET /users/{username}/digitalwallet/{label}
            """
            uri = self.url + "users/"+username+"/digitalwallet/"+label
            uri = uri + build_query_string(query_params)
            return self.session.get(uri)

    def RegisterDigitalWallet(self, data, username, query_params=None):
        """
        Create new bank account
        It is method for POST /users/{username}/digitalwallet
        """
        uri = self.url + "users/"+username+"/digitalwallet"
        uri = uri + build_query_string(query_params)
        return self.session.post(uri, json=data)

    def UpdateUserDigitalWallet(self, data, label, username, query_params=None):
        """
        Update an existing bankaccount and label.
        It is method for PUT /users/{username}/digitalwallet/{label}
        """
        uri = self.url + "users/"+username+"/digitalwallet/"+label
        uri = uri + build_query_string(query_params)
        return self.session.put(uri, json=data)

    def DeleteUserDigitalWallet(self, label, username, query_params=None):
        """
        Delete a BankAccount
        It is method for DELETE /users/{username}/digitalwallet/{label}
        """
        uri = self.url + "users/"+username+"/digitalwallet/"+label
        uri = uri + build_query_string(query_params)
        return self.session.delete(uri)


    def GetUserPublicKeys(self, username, query_params=None):
        """
        It is method for GET /users/{username}/publickeys
        """
        uri = self.url + "users/"+username+"/publickeys"
        uri = uri + build_query_string(query_params)
        return self.session.get(uri)

    def GetUserPublicKeyByLabel(self, label, username, query_params=None):
            """
            It is method for GET /users/{username}/publickeys/{label}
            """
            uri = self.url + "users/"+username+"/publickeys/"+label
            uri = uri + build_query_string(query_params)
            return self.session.get(uri)

    def RegisterUserPublicKey(self, data, username, query_params=None):
        """
        Create new bank account
        It is method for POST /users/{username}/publickeys
        """
        uri = self.url + "users/"+username+"/publickeys"
        uri = uri + build_query_string(query_params)
        return self.session.post(uri, json=data)

    def UpdateUserPublicKey(self, data, label, username, query_params=None):
        """
        Update an existing bankaccount and label.
        It is method for PUT /users/{username}/publickeys/{label}
        """
        uri = self.url + "users/"+username+"/publickeys/"+label
        uri = uri + build_query_string(query_params)
        return self.session.put(uri, json=data)

    def DeleteUserPublicKey(self, label, username, query_params=None):
        """
        Delete a BankAccount
        It is method for DELETE /users/{username}/publickeys/{label}
        """
        uri = self.url + "users/"+username+"/publickeys/"+label
        uri = uri + build_query_string(query_params)
        return self.session.delete(uri)


    def GetUserPhoneNumbers(self, username, query_params=None):
        """
        It is method for GET /users/{username}/phonenumbers
        """
        uri = self.url + "/users/"+username+"/phonenumbers"
        uri = uri + build_query_string(query_params)
        return self.session.get(uri)

    def RegisterNewUserPhonenumber(self, data, username, query_params=None):
        """
        Register a new phonenumber
        It is method for POST /users/{username}/phonenumbers
        """
        uri = self.url+"users/"+username+"/phonenumbers"
        return self.session.post(uri, json=data)

    def GetUserPhonenumberByLabel(self, label, username, query_params=None):
        """
        It is method for GET /users/{username}/phonenumbers/{label}
        """
        uri = self.url + "/users/"+username+"/phonenumbers/"+label
        uri = uri + build_query_string(query_params)
        return self.session.get(uri)

    def UpdateUserPhonenumber(self, data, label, username, query_params=None):
        """
        Update the label and/or value of an existing phonenumber.
        It is method for PUT /users/{username}/phonenumbers/{label}
        """
        uri = self.url + "users/"+username+"/phonenumbers/"+label
        return self.session.put(uri, json=data)

    def DeleteUserPhonenumber(self, label, username, query_params=None):
        """
        Removes a phonenumber
        It is method for DELETE /users/{username}/phonenumbers/{label}
        """
        uri = self.url + "users/"+username+"/phonenumbers/"+label
        uri = uri + build_query_string(query_params)
        return self.session.delete(uri)

    def ValidatePhonenumber(self,label, username, query_params=None):
        """
        Sends validation text to phone numbers
        It is method for POST /users/{username}/phonenumbers/{label}/validate
        """
        uri = self.url + "users/"+username+"/phonenumbers/"+label+"/validate"
        return self.session.post(uri)

    def VerifyPhoneNumber(self, data, label, username, query_params=None):
        """
        Verifies a phone number
        It is method for PUT /users/{username}/phonenumbers/{label}/validate
        """
        uri = self.url + "users/"+username+"/phonenumbers/"+label+"/validate"
        uri = uri + build_query_string(query_params)
        return self.session.put(uri, json=data)

    def GetUserContracts(self, username, query_params=None):
        """
        Get the contracts where the user is 1 of the parties. Order descending by date.
        It is method for GET /users/{username}/contracts
        """
        uri = self.url + "users/"+username+"/contracts"
        uri = uri + build_query_string(query_params)
        return self.session.get(uri)

    def CreateUserContract(self, data, username, query_params=None):
        """
        Create a new contract.
        It is method for POST /users/{username}/contracts
        """
        uri = self.url + "users/"+username+"/contracts"
        uri = uri + build_query_string(query_params)
        return self.session.post(uri, json=data)

    def GetAllAuthorizations(self, username, query_params=None):
        """
        Get the list of authorizations.
        It is method for GET /users/{username}/authorizations
        """
        uri = self.url + "users/"+username+"/authorizations"
        uri = uri + build_query_string(query_params)
        return self.session.get(uri)

    def GetAuthorization(self, grantedTo, username, query_params=None):
        """
        Get the authorization for a specific organization.
        It is method for GET /users/{username}/authorizations/{grantedTo}
        """
        uri = self.url + "/users/"+username+"/authorizations/"+grantedTo
        uri = uri + build_query_string(query_params)
        return self.session.get(uri)

    def UpdateAuthorization(self, data, grantedTo, username, query_params=None):
        """
        Modify which information an organization is able to see.
        It is method for PUT /users/{username}/authorizations/{grantedTo}
        """
        uri = self.url + "users/"+username+"/authorizations/"+grantedTo
        uri = uri + build_query_string(query_params)
        return self.session.put(uri, json=data)

    def DeleteAuthorization(self, grantedTo, username, query_params=None):
        """
        Remove the authorization for an organization, the granted organization will no longer have access the user's information.
        It is method for DELETE /users/{username}/authorizations/{grantedTo}
        """
        uri = self.url + "users/"+username+"/authorizations/"+grantedTo
        uri = uri + build_query_string(query_params)
        return self.session.delete(uri)





    def CreateNewOrganization(self, data, query_params=None):
        """
        Create a new organization. 1 user should be in the owners list. Validation is performed to check if the securityScheme allows management on this user.
        It is method for POST /organizations
        """
        uri = self.url + "organizations"
        uri = uri + build_query_string(query_params)
        return self.session.post(uri, json=data)

    def CreateNewSubOrganization(self, data, globalid, query_params=None):
        """
        Create a new suborganization.
        It is method for POST /organizations/{globalid}
        """
        uri = self.url + "organizations/"+globalid
        uri = uri + build_query_string(query_params)
        return self.session.post(uri, json=data)

    def UpdateOrganization(self, data, globalid, query_params=None):
        """
        Update organization info
        It is method for PUT /organizations/{globalid}
        """
        uri = self.url + "organizations/"+globalid
        uri = uri + build_query_string(query_params)
        return self.session.put(uri, json=data)

    def DeleteOrganization(self, globalid, query_params=None):
        """
        Delete organization by globalid.
        It is method for DELETE /organizations/{globalid}
        """
        uri = self.url + "organizations/"+globalid
        uri = uri + build_query_string(query_params)
        return self.session.delete(uri)

    def GetOrganization(self, globalid, query_params=None):
        """
        Get organization info
        It is method for GET /organizations/{globalid}
        """
        uri = self.url + "organizations/"+globalid
        uri = uri + build_query_string(query_params)
        return self.session.get(uri)

    def GetOrganizationUsers(self, globalid, query_params=None):
        """
        Get organization info
        It is method for GET /organizations/{globalid}/users
        """
        uri = self.url + "organizations/"+globalid+"/users"
        uri = uri + build_query_string(query_params)
        return self.session.get(uri)

    def GetOrganizationTree(self, globalid, query_params=None):
        """
        It is method for GET /organizations/{globalid}/tree
        """
        uri = self.url + "organizations/"+globalid+"/tree"
        uri = uri + build_query_string(query_params)
        return self.session.get(uri)

    def AddOrganizationMember(self, data, globalid, query_params=None):
        """
        Assign a member to organization.
        It is method for POST /organizations/{globalid}/members
        """
        uri = self.url + "organizations/"+globalid+"/members"
        uri = uri + build_query_string(query_params)
        return self.session.post(uri, json=data)

    def AddOrganizationOrgmember(self, data, globalid, query_params=None):
        """
        Assign a member to organization.
        It is method for POST /organizations/{globalid}/orgmembers
        """
        uri = self.url + "organizations/"+globalid+"/orgmembers"
        uri = uri + build_query_string(query_params)
        return self.session.post(uri, json=data)

    def AddOrganizationOrgowner(self, data, globalid, query_params=None):
        """
        Assign a member to organization.
        It is method for POST /organizations/{globalid}/orgowners
        """
        uri = self.url + "organizations/"+globalid+"/orgowners"
        uri = uri + build_query_string(query_params)
        return self.session.post(uri, json=data)

    def UpdateOrganizationMemberShip(self, data, globalid, query_params=None):
        """
        Update an organization membership
        It is method for PUT /organizations/{globalid}/members
        """
        uri = self.url + "organizations/"+globalid+"/members"
        return self.session.put(uri, json=data)

    def RemoveOrganizationMember(self, username, globalid, query_params=None):
        """
        Remove a member from organization
        It is method for DELETE /organizations/{globalid}/members/{username}
        """
        uri = self.url + "organizations/"+globalid+"/members/"+username
        uri = uri + build_query_string(query_params)
        return self.session.delete(uri)

    def AddOrganizationOwner(self, data, globalid, query_params=None):
        """
        Invite a user to become owner of an organization.
        It is method for POST /organizations/{globalid}/owners
        """
        uri = self.url + "organizations/"+globalid+"/owners"
        uri = uri + build_query_string(query_params)
        return self.session.post(uri, json=data)

    def RemoveOrganizationOwner(self, username, globalid, query_params=None):
        """
        Remove an owner from organization
        It is method for DELETE /organizations/{globalid}/owners/{username}
        """
        uri = self.url + "organizations/"+globalid+"/owners/"+username
        uri = uri + build_query_string(query_params)
        return self.session.delete(uri)

    def GetOrganizationContracts(self, globalid, query_params=None):
        """
        Get the contracts where the organization is 1 of the parties. Order descending by date.
        It is method for GET /organizations/{globalid}/contracts
        """
        uri = self.url + "/organizations/"+globalid+"/contracts"
        uri = uri + build_query_string(query_params)
        return self.session.get(uri)

    def CreateOrganizationContracts(self, data, globalid, query_params=None):
        """
        Create a new contract.
        It is method for POST /organizations/{globalid}/contracts
        """
        uri = self.url + "organizations/"+globalid+"/contracts"
        uri = uri + build_query_string(query_params)
        return self.session.post(uri, json=data)

    def GetPendingOrganizationInvitations(self, globalid, query_params=None):
        """
        Get the list of pending invitations for users to join this organization.
        It is method for GET /organizations/{globalid}/invitations
        """
        uri = self.url + "organizations/"+globalid+"/invitations"
        uri = uri + build_query_string(query_params)
        return self.session.get(uri)

    def RemovePendingOrganizationInvitation(self, username, globalid, query_params=None):
        """
        Cancel a pending invitation.
        It is method for DELETE /organizations/{globalid}/invitations/{username}
        """
        uri = self.url + "organizations/"+globalid+"/invitations/"+username
        uri = uri + build_query_string(query_params)
        return self.session.delete(uri)

    def GetOrganizationAPIKeyLabels(self, globalid, query_params=None):
        """
        Get the list of active api keys.
        It is method for GET /organizations/{globalid}/apikeys
        """
        uri = self.url + "/organizations/"+globalid+"/apikeys"
        uri = uri + build_query_string(query_params)
        return self.session.get(uri)

    def CreateNewOrganizationAPIKey(self, data, globalid, query_params=None):
        """
        Create a new API Key, a secret itself should not be provided, it will be generated serverside.
        It is method for POST /organizations/{globalid}/apikeys
        """
        uri = self.url + "organizations/"+globalid+"/apikeys"
        uri = uri + build_query_string(query_params)
        return self.session.post(uri, json=data)

    def GetOrganizationAPIKey(self, label, globalid, query_params=None):
        """
        It is method for GET /organizations/{globalid}/apikeys/{label}
        """
        uri = self.url + "/organizations/"+globalid+"/apikeys/"+label
        uri = uri + build_query_string(query_params)
        return self.session.get(uri)

    def UpdateOrganizationAPIKey(self, data, label, globalid, query_params=None):
        """
        Updates the label or other properties of a key.
        It is method for PUT /organizations/{globalid}/apikeys/{label}
        """
        uri = self.url + "organizations/"+globalid+"/apikeys/"+label
        uri = uri + build_query_string(query_params)
        return self.session.put(uri, json=data)

    def DeleteOrganizationAPIKey(self, label, globalid, query_params=None):
        """
        Removes an API key
        It is method for DELETE /organizations/{globalid}/apikeys/{label}
        """
        uri = self.url + "organizations/"+globalid+"/apikeys/"+label
        uri = uri + build_query_string(query_params)
        return self.session.delete(uri)

    def GetOrganizationRegistries(self, globalid, query_params=None):
        """
        Get the list of active registry.
        It is method for GET /organizations/{globalid}/registry
        """
        uri = self.url + "organizations/"+globalid+"/registry"
        uri = uri + build_query_string(query_params)
        return self.session.get(uri)

    def GetOrganizationRegistry(self, key, globalid, query_params=None):
        """
        Get the list of active registry.
        It is method for GET /organizations/{globalid}/registry
        """
        uri = self.url + "organizations/"+globalid+"/registry/"+key
        uri = uri + build_query_string(query_params)
        return self.session.get(uri)

    def CreateOrganizationRegistry(self, data, globalid, query_params=None):
        """
        Creates a new registry
        It is method for POST /organizations/{globalid}/registry
        """
        uri = self.url + "organizations/"+globalid+"/registry"
        uri = uri + build_query_string(query_params)
        return self.session.post(uri, json=data)

    def DeleteOrganizaitonRegistry(self, key, globalid,  query_params=None):
        """
        Removes a registry
        It is method for DELETE /organizations/{globalid}/registry/{key}
        """
        uri = self.url + "organizations/"+globalid+"/registry/"+key
        uri = uri + build_query_string(query_params)
        return self.session.delete(uri)

    def CreateOrganizationDNS(self, data, globalid,  query_params=None):
        """
        Creates a new DNS name associated with an organization
        It is method for POST /organizations/{globalid}/dns
        """
        uri = self.url + "organizations/"+globalid+"/dns"
        uri = uri + build_query_string(query_params)
        return self.session.post(uri, json=data)

    def UpdateOrganizationDNS(self, data, dnsname, globalid, query_params=None):
        """
        Updates an existing DNS name associated with an organization
        It is method for PUT /organizations/{globalid}/dns/{dnsname}
        """
        uri = self.url + "organizations/"+globalid+"/dns/"+dnsname
        uri = uri + build_query_string(query_params)
        return self.session.put(uri, json=data)

    def DeleteOrganizaitonDNS(self, dnsname, globalid, query_params=None):
        """
        Removes a DNS name
        It is method for DELETE /organizations/{globalid}/dns/{dnsname}
        """
        uri = self.url + "/organizations/"+globalid+"/dns/"+dnsname
        uri = uri + build_query_string(query_params)
        return self.session.delete(uri)

    def UpdateOrganizationOrgMemberShip(self, data, globalid, query_params=None):
        """
        Update the membership status of an organization
        It is method for PUT /organizations/{globalid}/orgmembers
        """
        uri = self.url + "organizations/"+globalid+"/orgmembers"
        return self.session.put(uri, json=data)

    def SetOrgMember(self, data, globalid, query_params=None):
        """
        Add another organization as a member of this one
        It is method for POST /organizations/{globalid}/orgmembers
        """
        uri = self.url + "organizations/"+globalid+"/orgmembers"
        return self.session.post(uri, json=data)

    def DeleteOrgMember(self, globalid2, globalid, query_params=None):
        """
        Remove an organization as a member
        It is method for DELETE /organizations/{globalid}/orgmembers/{globalid2}
        """
        uri = self.url + "organizations/"+globalid+"/orgmembers/"+globalid2
        return self.session.delete(uri)

    def SetOrgOwner(self, data, globalid, query_params=None):
        """
        Add another organization as an owner of this one
        It is method for POST /organizations/{globalid}/orgowners
        """
        uri = self.url + "organizations/"+globalid+"/orgowners"
        return self.session.post(uri, json=data)

    def DeleteOrgOwner(self, globalid2, globalid, query_params=None):
        """
        Remove an organization as an owner
        It is method for DELETE /organizations/{globalid}/orgowners/{globalid2}
        """
        uri = self.url + "organizations/"+globalid+"/orgowners/"+globalid2
        return self.session.delete(uri)


    def GetOrganizationLogo(self, globalid):
        """
        Get the Logo from an organization
        It is method for GET organizations/{globalid}/logo
        """
        uri = self.url + "organizations/"+globalid+"/logo"
        return self.session.get(uri)


    def SetOrganizationLogo(self, data, globalid):
        """
        Set the organization Logo for the organization
        It is method for PUT organizations/{globalid}/logo
        """
        uri = self.url + "organizations/"+globalid+"/logo"
        return self.session.put(uri, json=data)

    def IncludeSuborgsof(self, data, globalid):
        """
        Add an orgmember or orgowner organization to the includesuborgsof list
        It is method for POST /organizations/{globalid}/orgmembers/includesuborgs
        """
        uri = self.url + "organizations/"+globalid+"/orgmembers/includesuborgs"
        return self.session.post(uri, json=data)

    def RemoveIncludeSuborgsof(self, globalid, orgmember):
        """
        Remove an orgmember or orgowner organization to the includesuborgsof list
        It is method for DELETE /organizations/{globalid}/orgmembers/includesuborgs/{orgmember}
        """
        uri = self.url + "organizations/"+globalid+"/orgmembers/includesuborgs/"+orgmember
        return self.session.delete(uri)

    def DeleteOrganizationLogo(self, globalid):
        """
        Removes the Logo from an organization
        It is method for DELETE organizations/{globalid}/logo
        """
        uri = self.url + "organizations/"+globalid+"/logo"
        return self.session.delete(uri)


    def CreateCompany(self, data, query_params=None):
        """
        Register a new company
        It is method for POST companies
        """
        uri = self.url + "companies"
        uri = uri + build_query_string(query_params)
        return self.session.post(uri, json=data)

    def GetCompanyList(self, query_params=None):
        """
        Get companies. Authorization limits are applied to requesting user.
        It is method for GET companies
        """
        uri = self.url + "companies"
        uri = uri + build_query_string(query_params)
        return self.session.get(uri)

    def GetCompany(self, globalId, query_params=None):
        """
        Get organization info
        It is method for GET companies/{globalId}
        """
        uri = self.url + "companies/"+globalId
        uri = uri + build_query_string(query_params)
        return self.session.get(uri)

    def UpdateCompany(self, data, globalId, query_params=None):
        """
        Update existing company. Updating ``globalId`` is not allowed.
        It is method for PUT companies/{globalId}
        """
        uri = self.url + "companies/"+globalId
        uri = uri + build_query_string(query_params)
        return self.session.put(uri, json=data)

    def GetCompanyContracts(self, globalId, query_params=None):
        """
        Get the contracts where the organization is 1 of the parties. Order descending by date.
        It is method for GET companies/{globalId}/contracts
        """
        uri = self.url + "companies/"+globalId+"/contracts"
        uri = uri + build_query_string(query_params)
        return self.session.get(uri)

    def CreateCompanyContract(self, data, globalId, query_params=None):
        """
        Create a new contract.
        It is method for POST companies/{globalId}/contracts
        """
        uri = self.url + "companies/"+globalId+"/contracts"
        uri = uri + build_query_string(query_params)
        return self.session.post(uri, json=data)

    def GetCompanyInfo(self, globalId, query_params=None):
        """
        It is method for GET companies/{globalId}/info
        """
        uri = self.url + "companies/"+globalId+"/info"
        uri = uri + build_query_string(query_params)
        return self.session.get(uri)

    def companies_byGlobalId_validate_get(self, globalId, query_params=None):
        """
        It is method for GET companies/{globalId}/validate
        """
        uri = self.url + "companies/"+globalId+"/validate"
        uri = uri + build_query_string(query_params)
        return self.session.get(uri)

    def GetContract(self, contractId, query_params=None):
        """
        Get a contract
        It is method for GET /contracts/{contractId}
        """
        uri = self.url + "contracts/"+contractId
        uri = uri + build_query_string(query_params)
        return self.session.get(uri)

    def SignContract(self, data, contractId, query_params=None):
        """
        Sign a contract
        It is method for POST /contracts/{contractId}/signatures
        """
        uri = self.url + "contracts/"+contractId+"/signatures"
        uri = uri + build_query_string(query_params)
        return self.session.post(uri, json=data)

    def GetRegistries(self, username, query_params=None):
        """
        It is method for GET /users/{username}/registry
        """
        uri = self.url + "users/"+username+"/registry"
        uri = uri + build_query_string(query_params)
        return self.session.get(uri)

    def GetRegistry(self, label, username, query_params=None):
        """
        It is method for GET /users/{username}/registry
        """
        uri = self.url + "users/"+username+"/registry"+label
        uri = uri + build_query_string(query_params)
        return self.session.get(uri)

    def CreateNewRegistry(self, data, username, query_params=None):
        """
        It is method for GET /users/{username}/registry
        """
        uri = self.url + "users/"+username+"/registry"
        uri = uri + build_query_string(query_params)
        return self.session.post(uri, json=data)

    def DeleteRegistry(self, key, username, query_params=None):
        """
        It is method for DELETE /users/{username}/registry
        """
        uri = self.url + "users/"+username+"/registry"+key
        uri = uri + build_query_string(query_params)
        return self.session.delete(uri)

    def GetTwofamethods(self, username, query_params=None):
        """
        It is method for GET /users/{username}/twofamethods
        """
        uri = self.url + "users/"+username+"/twofamethods"
        uri = uri + build_query_string(query_params)
        return self.session.get(uri)

    def GetTotp(self, username, query_params=None):
        """
        It is method for GET /users/{username}/totp
        """
        uri = self.url + "users/"+username+"/totp"
        uri = uri + build_query_string(query_params)
        return self.session.get(uri)

    def EditTotp(self, data, username, query_params=None):
        """
        It is method for POST /users/{username}/totp
        """
        uri = self.url + "users/"+username+"/totp"
        uri = uri + build_query_string(query_params)
        return self.session.post(uri, json=data)

    def DeleteTotp(self, username, query_params=None):
        """
        It is method for Delete /users/{username}/totp
        """
        uri = self.url + "users/"+username+"/totp"
        uri = uri + build_query_string(query_params)
        return self.session.delete(uri)

    def AddOrganizationDescription(self, data, globalid, query_params=None):
        """
        It is method for POST organizations/{globalid}/description
        """
        uri = self.url + "organizations/"+globalid+"/description"
        uri = uri + build_query_string(query_params)
        return self.session.post(uri, json=data)

    def UpdateOrganizationDescription(self, data, globalid, query_params=None):
        """
        It is method for PUT organizations/{globalid}/description
        """
        uri = self.url + "organizations/"+globalid+"/description"
        uri = uri + build_query_string(query_params)
        return self.session.put(uri, json=data)

    def GetOrganizationDescription(self, langkey, globalid, query_params=None):
        """
        It is method for GET organizations/{globalid}/description/{langkey}
        """
        uri = self.url + "organizations/"+globalid+"/description/"+langkey
        uri = uri + build_query_string(query_params)
        return self.session.get(uri)

    def GetOrganizationDescriptionWithfallback(self, langkey, globalid, query_params=None):
        """
        It is method for GET organizations/{globalid}/description/{langkey}/withfallback
        """
        uri = self.url + "organizations/"+globalid+"/description/"+langkey+"/withfallback"
        uri = uri + build_query_string(query_params)
        return self.session.get(uri)


    def DeleteOrganizationDescription(self, langkey, globalid, query_params=None):
        """
        It is method for DELETE organizations/{globalid}/description/{langkey}
        """
        uri = self.url + "organizations/"+globalid+"/description/"+langkey
        uri = uri + build_query_string(query_params)
        return self.session.delete(uri)

    def GetTwoFA(self, globalid, query_params=None):
        """
        It is method for GET organizations/{globalid}/2fa/validity
        """
        uri = self.url + "organizations/"+globalid+"/2fa/validity"
        uri = uri + build_query_string(query_params)
        return self.session.get(uri)

    def UpdateTwoFA(self, data, globalid, query_params=None):
        """
        It is method for PUT organizations/{globalid}/2fa/validity
        """
        uri = self.url + "organizations/"+globalid+"/2fa/validity"
        uri = uri + build_query_string(query_params)
        return self.session.put(uri, json=data)

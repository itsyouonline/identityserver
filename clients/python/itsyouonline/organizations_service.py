
from .APIKeyLabel import APIKeyLabel
from .Contract import Contract
from .DnsAddress import DnsAddress
from .Error import Error
from .GetOrganizationUsersResponseBody import GetOrganizationUsersResponseBody
from .IsMember import IsMember
from .JoinOrganizationInvitation import JoinOrganizationInvitation
from .LocalizedInfoText import LocalizedInfoText
from .Organization import Organization
from .OrganizationAPIKey import OrganizationAPIKey
from .OrganizationLogo import OrganizationLogo
from .OrganizationTreeItem import OrganizationTreeItem
from .RegistryEntry import RegistryEntry
from .ValidityTime import ValidityTime
from .api_response import APIResponse
from .unhandled_api_error import UnhandledAPIError
from .unmarshall_error import UnmarshallError


class OrganizationsService:
    def __init__(self, client):
        self.client = client

    def Get2faValidityTime(self, globalid, headers=None, query_params=None, content_type="application/json"):
        """
        Get the 2FA validity time for the organization, in seconds
        It is method for GET /organizations/{globalid}/2fa/validity
        """
        uri = self.client.base_url + "/organizations/" + globalid + "/2fa/validity"
        resp = self.client.get(uri, None, headers, query_params, content_type)
        try:
            if resp.status_code == 200:
                return APIResponse(data=ValidityTime(resp.json()), response=resp)

            message = 'unknown status code={}'.format(resp.status_code)
            raise UnhandledAPIError(response=resp, code=resp.status_code,
                                    message=message)
        except ValueError as msg:
            raise UnmarshallError(resp, msg)
        except UnhandledAPIError as uae:
            raise uae
        except Exception as e:
            raise UnmarshallError(resp, e.message)

    def Set2faValidityTime(self, data, globalid, headers=None, query_params=None, content_type="application/json"):
        """
        Update the 2FA validity time for the organization
        It is method for PUT /organizations/{globalid}/2fa/validity
        """
        uri = self.client.base_url + "/organizations/" + globalid + "/2fa/validity"
        return self.client.put(uri, data, headers, query_params, content_type)

    def DeleteOrganizationAPIKey(
            self,
            label,
            globalid,
            headers=None,
            query_params=None,
            content_type="application/json"):
        """
        Removes an API key
        It is method for DELETE /organizations/{globalid}/apikeys/{label}
        """
        uri = self.client.base_url + "/organizations/" + globalid + "/apikeys/" + label
        return self.client.delete(uri, None, headers, query_params, content_type)

    def GetOrganizationAPIKey(self, label, globalid, headers=None, query_params=None, content_type="application/json"):
        """
        Get an api key from an organization
        It is method for GET /organizations/{globalid}/apikeys/{label}
        """
        uri = self.client.base_url + "/organizations/" + globalid + "/apikeys/" + label
        resp = self.client.get(uri, None, headers, query_params, content_type)
        try:
            if resp.status_code == 200:
                return APIResponse(data=OrganizationAPIKey(resp.json()), response=resp)

            message = 'unknown status code={}'.format(resp.status_code)
            raise UnhandledAPIError(response=resp, code=resp.status_code,
                                    message=message)
        except ValueError as msg:
            raise UnmarshallError(resp, msg)
        except UnhandledAPIError as uae:
            raise uae
        except Exception as e:
            raise UnmarshallError(resp, e.message)

    def UpdateOrganizationAPIKey(
            self,
            data,
            label,
            globalid,
            headers=None,
            query_params=None,
            content_type="application/json"):
        """
        Updates the label or other properties of a key.
        It is method for PUT /organizations/{globalid}/apikeys/{label}
        """
        uri = self.client.base_url + "/organizations/" + globalid + "/apikeys/" + label
        return self.client.put(uri, data, headers, query_params, content_type)

    def GetOrganizationAPIKeyLabels(self, globalid, headers=None, query_params=None, content_type="application/json"):
        """
        Get the list of active api keys.
        It is method for GET /organizations/{globalid}/apikeys
        """
        uri = self.client.base_url + "/organizations/" + globalid + "/apikeys"
        resp = self.client.get(uri, None, headers, query_params, content_type)
        try:
            if resp.status_code == 200:
                resps = []
                for elem in resp.json():
                    resps.append(APIKeyLabel(elem))
                return APIResponse(data=resps, response=resp)

            message = 'unknown status code={}'.format(resp.status_code)
            raise UnhandledAPIError(response=resp, code=resp.status_code,
                                    message=message)
        except ValueError as msg:
            raise UnmarshallError(resp, msg)
        except UnhandledAPIError as uae:
            raise uae
        except Exception as e:
            raise UnmarshallError(resp, e.message)

    def CreateNewOrganizationAPIKey(
            self,
            data,
            globalid,
            headers=None,
            query_params=None,
            content_type="application/json"):
        """
        Create a new API Key, a secret itself should not be provided, it will be generated serverside.
        It is method for POST /organizations/{globalid}/apikeys
        """
        uri = self.client.base_url + "/organizations/" + globalid + "/apikeys"
        resp = self.client.post(uri, data, headers, query_params, content_type)
        try:
            if resp.status_code == 201:
                return APIResponse(data=OrganizationAPIKey(resp.json()), response=resp)

            message = 'unknown status code={}'.format(resp.status_code)
            raise UnhandledAPIError(response=resp, code=resp.status_code,
                                    message=message)
        except ValueError as msg:
            raise UnmarshallError(resp, msg)
        except UnhandledAPIError as uae:
            raise uae
        except Exception as e:
            raise UnmarshallError(resp, e.message)

    def GetOrganizationContracts(self, globalid, headers=None, query_params=None, content_type="application/json"):
        """
        Get the contracts where the organization is 1 of the parties. Order descending by date.
        It is method for GET /organizations/{globalid}/contracts
        """
        uri = self.client.base_url + "/organizations/" + globalid + "/contracts"
        return self.client.get(uri, None, headers, query_params, content_type)

    def CreateOrganizationContracty(
            self,
            data,
            globalid,
            headers=None,
            query_params=None,
            content_type="application/json"):
        """
        Create a new contract.
        It is method for POST /organizations/{globalid}/contracts
        """
        uri = self.client.base_url + "/organizations/" + globalid + "/contracts"
        resp = self.client.post(uri, data, headers, query_params, content_type)
        try:
            if resp.status_code == 201:
                return APIResponse(data=Contract(resp.json()), response=resp)

            message = 'unknown status code={}'.format(resp.status_code)
            raise UnhandledAPIError(response=resp, code=resp.status_code,
                                    message=message)
        except ValueError as msg:
            raise UnmarshallError(resp, msg)
        except UnhandledAPIError as uae:
            raise uae
        except Exception as e:
            raise UnmarshallError(resp, e.message)

    def GetDescriptionWithFallback(
            self,
            langkey,
            globalid,
            headers=None,
            query_params=None,
            content_type="application/json"):
        """
        Get the description for an organization for this langkey, try to use the English is there is no description for this langkey
        It is method for GET /organizations/{globalid}/description/{langkey}/withfallback
        """
        uri = self.client.base_url + "/organizations/" + globalid + "/description/" + langkey + "/withfallback"
        resp = self.client.get(uri, None, headers, query_params, content_type)
        try:
            if resp.status_code == 200:
                return APIResponse(data=LocalizedInfoText(resp.json()), response=resp)

            message = 'unknown status code={}'.format(resp.status_code)
            raise UnhandledAPIError(response=resp, code=resp.status_code,
                                    message=message)
        except ValueError as msg:
            raise UnmarshallError(resp, msg)
        except UnhandledAPIError as uae:
            raise uae
        except Exception as e:
            raise UnmarshallError(resp, e.message)

    def DeleteDescription(self, langkey, globalid, headers=None, query_params=None, content_type="application/json"):
        """
        Delete the description for this organization for a given language key
        It is method for DELETE /organizations/{globalid}/description/{langkey}
        """
        uri = self.client.base_url + "/organizations/" + globalid + "/description/" + langkey
        return self.client.delete(uri, None, headers, query_params, content_type)

    def GetDescription(self, langkey, globalid, headers=None, query_params=None, content_type="application/json"):
        """
        Get the description for an organization for this langkey
        It is method for GET /organizations/{globalid}/description/{langkey}
        """
        uri = self.client.base_url + "/organizations/" + globalid + "/description/" + langkey
        resp = self.client.get(uri, None, headers, query_params, content_type)
        try:
            if resp.status_code == 200:
                return APIResponse(data=LocalizedInfoText(resp.json()), response=resp)

            message = 'unknown status code={}'.format(resp.status_code)
            raise UnhandledAPIError(response=resp, code=resp.status_code,
                                    message=message)
        except ValueError as msg:
            raise UnmarshallError(resp, msg)
        except UnhandledAPIError as uae:
            raise uae
        except Exception as e:
            raise UnmarshallError(resp, e.message)

    def SetDescription(self, data, globalid, headers=None, query_params=None, content_type="application/json"):
        """
        Set the description for this organization for a given language key
        It is method for POST /organizations/{globalid}/description
        """
        uri = self.client.base_url + "/organizations/" + globalid + "/description"
        resp = self.client.post(uri, data, headers, query_params, content_type)
        try:
            if resp.status_code == 201:
                return APIResponse(data=LocalizedInfoText(resp.json()), response=resp)

            message = 'unknown status code={}'.format(resp.status_code)
            raise UnhandledAPIError(response=resp, code=resp.status_code,
                                    message=message)
        except ValueError as msg:
            raise UnmarshallError(resp, msg)
        except UnhandledAPIError as uae:
            raise uae
        except Exception as e:
            raise UnmarshallError(resp, e.message)

    def UpdateDescription(self, data, globalid, headers=None, query_params=None, content_type="application/json"):
        """
        Update the description for this organization for a given language key
        It is method for PUT /organizations/{globalid}/description
        """
        uri = self.client.base_url + "/organizations/" + globalid + "/description"
        resp = self.client.put(uri, data, headers, query_params, content_type)
        try:
            if resp.status_code == 200:
                return APIResponse(data=LocalizedInfoText(resp.json()), response=resp)

            message = 'unknown status code={}'.format(resp.status_code)
            raise UnhandledAPIError(response=resp, code=resp.status_code,
                                    message=message)
        except ValueError as msg:
            raise UnmarshallError(resp, msg)
        except UnhandledAPIError as uae:
            raise uae
        except Exception as e:
            raise UnmarshallError(resp, e.message)

    def DeleteOrganizationDns(
            self,
            dnsname,
            globalid,
            headers=None,
            query_params=None,
            content_type="application/json"):
        """
        Removes a DNS name associated with an organization
        It is method for DELETE /organizations/{globalid}/dns/{dnsname}
        """
        uri = self.client.base_url + "/organizations/" + globalid + "/dns/" + dnsname
        return self.client.delete(uri, None, headers, query_params, content_type)

    def UpdateOrganizationDns(
            self,
            data,
            dnsname,
            globalid,
            headers=None,
            query_params=None,
            content_type="application/json"):
        """
        Updates an existing DNS name associated with an organization
        It is method for PUT /organizations/{globalid}/dns/{dnsname}
        """
        uri = self.client.base_url + "/organizations/" + globalid + "/dns/" + dnsname
        return self.client.put(uri, data, headers, query_params, content_type)

    def CreateOrganizationDns(self, data, globalid, headers=None, query_params=None, content_type="application/json"):
        """
        Creates a new DNS name associated with an organization
        It is method for POST /organizations/{globalid}/dns
        """
        uri = self.client.base_url + "/organizations/" + globalid + "/dns"
        resp = self.client.post(uri, data, headers, query_params, content_type)
        try:
            if resp.status_code == 201:
                return APIResponse(data=DnsAddress(resp.json()), response=resp)

            message = 'unknown status code={}'.format(resp.status_code)
            raise UnhandledAPIError(response=resp, code=resp.status_code,
                                    message=message)
        except ValueError as msg:
            raise UnmarshallError(resp, msg)
        except UnhandledAPIError as uae:
            raise uae
        except Exception as e:
            raise UnmarshallError(resp, e.message)

    def RemovePendingOrganizationInvitation(
            self,
            username,
            globalid,
            headers=None,
            query_params=None,
            content_type="application/json"):
        """
        Cancel a pending invitation.
        It is method for DELETE /organizations/{globalid}/invitations/{username}
        """
        uri = self.client.base_url + "/organizations/" + globalid + "/invitations/" + username
        return self.client.delete(uri, None, headers, query_params, content_type)

    def GetInvitations(self, globalid, headers=None, query_params=None, content_type="application/json"):
        """
        Get the list of pending invitations for users to join this organization.
        It is method for GET /organizations/{globalid}/invitations
        """
        uri = self.client.base_url + "/organizations/" + globalid + "/invitations"
        resp = self.client.get(uri, None, headers, query_params, content_type)
        try:
            if resp.status_code == 200:
                resps = []
                for elem in resp.json():
                    resps.append(JoinOrganizationInvitation(elem))
                return APIResponse(data=resps, response=resp)

            message = 'unknown status code={}'.format(resp.status_code)
            raise UnhandledAPIError(response=resp, code=resp.status_code,
                                    message=message)
        except ValueError as msg:
            raise UnmarshallError(resp, msg)
        except UnhandledAPIError as uae:
            raise uae
        except Exception as e:
            raise UnmarshallError(resp, e.message)

    def DeleteOrganizationLogo(self, globalid, headers=None, query_params=None, content_type="application/json"):
        """
        Removes the Logo from an organization
        It is method for DELETE /organizations/{globalid}/logo
        """
        uri = self.client.base_url + "/organizations/" + globalid + "/logo"
        return self.client.delete(uri, None, headers, query_params, content_type)

    def GetOrganizationLogo(self, globalid, headers=None, query_params=None, content_type="application/json"):
        """
        Get the Logo from an organization
        It is method for GET /organizations/{globalid}/logo
        """
        uri = self.client.base_url + "/organizations/" + globalid + "/logo"
        resp = self.client.get(uri, None, headers, query_params, content_type)
        try:
            if resp.status_code == 200:
                return APIResponse(data=OrganizationLogo(resp.json()), response=resp)

            message = 'unknown status code={}'.format(resp.status_code)
            raise UnhandledAPIError(response=resp, code=resp.status_code,
                                    message=message)
        except ValueError as msg:
            raise UnmarshallError(resp, msg)
        except UnhandledAPIError as uae:
            raise uae
        except Exception as e:
            raise UnmarshallError(resp, e.message)

    def SetOrganizationLogo(self, data, globalid, headers=None, query_params=None, content_type="application/json"):
        """
        Set the organization Logo for the organization
        It is method for PUT /organizations/{globalid}/logo
        """
        uri = self.client.base_url + "/organizations/" + globalid + "/logo"
        resp = self.client.put(uri, data, headers, query_params, content_type)
        try:
            if resp.status_code == 200:
                return APIResponse(data=OrganizationLogo(resp.json()), response=resp)

            message = 'unknown status code={}'.format(resp.status_code)
            raise UnhandledAPIError(response=resp, code=resp.status_code,
                                    message=message)
        except ValueError as msg:
            raise UnmarshallError(resp, msg)
        except UnhandledAPIError as uae:
            raise uae
        except Exception as e:
            raise UnmarshallError(resp, e.message)

    def RemoveOrganizationMember(
            self,
            username,
            globalid,
            headers=None,
            query_params=None,
            content_type="application/json"):
        """
        Remove a member from an organization.
        It is method for DELETE /organizations/{globalid}/members/{username}
        """
        uri = self.client.base_url + "/organizations/" + globalid + "/members/" + username
        return self.client.delete(uri, None, headers, query_params, content_type)

    def AddOrganizationMember(self, data, globalid, headers=None, query_params=None, content_type="application/json"):
        """
        Invite someone to become member of an organization.
        It is method for POST /organizations/{globalid}/members
        """
        uri = self.client.base_url + "/organizations/" + globalid + "/members"
        resp = self.client.post(uri, data, headers, query_params, content_type)
        try:
            if resp.status_code == 201:
                return APIResponse(data=JoinOrganizationInvitation(resp.json()), response=resp)

            message = 'unknown status code={}'.format(resp.status_code)
            raise UnhandledAPIError(response=resp, code=resp.status_code,
                                    message=message)
        except ValueError as msg:
            raise UnmarshallError(resp, msg)
        except UnhandledAPIError as uae:
            raise uae
        except Exception as e:
            raise UnmarshallError(resp, e.message)

    def UpdateOrganizationMemberShip(
            self,
            data,
            globalid,
            headers=None,
            query_params=None,
            content_type="application/json"):
        """
        Update an organization membership
        It is method for PUT /organizations/{globalid}/members
        """
        uri = self.client.base_url + "/organizations/" + globalid + "/members"
        resp = self.client.put(uri, data, headers, query_params, content_type)
        try:
            if resp.status_code == 200:
                return APIResponse(data=Organization(resp.json()), response=resp)

            message = 'unknown status code={}'.format(resp.status_code)
            raise UnhandledAPIError(response=resp, code=resp.status_code,
                                    message=message)
        except ValueError as msg:
            raise UnmarshallError(resp, msg)
        except UnhandledAPIError as uae:
            raise uae
        except Exception as e:
            raise UnmarshallError(resp, e.message)

    def RejectOrganizationInvite(
            self,
            invitingorg,
            role,
            globalid,
            headers=None,
            query_params=None,
            content_type="application/json"):
        """
        Reject the invite for one of your organizations
        It is method for DELETE /organizations/{globalid}/organizations/{invitingorg}/roles/{role}
        """
        uri = self.client.base_url + "/organizations/" + globalid + "/organizations/" + invitingorg + "/roles/" + role
        return self.client.delete(uri, None, headers, query_params, content_type)

    def AcceptOrganizationInvite(
            self,
            data,
            invitingorg,
            role,
            globalid,
            headers=None,
            query_params=None,
            content_type="application/json"):
        """
        Accept the invite for one of your organizations
        It is method for POST /organizations/{globalid}/organizations/{invitingorg}/roles/{role}
        """
        uri = self.client.base_url + "/organizations/" + globalid + "/organizations/" + invitingorg + "/roles/" + role
        resp = self.client.post(uri, data, headers, query_params, content_type)
        try:
            if resp.status_code == 201:
                return APIResponse(data=JoinOrganizationInvitation(resp.json()), response=resp)

            message = 'unknown status code={}'.format(resp.status_code)
            raise UnhandledAPIError(response=resp, code=resp.status_code,
                                    message=message)
        except ValueError as msg:
            raise UnmarshallError(resp, msg)
        except UnhandledAPIError as uae:
            raise uae
        except Exception as e:
            raise UnmarshallError(resp, e.message)

    def RemoveIncludeSubOrgsOf(
            self,
            orgmember,
            globalid,
            headers=None,
            query_params=None,
            content_type="application/json"):
        """
        Remove an orgmember or orgowner organization to the includesuborgsof list
        It is method for DELETE /organizations/{globalid}/orgmembers/includesuborgs/{orgmember}
        """
        uri = self.client.base_url + "/organizations/" + globalid + "/orgmembers/includesuborgs/" + orgmember
        return self.client.delete(uri, None, headers, query_params, content_type)

    def AddIncludeSubOrgsOf(self, data, globalid, headers=None, query_params=None, content_type="application/json"):
        """
        Add an orgmember or orgowner organization to the includesuborgsof list
        It is method for POST /organizations/{globalid}/orgmembers/includesuborgs
        """
        uri = self.client.base_url + "/organizations/" + globalid + "/orgmembers/includesuborgs"
        resp = self.client.post(uri, data, headers, query_params, content_type)
        try:
            if resp.status_code == 201:
                return APIResponse(data=Organization(resp.json()), response=resp)

            message = 'unknown status code={}'.format(resp.status_code)
            raise UnhandledAPIError(response=resp, code=resp.status_code,
                                    message=message)
        except ValueError as msg:
            raise UnmarshallError(resp, msg)
        except UnhandledAPIError as uae:
            raise uae
        except Exception as e:
            raise UnmarshallError(resp, e.message)

    def DeleteOrgMember(self, globalid2, globalid, headers=None, query_params=None, content_type="application/json"):
        """
        Remove an organization as a member
        It is method for DELETE /organizations/{globalid}/orgmembers/{globalid2}
        """
        uri = self.client.base_url + "/organizations/" + globalid + "/orgmembers/" + globalid2
        return self.client.delete(uri, None, headers, query_params, content_type)

    def SetOrgMember(self, data, globalid, headers=None, query_params=None, content_type="application/json"):
        """
        Add another organization as a member of this one
        It is method for POST /organizations/{globalid}/orgmembers
        """
        uri = self.client.base_url + "/organizations/" + globalid + "/orgmembers"
        return self.client.post(uri, data, headers, query_params, content_type)

    def UpdateOrganizationOrgMemberShip(
            self,
            data,
            globalid,
            headers=None,
            query_params=None,
            content_type="application/json"):
        """
        Update the membership status of an organization
        It is method for PUT /organizations/{globalid}/orgmembers
        """
        uri = self.client.base_url + "/organizations/" + globalid + "/orgmembers"
        resp = self.client.put(uri, data, headers, query_params, content_type)
        try:
            if resp.status_code == 200:
                return APIResponse(data=Organization(resp.json()), response=resp)

            message = 'unknown status code={}'.format(resp.status_code)
            raise UnhandledAPIError(response=resp, code=resp.status_code,
                                    message=message)
        except ValueError as msg:
            raise UnmarshallError(resp, msg)
        except UnhandledAPIError as uae:
            raise uae
        except Exception as e:
            raise UnmarshallError(resp, e.message)

    def DeleteOrgOwner(self, globalid2, globalid, headers=None, query_params=None, content_type="application/json"):
        """
        Remove an organization as an owner
        It is method for DELETE /organizations/{globalid}/orgowners/{globalid2}
        """
        uri = self.client.base_url + "/organizations/" + globalid + "/orgowners/" + globalid2
        return self.client.delete(uri, None, headers, query_params, content_type)

    def SetOrgOwner(self, data, globalid, headers=None, query_params=None, content_type="application/json"):
        """
        Add another organization as an owner of this one
        It is method for POST /organizations/{globalid}/orgowners
        """
        uri = self.client.base_url + "/organizations/" + globalid + "/orgowners"
        return self.client.post(uri, data, headers, query_params, content_type)

    def RemoveOrganizationOwner(
            self,
            username,
            globalid,
            headers=None,
            query_params=None,
            content_type="application/json"):
        """
        Remove an owner from organization
        It is method for DELETE /organizations/{globalid}/owners/{username}
        """
        uri = self.client.base_url + "/organizations/" + globalid + "/owners/" + username
        return self.client.delete(uri, None, headers, query_params, content_type)

    def AddOrganizationOwner(self, data, globalid, headers=None, query_params=None, content_type="application/json"):
        """
        Invite someone to become owner of an organization.
        It is method for POST /organizations/{globalid}/owners
        """
        uri = self.client.base_url + "/organizations/" + globalid + "/owners"
        resp = self.client.post(uri, data, headers, query_params, content_type)
        try:
            if resp.status_code == 201:
                return APIResponse(data=JoinOrganizationInvitation(resp.json()), response=resp)

            message = 'unknown status code={}'.format(resp.status_code)
            raise UnhandledAPIError(response=resp, code=resp.status_code,
                                    message=message)
        except ValueError as msg:
            raise UnmarshallError(resp, msg)
        except UnhandledAPIError as uae:
            raise uae
        except Exception as e:
            raise UnmarshallError(resp, e.message)

    def DeleteOrganizationRegistryEntry(
            self,
            key,
            globalid,
            headers=None,
            query_params=None,
            content_type="application/json"):
        """
        Removes a RegistryEntry from the organization's registry
        It is method for DELETE /organizations/{globalid}/registry/{key}
        """
        uri = self.client.base_url + "/organizations/" + globalid + "/registry/" + key
        return self.client.delete(uri, None, headers, query_params, content_type)

    def GetOrganizationRegistryEntry(
            self,
            key,
            globalid,
            headers=None,
            query_params=None,
            content_type="application/json"):
        """
        Get a RegistryEntry from the organization's registry.
        It is method for GET /organizations/{globalid}/registry/{key}
        """
        uri = self.client.base_url + "/organizations/" + globalid + "/registry/" + key
        resp = self.client.get(uri, None, headers, query_params, content_type)
        try:
            if resp.status_code == 200:
                return APIResponse(data=RegistryEntry(resp.json()), response=resp)

            message = 'unknown status code={}'.format(resp.status_code)
            raise UnhandledAPIError(response=resp, code=resp.status_code,
                                    message=message)
        except ValueError as msg:
            raise UnmarshallError(resp, msg)
        except UnhandledAPIError as uae:
            raise uae
        except Exception as e:
            raise UnmarshallError(resp, e.message)

    def ListOrganizationRegistry(self, globalid, headers=None, query_params=None, content_type="application/json"):
        """
        Lists the RegistryEntries in an organization's registry.
        It is method for GET /organizations/{globalid}/registry
        """
        uri = self.client.base_url + "/organizations/" + globalid + "/registry"
        resp = self.client.get(uri, None, headers, query_params, content_type)
        try:
            if resp.status_code == 200:
                resps = []
                for elem in resp.json():
                    resps.append(RegistryEntry(elem))
                return APIResponse(data=resps, response=resp)

            message = 'unknown status code={}'.format(resp.status_code)
            raise UnhandledAPIError(response=resp, code=resp.status_code,
                                    message=message)
        except ValueError as msg:
            raise UnmarshallError(resp, msg)
        except UnhandledAPIError as uae:
            raise uae
        except Exception as e:
            raise UnmarshallError(resp, e.message)

    def AddOrganizationRegistryEntry(
            self,
            data,
            globalid,
            headers=None,
            query_params=None,
            content_type="application/json"):
        """
        Adds a RegistryEntry to the organization's registry, if the key is already used, it is overwritten.
        It is method for POST /organizations/{globalid}/registry
        """
        uri = self.client.base_url + "/organizations/" + globalid + "/registry"
        resp = self.client.post(uri, data, headers, query_params, content_type)
        try:
            if resp.status_code == 201:
                return APIResponse(data=RegistryEntry(resp.json()), response=resp)

            message = 'unknown status code={}'.format(resp.status_code)
            raise UnhandledAPIError(response=resp, code=resp.status_code,
                                    message=message)
        except ValueError as msg:
            raise UnmarshallError(resp, msg)
        except UnhandledAPIError as uae:
            raise uae
        except Exception as e:
            raise UnmarshallError(resp, e.message)

    def DeleteRequiredScope(
            self,
            requiredscope,
            globalid,
            headers=None,
            query_params=None,
            content_type="application/json"):
        """
        Deletes a required scope
        It is method for DELETE /organizations/{globalid}/requiredscopes/{requiredscope}
        """
        uri = self.client.base_url + "/organizations/" + globalid + "/requiredscopes/" + requiredscope
        return self.client.delete(uri, None, headers, query_params, content_type)

    def UpdateRequiredScope(
            self,
            data,
            requiredscope,
            globalid,
            headers=None,
            query_params=None,
            content_type="application/json"):
        """
        Updates a required scope
        It is method for PUT /organizations/{globalid}/requiredscopes/{requiredscope}
        """
        uri = self.client.base_url + "/organizations/" + globalid + "/requiredscopes/" + requiredscope
        return self.client.put(uri, data, headers, query_params, content_type)

    def AddRequiredScope(self, data, globalid, headers=None, query_params=None, content_type="application/json"):
        """
        Adds a required scope
        It is method for POST /organizations/{globalid}/requiredscopes
        """
        uri = self.client.base_url + "/organizations/" + globalid + "/requiredscopes"
        return self.client.post(uri, data, headers, query_params, content_type)

    def GetOrganizationTree(self, globalid, headers=None, query_params=None, content_type="application/json"):
        """
        Tree structure of all suborganizations
        It is method for GET /organizations/{globalid}/tree
        """
        uri = self.client.base_url + "/organizations/" + globalid + "/tree"
        resp = self.client.get(uri, None, headers, query_params, content_type)
        try:
            if resp.status_code == 200:
                resps = []
                for elem in resp.json():
                    resps.append(OrganizationTreeItem(elem))
                return APIResponse(data=resps, response=resp)

            message = 'unknown status code={}'.format(resp.status_code)
            raise UnhandledAPIError(response=resp, code=resp.status_code,
                                    message=message)
        except ValueError as msg:
            raise UnmarshallError(resp, msg)
        except UnhandledAPIError as uae:
            raise uae
        except Exception as e:
            raise UnmarshallError(resp, e.message)

    def UserIsMember(self, username, globalid, headers=None, query_params=None, content_type="application/json"):
        """
        Checks if the user has memberschip rights on the organization
        It is method for GET /organizations/{globalid}/users/ismember/{username}
        """
        uri = self.client.base_url + "/organizations/" + globalid + "/users/ismember/" + username
        resp = self.client.get(uri, None, headers, query_params, content_type)
        try:
            if resp.status_code == 200:
                return APIResponse(data=IsMember(resp.json()), response=resp)

            message = 'unknown status code={}'.format(resp.status_code)
            raise UnhandledAPIError(response=resp, code=resp.status_code,
                                    message=message)
        except ValueError as msg:
            raise UnmarshallError(resp, msg)
        except UnhandledAPIError as uae:
            raise uae
        except Exception as e:
            raise UnmarshallError(resp, e.message)

    def GetOrganizationUsers(self, globalid, headers=None, query_params=None, content_type="application/json"):
        """
        Get all users from this organization, not including suborganizations.
        It is method for GET /organizations/{globalid}/users
        """
        uri = self.client.base_url + "/organizations/" + globalid + "/users"
        resp = self.client.get(uri, None, headers, query_params, content_type)
        try:
            if resp.status_code == 200:
                return APIResponse(data=GetOrganizationUsersResponseBody(resp.json()), response=resp)

            message = 'unknown status code={}'.format(resp.status_code)
            raise UnhandledAPIError(response=resp, code=resp.status_code,
                                    message=message)
        except ValueError as msg:
            raise UnmarshallError(resp, msg)
        except UnhandledAPIError as uae:
            raise uae
        except Exception as e:
            raise UnmarshallError(resp, e.message)

    def DeleteOrganization(self, globalid, headers=None, query_params=None, content_type="application/json"):
        """
        Deletes an organization and all data linked to it (join-organization-invitations, oauth_access_tokens, oauth_clients, logo)
        It is method for DELETE /organizations/{globalid}
        """
        uri = self.client.base_url + "/organizations/" + globalid
        return self.client.delete(uri, None, headers, query_params, content_type)

    def GetOrganization(self, globalid, headers=None, query_params=None, content_type="application/json"):
        """
        Get organization info
        It is method for GET /organizations/{globalid}
        """
        uri = self.client.base_url + "/organizations/" + globalid
        resp = self.client.get(uri, None, headers, query_params, content_type)
        try:
            if resp.status_code == 200:
                return APIResponse(data=Organization(resp.json()), response=resp)

            message = 'unknown status code={}'.format(resp.status_code)
            raise UnhandledAPIError(response=resp, code=resp.status_code,
                                    message=message)
        except ValueError as msg:
            raise UnmarshallError(resp, msg)
        except UnhandledAPIError as uae:
            raise uae
        except Exception as e:
            raise UnmarshallError(resp, e.message)

    def CreateNewSubOrganization(
            self,
            data,
            globalid,
            headers=None,
            query_params=None,
            content_type="application/json"):
        """
        Create a new suborganization.
        It is method for POST /organizations/{globalid}
        """
        uri = self.client.base_url + "/organizations/" + globalid
        resp = self.client.post(uri, data, headers, query_params, content_type)
        try:
            if resp.status_code == 201:
                return APIResponse(data=Organization(resp.json()), response=resp)

            message = 'unknown status code={}'.format(resp.status_code)
            raise UnhandledAPIError(response=resp, code=resp.status_code,
                                    message=message)
        except ValueError as msg:
            raise UnmarshallError(resp, msg)
        except UnhandledAPIError as uae:
            raise uae
        except Exception as e:
            raise UnmarshallError(resp, e.message)

    def CreateNewOrganization(self, data, headers=None, query_params=None, content_type="application/json"):
        """
        Create a new organization. 1 user should be in the owners list. Validation is performed to check if the securityScheme allows management on this user.
        It is method for POST /organizations
        """
        uri = self.client.base_url + "/organizations"
        resp = self.client.post(uri, data, headers, query_params, content_type)
        try:
            if resp.status_code == 201:
                return APIResponse(data=Organization(resp.json()), response=resp)

            message = 'unknown status code={}'.format(resp.status_code)
            raise UnhandledAPIError(response=resp, code=resp.status_code,
                                    message=message)
        except ValueError as msg:
            raise UnmarshallError(resp, msg)
        except UnhandledAPIError as uae:
            raise uae
        except Exception as e:
            raise UnmarshallError(resp, e.message)

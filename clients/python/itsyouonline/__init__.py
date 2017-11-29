import requests


from .AddApiKeyReqBody import AddApiKeyReqBody
from .AddIncludeSubOrgsOfReqBody import AddIncludeSubOrgsOfReqBody
from .AddIncludeSubOrgsOfReqBody import AddIncludeSubOrgsOfReqBody
from .AddOrganizationMemberReqBody import AddOrganizationMemberReqBody
from .Address import Address
from .Authorization import Authorization
from .AuthorizationMap import AuthorizationMap
from .Avatar import Avatar
from .BankAccount import BankAccount
from .Company import Company
from .Contract import Contract
from .ContractSigningRequest import ContractSigningRequest
from .CreateUserBankAccountReqBody import CreateUserBankAccountReqBody
from .CreateUserBankAccountReqBody import CreateUserBankAccountReqBody
from .DigitalAssetAddress import DigitalAssetAddress
from .DigitalWalletAuthorization import DigitalWalletAuthorization
from .DnsAddress import DnsAddress
from .EmailAddress import EmailAddress
from .EnumJoinOrganizationInvitationMethod import EnumJoinOrganizationInvitationMethod
from .EnumJoinOrganizationInvitationRole import EnumJoinOrganizationInvitationRole
from .EnumJoinOrganizationInvitationStatus import EnumJoinOrganizationInvitationStatus
from .Error import Error
from .FacebookAccount import FacebookAccount
from .GetNotificationsReqBody import GetNotificationsReqBody
from .GetOrganizationUsersResponseBody import GetOrganizationUsersResponseBody
from .GetTotpsecretReqBody import GetTotpsecretReqBody
from .GetTwoFamethodsReqBody import GetTwoFamethodsReqBody
from .GetUserOrganizationsReqBody import GetUserOrganizationsReqBody
from .GithubAccount import GithubAccount
from .JoinOrganizationInvitation import JoinOrganizationInvitation
from .KeyData import KeyData
from .KeyStoreKey import KeyStoreKey
from .Label import Label
from .LocalizedInfoText import LocalizedInfoText
from .Member import Member
from .Membership import Membership
from .MissingScopes import MissingScopes
from .Organization import Organization
from .OrganizationAPIKey import OrganizationAPIKey
from .OrganizationTreeItem import OrganizationTreeItem
from .OrganizationUser import OrganizationUser
from .Ownerof import Ownerof
from .Party import Party
from .Phonenumber import Phonenumber
from .PublicKey import PublicKey
from .RegistryEntry import RegistryEntry
from .RequiredScope import RequiredScope
from .See import See
from .SeeVersion import SeeVersion
from .SeeView import SeeView
from .SetOrgMemberReqBody import SetOrgMemberReqBody
from .SetOrgOwnerReqBody import SetOrgOwnerReqBody
from .SetOrganizationLogoReqBody import SetOrganizationLogoReqBody
from .SetupTotpreqBody import SetupTotpreqBody
from .Signature import Signature
from .UpdateApikeyReqBody import UpdateApikeyReqBody
from .UpdateOrganizationApikeyReqBody import UpdateOrganizationApikeyReqBody
from .UpdateOrganizationOrgMemberShipReqBody import UpdateOrganizationOrgMemberShipReqBody
from .UpdatePasswordReqBody import UpdatePasswordReqBody
from .UpdateUserNameReqBody import UpdateUserNameReqBody
from .User import User
from .UserAPIKey import UserAPIKey
from .UserIsMemberReqBody import UserIsMemberReqBody
from .ValidatePhonenumberReqBody import ValidatePhonenumberReqBody
from .VerifyPhoneNumberReqBody import VerifyPhoneNumberReqBody
from .companyview import companyview
from .userview import userview

from .client import Client as APIClient

from .oauth2_client_oauth_2_0 import Oauth2ClientOauth_2_0

class Client:
    def __init__(self, base_uri="https://itsyou.online/api"):
        self.api = APIClient(base_uri)
        
        self.oauth2_client_oauth_2_0 = Oauth2ClientOauth_2_0()
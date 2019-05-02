## Customize the Authorization Code Flow

### Show an organization logo on the login/register screen

When you use the authorization code flow to authenticate your users using Itsyou.Online, you can provide a better user experience by showing your logo and customize the text on the login page.

Go to the settings page of an organization:

![Organization Settings](OrganizationSettingsTab.png)

where you can add an organization logo and change the text shown on the login page by modifying the organization description text.

When a user is asked to login, this logo and text are added to the login/register page:

![Branded login page](BrandedLoginPage.png)


### Choose a different default language

When an external site uses ItsYou.online using the authorization code flow, it can add the `lang` query parameter to change the default language if a user has not explicitly changed it. Possible values are `en` and `nl`, if no `lang` query parameter is supplied, English is taken.


### Choose to force 2-Factor authentication or not 

Two factor Authentication is forced by default.

![](./2faoptional.png)

- When enabled (Default behavior)
    - User will be forced to use `2-FA` to authenticate through this specific organization.
    - If user had not setup `2-FA` to their account, they will be asked to do so first time they authenticate through this specific organization

- When disabled (NOT RECOMMENDED)
    - It's up to `User Settings` who authenticate through this specific organization.
     ![](user2fa.png)

        - If user had enabled `2-FA` for their account, then `2-FA` will be used even if organization does not force it.
        - If user had enabled the option `Skip 2-FA for organizations when not required`, at then user will not be required to use `2-FA` with organizations that do not force `2-FA`

### Configuring the frequency of the 2FA challenge


When logging in to an external site using Itsyou.Online, a successful 2 factor authentication will gain a validity period, for which no further 2FA's are required. This 2FA validity is bound to the external site. As long as the user does not provide an invalid password, and the validity period hasn't expired, the 2FA step is not required for logging in. As soon as an invalid password is provided, the validity of the 2FA, if one is still active, is revoked. When no active validity for the user is detected, they will have to do the 2FA step, and will acquire a new validity period for their successful authentication. The default validity period duration is 7 days.

Currently, it is only possible to view or modify the validity period using the `organizations/{globalid}/2fa/validity` api. The validity period is expressed in seconds. The api supports both **GET** requests to retrieve the validity duration, and **PUT** requests to change the validity duration. Note that the validity period should be between 0 and 2678400 (31 days).

Example to retrieve and modify the validity period of an organization with globalid `mycompany`:

1. Inspect the validity duration
```
GET https://itsyou.online/api/organizations/mycompany/2fa/validity
```
The following information is returned in the response body:
```json
{
    "secondsvalidity":  604800
}
```
At this moment, the validity duration for a successful 2FA login is 604800 seconds (7 days, the default).

2. Change the validity duration
```
PUT https://itsyou.online/api/organization/mycompany/2fa/validity
```
In the body of the request, we specify the new duration, which we will set to 86400 (1 day).
```json
{
    "secondsvalidity": 86400
}
```
Also note that an access token will have to be specified, either by appending it to the request url, or by setting it in the Authorization header.

### Show the register screen instead of the login screen

If you are think the user has no account with ItsYou.online yet, you can supply the `prefer=register` queryparameter in the oauth flow. This will show the user the register screen instead of the login screen if we do do not detect a previous login (this is registered in the local storage).

### Make 2 Factor Authentication optional for organization

By default, any newly created organization has the option 
`Force Two Factor Authentication` set to true in `organization settings`
which means, that user is required have `2 Factor Authentication` set up 
before they can login throguh this organization.
For users who don't have `2 Factor Authentication` setup, they will be required to verify a phone number the 1st tmie they login with that organization.
For new users, they can't skip the `2 Factor authentication` step if they are trying to registering new user while using this organization.

If this option `Force Two Factor Authentication` is disabled, then `Two Factor Authentication` now is optional and up to user to use it or not.
Users can use `Two Factor Authentication` everywhere but in same time skip it for organizations that don not require it by activating this option in their account settings
`Skip 2 Factor Auth for organizations when not required`This option is not enabled by default and thus by default user is forced to use 2 Factor Auth for all organizations



## External site 2 factor authentication

When logging in to an external site using Itsyou.Online, a successful 2 factor authentication will gain a validity period, for which no further 2FA's are required. This 2FA validity is bound to the external site. As long as the user does not provide an invalid password, and the validity period hasn't expired, the 2FA step is not required for logging in. As soon as an invalid password is provided, the validity of the 2FA, if one is still active, is revoked. When no active validity for the user is detected, they will have to do the 2FA step, and will acquire a new validity period for their successful authentication. The default validity period duration is 7 days

### Configuring the validity period

Currently, it is only possible to view or modify the validity period using the `organizations/{globalid}/2fa/validity` api. The validity period is expressed in seconds. The api suports both **GET** requests to retrieve the validity duration, and **PUT** requests to change the validity duration. Note that the validity period should be anywhere between 0 and 2678400 (31 days).

Lets take a look at an example, where we will attempt to retrieve and modify the validity period of an organization with globalid `mycompany`.

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

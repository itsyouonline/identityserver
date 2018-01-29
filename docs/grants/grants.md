# grants

## Why

While being a member of an organization is a good way to handle long standing authorizations, it is bulky, and not intended to give individuals "rights". A person is added to an organization, and applications can limit access to people in this specific organization. The problem is that organization membership is, in a way, part of the user's identity, and he thus must first accept an invitation before being added to the organization. As a result, he must log in again before being granted access to the protected resource, and as such this does not work for apps with real time requirements.

## What

Grants aim to solve this problem by allowing oranizations to add custom defined values to an authorization from a user. Adding this authorization does not require any interaction on the part of the user (he does not need to accept the grant). Once given, a grant can be freely changed, listed, or removed by the organization, or even included in future access tokens or JWT's. This means that grants allow for single-side, real time access management.

Grants on access tokens and JWT are added as `grant:...`.

When adding grants, only characters alphanumerical characters, dashes and underscores are accepted (`a-Z, -, _, ., 0-9`).

## Remarks

Operations involving grants are only possible for users who have an existing authorization for your organization. If this is not the case, a `403 Forbidden` will be returned. Removal of an authorization (by a user) does not lead to removal of the grants, they will remain and will be present if a new authorization is given

In order to avoid doing excessive API calls, it is recommended to use grants in combination with a refresheable JWT. After giving a grant, the JWT can be refreshed and you can request that Itsyou.Online adds the grant to the JWT. This eliminates the need for a separate request to list the grants a user has.

Furthermore, note that grants are always added to the scopes on an access_token, but for a JWT the queryparameter `add_grants=true` needs to be added. This prevents an organization from accidentally passing a JWT with application defined grants to another third party.

## Example

Lets assume our application has a website, integrated with Itsyou.Online. On this website are 3 pages wich require sign in. The first page is accessible by anyone who has signed in through Itsyou.Online. On this page, a user can perform actions such as making a purchase in a web shop. 

After making a purchase, users should be able to access page 2 and 3, where they are offered services such as customer support, the ability to comment on the provided service etc...

Since authorization and authentication is provided by Itsyou.Online, it makes sense to also store who can access page 2 and 3 there. Without grants, the application would need to invite the user to the correct Itsyou.Online organization. The user would then need to login on Itsyou.Online, accept the invite, and log out and back in again on the applications website. But with grants, the application can simpy make an api call to Itsyou.Online, and give the person the `haspurchased` grant. The application can now do api calls to get all grants for this user, or inspect a new access token / JWT, and identifies that the user has gotten this grant (by presense of the `grant:haspurchased` scope), and thus has previously made a purchase. As a result, he can access the additional pages.


## Limitations

Currently, per user, an organization can only add a max of 50 grants. Also, each grant can only have a length of 100 bytes (characters) maximum .

As mentioned previously, only the characters `a-z`, `-`, `_`, `.` and `0-9` are allowed at the moment when creating a grant.


## API Endpoints

`GET /api/organizations/{globalid}/grants/{username}`: Lists all grants given to a user

`DELETE /api/organizations/{globalid}/grants/{username}`: Deletes all grants given to a user

`POST /api/organizations/{globalid}/grants/{username}`: Ensure this user has this grant

`PUT /api/organizations/{globalid}/grants/{username}`: Update a grant to a new one, if present

`DELETE /api/organizations/{globalid}/grants/{username}/{grant}`: Ensure this grant is not present on the user

`GET /organization/{globalid}/grants/havegrant/{grant}`: Returns a list of user identifiers for users with a given grant


Full documentation about these endpoints can be found in [the api docs](https://itsyou.online/apidocumentation)




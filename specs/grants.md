#Grants

## Problem

Currently the only wat for an organization or application to handle fine grained authorization or classification of users is to use suborganizations or store this classification itself. Neither of them are practical or a one fits all solution.
Suborganizations are a handy concept bu not really lightweight sine a user needs to accept membership for a right or authorization an app wants to giveHence the concept of  ** grants**. These are scopes an organization can stick on a user (usrname/email/phone nmber... and when u user authenticates using an oauth flow, the applications receives these grants as ** grant:name ** scopes, the users can olso be listed on a per grant basis thtough the api allowing a lightweight grouping and listing

## Todo: detailed spec to be later moved to the documentation

### High level description

Grants are application defined properties a third-party can stick on a user. Grants are only restricted in length and number. Besides that, an app can define any grant it wants for users. This allows the application more freedom to store small amounts of data. It is essentially a replacement for the current (ab)use of subgroups to manage access rights. While these subgroups work fine for actual members of organizations, an invite to a subgroup needs to be accepted first before being able to authorize the memberof scope for it, thus this system is not able to manage access rights in real time, and requires a client - server operation instead of only a server operation. 

In order to make the best use of grants, the app should workt with refresheable JWT's. When a grant gets added, the app can then simply call the JWT refresh endpoint. This will then return an updated JWT, including the new grant. In this way, a user can be given specific in app authorizations (e.g. folder access), without passing through ItsYou.online to accept an invite for a suborganization.

It should be noted that JWT's could be transfered to other servers (not under controll of the original authorized party).
Because you might not want to disclose the grants you use to these people, or because you want to reduce the size of the 
JWT which is send, grants are not automatically added to the JWT. To add them, the queryparameter `add_grants=true` must be 
provided when retrieving the JWT from Itsyou.Online. Adding grants can be done at any time, e.g. it is possible to only
add them after a JWT is refreshed (if they are not present in the original JWT for example) 

### Technical spec

- a grant is a scope of the form `grant:...`
- grants are stored on an authorization
- since grants are application defined, their length is limited (say 100 bytes to start off)
- to avoid spamming of the database, the amount of grants is limited (say 50 grants max per user to start)
- an app can not add a grant more than once to the same user
- although grants are actually just scopes, they can not be requested or authorized.
- existing functions to deal with authorizations and scopes must be checked to only verify `user:...` scopes and not grants
- a small set op api endpoints should be created to expose searching / sorting based on grants. To be defined

## Define api endpoints

Because these calls are only usefull for the organization and less so for the user (and they definetely should not be able to add
grants themselves), these endpoints are all under the `/api/organizations/{globalid}` endpoint.

1. GET `/api/organizations/{globalid}/grants/{username}`:

    Returns all the grants added to the authorization for the user

2. DELETE `/api/organizations/{globalid}/grants/{username}/{grant}`:

    Deletes the specified grant for the user

3. DELETE `/api/organizations/{globalid}/grants/{username}`:

    Deletes all grants for the user

4. POST: `/api/organizations/{globalid}/grants`:

    Creates a new grant for the user

    Body:

        user: string
        grant: string

5. PUT: `/api/organizations/{globalid}/grants`:

    Updates a grant to a new one for a user

    Body:

        user: string
        oldgrant: string
        newgrant: string

6. GET: `/api/organizations/{globalid}/grants/havegrant/{grant}`:

    Lists all users who have a specified tag added onto them.




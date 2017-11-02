#Grants

## Problem

Currently the only wat for an organization or application to handle fine grained authorization or classification of users is to use suborganizations or store this classification itself. Neither of them are practical or a one fits all solution.
Suborganizations are a handy concept bu not really lightweight sine a user needs to accept membership for a right or authorization an app wants to giveHence the concept of  ** grants**. These are scopes an organization can stick on a user (usrname/email/phone nmber... and when u user authenticates using an oauth flow, the applications receives these grants as ** grant:name ** scopes, the users can olso be listed on a per grant basis thtough the api allowing a lightweight grouping and listing

## Todo: detailed spec to be later moved to the documentation

### High level description

Grants are application defined properties a third-party can stick on a user. Obviously, grants can only be added for users who are sharing information with the third-party (even if we allow grants to be added for arbitrary users, they won't be available unless info is shared anyway). Grants are only restricted in length and number. Besides that, an app can define any grant it wants for users. This allows the application more freedom to store small amounts of data.

In order to make the best use of grants, the app should workt with refresheable JWT's. When a grant gets added, the app can then simply call the JWT refresh endpoint. This will then return an updated JWT, including the new grant. In this way, a user can be given specific in app authorizations (e.g. folder access), without passing through ItsYou.online to accept an invite for a suborganization. 

### Techincal spec

- a grant is a scope of the form `grant:...`
- grants are stored on an authorization
- since grants are application defined, their length is limited (say 100 bytes to start off)
- to avoid spamming of the database, the amount of grants is limited (say 50 grants max per user to start)
- an app can not add a grant more than once to the same user
- although grants are actually just scopes, they can not be requested or authorized.
- existing functions to deal with authorizations and scopes must be checked to only verify `user:...` scopes and not grants
- a small set op api endpoints should be created to expose searching / sorting based on grants. To be defined

### TODO: Define api endpoints

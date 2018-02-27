# OIDC (OpenID Connect) support

OpenID Connect 1.0 is a simple identity layer on top of the OAuth 2.0 protocol. It allows Clients to verify the identity of the End-User based on the authentication performed by an Authorization Server, as well as to obtain basic profile information about the End-User in an interoperable and REST-like manner.

From the `access token` request an additional token is returned with the `access token` namely an `ID token`.
The `ID token` is returned in the `id_token` field of the JSON object returned from request.

To get an OIDC response from the `access token` request, the oidc scope: `openid`, should be provided to the scope in the `authorization` request, made before the `access token` request in the OAuth 2 flow.

## ID Token

The `ID token` resembles the concept of an identity card, in a standard JWT format, signed by the Identity Provider (itsyou.online).

A standard ID token has the following claims.

* sub: [required] subject: (unique) identity of the user
* iss: [required] issuer: the issuing authority
* aud: [required] The intended recipient of the token (currently set to: `ALL_AUDIENCES`)
* iat: [required] Issued at timestamp
* exp: [required] Expiration timestamp
* nonce: [optional]
* auth_time: [optional] time of authentication
* acr: [optional] [Authentication Context Class Reference](http://openid.net/specs/openid-connect-eap-acr-values-1_0.html#acrValues)

### Additional Claims

Additional claims can be added to the ID token.
OpenID Connect specifies a set of standard claims, or user attributes. They are intended to supply the client app with consented user details such as email, name and picture, upon request.

The identityserver's implementation supports the following claims:
|Scope Value|Associated claims|
|-|-|-|
|email| email, email_verified|
|profile|name, family_name, given_name|
|phone|phone_number, phone_number_verified|

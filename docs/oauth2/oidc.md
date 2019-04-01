# OIDC (OpenID Connect) support

[OpenID Connect 1.0](https://openid.net/specs/openid-connect-core-1_0.html) is a simple identity layer on top of the OAuth 2.0 protocol. It allows Clients to verify the identity of the End-User based on the authentication performed by an Authorization Server, as well as to obtain basic profile information about the End-User in an interoperable and REST-like manner.

From the `access token` request an additional token is returned with the `access token` namely an `ID token`.
The `ID token` is returned in the `id_token` field of the JSON object returned from request.

To get an OIDC response from the `access token` request, the oidc scope: `openid`, should be provided to the scope in the `authorization` request, made before the `access token` request in the OAuth 2 flow.

## ID Token

The `ID token` resembles the concept of an identity card, in a standard JWT format, signed by the Identity Provider (itsyou.online).

A standard ID token has the following claims.

* sub: [required] subject: (unique) identity of the user
* iss: [required] issuer: the issuing authority
* aud: [required] The intended recipient of the token (ClientID of Relying Party)
* iat: [required] Issued at timestamp
* exp: [required] Expiration timestamp
* nonce: [optional]
* auth_time: [optional] time of authentication
* acr: [optional] [Authentication Context Class Reference](http://openid.net/specs/openid-connect-eap-acr-values-1_0.html#acrValues)

### Additional Claims

Additional claims can be added to the ID token.
OpenID Connect specifies an optional set of standard claims, or user attributes. They are not supported in favor of using the existing scopes, similar to the [JWT](jwt.md#Storing-the-actual-values-of-scopes-in-JWT) implementation of adding values of scopes to the JWT.

The oidc implementation supports the following scopes with their respective claims:  

| Scope value | Associated claims |
|-|-|
| user:name | user:name |
|user:email[:label]|user:email[:label]|
|user:validated:email[:label]| user:validated:email[:label]|
|user:phone[:label]|user:phone[:label]|
|user:validated:phone[:label]|user:validated:phone[:label]|

Example of returned ID token:
```js
// requested scopes: "openid,user:name,user:validated:email,user:validated:phone"
{
  "aud": "test-organization",
  "exp": 1520000016,
  "iat": 1519913616,
  "iss": "itsyouonline",
  "sub": "user_name_1",
  "user:name": "user name",
  "user:validated:email": "user.name@example.com",
  "user:validated:phone": "+32123456789"
}
```

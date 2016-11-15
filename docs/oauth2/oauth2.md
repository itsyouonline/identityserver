## Authorization grant types

OAuth2 defines four grant types, each of which is useful in different cases:

1. Authorization Code: used with server-side Applications
2. Implicit: used with Mobile Apps or Web Applications (applications that run on the user's device)
3. Resource Owner Password Credentials: used with trusted Applications, such as those owned by the service itself
4. Client Credentials: used with Applications API access

Currently the **authorization code** and **client credentials** grant types are supported.


## Authorization Code Flow
The authorization code grant type is the most commonly used because it is optimized for server-side applications where Client Secret confidentiality can be maintained. This is a redirection-based flow, which means that the application must be capable of interacting with the user-agent (i.e. the user's web browser) and receiving API authorization codes that are routed through the user-agent.

### Prerequisite: clientid and client secret

In order to acquire an oauth access token, a client id and client secret are required.

In itsyou.online, organizations map to clients in the oauth2 terminology and the organization's globalid is used as the clientid. Client secrets can be created through the UI or through the `organizations/{globalid}/apikeys` api.

![AuthorizationCodeFlow](https://rawgit.com/itsyouonline/identityserver/master/docs/oauth2/AuthorizationCodeFlow.svg)

### Step 1: Authorization Code Link

First, the user is given an authorization code link that looks like the following:

```
https://itsyou.online/v1/oauth/authorize?response_type=code&client_id=CLIENT_ID&redirect_uri=CALLBACK_URL&scope=user:name&state=STATE
```

* https://itsyou.online/v1/oauth/authorize: the API authorization endpoint
* client_id=client_id

    the application's client ID
* redirect_uri=CALLBACK_URL

    The redirect_uri parameter is required. The redirect URL's host and port must exactly match the callback URL and the redirect URL's path must reference a subdirectory of the callback URL.
    The redirect_uri *must* start with a scheme indicator (`scheme://`).


* response_type=code

    specifies that your application is requesting an authorization code grant
* scope=user:name

    specifies the level of access that the application is requesting

    here we specify the scope "user:name", see the [available scopes](availableScopes.md) for all other supported scopes.


* state=STATE

    A random string. It is used to protect against csrf attacks.

### Step 2: User Authorizes Application

When the user clicks the link, they must first log in to the service, to authenticate their identity (unless they are already logged in). Then they will be prompted by the service to authorize or deny the application access to the requested information.

### Step 3: Application Receives Authorization Code

After the user authorizes the application some of it's information, itsyou.online redirects the user-agent to the application redirect URI, which was specified during the client registration, along with an authorization code and a state parameter passed in step 1. If the state parameters don't match, the request has been created by a third party and the process should be aborted.
The redirect would look something like this (assuming the application is "petshop.com"):

```
https://petshop.com/callback?code=AUTHORIZATION_CODE&state=STATE
```

### Step 4: Application Requests Access Token

The application requests an access token from the API, by passing the authorization code along with authentication details, including the client secret, to the API token endpoint and the state. Here is an example POST request to itsyou.online's token endpoint:

```
POST https://itsyou.online/v1/oauth/access_token?client_id=CLIENT_ID&client_secret=CLIENT_SECRET&code=AUTHORIZATION_CODE&redirect_uri=CALLBACK_URL&state=STATE
```

The redirect_uri must match the redirect_uri passed in the access_code request and the callback URI registered in the api key. The redirect URL's host and port must exactly match the callback URL and the redirect URL's path must reference a subdirectory of the callback URL. The state must match the state received with the authorization code

* response_type=code

### Step 5: Application Receives Access Token

If the authorization is valid, the API will send a response containing the access token (and optionally, a refresh token) to the application. The entire response will look something like this:

```
{"access_token":"ACCESS_TOKEN","token_type":"bearer","expires_in":86400,"refresh_token":"REFRESH_TOKEN","scope":"read","info":{"username":"bob"}}
```
Now the application is authorized.
It may use the token to access the user's account via the service API, limited to the scope of access, until the token expires or is revoked.
If a refresh token was issued, it may be used to request new access tokens if the original token has expired.


### Use the access token to access the API

The access token allows you to make requests to the API on a behalf of a user.

```
GET https://itsyou.online/api/users/bob/info?access_token=...
```
You can pass the token in the query params like shown above, but a cleaner approach is to include it in the Authorization header

```
Authorization: token OAUTH-TOKEN
```
For example, in curl you can set the Authorization header like this:

```
curl -H "Authorization: token OAUTH-TOKEN" https://itsyou.online/api/users/bob/info
```

### Customize the user experience

Small customizations can be configured such as an organization logo and 2 factor authentication validity.

Details are described in the [Customize Authorization Code Flow documentation](CustomizeAuthorizationCodeFlow.md).

## Client Credentials Flow


The client credentials grant type can be used in two scenario's:
1. An application linked to an **organization** to access its own account.
2. An application linked to a **user** to access the user's account

Examples of when this might be useful include if an application wants to invite someone to an organization or access other organization data using the API.

### Prerequisite: client_id and client_secret

In order to acquire an oauth access token, a client id and client secret are required.

#### Organization api key
In itsyou.online, organizations map to clients in the oauth2 terminology and the organization's globalid is used as the client_id. Client secrets can be created through the UI by going to the organization's settings page or through the `api/organizations/{globalid}/apikeys` api.

In order to use an apikey in a client credentials flow, *enable client credentials flow* must be set on the apikey (through the api or in the apikey detail dialog in the UI).

#### User api key
It is possible to create api keys to access a user's data through the api instead of through the UI.
In the UI, create a new api key in the user's settings page. This will generate an *application id* and a *secret*. Use the *application id* as client_id and the *secret* as client_secret.

Api keys can also be created through the `api/users/{username}/apikeys` api or by creating an user through the api ( POST on `/api/users`)


### Acquire an access token

The application requests an access token by sending its credentials, its client_id and client secret, to the authorization server. An example POST request might look like the following:

```
https://itsyou.online/v1/oauth/access_token?grant_type=client_credentials&client_id=CLIENT_ID&client_secret=CLIENT_SECRET
```

* client_id=CLIENT_ID

    the organization's globalid or the application id as described above

* client_secret=CLIENT_SECRET

    the api key secret

* grant_type=client_credentials

    specifies that your application is requesting an access token using the client credentials flow



If the application credentials check out, the authorization server returns an access token to the application.


### Use the access token to access the API

The access token allows you to make requests to the API like described in the authorization code grant type above. When an organization api key is used, the requests are on behalf of the organization instead of on behalf of a user.

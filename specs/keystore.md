# IYO Keystore

Every user has a private keystore linked to his account, where organizations can save public keys.
The keystore is not exposed to the user, except for when authorization is given to read and push to a keystore.
Keystores are linked to both the user and the organization: an organization can only read and write to its own keystore, which is seperate for every organization

To gain access to a user's keystore, an oauth authorization code grant flow login is required, where the organization requests the `user:keystore` scope.
This scope gives both read and write access to the keystore.

With this scope, 3 API endpoints are made available:
 - `GET` `/users/{username}/keystore`: gets all the keys of saved in the keystore of the user identified by `{username}`
 - `GET` `/users/{username}/keystore/{label}`: gets the key saved with `{label}` in the keystore of the user identified by `{username}`
 - `POST` `/users/{username}/keystore`: saves a key in the keystore of the user identified by `{username}`
 
 A key body or response is defined as:
 
 ```golang
 type KeyStoreKey struct { 
  Key      string  `json:"key"` 
  Globalid string  `json:"globalid"` 
  Username string  `json:"username"` 
  Label    string  `json:"label"` 
  KeyData  KeyData `json:"keydata"` 
} 

type KeyData struct { 
  TimeStamp db.DateTime `json:"timestamp"` 
  Comment   string      `json:"comment"` 
  Algorithm string      `json:"algorithm"` 
} 
```

Globalid, Username and KeyData.TimeStamp are optional and overwritten by the server when the save API is called.
Comment is optional.
The combination `Globalid`, `Username`, `Label` must be unique


 Keys saved in the keystore can not be deleted or modified after being uploaded.

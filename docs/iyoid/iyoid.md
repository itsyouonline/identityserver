# IyoID

Third party applications, who rely on Itsyou.Online to handle the authorization and authentication of users, can also store and/or proces data themselves. While users are identified by a user identifier (a validated email address or phone number), this could leak data to unauthorized personel. To solve this, the concept of an IyoID is introduced.

Essentially, IyoID's are just random strings. An application can request such an ID from Itsyou.Online. In turn, a mapping is created between the authorized party (`azp`), and the user. the authorized party is the `clientid` of an api key: either an organization, or a personal api key.

Only the party identified by the `azp` can then ask Itsyou.Online to return a valid identifier for this party. In this way, information about a user can be stored e.g. on a blockchain. The organization who first stored the info can then find out who is represented by the IyoID.

It should be noted that the amount of ID's which can be requested by the application for a certain user is limited (currently this limit is 25). Also note that these ID's are only deleted in case the organization itself is removed, to prevent them from being leaked.

## API Endpoints:

`GET /api/users/{username}/identifiers`: Lists all generated identifiers for this user. Note that this only lists identifiers with the callers clientID

`POST /api/users/{username}/identifiers`: Create a new identifier for this user, only accessible with the callers clientid

`GET /api/users/identifiers/{identifier}`: Look the user behind an identifier

Complete api info can be found in [the api docs](https://itsyou.online/apidocumentation)
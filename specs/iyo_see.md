# iyo see

Links to documents that can be signed. (e.g. TOC) 

First version:
- UI:
  - access via `itsyou.online/#/organization/{globalid}/see`
  - all different items, sorted by date
  - by default only latest version is shown
  - full history can be viewed on the detail page of the document (`itsyou.online/#/organization/{globalid}/see/{uniqueid}/detail`)
  - checkmark if version is signed
  - short description rendered by default - user can expand this to full description
  
- API functions:
  - To enable the organization to read from and write to iyo see, the `user:see` scope must be obtained
  - A user can only read the `see` objects stored by organizations
  - The oranization can upload new `see` objects
  - If a `see` object is updated, a new version will be created, the old one is kept
  - `See` items can't be changed after being uploaded, only a new version can be uploaded
  - The items can only be signed after being uploaded without creating a new version
  
- API endpoints:
  - `GET` `/users/{username}/see/{globalid}`: Lists all see items for this user from this organization
  - `POST` `/users/{username}/see/{globalid}`: Create a new see object (errors if a duplicate user - organization - uniqueid exists)
  - `PUT` `/users/{username}/see/{globalid}/{uniqueid}`: Updates a see object by creating a new version (errors if the user - organization - uniqueid combination does not exist)
  - `GET` `/users/{username}/see/{globalid}/{uniqueid}`: Gets the detail of a single see item. Query parameter `version`:
    - `all`: returns the full version history of the item
    - `version_index`: an integer pointing to a verion in the history of the item. Not found if index does not exist
    - `latest` (default): only return the latest version of the item
  - `PUT` `/users/{username}/see/{globalid}/{uniqueid}/sign/{version}`: Updates the signature of a see item. Does not create a new version in the history.  Can only be called once (returns conflict status error on successive calls). 
  
The body of the see API endpoints is defined as:

```golang
type SeeView struct {
	Username                 string       `json:"username"`
	Globalid                 string       `json:"globalid"`
	Uniqueid                 string       `json:"uniqueid"`
	Version                  int          `json:"version"`
	Category                 string       `json:"category"`
	Link                     string       `json:"link"`
	ContentType              string       `json:"content_type"`
	MarkdownShortDescription string       `json:"markdown_short_description"`
	MarkdownFullDescription  string       `json:"markdown_full_description"`
	CreationDate             *db.DateTime `json:"creation_date"`
	StartDate                *db.DateTime `json:"start_date,omitempty"`
	EndDate                  *db.DateTime `json:"end_date,omitempty"`
	KeyStoreLabel            string       `json:"keystore_label"`
	Signature                string       `json:"signature"`
}
```

When retrieving the details of a single see item (only for `GET` `/users/{username}/see/{globalid}/{uniqueid}`):

```golang
type See struct {
	Username string        `json:"username"`
	Globalid string        `json:"globalid"`
	Uniqueid string        `json:"uniqueid"`
	Versions []SeeVersion  `json:"versions"`
}

type SeeVersion struct {
	Version                  int          `json:"version" bson:"-"`
	Category                 string       `json:"category"`
	Link                     string       `json:"link"`
	ContentType              string       `json:"content_type"`
	MarkdownShortDescription string       `json:"markdown_short_description"`
	MarkdownFullDescription  string       `json:"markdown_full_description"`
	CreationDate             *db.DateTime `json:"creation_date"`
	StartDate                *db.DateTime `json:"start_date,omitempty" bson:"startdate,omitempty"`
	EndDate                  *db.DateTime `json:"end_date,omitempty" bson:"enddate,omitempty"`
	KeyStoreLabel            string       `json:"keystore_label"`
	Signature                string       `json:"signature"`
}
```


Username, Globalid, CreationDate are overwritten when the save api is called.
KeyStoreLabel points to the label of the key stored in the keystore. 


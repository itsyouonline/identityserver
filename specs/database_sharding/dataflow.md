# Normal data flow

The database, represented to the API by the `mongos` instance, receives a `bson document`
to be stored in the database. First the `mongos` determines whether this document
is in a sharded collection or not. In case it is an unsharded collection, the document
is send to the primary shard, where it is stored. If the collection is sharded, the
`mongos` inspects the document and looks for the shard key. It then sends the document
to the shard specified by the sharding rules for the shard key, where it is stored.

In short, the only place where any kind of persistent storage of the actual data
(aside from the configuration of the database) occurs is on the shards.

Example:

A user registers on `ItsYou.online`. After filling in the required fields and submitting,
the API constructs `bson documents` containing the provided data. Furthermore, from
the information provided, the API has identified this user as a European user (Criteria
to be determined). On the documents stored in sharded collections, the API adds the right
shard key (e.g.: `Country: EU`). The API then transmits the documents to a database interface
and instructs these documents to be saved. Transparent to the API, the documents are send
to the `mongos` instance. It will send any documents for unsharded collections to the
primary shard. All other documents are inspected for the shard key and send to their
respective shards for storage. In th case of this example, the shard key (`Country`) will have the
`EU` value, thus the `mongos` will send the documents to shards tagged as `EU` (shard tags
are maintained in the config). Ultimately, the documents for the sharded collection from our
small example are only stored on the European set, for it is the only shard tagged as `EU`  

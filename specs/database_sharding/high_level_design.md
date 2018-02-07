# Sharded database design

Currently the database runs as a single replica set. Each member of the set will
be assigned the role of shard server, thereby turning the set into the first shard.
Next to this, a config replica set is to be deployed. This set will store the
configuration of the database cluster and must be reachable by every `mongos` instance
and every shard.

Next to the front end, a single `mongos` instance is deployed. This will act as the
new database access point. A `mongos` must connect to the config set and every shard.
The existing shard (created from the existing database replica set) is elected as the
primary shard: this shard will hold, next to the parts of the sharded collections that
will stay there, all unsharded collections in their entirety. Any other shards only
contain the data from the sharded collections, that has been designated to be stored
on the respective shards by the shard key.

The shard key is an additional field that will be introduced on the sharded collections.
Its value will match the tag assigned to the shard where this data object needs to be stored.

Additional shards will be deployed as 3-member replica sets.
Once the shard is connected to the cluster, the `mongos` balancers are enabled
to automatically redistribute the data from the sharded collections. The `mongos` instances
also take care of directing new writes to the correct shard, as per the shard key.

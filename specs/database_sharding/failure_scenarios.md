# Failure scenarios

A small overview of some possible (but not necessarily all) failure scenarios. Note
that only those directly involving the database are considered, as other components
are beyond the scope of this document.

Failures of the setup have different results depending on which components they effect:

- `Mongos`: As the mongos instance(s) is (are) the connection endpoints for the front end to
  the database, unavailability means that the connected front end cannot reach the database
  anymore.


- `Primary shard`: The primary shard is configured as a 3 member replica set, therefore
we separate multiple scenarios:
  - A single database server within the replica set fails or is otherwise unreachable.
  As long as the other 2 database servers remain online, the replica set and therefore the
  shard remains intact. When the downed node comes back online, it is reconnected in the
  set, and moved in a recover state until it has finished syncing all data with the set primary.
  Should it be the set primary that goes offline, the other 2 nodes will start a primary
  election to appoint a new primary. Once the old primary comes back online, it syncs with
  the new primary and continues as a secondary thereafter until it is possibly appointed
  as the primary again once a new primary election is triggered.

  - Two database servers fail: As there is no longer a strict majority of the set
  members online and able to participate in a primary election, this is considered
  a total failure. See all three members fail.

  - All three members fail: This is a total failure of the set and therefore the shard.
  If api's are called which read or write unsharded collections, internal server
  error is returned. Should data be read only from sharded collections, the api
  will hang indefinitely. Write operations to sharded collections still go through,
  but the api will likely return a status 500 or hang as api's generally perform reads first
  to verify data integrity.

- `Secondary shards`: Like the primary, these are deployed as a three member replica
set, so we have the same scenarios:
  - Single server failure: The replica and thus the shard continue operations as
  expected, see single server failure in the primary shard.

  - Two server failure: Total set and thus shard failure.

  - Three server failure (total set/shard failure): If a read operation is performed on
  a sharded collection, it hangs indefinitely. Write operations to sharded collections
  can go though to the non failing shards, as the shard key is known and thus the
  `mongos` instances can perform a targeted shard access. Reads and writes from unsharded
  collections work as normal, as the system understands that these collections aren't located
  here.

- `Config set`: The config set persists all the configurations for the database. It must
be reachable by all shards and every `mongos` instance.
  - Single server failure: The set continues operations as expected, see above

  - Two server failure: Total set and thus shard failure.

  - Three server failure (total set/shard failure): As the config set holds all the
  metadata about the sharded collections, reads and writes to those are unavailable
  and will (eventually) result in an internal server error. Data can still be read from
  and written to unsharded collections.

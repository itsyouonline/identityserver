# Database sharding

This directory will contain the markdown documents for the database sharding so
they can easily be reffered to. The database backend is `Mongodb`. Next to the
architecture of the database itself, some backend changes will also be required.
Most notably, the collections which will be sharded will require an additional field
to be used as the shard key.

For a high level overview of the proposed changes, see
[high_level_design.md](high_level_design.md)

Normal operation data flow: [dataflow.md](dataflow.md)

List of failure scenarios related to a sharded db setup:
[failure_scenarios.md](failure_scenarios.md)


To deploy an example setup, [a shell script is provided](deploy_shards.sh). This
script deploys: a 3 member config replica set, a 3 member primary shard replica set,
and a 3 member secondary shard replica set. It also deploys a `mongos` instance, with
port `27017` exposed on localhost, linked to the created shards. If a `dump` directory
is present in the same directory as the script file, it will try to restore the dump.
The script then modifies the database and sets up the user collection for sharding.
To remove the dockers, run the script with `clean` as first argument
 (`./deploy_shards clean`). WARNING: if the script loads a dump into the database,
 you will be left with some big dangling volumes after remvoing the containers.
 These can be cleaned with `docker volume prune`.

 The containers can alo be started or stopped without removing by passing in `start`
 or `stop` as the first argument. To connect to the database, make sure to comment
 line 25 in [db/user/db.go](../../db/user/db.go), as the shard key won't allow for
 another unique index

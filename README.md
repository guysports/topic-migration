# Event Streams Topic Migration

This is an example of how to migrate topic definitions programmatically from a source event streams service instance to a target event streams service instance

## Build 
Build the binary by issuing a `go build` command

## Execute
1. Set the `SOURCE_API_KEY` and `TARGET_API_KEY` environment variables to the service credentials api key of the service instances having the topics details migrated from and to.
2. Ensure the ServiceID of the target cluster has a `Manager` role policy
```
bx iam service-policy-create <serviceID> --roles Manager --service-name messagehub --service-instance <service_instance_name>
```
3. Run the command to migrate (note the --dry-run command will just simulate the migration, omit to actually perform the migration)
```
./topic-migration migrate --source-broker-list={"kafka01:9093","kafka02:9093"} --target-broker-list={"broker-0:9093","broker-1:9093"} --dry-run
```
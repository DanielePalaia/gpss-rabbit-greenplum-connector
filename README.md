This is a simple gpss client for greenplum:

https://gpdb.docs.pivotal.io/5160/greenplum-stream/api/dev_client.html

it takes a .csv file a schema definition and ingest it inside a greenplum table using the gpss protocol.

a gpss server must be run before using this cient.

The software uses grpc.

gprc contract definition can be found in proto/gpss.proto

From this file it's possible to automatic generate code (in this case Go code) with this 
command:

protoc --go_out=plugins=grpc:. *.proto

Code is already generated for you and can be found in /gpss/gpss.pb.go

Code must be compiled before execution using the Go compiler:

go build client.go

After this the properties file needs to be filled with this explicative info:

* GpssAddress=10.91.51.23:50007     
* GreenplumAddress=10.91.51.23
* GreenplumPort=5533
* GreenplumUser=gpadmin
* GreenplumPasswd=****
* Database=test
* SchemaName=public
* TableName=mytest
* columnfile=/Users/dpalaia/GO/src/gpss-client/test/columns.csv
* datafile=/Users/dpalaia/GO/src/gpss-client/test/data.csv


columnfile must contain the schema table of TableName with all fields separated by coma ex (col1,col2,col3,col4)
datafile contains the .csv file to ingest.





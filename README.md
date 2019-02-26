This software is intended to be a simple connector rabbitmq-greenplum, similar to the gpsscli which is supporting kafka.

It is based on gpss (greenplum streaming server) so will work just with greenplum 5.16 or above.
https://gpdb.docs.pivotal.io/5160/greenplum-stream/overview.html

The connector will attach to a rabbitmq queue specified at configuration time will batch a certain amount of elements specified and will ask the gpss server to push them on a greenplum table.

For the moment the connector is supporting just json format (elements must be strings of json in rabbitmq and the resulting greenplum table needs to be a table with a json field).

These are the steps to run the software:

Prerequisites:

1) Activate the gpss extension on the greenplum database you want to use (for example test)
   test=# CREATE EXTENSION gpss;
   
2) create a table inside this database with a json field on it (for example mytest3)
   test=# create table mytest3(data json);
   
3) Run a gpss server with the right configuration (ex):
  gpss ./gpsscfg1.json --log-dir ./gpsslogs
  where gpsscfg1.json is 
  {
    "ListenAddress": {
        "Host": "",
        "Port": 50007,
        "SSL": false
    },
    "Gpfdist": {
        "Host": "",
        "Port": 8086
    }
}

4) download, install and run a rabbitmq broker
./rabbitmq-server

5) Create a rabbitmq transient queue with the rabbitmq UI interface you want the connector to connect (es gpss):
  ![Screenshot](queue.png)
  
Running the application:

1) The application is written in GO. If you are using MacOs then you can directly use the binary version inside /bin of this project called: gpss-rabbit-greenplum-connect otherwise you must compile it with the GO compiler

2) Use the file properties.ini (that should be place in the same directory of the binary in order to instruct the program with this properties

GpssAddress=10.91.51.23:50007
GreenplumAddress=10.91.51.23
GreenplumPort=5533
GreenplumUser=gpadmin
GreenplumPasswd=**** 
Database=test
SchemaName=public
TableName=mytest3
rabbit=amqp://guest:guest@localhost:5672/
queue=gpss
batch=50000  

queue is the rabbitmq queue name while batch is the amount of batching that the rabbit-greenplum connector must take before pushing the data into greenplum.

3) Run the connector:
./gpss-rabbit-greenplum-connect 
  ![Screenshot](connector.png)

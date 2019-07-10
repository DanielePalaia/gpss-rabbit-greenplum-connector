# Summary
This software is intended to be a simple (non production ready) connector rabbitmq-greenplum using Greenplum Streaming Service functionalities. </br>
It is written in GO and it uses the following technologies: RabbitMQ, GO, GRPC, Greenplum GPSS. </br>

The following reading can help you to better understand the software:

**RabbitMQ:** </br>
https://www.rabbitmq.com/ </br>
**GRPC:**  </br>
https://grpc.io/ </br>
**Greenplum GPSS:**</br>
https://gpdb.docs.pivotal.io/5160/greenplum-stream/overview.html</br>
https://gpdb.docs.pivotal.io/5160/greenplum-stream/api/dev_client.html</br>

![Screenshot](./pics/image2.png)

The connector will attach to a rabbitmq queue specified at configuration time will then batch a certain amount of elements specified and finally will ask the gpss server to push them on a greenplum table. </br>

Also it supports persistency so in case the connector is stopped before it reaches the minimum batch elements to send a request to GPSS to ingest items, when restarted it will load again the lost items. </br>

These are the steps to run the software:

## Prerequisites

1. **Activate the gpss extension on the greenplum database you want to use (for example test)**
   
      ```
      test=# CREATE EXTENSION gpss;
      ```
   
2. **Create the Greenplum table to be ingested**
   
      ```
      test=# create table companies(id varchar 200, city varchar 200, foundation timestamp, description text, data json);
      ```

   ![Screenshot](./pics/definition.png)
   
3. **Run a gpss server with the right configuration (ex):**
  
      gpss ./gpsscfg1.json --log-dir ./gpsslogs
      where gpsscfg1.json 
  
      ```
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
      ```

4. **download, install and run a rabbitmq broker**

      ./rabbitmq-server

5. **Create a rabbitmq durable queue with the rabbitmq UI interface you want the connector to connect (es gpss)**

  ![Screenshot](./pics/queue.png)<br/>
  
## Running the application

1. **Find binaries** 
      The application is written in GO. If you are using MacOs then you can directly use the binary version inside ./bin/osx       of this project called: gpss-rabbit-greenplum-connect otherwise you must compile it with the GO compiler<br/>

2. **Setting property file**    
      Use the file properties.ini (that should be place in the same directory of the binary in order to instruct the program        with this properties
      
         ```
         GpssAddress=10.91.51.23:50007
         GreenplumAddress=10.91.51.23
         GreenplumPort=5533
         GreenplumUser=gpadmin
         GreenplumPasswd= 
         Database=test
         SchemaName=public
         TableName=mytest3
         rabbit=amqp://guest:guest@localhost:5672/
         queue=gpss
         batch=50000 
         mode=1
         ```
         
      queue is the rabbitmq queue name while batch is the amount of messages that the rabbit-greenplum connector must             receive     before pushing the data into greenplum.<br/>
      If mode is set to 1 the items batched will be saved on a disk file so in case of crash or network issue at the next         restart the connector will be automatically able to recover this info again<br/>

3. **Run the connector**
```
./gpss-rabbit-greenplum-connect 
Danieles-MBP:bin dpalaia$ ./gpss-rabbit-greenplum-connector 
connecting to grpc server
connected
2019/02/26 17:01:30  [*] Waiting for messages. To exit press CTRL+C
```

4. **Populate the queue with the UI interface (Publish command) **
![Screenshot](./pics/queue3.png)

Every line correspond to the respective table field.

5. **Insert elements as specified in batches property** 
      Once you publish more messages than the batch value you should then see the table populated and you can restart             publishing.<br/>

6. **Try producer client**
      In order to make tests easy I also developed a simple consumer inside rabbit-client, you can find a binary for macos         always inside bin.<br/>
      If you run<br/>
      ./rabbit-client<br/>
      
he will take the same configuration that is inside properties.ini and will start to fire messages inside the same queue.

## Unit testing
A functional test is provided, it takes the parameters specified in ./properties insert batch elements inside the queue specified and then is checking that these elements have been inserted correctly.
To work properly it needs the table and the rabbitmq queue to be initially empty.
The table should be the same as the example and so:

```
  test=# create table companies(id varchar 200, city varchar 200, foundation timestamp, description text, data json);
```

  Then you can just go test -v ./... to let the test start
  
## Compiling and Installing the application </br> 

The application is written in GO. Binary for MacosX and Linux are already provided inside the /bin folder. <br/>
If you need to compile and install it you need to download a GO compiler (ex for Linux - ubuntu) </br>

1. sudo apt-get install golang-go <br>
2. export GOPATH=/home/user/GO <br>
3. create a directory src inside GO and go there </br>
4. git clone https://github.com/DanielePalaia/gpss-rabbit-greenplum-connector and enter the project</br>
5. go get github.com/golang/protobuf/proto </br>
   go get github.com/streadway/amqp </br>
   go get google.golang.org/grpc </br>
   cp -fR ./gpss /home/user/GO/src/gpssclient </br>
6. go install gss-rabbit-greenplum-connector and you will find your binary in GOPATH/bin </br> </br>

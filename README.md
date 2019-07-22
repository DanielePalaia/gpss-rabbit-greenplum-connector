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
      test=# create table companies(id varchar (200), city varchar (200), foundation timestamp, description text, data json);
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

      https://www.rabbitmq.com/download.html </br>
      then run the broker
      ./rabbitmq-server

5. **Create a rabbitmq durable queue with the rabbitmq UI interface you want the connector to connect (es gpss)**

  ![Screenshot](./pics/queue.png)<br/>
  
## Running the application

1. **Find binaries** 
      The application is written in GO. If you are using MacOs or Linux64 then you can directly use the binary version inside ./bin/osx and ./bin/linux of this project called: gpss-rabbit-greenplum-connect otherwise you must compile it with the GO compiler<br/>

2. **Setting property file**    
      Use the file properties.ini (that should be place in the same directory of the binary in order to instruct the program        with this properties
      
         GpssAddress=10.91.51.23:50007
         GreenplumAddress=10.91.51.23
         GreenplumPort=5432
         GreenplumUser=gpadmin
         GreenplumPasswd= 
         Database=test
         SchemaName=public
         TableName=mytest3
         rabbit=amqp://guest:guest@localhost:5672/
         queue=gpss
         batch=50000 
         mode=1
                  
      queue is the rabbitmq queue name while batch is the amount of messages that the rabbit-greenplum connector must             receive     before pushing the data into greenplum.<br/>
      If mode is set to 1 the items batched will be saved on a disk file so in case of crash or network issue at the next         restart the connector will be automatically able to recover this info again<br/>

3. **Run the connector**
```
./gpss-rabbit-greenplum-connector 
Danieles-MBP:bin dpalaia$ ./gpss-rabbit-greenplum-connector 
connecting to grpc server
connected
2019/02/26 17:01:30  [*] Waiting for messages. To exit press CTRL+C
```

4. **Populate the queue with the UI interface (Publish command)**
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

  Then you can just go test -v ./... to let the test start </br>
  Useful commands in Go are also: </br>
  
  ```
  Danieles-MBP:gpss-rabbit-greenplum-connector dpalaia$ go test -coverprofile=coverage.out 
   2019/07/21 14:31:37 Properties read: Connecting to the Grpc server specified
   2019/07/21 14:31:37 Connected to the grpc server
   2019/07/21 14:31:37  [*] Waiting for messages. To exit press CTRL+C
   2019/07/21 14:31:37 Batch reached: I'm sending request to write to gpss/gprc server
   2019/07/21 14:31:37 connecting to a greenplum database
   2019/07/21 14:31:37 Beginning to write to greenplum
   2019/07/21 14:31:37 table informations
   2019/07/21 14:31:37 prepare for writing
   Result:  SuccessCount:3 
   2019/07/21 14:31:37 disconnecting to a greenplum database
   PASS
   coverage: 64.1% of statements
   ok  	_/Users/dpalaia/go/src/gpss-rabbit-greenplum-connector	0.461s
   ```
   to see the degree of coverage in this case 64.1% and </br>
   
   go tool cover -html=coverage.out  </br>
   it will open a .html page to see the code covered and the code not covered </br>
   
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

## Future development

1) Better loggings and put loggings in a central .log file
2) Transform the connector as a web service adding also a Dockerfile in order to push it on Dockerhub and be able to be deployed on Kubernetes. Also adding a simple web-interface like Swagger.
3) Better code review


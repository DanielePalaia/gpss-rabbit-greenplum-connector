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

GPSS is able to receive requests from different clients (Kafka, Greenplum-Informatica connector) as shown in the pic and proceed with ingestion process. We are adding support for RabbitMQ
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

      The table can be generic with any number of coloumns and data types. The important thing is that the input to ingest         will be coherent with the table definition. Let's try this table as example:
   
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
      then run the broker </br>
      ./rabbitmq-server </br>
      Then go with your browser to the rabbit web UI in: </br>
      http://localhost:15672/ </br></br>
      ![Screenshot](./pics/connection.png)<br/>
      and log with guest/guest (default)
      
      
5. **Create a rabbitmq durable queue with the rabbitmq UI interface you want the connector to connect (es gpss)**

  ![Screenshot](./pics/queue.png)<br/>
  
## Running the application

1. **Find binaries** </br>
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

5. **Insert elements as specified in batches property** </br>
      Once you publish more messages than the batch value you should then see the table populated and you can restart             publishing.<br/>


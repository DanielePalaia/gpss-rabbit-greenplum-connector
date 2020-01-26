# Summary
This software is intended to be a simple (non production ready) processor for apache nifi server, using Greenplum Streaming Service functionalities. </br>
It is written in Java and it uses the following technologies: Apache nifi, java, GRPC, Greenplum GPSS. </br>
At the moment it is just supportin .json. The processor is receiving .json entries from a nifi relashion and ingest a Greenplum table.</br> 

The following reading can help you to better understand the software:

**Apache Nifi:** </br>
https://nifi.apache.org/ </br>
**GRPC:**  </br>
https://grpc.io/ </br>
**Greenplum GPSS:**</br>
https://gpdb.docs.pivotal.io/5160/greenplum-stream/overview.html</br>
https://gpdb.docs.pivotal.io/5160/greenplum-stream/api/dev_client.html</br>

These are the steps to run the software:

## Prerequisites

1. **Activate the gpss extension on the greenplum database you want to use (for example test)**
   
      ```
      test=# CREATE EXTENSION gpss;
      ```
   
2. **Create the Greenplum table to be ingested**

      The table can be generic with any number of coloumns and data types. The important thing is that the input to ingest         will be coherent with the table definition. Let's try this table as example:
   
      ```
      test=# create table test(id varchar (data json);
      ```
   
3. **Run a gpss server with the right configuration (ex):**
  
      gpss ./gpsscfg1.json --log-dir ./gpsslogs
      where gpsscfg1.json 
  
      ```
      {
         "ListenAddress": {
            "Host": "",
            "Port": 8088,
            "SSL": false
         },
         "Gpfdist": {
            "Host": "",
            "Port": 8086
         }
      }
      ```

4. **download, install and start nifi**

![Screenshot](./pics/fourth)
  
## Deploy and test the nifi processor

1. **Copy the .nar file** </br>

2. **restart nifi** </br>

3. **insert the processor in the nifi UI** </br>


2. **Setting property of the processor**  </br>   
    

3. **Add a tester** </br>  

4. **See the greenplum table populated </br>  






package main

import (
	"context"
	"flag"
	"fmt"
	"gpssclient/gpss"
	"log"

	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc"
)

type gpssClient struct {
	client           gpss.GpssClient
	session          *gpss.Session
	conn             *grpc.ClientConn
	GpssAddress      string
	GreenplumAddress string
	GreenplumPort    int32
	GreenplumUser    string
	GreenplumPasswd  string
	Database         string
	SchemaName       string
	TableName        string
}

func MakeGpssClient(gpssAddress string, greenplumAddress string, greenplumPort int32, user string, password string, database string, schema string, table string) *gpssClient {

	gpssClient := new(gpssClient)
	gpssClient.GpssAddress = gpssAddress
	gpssClient.GreenplumAddress = greenplumAddress
	gpssClient.GreenplumPort = greenplumPort
	gpssClient.GreenplumUser = user
	gpssClient.GreenplumPasswd = password
	gpssClient.Database = database
	gpssClient.SchemaName = schema
	gpssClient.TableName = table

	return gpssClient
}

func (client *gpssClient) ConnectToGrpcServer() {

	fmt.Println("connecting to grpc server")
	/* Connecting to the grpc server */
	serverAddr := flag.String("server_addr", client.GpssAddress, "The server address in the format of host:port")
	var err error
	client.conn, err = grpc.Dial(*serverAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	client.client = gpss.NewGpssClient(client.conn)

	fmt.Println("connected")

}

func (client *gpssClient) ConnectToGreenplumDatabase() {

	fmt.Println("connecting to a greenplum database")
	connReq := gpss.ConnectRequest{Host: client.GreenplumAddress, Port: client.GreenplumPort, Password: client.GreenplumPasswd, Username: client.GreenplumUser, DB: client.Database}
	var err error
	client.session, err = client.client.Connect(context.Background(), &connReq)
	if err != nil {
		log.Fatalf("fail to connect to database: %v", err)
	}

}

func (client *gpssClient) DisconnectToGreenplumDatabase() {

	fmt.Println("disconnecting to a greenplum database")
	var err error
	_, err = client.client.Disconnect(context.Background(), client.session)
	if err != nil {
		log.Fatalf("fail to connect to database: %v", err)
	}

}

func (client *gpssClient) prepareForWriting() {

	// Prepare for writing
	columns := make([]string, 1)

	columns[0] = "data"

	insertOption := gpss.InsertOption{ErrorLimitCount: 25, ErrorLimitPercentage: 25, TruncateTable: false, InsertColumns: columns}
	openRequestInsertOption := gpss.OpenRequest_InsertOption{InsertOption: &insertOption}
	openRequest := gpss.OpenRequest{Session: client.session, SchemaName: client.SchemaName, TableName: client.TableName, Timeout: 5, Option: &openRequestInsertOption}
	_, err := client.client.Open(context.Background(), &openRequest)
	if err != nil {
		log.Fatalf("fail to open request to write: %v", err)
	}

}

func (client *gpssClient) WriteToGreenplum(buffer []string) {

	client.prepareForWriting()
	rowData := make([]*gpss.RowData, len(buffer))
	for i, line := range buffer {
		//fmt.Println("value: ", line)
		dbValue := make([]*gpss.DBValue, 1)
		dbValue[0] = new(gpss.DBValue)
		dbValue[0].DBType = &gpss.DBValue_StringValue{StringValue: string(line)}

		rowLine := gpss.Row{Columns: dbValue}
		rowData[i] = new(gpss.RowData)
		//databytes, _ := rowLine.Descriptor()
		rowData[i].Data, _ = proto.Marshal(&rowLine)
	}

	req := gpss.WriteRequest{Session: client.session, Rows: rowData}
	_, err := client.client.Write(context.Background(), &req)
	if err != nil {
		log.Fatalf("fail to open request to write: %v", err)
	}
}

/*
func (client *gpssClient) WriteToGreenplum(columnFile string, valueFile string) {

	client.prepareForWriting(columnFile)

	lines := client.readCsvFile(valueFile)
	fmt.Println("number of lines", len(lines))
	rowData := make([]*gpss.RowData, len(lines))
	for i, line := range lines {
		fmt.Println("number of words", len(line))
		dbValue := make([]*gpss.DBValue, len(line))
		for j, column := range line {
			fmt.Println("value: ", column)
			dbValue[j] = new(gpss.DBValue)
			dbValue[j].DBType = &gpss.DBValue_StringValue{StringValue: string(column)}
		}

		rowLine := gpss.Row{Columns: dbValue}
		rowData[i] = new(gpss.RowData)
		//databytes, _ := rowLine.Descriptor()
		rowData[i].Data, _ = proto.Marshal(&rowLine)
	}

	req := gpss.WriteRequest{Session: client.session, Rows: rowData}
	_, err := client.client.Write(context.Background(), &req)
	if err != nil {
		log.Fatalf("fail to open request to write: %v", err)
	}

}
*/

func (client *gpssClient) CloseRequest() {

	closeRequest := gpss.CloseRequest{Session: client.session}
	tStats, err := client.client.Close(context.Background(), &closeRequest)
	if err != nil {
		log.Fatalf("fail close write to database: %v", err)
	}

	fmt.Println("Result: ", tStats.String())

}

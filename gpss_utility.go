package main

import (
	"context"
	"flag"
	"fmt"
	"gpssclient/gpss"
	"log"
	"strconv"
	"strings"

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

	/* Connecting to the grpc server */
	serverAddr := flag.String("server_addr", client.GpssAddress, "The server address in the format of host:port")
	var err error
	client.conn, err = grpc.Dial(*serverAddr, grpc.WithInsecure(), grpc.WithMaxMsgSize(121958008))
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	client.client = gpss.NewGpssClient(client.conn)

}

func (client *gpssClient) ConnectToGreenplumDatabase() {

	log.Printf("connecting to a greenplum database")
	connReq := gpss.ConnectRequest{Host: client.GreenplumAddress, Port: client.GreenplumPort, Password: client.GreenplumPasswd, Username: client.GreenplumUser, DB: client.Database}
	var err error
	client.session, err = client.client.Connect(context.Background(), &connReq)
	if err != nil {
		log.Fatalf("fail to connect to database: %v", err)
	}

}

func (client *gpssClient) DisconnectToGreenplumDatabase() {

	log.Printf("disconnecting to a greenplum database")
	var err error
	_, err = client.client.Disconnect(context.Background(), client.session)
	if err != nil {
		log.Fatalf("fail to connect to database: %v", err)
	}

}

func (client *gpssClient) DescribeTable() *gpss.Columns {

	log.Printf("table informations")
	describeReq := gpss.DescribeTableRequest{Session: client.session, TableName: client.TableName, SchemaName: client.SchemaName}

	columns, _ := client.client.DescribeTable(context.Background(), &describeReq)

	return columns

}

func (client *gpssClient) prepareForWriting(cols *gpss.Columns) {

	log.Printf("prepare for writing")
	// Prepare for writing
	columns := make([]string, len(cols.Columns))
	// Looping over table columns as defined in greenplum
	for i, col := range cols.Columns {
		//fmt.Println("field: " + col.Name)
		columns[i] = col.Name
	}

	insertOption := gpss.InsertOption{ErrorLimitCount: 25, ErrorLimitPercentage: 25, TruncateTable: false, InsertColumns: columns}
	openRequestInsertOption := gpss.OpenRequest_InsertOption{InsertOption: &insertOption}
	openRequest := gpss.OpenRequest{Session: client.session, SchemaName: client.SchemaName, TableName: client.TableName, Timeout: 5, Option: &openRequestInsertOption}
	_, err := client.client.Open(context.Background(), &openRequest)
	if err != nil {
		log.Fatalf("fail to open request to write: %v", err)
	}

}

func convertType(field string, databaseType string) *gpss.DBValue {
	//	*DBValue_Int32Value
	//	*DBValue_Int64Value
	//	*DBValue_Float32Value
	//	*DBValue_Float64Value
	//	*DBValue_StringValue
	//	*DBValue_BytesValue
	//	*DBValue_TimeStampValue
	//	*DBValue_NullValue
	//fmt.Println("field is: " + field + "length: " + string(len(field)))
	dbValue := new(gpss.DBValue)
	// Consider NULL value
	if field == "NULL" {
		dbValue.DBType = &gpss.DBValue_NullValue{}
		return dbValue
	}
	if strings.Contains(databaseType, "int") || strings.Contains(databaseType, "serial") {
		if n, err := strconv.Atoi(field); err == nil {
			dbValue.DBType = &gpss.DBValue_Int64Value{Int64Value: int64(n)}
		} else {
			fmt.Println(field, "is not an integer.")
		}
	} else if strings.Contains(databaseType, "float") || strings.Contains(databaseType, "numeric") || strings.Contains(databaseType, "decimal") {
		if n, err := strconv.ParseFloat(field, 64); err == nil {
			dbValue.DBType = &gpss.DBValue_Float64Value{Float64Value: float64(n)}
		}
	} else {
		dbValue.DBType = &gpss.DBValue_StringValue{StringValue: field}
	} /*else if strings.Contains(databaseType, "timestamp") {
		dbValue.DBType = &gpss.DBValue_TimeStampValue{TimeStampValue: string(field)}
	} */

	return dbValue
}

func (client *gpssClient) WriteToGreenplum(buffer []string) {

	log.Printf("Beginning to write to greenplum")
	// Take table information from greenplum
	cols := client.DescribeTable()

	// Prepare table for writing
	client.prepareForWriting(cols)

	// For every entry buffered in memory check for fields
	rowData := make([]*gpss.RowData, len(buffer))
	for i, line := range buffer {

		dbValue := make([]*gpss.DBValue, len(cols.Columns))

		// From an entry composed by lines return a list
		fields := strings.Split(strings.Replace(line, "\r\n", "\n", -1), "\n")

		for j, field := range fields {

			dbValue[j] = convertType(field, cols.Columns[j].DatabaseType)

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

	client.closeRequest()
}

func (client *gpssClient) closeRequest() {

	closeRequest := gpss.CloseRequest{Session: client.session}
	tStats, err := client.client.Close(context.Background(), &closeRequest)
	if err != nil {
		log.Fatalf("fail close write to database: %v", err)
	}

	fmt.Println("Result: ", tStats.String())

}

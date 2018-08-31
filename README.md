To run a curl command that inserts data from a CSV in S3 to Dynamo DB use the following command :

curl -XPOST http://localhost:8001/process -d '{"ColumnDefinitions":[{"CSVColumnIndex":-1,"InsertUUID":true,"DynamoColumnName":"id","DynamoDataType":"string"},{"CSVColumnIndex":0,"InsertUUID":false,"DynamoColumnName":"colour","DynamoDataType":"string"},{"CSVColumnIndex":1,"InsertUUID":false,"DynamoColumnName":"anumber","DynamoDataType":"integer"},{"CSVColumnIndex":2,"InsertUUID":false,"DynamoColumnName":"chocolate","DynamoDataType":"string"}],"DynamoConfig":{"TableName":"Dynago","Hash":"id","Sort":"name","MaximumCapacity":200,"MaximumPercentageCapacity":80,"StartCapacity":20},"S3Config":{"Region":"us-east-1","Bucket":"dynagouseast1","Key":"testdata.csv"}}'

To See the data in the database using the command line :

```
aws dynamodb scan --table-name DynagoDb
```

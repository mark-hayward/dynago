package main


type S3Config struct {
	Region string
	Bucket string
	Key    string
}

type DynamoConfig struct {
	TableName                string
	Hash                     string
	Sort                     string
	MaximumCapacity          int
	MaximumPercentageCapacity int
	StartCapacity            int
}

type ColumnDefinition struct {
	CSVColumnIndex   int
	InsertUUID       bool
	DynamoColumnName string
	DynamoDataType   string
}

type InsertTask struct {
	S3Config
	ColumnDefinitions []ColumnDefinition
	DynamoConfig
}

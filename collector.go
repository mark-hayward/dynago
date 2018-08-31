package main

import (
	//		"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"time"
	"github.com/aws/aws-sdk-go-v2/service/s3/s3manager"
	"bytes"
	"encoding/csv"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/aws/endpoints"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/twinj/uuid"
	"net/http"
	"encoding/json"
	"log"
	"io/ioutil"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
)




type WorkRequest struct {
	record		*dynamodb.PutItemInput
}

type JSONRequest struct {
	ColumnDefinitions []ColumnDefinition
	DynamoConfig    DynamoConfig
	S3Config    S3Config
	Name 	string
}



// A buffered channel that we can send work requests on.
var WorkQueue = make(chan WorkRequest, 100)

	//dynamoConfig := &DynamoConfig{
	//	TableName: "DynagoDb",
	//	Hash:"id",
	//	Sort:"somethingelse",
	//	MaximumCapacity:2000,
	//	MaximumPercentageCapacity: 80,
	//	StartCapacity: 50,
	//}

func CloudwatchWorker(rate chan<- time.Duration, d DynamoConfig) {
	fmt.Println(rate)
	// Get Current Write Capacity
	sleepTime := time.Duration(0)

	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		panic("unable to load SDK config, " + err.Error())
	}
	cfg.Region = endpoints.UsEast1RegionID
	svc := cloudwatch.New(cfg)
	ticker := time.NewTicker(time.Minute * 1)
	for tick := range ticker.C {
		fmt.Println(tick)

 		if d.MaximumPercentageCapacity < 100 {

			// Get Current Throughput
			currentCapacity := float64(5)

			// Send the request, and get the response or error back
			now := time.Now()

			// TODO Maybe use now.Sub?
			prev := now.Add(time.Duration(60) * time.Minute * -1)

			rcr := svc.GetMetricStatisticsRequest(&cloudwatch.GetMetricStatisticsInput{
				MetricName: aws.String("ProvisionedReadCapacityUnits"), //ProvisionedWriteCapacityUnits
				Namespace:  aws.String("AWS/DynamoDB"),
				Dimensions: []cloudwatch.Dimension{
					{
						Name:  aws.String("TableName"),
						Value: aws.String(d.TableName),
					},
				},
				Period:     aws.Int64(60),
				StartTime:  &prev,
				EndTime:    &now,
				Statistics: []cloudwatch.Statistic{"Sum"},
			})

			rc, err := rcr.Send()
			if err != nil {
				panic("failed to describe table, " + err.Error())
			}
			if rc.Datapoints != nil {
				currentCapacity = *rc.Datapoints[0].Sum
			}

			currentMaxPercentageThroughput := float64(d.MaximumPercentageCapacity) * currentCapacity / 100
			// currentMaxPercentageThroughput is now the maximum speed if percentage is declared
			// t.DynamoConfig.MaximumCapacity is the maximum capacity
			// Which is lower?

			// If the current throughput speed is greater than the Maximum Capacity then
			// MaximumCapacity is the lower.
			fmt.Println("currentMaxPercentageThroughput")
			fmt.Println(currentMaxPercentageThroughput)

			if d.MaximumCapacity <= int(currentMaxPercentageThroughput) {
				fmt.Println("Max Capacity Less than Current")
				fmt.Println(d.MaximumCapacity)
				fmt.Println(currentMaxPercentageThroughput)
				sleepTime = 1 * time.Second / time.Duration(d.MaximumCapacity)
			} else {
				sleepTime = 1 * time.Second / time.Duration(currentMaxPercentageThroughput)
			}
			rate <- sleepTime
		} else {
			i := float32(1) / float32(d.MaximumCapacity)
			rate <- time.Duration(i) * time.Second
		}
	}

}

func Collector(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	log.Println(string(body))
	var t JSONRequest
	err = json.Unmarshal(body, &t)
	if err != nil {
		panic(err)
	}
	log.Println(t.Name)

	rate := make(chan time.Duration)

	go CloudwatchWorker(rate, t.DynamoConfig)

	// TODO Could the default configuration be set at startup? If it needs to be configurable, perhaps is needs to be
	// passed in as a POST parameter (perhaps with its own handler function)
	dynamoCfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		panic("unable to load SDK config, " + err.Error())
	}

	// TODO This could be made configurable.
	dynamoCfg.Region = endpoints.UsEast1RegionID
	NWorkers := 4

	StartDispatcher(NWorkers, dynamoCfg)

	// TODO This could also be set at startup, or made configurable with its own handler function.
	s3Cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		panic("unable to load SDK config, " + err.Error())
	}

	// Set the AWS Region that the S3 service clients should use
	// region usually be regionUsEast1RegionID
	s3Cfg.Region = endpoints.UsEast1RegionID

	svc := s3.New(s3Cfg)

	downloader := s3manager.NewDownloaderWithClient(svc, func(d *s3manager.Downloader) {
		d.PartSize = 1024 * 1024 * 64 // 64MB per part
		d.Concurrency = 10
	})
	buff := &aws.WriteAtBuffer{}

	numBytes, err := downloader.Download(buff,
		&s3.GetObjectInput{
			Bucket: aws.String(t.S3Config.Bucket),
			Key:    aws.String(t.S3Config.Key),
		})

	if err != nil {
		panic(err)
	}

	// Split buffer into new lines
	// Is this how I get a string from this buff?
	if numBytes == 1 {
		fmt.Println("shit file")
	}
	csvFile := bytes.NewReader(buff.Bytes())


	//Now I need to get I want to split the CSV into lines. How do I do this and append any trailing characters to the next buffer?
	lines, _ := csv.NewReader(csvFile).ReadAll()

	// This should be one second (1000ms) divided by the amount of write throughput.
	// In this example, we have a write throughput of 5 items per second. Answer should be 200ms
	// If val
	sleepTime :=  time.Duration(0)

	for _, line := range lines {

		// split record into columns:
		//var m map[string]interface{}
		m := make(map[string]interface{})
		for _, definition := range t.ColumnDefinitions {
			if definition.InsertUUID {
				m[definition.DynamoColumnName] = uuid.NewV4().String()
			} else {
				m[definition.DynamoColumnName] = line[definition.CSVColumnIndex]
			}
		}
		a,err := dynamodbattribute.MarshalMap(m)
		if err != nil {
			fmt.Println(err)
		}
		// Need a way to listen to a channel to get new wait times when changed
		// decode the csv line

		//s:=data["asset_id"].(string)
		params := &dynamodb.PutItemInput{
			Item: a,
			TableName: aws.String(t.DynamoConfig.TableName), // Required
		}
		fmt.Println(params)
		work := WorkRequest{record: params}
		// Push the line into the queue.
		WorkQueue <- work
		fmt.Println("Dynamo DB request queued")
		select {
		case v := <-rate:
			sleepTime = v
			fmt.Println(sleepTime)
		default:
			//do nowt
		}
		time.Sleep(sleepTime)
	}

	return
}

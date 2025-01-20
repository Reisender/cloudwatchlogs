package main

import (
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	cwl "github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"

	"github.com/urfave/cli/v2"
)

const pageSize = 100

func main() {
	app := &cli.App{
		Name:  "getlogstream",
		Usage: "stream out all events from a log stream",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Required: true,
				Name:     "log-group",
				Usage:    "the log group that the stream is a part of.",
			},
			&cli.StringFlag{
				Required: true,
				Name:     "log-stream",
				Usage:    "the log stream to get events from",
			},
			&cli.StringFlag{
				Name:    "profile",
				Usage:   "the aws profile to use",
				EnvVars: []string{"AWS_PROFILE"},
			},
		},
		Action: run,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {
	ctx := c.Context

	// Load the Shared AWS Configuration (~/.aws/config
	options := [](func(*config.LoadOptions) error){}
	if len(c.String("profile")) > 0 {
		options = append(options, config.WithSharedConfigProfile(c.String("profile")))
	}
	cfg, err := config.LoadDefaultConfig(ctx, options...)
	if err != nil {
		log.Fatal(err)
	}

	// Create an Amazon S3 service client
	client := cwl.NewFromConfig(cfg)

	var nextToken *string
	var output *cwl.GetLogEventsOutput
	eventCount := 1

	for eventCount > 0 {
		// Get the first page of results for ListObjectsV2 for a bucket
		err := retry(3, func() error {
			var err error
			output, err = client.GetLogEvents(ctx, &cwl.GetLogEventsInput{
				LogGroupName:  aws.String(c.String("log-group")),
				LogStreamName: aws.String(c.String("log-stream")),
				Limit:         aws.Int32(pageSize),
				StartFromHead: aws.Bool(true),
				NextToken:     nextToken,
			})
			return err
		})
		if err != nil {
			return err
		}

		nextToken = output.NextForwardToken
		eventCount = len(output.Events)

		for _, evt := range output.Events {
			fmt.Println(*evt.Message)
		}
	}

	return nil
}

func retry(times int, try func() error) error {
	var err error
	tried := 0
	for times >= tried {
		if err = try(); err == nil {
			return nil
		}
		tried++
	}
	return err
}

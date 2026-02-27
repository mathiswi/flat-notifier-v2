package main

import (
	"os"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsevents"
	"github.com/aws/aws-cdk-go/awscdk/v2/awseventstargets"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslogs"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type FlatNotifierStackProps struct {
	awscdk.StackProps
}

func NewFlatNotifierStack(scope constructs.Construct, id string, props *FlatNotifierStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	// DynamoDB table
	table := awsdynamodb.NewTable(stack, jsii.String("EbayTable"), &awsdynamodb.TableProps{
		TableName:    jsii.String("ebayTable"),
		PartitionKey: &awsdynamodb.Attribute{Name: jsii.String("flatId"), Type: awsdynamodb.AttributeType_STRING},
		BillingMode:  awsdynamodb.BillingMode_PAY_PER_REQUEST,
	})

	// Lambda function
	fn := awslambda.NewFunction(stack, jsii.String("FlatNotifierFn"), &awslambda.FunctionProps{
		Runtime:      awslambda.Runtime_PROVIDED_AL2023(),
		Handler:      jsii.String("bootstrap"),
		Architecture: awslambda.Architecture_ARM_64(),
		Code:         awslambda.Code_FromAsset(jsii.String("../deploy"), nil),
		MemorySize:   jsii.Number(128),
		Timeout:      awscdk.Duration_Seconds(jsii.Number(30)),
		Environment: &map[string]*string{
			"DYNAMODB_TABLE":  table.TableName(),
			"OVERVIEW_URL":    jsii.String(envOrDefault("OVERVIEW_URL", "")),
			"TELEGRAM_TOKEN":  jsii.String(envOrDefault("TELEGRAM_TOKEN", "")),
			"TELEGRAM_CHAT_ID": jsii.String(envOrDefault("TELEGRAM_CHAT_ID", "")),
			"DISCORD_TOKEN":   jsii.String(envOrDefault("DISCORD_TOKEN", "")),
			"DISCORD_USER_ID": jsii.String(envOrDefault("DISCORD_USER_ID", "")),
		},
		LogRetention: awslogs.RetentionDays_TWO_WEEKS,
	})

	// Grant DynamoDB access
	table.GrantReadWriteData(fn)

	// EventBridge schedule: every 5 minutes
	rule := awsevents.NewRule(stack, jsii.String("Every5MinRule"), &awsevents.RuleProps{
		Schedule: awsevents.Schedule_Rate(awscdk.Duration_Minutes(jsii.Number(5))),
	})
	rule.AddTarget(awseventstargets.NewLambdaFunction(fn, nil))

	return stack
}

func envOrDefault(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	NewFlatNotifierStack(app, "FlatNotifierStack", &FlatNotifierStackProps{
		awscdk.StackProps{
			Env: &awscdk.Environment{
				Region: jsii.String("eu-central-1"),
			},
		},
	})

	app.Synth(nil)
}

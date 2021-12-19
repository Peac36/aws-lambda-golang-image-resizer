## AWS Lambda Golang Image Resizer

A simple image resizer using [Golang](https://go.dev/) and [AWS Lambda](https://aws.amazon.com/lambda/).

#### Environment variable

* `OUTPUT_BUCKET` - output bucket
* `SIZES` - a JSON that containers resized image size,width and save directory.

Examples:

`OUTPUT_BUCKET`: `my-bucket`

`SIZES`:`[{"OutputDirectory":"/thumbnail/","SizeWidth":600,"SizeHeight":300},{"OutputDirectory":"/medium/","SizeWidth":1920,"SizeHeight":1020}]`

#### How to build

Simplity execute the `build.sh` script

#### Setup

1. Create a S3
2. Create a SQS
3. Create the following policies:

    * SQS Policy
        ```
        {
            "Version": "2012-10-17",
            "Statement": [
                {
                    "Sid": "SQSRole",
                    "Effect": "Allow",
                    "Action": [
                        "sqs:DeleteMessage",
                        "sqs:ChangeMessageVisibility",
                        "sqs:ReceiveMessage",
                        "sqs:GetQueueAttributes"
                    ],
                    "Resource": "<%your-sqs-arn%>"
                }
            ]
        }
        ```
    * S3 Policy
        ```
                {
            "Version": "S3-role",
            "Statement": [
                {
                    "Sid": "VisualEditor0",
                    "Effect": "Allow",
                    "Action": [
                        "s3:PutObject",
                        "s3:GetObjectAcl",
                        "s3:GetObject",
                        "s3:PutObjectTagging",
                        "s3:PutObjectAcl"
                    ],
                    "Resource": "<%your-s3-bucket-arn%>/*"
                }
            ]
        }
        ```
4. Create a role with following policies: `SQS Role`, `S3 Role` and `AWSLambdaBasicExecutionRole`.
5. Create a S3 event notification to post to SQS whenever a file is uploaded
6. Apply following policy to your SQS(`Access policy` tab)
    ```
    {
        "Version": "2012-10-17",
        "Id": "golang-resizer",
        "Statement": [
            {
            "Sid": "golang-resizer",
            "Effect": "Allow",
            "Principal": {
                "Service": "s3.amazonaws.com"
            },
            "Action": "SQS:SendMessage",
            "Resource": "<%your-sqs-arn%>",
            "Condition": {
                "ArnLike": {
                "aws:SourceArn": "<%your-s3-bucket-arn%>"
                }
            }
            }
        ]
    }
    ```

7. Create a Lambda function with role from step 4, handler `main` and trigger with SQS queue created from step 2.
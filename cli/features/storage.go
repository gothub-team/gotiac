package features

import (
	"context"
	"fmt"
	"os"

	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/s3"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optdestroy"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optup"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func Storage() {
	// to destroy our program, we can run `go run main.go destroy`
	destroy := false
	argsWithoutProg := os.Args[1:]
	if len(argsWithoutProg) > 0 {
		if argsWithoutProg[0] == "destroy" {
			destroy = true
		}
	}

	// define our program that creates our pulumi resources.
	// we refer to this style as "inline" pulumi programs where both program + automation can be compiled in the same
	// binary. no need for separate projects.
	deployFunc := func(ctx *pulumi.Context) error {
		// similar go git_repo_program, our program defines a s3 website.
		// here we create the bucket
		siteBucket, err := s3.NewBucket(ctx, "s3-website-bucket", &s3.BucketArgs{
			Website: s3.BucketWebsiteArgs{
				IndexDocument: pulumi.String("index.html"),
			},
		})
		if err != nil {
			return err
		}

		// we define and upload our HTML inline.
		indexContent := `<html><head>
		<title>Hello S3</title><meta charset="UTF-8">
	</head>
	<body><p>Hello, world!</p><p>Made with ❤️ with <a href="https://pulumi.com">Pulumi</a></p>
	</body></html>
`
		// upload our index.html
		if _, err := s3.NewBucketObject(ctx, "index", &s3.BucketObjectArgs{
			Bucket:      siteBucket.ID(), // reference to the s3.Bucket object
			Content:     pulumi.String(indexContent),
			Key:         pulumi.String("index.html"),               // set the key of the object
			ContentType: pulumi.String("text/html; charset=utf-8"), // set the MIME type of the file
		}); err != nil {
			return err
		}

		// Allow public ACLs for the bucket
		accessBlock, err := s3.NewBucketPublicAccessBlock(ctx, "public-access-block", &s3.BucketPublicAccessBlockArgs{
			Bucket:          siteBucket.ID(),
			BlockPublicAcls: pulumi.Bool(false),
		})
		if err != nil {
			return err
		}

		// Set the access policy for the bucket so all objects are readable.
		if _, err := s3.NewBucketPolicy(ctx, "bucketPolicy", &s3.BucketPolicyArgs{
			Bucket: siteBucket.ID(), // refer to the bucket created earlier
			Policy: pulumi.Any(map[string]interface{}{
				"Version": "2012-10-17",
				"Statement": []map[string]interface{}{
					{
						"Effect":    "Allow",
						"Principal": "*",
						"Action": []interface{}{
							"s3:GetObject",
						},
						"Resource": []interface{}{
							pulumi.Sprintf("arn:aws:s3:::%s/*", siteBucket.ID()), // policy refers to bucket name explicitly
						},
					},
				},
			}),
		}, pulumi.DependsOn([]pulumi.Resource{accessBlock})); err != nil {
			return err
		}

		// export the website URL
		ctx.Export("websiteUrl", siteBucket.WebsiteEndpoint)
		return nil
	}

	ctx := context.Background()

	projectName := "inlineS3Project"
	// we use a simple stack name here, but recommend using auto.FullyQualifiedStackName for maximum specificity.
	stackName := "dev"
	// stackName := auto.FullyQualifiedStackName("myOrgOrUser", projectName, stackName)

	// create or select a stack matching the specified name and project.
	// this will set up a workspace with everything necessary to run our inline program (deployFunc)
	s, err := auto.UpsertStackInlineSource(ctx, stackName, projectName, deployFunc)
	if err != nil {
		fmt.Printf("Failed to set up a workspace: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Created/Selected stack %q\n", stackName)

	w := s.Workspace()

	fmt.Println("Installing the AWS plugin")

	// for inline source programs, we must manage plugins ourselves
	err = w.InstallPlugin(ctx, "aws", "v4.0.0")
	if err != nil {
		fmt.Printf("Failed to install program plugins: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Successfully installed AWS plugin")

	// set stack configuration specifying the AWS region to deploy
	s.SetConfig(ctx, "aws:region", auto.ConfigValue{Value: "eu-central-1"})

	fmt.Println("Successfully set config")
	fmt.Println("Starting refresh")

	_, err = s.Refresh(ctx)
	if err != nil {
		fmt.Printf("Failed to refresh stack: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Refresh succeeded!")

	if destroy {
		fmt.Println("Starting stack destroy")

		// wire up our destroy to stream progress to stdout
		stdoutStreamer := optdestroy.ProgressStreams(os.Stdout)

		// destroy our stack and exit early
		_, err := s.Destroy(ctx, stdoutStreamer)
		if err != nil {
			fmt.Printf("Failed to destroy stack: %v", err)
		}
		fmt.Println("Stack successfully destroyed")
		os.Exit(0)
	}

	fmt.Println("Starting update")

	// wire up our update to stream progress to stdout
	stdoutStreamer := optup.ProgressStreams(os.Stdout)

	// run the update to deploy our s3 website
	res, err := s.Up(ctx, stdoutStreamer)
	if err != nil {
		fmt.Printf("Failed to update stack: %v\n\n", err)
		os.Exit(1)
	}

	fmt.Println("Update succeeded!")

	// get the URL from the stack outputs
	url, ok := res.Outputs["websiteUrl"].Value.(string)
	if !ok {
		fmt.Println("Failed to unmarshall output URL")
		os.Exit(1)
	}

	fmt.Printf("URL: %s\n", url)
}

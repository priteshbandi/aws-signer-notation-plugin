package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, event AWSEvent) error {
	if event.Version != "0" {
		return fmt.Errorf("unsupported event version detected: %s", event.Version)
	}

	if event.Source != "aws.ecr" || event.DetailType != "AWS API Call via CloudTrail" || event.Detail.EventName != "PutImage" {
		return fmt.Errorf("unsupported event source: %s, detail type: %s, or event name: %s", event.Source, event.DetailType, event.Detail.EventName)
	}

	image := event.Detail.ResponseElements.Image
	imageURI := fmt.Sprintf("%s.dkr.ecr.%s.amazonaws.com/%s@%s", image.RegistryId, event.Region, image.RepositoryName, image.ImageId.ImageDigest)

	if strings.Contains(image.ImageManifest, "application/vnd.cncf.notary.signature") {
		log.Printf("%s ia already signed, skipping signing...", imageURI)
		return nil
	}

	trustedRoleArn := os.Getenv("AWS_TRUSTED_IAM_ROLE_ARN")
	if trustedRoleArn == "" {
		return fmt.Errorf("AWS_TRUSTED_IAM_ROLE_ARN environment variable not set")
	}
	if event.Detail.UserIdentity.Arn != trustedRoleArn {
		log.Printf("%s pushed using non-trusted IAM role '%s', skipping signing...", imageURI, event.Detail.UserIdentity.Arn)
		return nil
	}

	signingProfileArn := os.Getenv("AWS_SIGNER_PROFILE_ARN")
	if signingProfileArn == "" {
		return fmt.Errorf("AWS_SIGNER_PROFILE_ARN environment variable not set")
	}

	log.Printf("signing image %s using %s signing profile", imageURI, signingProfileArn)
	userMetadata := map[string]string{"tagAtPush": image.ImageId.ImageTag} // Optional, add if you want to add metadata to the signature, else use nil

	// signing
	signer, err := NewNotationSigner(ctx, event.Region)
	if err != nil {
		panic(err)
	}
	err = signer.Sign(ctx, signingProfileArn, imageURI, userMetadata)
	if err != nil {
		return err
	}

	fmt.Printf("Sucessfully signed artifact: %s.", imageURI)
	return nil
}

func main() {
	lambda.Start(handler)
}

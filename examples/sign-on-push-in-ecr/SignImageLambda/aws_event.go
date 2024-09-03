package main

import (
	"encoding/json"
	"time"
)

type AWSEvent struct {
	Version    string                  `json:"version"`
	Id         string                  `json:"id"`
	Source     string                  `json:"source"`
	Account    string                  `json:"account"`
	Time       time.Time               `json:"time"`
	Region     string                  `json:"region"`
	Resources  []string                `json:"resources"`
	DetailType string                  `json:"detail-type"`
	Detail     AWSAPICallViaCloudTrail `json:"detail"`
}

type AWSAPICallViaCloudTrail struct {
	ResponseElements ResponseElements `json:"responseElements"`
	UserIdentity     UserIdentity     `json:"userIdentity"`
	EventID          string           `json:"eventID"`
	AwsRegion        string           `json:"awsRegion"`
	EventVersion     string           `json:"eventVersion"`
	EventSource      string           `json:"eventSource"`
	Resources        []Resources      `json:"resources"`
	EventType        string           `json:"eventType"`
	RequestID        string           `json:"requestID"`
	EventName        string           `json:"eventName"`
}

type ResponseElements struct {
	Image Image `json:"image"`
}

type Image struct {
	RegistryId     string  `json:"registryId"`
	RepositoryName string  `json:"repositoryName"`
	ImageId        ImageId `json:"imageId"`
	ImageManifest  string  `json:"imageManifest"`
}

type ImageId struct {
	ImageTag               string `json:"imageTag"`
	ImageDigest            string `json:"imageDigest"`
	ImageManifestMediaType string `json:"imageManifestMediaType"`
}

type UserIdentity struct {
	PrincipalId string `json:"principalId"`
	Arn         string `json:"arn"`
}

type Resources struct {
	ARN string `json:"ARN"`
}

func UnmarshalEvent(inputStream []byte) (AWSEvent, error) {
	var outputEvent AWSEvent
	err := json.Unmarshal(inputStream, &outputEvent)
	if err != nil {
		return AWSEvent{}, err
	}

	return outputEvent, nil
}

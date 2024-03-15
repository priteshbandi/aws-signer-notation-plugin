//  Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
//
//  Licensed under the Apache License, Version 2.0 (the "License"). You may
//  not use this file except in compliance with the License. A copy of the
//  License is located at
//
// 	http://aws.amazon.com/apache2.0
//
//  or in the "license" file accompanying this file. This file is distributed
//  on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
//  express or implied. See the License for the specific language governing
//  permissions and limitations under the License.

package client

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/signer"
)

// Interface facilitates unit testing
type Interface interface {
	SignPayload(ctx context.Context, params *signer.SignPayloadInput, optFns ...func(*signer.Options)) (*signer.SignPayloadOutput, error)
	GetRevocationStatus(ctx context.Context, params *signer.GetRevocationStatusInput, optFns ...func(*signer.Options)) (*signer.GetRevocationStatusOutput, error)
}
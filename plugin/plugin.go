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

// Package plugin implements the interface [github.com/notaryproject/notation-plugin-framework-go/plugin],
// enabling its use as a library in the notation-go package and to generate executable
package plugin

import (
	"context"

	"github.com/aws/aws-signer-notation-plugin/internal/client"
	"github.com/aws/aws-signer-notation-plugin/internal/signer"
	"github.com/aws/aws-signer-notation-plugin/internal/verifier"
	"github.com/aws/aws-signer-notation-plugin/internal/version"
	"github.com/notaryproject/notation-plugin-framework-go/plugin"
)

// AWSSignerPlugin provides functionality for signing and verification in accordance with the NotaryProject AWSSignerPlugin contract.
type AWSSignerPlugin struct {
	awssigner client.Interface
}

// NewAWSSigner creates new AWSSignerPlugin
func NewAWSSigner(s client.Interface) *AWSSignerPlugin {
	return &AWSSignerPlugin{awssigner: s}
}

// NewAWSSignerForCLI creates a new AWSSignerPlugin and is intended solely for generating executables.
func NewAWSSignerForCLI() *AWSSignerPlugin {
	return &AWSSignerPlugin{}
}

// VerifySignature performs the extended verification of signature by optionally calling AWS Signer.
func (sp *AWSSignerPlugin) VerifySignature(ctx context.Context, req *plugin.VerifySignatureRequest) (*plugin.VerifySignatureResponse, error) {
	if req == nil {
		return nil, plugin.NewValidationError("verifySignature req is nil")
	}
	if err := req.Validate(); err != nil {
		return nil, err
	}
	if err := sp.setSignerClientIfNotPresent(ctx, req.PluginConfig); err != nil {
		return nil, err
	}

	return verifier.New(sp.awssigner).Verify(ctx, req)
}

// GetMetadata returns the metadata information of the plugin.
func (sp *AWSSignerPlugin) GetMetadata(_ context.Context, _ *plugin.GetMetadataRequest) (*plugin.GetMetadataResponse, error) {
	return &plugin.GetMetadataResponse{
		Name:                      "com.amazonaws.signer.notation.plugin",
		Description:               "AWS Signer plugin for Notation",
		Version:                   version.GetVersion(),
		URL:                       "https://docs.aws.amazon.com/signer",
		SupportedContractVersions: []string{plugin.ContractVersion},
		Capabilities: []plugin.Capability{
			plugin.CapabilityEnvelopeGenerator,
			plugin.CapabilityTrustedIdentityVerifier,
			plugin.CapabilityRevocationCheckVerifier,
		},
	}, nil
}

// DescribeKey describes the key being used for signing. This method is not supported by AWS Signer's plugin.
func (sp *AWSSignerPlugin) DescribeKey(_ context.Context, _ *plugin.DescribeKeyRequest) (*plugin.DescribeKeyResponse, error) {
	return nil, plugin.NewUnsupportedError("DescribeKey operation")
}

// GenerateSignature generates the raw signature. This method is not supported by AWS Signer's plugin.
func (sp *AWSSignerPlugin) GenerateSignature(_ context.Context, _ *plugin.GenerateSignatureRequest) (*plugin.GenerateSignatureResponse, error) {
	return nil, plugin.NewUnsupportedError("GenerateSignature operation")
}

// GenerateEnvelope returns the signature envelope generated by calling AWS Signer.
func (sp *AWSSignerPlugin) GenerateEnvelope(ctx context.Context, req *plugin.GenerateEnvelopeRequest) (*plugin.GenerateEnvelopeResponse, error) {
	if req == nil {
		return nil, plugin.NewValidationError("generateEnvelope request is nil")
	}
	if err := req.Validate(); err != nil {
		return nil, err
	}
	if err := sp.setSignerClientIfNotPresent(ctx, req.PluginConfig); err != nil {
		return nil, err
	}

	return signer.New(sp.awssigner).GenerateEnvelope(ctx, req)
}

func (sp *AWSSignerPlugin) setSignerClientIfNotPresent(ctx context.Context, plConfig map[string]string) error {
	if sp.awssigner == nil {
		s, err := client.NewAWSSigner(ctx, plConfig)
		if err != nil {
			return err
		}
		sp.awssigner = s
	}
	return nil
}
// Copyright Aviatrix Systems, Inc.
// SPDX-License-Identifier: MPL-2.0

package clients

import (
	"context"
	"crypto/tls"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	configpb "github.com/AviatrixSystems/terraform-provider-avxcloud/internal/gen-protogo/aviatrix/config/v1alpha"
	configtf "github.com/AviatrixSystems/terraform-provider-avxcloud/internal/gen-terraform/aviatrix/config/v1alpha"
)

const (
	defaultApiAddr = "api.cloud.aviatrix.com"
)

type Config struct {
	Addr         string
	InsecureTLS  bool
	APIAccessKey string
}

type ClientGetter interface {
	configtf.CloudAccountClientGetter
	configtf.NetworkClientGetter
	configtf.NetworkInspectionClientGetter
	configtf.SmartGroupClientGetter
	configtf.WebGroupClientGetter
	configtf.DcfPolicyBlockClientGetter
	configtf.DcfPolicyListClientGetter
}

var _ ClientGetter = &Clients{}

type Clients struct {
	Addr     string
	apidConn *grpc.ClientConn
}

func (c *Clients) CloudAccountClient() configpb.CloudAccountServiceClient {
	return configpb.NewCloudAccountServiceClient(c.apidConn)
}

func (c *Clients) NetworkClient() configpb.NetworkServiceClient {
	return configpb.NewNetworkServiceClient(c.apidConn)
}

func (c *Clients) NetworkInspectionClient() configpb.NetworkInspectionServiceClient {
	return configpb.NewNetworkInspectionServiceClient(c.apidConn)
}

func (c *Clients) SmartGroupClient() configpb.SmartGroupServiceClient {
	return configpb.NewSmartGroupServiceClient(c.apidConn)
}

func (c *Clients) WebGroupClient() configpb.WebGroupServiceClient {
	return configpb.NewWebGroupServiceClient(c.apidConn)
}

func (c *Clients) DcfPolicyBlockClient() configpb.DcfPolicyBlockServiceClient {
	return configpb.NewDcfPolicyBlockServiceClient(c.apidConn)
}

func (c *Clients) DcfPolicyListClient() configpb.DcfPolicyListServiceClient {
	return configpb.NewDcfPolicyListServiceClient(c.apidConn)
}

type jwtCredentials struct {
	token      string
	requireTLS bool
}

func (j *jwtCredentials) RequireTransportSecurity() bool {
	return j.requireTLS
}

func (j *jwtCredentials) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"authorization": j.token,
	}, nil
}

func tlsCredentials(insecure bool) credentials.TransportCredentials {
	config := &tls.Config{
		// #nosec G402
		InsecureSkipVerify: insecure,
	}

	return credentials.NewTLS(config)
}

func New(ctx context.Context, config *Config) (*Clients, error) {
	clients := &Clients{
		Addr: defaultApiAddr,
	}

	if config.Addr != "" {
		clients.Addr = config.Addr
	}

	err := connectApid(
		ctx,
		clients,
		config,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Aviatrix Cloud endpoint: %w", err)
	}
	return clients, nil
}

func connectApid(ctx context.Context, c *Clients, config *Config) error {
	creds := &jwtCredentials{
		token:      config.APIAccessKey,
		requireTLS: !config.InsecureTLS,
	}

	tlsCreds := tlsCredentials(config.InsecureTLS)
	opts := []grpc.DialOption{
		grpc.WithPerRPCCredentials(creds),
		grpc.WithTransportCredentials(tlsCreds),
	}

	ApidConn, err := grpc.NewClient(c.Addr, opts...)
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	tflog.Info(ctx, "connected to Aviatrix Cloud endpoint", map[string]any{"address": c.Addr})
	c.apidConn = ApidConn
	return nil
}

func (c *Clients) Close() {
	if err := c.apidConn.Close(); err != nil {
		tflog.Warn(context.Background(), "failed to close connection to Aviatrix Cloud", map[string]any{
			"error": err.Error(),
		})
	}
}

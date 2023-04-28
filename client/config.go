package client

import (
	"crypto/tls"
	"errors"
	"net/url"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

var DefaultGrpcOpts = []grpc.DialOption{
	grpc.WithBlock(),
	grpc.WithKeepaliveParams(keepalive.ClientParameters{
		Time:                5 * time.Second,
		Timeout:             10 * time.Second,
		PermitWithoutStream: true,
	}),
	grpc.WithConnectParams(grpc.ConnectParams{
		Backoff: backoff.Config{
			BaseDelay:  100 * time.Millisecond,
			Multiplier: 1.6,
			Jitter:     0.2,
			MaxDelay:   3 * time.Second,
		},
		MinConnectTimeout: 3 * time.Second,
	}),
}

type CloudConfig struct {
	APIKey      string
	ClusterName string
}

type Config struct {
	Address       string // Remote address, "localhost:19530".
	Username      string // Username for auth.
	Password      string // Password for auth.
	DBName        string // DBName for this client.
	EnableTLSAuth bool   // Enable TLS Auth for transport security.

	Cloud *CloudConfig // configuration for cloud.

	DialOptions []grpc.DialOption // Dial options for GRPC.
}

func (c *Config) getDialOption() ([]grpc.DialOption, error) {
	options := c.DialOptions
	if c.DialOptions == nil {
		// Add default connection options.
		options = make([]grpc.DialOption, len(DefaultGrpcOpts))
		copy(options, DefaultGrpcOpts)
	}

	// Construct dial option.
	enableTLSAuth, err := c.parseEnableTLSAuth()
	if err != nil {
		return nil, err
	}
	if enableTLSAuth {
		options = append(options, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})))
	} else {
		options = append(options, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	// options = append(options,
	// 	grpc.WithChainUnaryInterceptor(grpc_retry.UnaryClientInterceptor(
	// 		grpc_retry.WithMax(6),
	// 		grpc_retry.WithBackoff(func(attempt uint) time.Duration {
	// 			return 60 * time.Millisecond * time.Duration(math.Pow(3, float64(attempt)))
	// 		}),
	// 		grpc_retry.WithCodes(codes.Unavailable, codes.ResourceExhausted)),
	// 		RetryOnRateLimitInterceptor(10, func(ctx context.Context, attempt uint) time.Duration {
	// 			return 10 * time.Millisecond * time.Duration(math.Pow(3, float64(attempt)))
	// 		}),
	// 	))

	// Construct username:password field
	if c.Username != "" && c.Password != "" {
		options = append(options,
			grpc.WithChainUnaryInterceptor(
				createAuthenticationUnaryInterceptor(c.Username, c.Password),
			),
			grpc.WithStreamInterceptor(createAuthenticationStreamInterceptor(c.Username, c.Password)),
		)
	} else if c.Cloud != nil {
		options = append(options,
			grpc.WithChainUnaryInterceptor(
				createCloudMetaUnaryInterceptor(c.Cloud),
			),
			grpc.WithStreamInterceptor(createCloudMetaStreamInterceptor(c.Cloud)),
		)
	}

	// Construct DBName field
	if c.DBName != "" {
		options = append(options,
			grpc.WithChainUnaryInterceptor(
				createDatabaseNameUnaryInterceptor(c.DBName),
			),
			grpc.WithStreamInterceptor(createDatabaseNameStreamInterceptor(c.DBName)),
		)
	}

	return options, nil
}

// Validate the config.
func (c *Config) validate() error {
	if c.Address == "" {
		return errors.New("empty remote address")
	}
	if _, err := c.parseRemoteAddr(); err != nil {
		return err
	}
	if _, err := c.parseEnableTLSAuth(); err != nil {
		return err
	}
	return nil
}

func (c *Config) parseRemoteAddr() (addr string, err error) {
	addr, _, err = parseURI(c.Address)
	return
}

func (c *Config) parseEnableTLSAuth() (enable bool, err error) {
	if c.EnableTLSAuth {
		return true, nil
	}
	_, secureHTTPS, err := parseURI(c.Address)
	return secureHTTPS, err
}

func parseURI(uri string) (string, bool, error) {
	hasPrefix := false
	inSecure := false
	if strings.HasPrefix(uri, "https://") {
		inSecure = true
		hasPrefix = true
	}

	if strings.HasPrefix(uri, "http://") {
		inSecure = false
		hasPrefix = true
	}

	if hasPrefix {
		url, err := url.Parse(uri)
		if err != nil {
			return "", inSecure, err
		}
		return url.Host, inSecure, nil
	}

	return uri, inSecure, nil
}

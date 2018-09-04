package aws

import (
	"context"
	"fmt"
	"net/url"
	"time"
	"unicode/utf8"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/pkg/errors"
	glooaws "github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1/plugins/aws"

	"github.com/solo-io/solo-kit/pkg/utils/contextutils"
	discovery "github.com/solo-io/solo-kit/projects/discovery/pkg"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1/plugins"
)

const (
	// expected map identifiers for secrets
	awsAccessKey = "access_key"
	awsSecretKey = "secret_key"
)

type AWSLambdaFuncitonDiscoveryFactory struct {
	PollingTime time.Duration
}

func (f *AWSLambdaFuncitonDiscoveryFactory) NewFunctionDiscovery(u *v1.Upstream) discovery.UpstreamFunctionDiscovery {
	return &AWSLambdaFuncitonDiscovery{
		timetowait: f.PollingTime,
		upstream:   u,
	}
}

type AWSLambdaFuncitonDiscovery struct {
	timetowait time.Duration
	upstream   *v1.Upstream
}

func (f *AWSLambdaFuncitonDiscovery) IsFunctional() bool {
	_, ok := f.upstream.UpstreamSpec.UpstreamType.(*v1.UpstreamSpec_Aws)
	return ok
}

func (f *AWSLambdaFuncitonDiscovery) DetectType(ctx context.Context, url *url.URL) (*plugins.ServiceSpec, error) {
	return nil, nil
}

// TODO: how to handle changes in secret or upstream (like the upstream ref)?
// perhaps the in param for the upstream should be a function? in func() *v1.Upstream
func (f *AWSLambdaFuncitonDiscovery) DetectFunctions(ctx context.Context, url *url.URL, secrets func() v1.SecretList, updatecb func(discovery.UpstreamMutator) error) error {
	for {
		// TODO: get backoff values from config?
		err := contextutils.NewExponentioalBackoff(contextutils.ExponentioalBackoff{}).Backoff(ctx, func(ctx context.Context) error {

			newfunctions, err := f.DetectFunctionsOnce(ctx, secrets)

			if err != nil {
				return err
			}

			err = updatecb(func(out *v1.Upstream) error {
				awsspec, ok := out.UpstreamSpec.UpstreamType.(*v1.UpstreamSpec_Aws)
				if !ok {
					return errors.New("not aws upstream")
				}
				awsspec.Aws.LambdaFunctions = newfunctions
				return nil
			})

			if err != nil {
				return errors.Wrap(err, "unable to update upstream")
			}
			return nil

		})
		if err != nil {
			if ctx.Err() != nil {
				return ctx.Err()
			}
			// ignore other erros as we would like to continue forever.
		}

		// sleep so we are not hogging
		if err := contextutils.Sleep(ctx, f.timetowait); err != nil {
			return err
		}
	}
}

func (f *AWSLambdaFuncitonDiscovery) DetectFunctionsOnce(ctx context.Context, secrets func() v1.SecretList) ([]*glooaws.LambdaFunctionSpec, error) {
	in := f.upstream
	awsspec, ok := in.UpstreamSpec.UpstreamType.(*v1.UpstreamSpec_Aws)

	if !ok {
		return nil, errors.New("not a lambda upstream spec")
	}
	lambdaSpec := awsspec.Aws
	awsSecrets, err := secrets().Find(in.Metadata.Namespace, lambdaSpec.SecretRef)
	if err != nil {
		return nil, errors.Wrapf(err, "secrets not found for secret ref %v", lambdaSpec.SecretRef)
	}

	accessKey, ok := awsSecrets.Data[awsAccessKey]
	if !ok {
		return nil, errors.Errorf("key %v missing from provided secret", awsAccessKey)
	}
	if accessKey != "" && !utf8.Valid([]byte(accessKey)) {
		return nil, errors.Errorf("%s not a valid string", awsAccessKey)
	}
	secretKey, ok := awsSecrets.Data[awsSecretKey]
	if !ok {
		return nil, errors.Errorf("key %v missing from provided secret", awsSecretKey)
	}
	if secretKey != "" && !utf8.Valid([]byte(secretKey)) {
		return nil, errors.Errorf("%s not a valid string", awsSecretKey)
	}

	sess, err := session.NewSession(aws.NewConfig().
		WithCredentials(credentials.NewStaticCredentials(accessKey, secretKey, "")))
	if err != nil {
		return nil, errors.Wrap(err, "unable to create AWS session")
	}
	svc := lambda.New(sess, &aws.Config{Region: aws.String(lambdaSpec.Region)})

	var newfunctions []*glooaws.LambdaFunctionSpec

	options := &lambda.ListFunctionsInput{FunctionVersion: aws.String("ALL")}
	err = svc.ListFunctionsPagesWithContext(ctx, options, func(results *lambda.ListFunctionsOutput, _ bool) bool {

		for _, f := range results.Functions {
			version := aws.StringValue(f.Version)
			name := aws.StringValue(f.FunctionName)

			logicalname := fmt.Sprintf("%s:%s", name, version)
			if version == "$LATEST" {
				logicalname = name
			}

			newfunctions = append(newfunctions, &glooaws.LambdaFunctionSpec{
				LambdaFunctionName: name,
				Qualifier:          version,
				LogicalName:        logicalname,
			})
		}

		return true
	})
	if err != nil {
		return nil, errors.Wrap(err, "unable to get list of functions from AWS")
	}

	return newfunctions, nil
}

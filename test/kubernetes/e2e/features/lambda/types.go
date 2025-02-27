package lambda

import (
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kgateway-dev/kgateway/v2/api/v1alpha1"
	"github.com/kgateway-dev/kgateway/v2/pkg/utils/fsutils"
)

var (
	setupManifest         = filepath.Join(fsutils.MustGetThisDir(), "testdata", "setup.yaml")
	lambdaBackendManifest = filepath.Join(fsutils.MustGetThisDir(), "testdata", "lambda-backend.yaml")
	lambdaAsyncManifest   = filepath.Join(fsutils.MustGetThisDir(), "testdata", "lambda-async.yaml")
	lambdaNamespace       = "lambda-test"

	lambdaBackendName       = "lambda-backend"
	lambdaBackendObjectMeta = metav1.ObjectMeta{
		Name:      lambdaBackendName,
		Namespace: lambdaNamespace,
	}
	lambdaBackend = &v1alpha1.Backend{ObjectMeta: lambdaBackendObjectMeta}

	gatewayName       = "lambda-gateway"
	gatewayObjectMeta = metav1.ObjectMeta{
		Name:      gatewayName,
		Namespace: lambdaNamespace,
	}
)

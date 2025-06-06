package install

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	appsv1api "k8s.io/api/apps/v1"
	corev1api "k8s.io/api/core/v1"
	apiextv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apiextv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	v1crds "github.com/vmware-tanzu/velero/config/crd/v1/crds"
	"github.com/vmware-tanzu/velero/pkg/test"
)

func TestInstall(t *testing.T) {
	dc := &test.FakeDynamicClient{}
	dc.On("Create", mock.Anything).Return(&unstructured.Unstructured{}, nil)

	factory := &test.FakeDynamicFactory{}
	factory.On("ClientForGroupVersionResource", mock.Anything, mock.Anything, mock.Anything).Return(dc, nil)

	c := fake.NewClientBuilder().WithObjects(
		&apiextv1.CustomResourceDefinition{
			ObjectMeta: metav1.ObjectMeta{
				Name: "backuprepositories.velero.io",
			},

			Status: apiextv1.CustomResourceDefinitionStatus{
				Conditions: []apiextv1.CustomResourceDefinitionCondition{
					{
						Type:   apiextv1.Established,
						Status: apiextv1.ConditionTrue,
					},
					{
						Type:   apiextv1.NamesAccepted,
						Status: apiextv1.ConditionTrue,
					},
				},
			},
		},
	).Build()

	resources := &unstructured.UnstructuredList{}
	require.NoError(t, appendUnstructured(resources, v1crds.CRDs[0]))
	require.NoError(t, appendUnstructured(resources, Namespace("velero")))

	assert.NoError(t, Install(factory, c, resources, os.Stdout))
}

func Test_crdsAreReady(t *testing.T) {
	c := fake.NewClientBuilder().WithObjects(
		&apiextv1beta1.CustomResourceDefinition{
			ObjectMeta: metav1.ObjectMeta{
				Name: "backuprepositories.velero.io",
			},

			Status: apiextv1beta1.CustomResourceDefinitionStatus{
				Conditions: []apiextv1beta1.CustomResourceDefinitionCondition{
					{
						Type:   apiextv1beta1.Established,
						Status: apiextv1beta1.ConditionTrue,
					},
					{
						Type:   apiextv1beta1.NamesAccepted,
						Status: apiextv1beta1.ConditionTrue,
					},
				},
			},
		},
	).Build()

	crd := &apiextv1beta1.CustomResourceDefinition{
		TypeMeta: metav1.TypeMeta{
			Kind:       "CustomResourceDefinition",
			APIVersion: "v1beta1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "backuprepositories.velero.io",
		},
	}
	obj, err := runtime.DefaultUnstructuredConverter.ToUnstructured(crd)
	require.NoError(t, err)

	crds := []*unstructured.Unstructured{
		{
			Object: obj,
		},
	}

	ready, err := crdsAreReady(c, crds)
	require.NoError(t, err)
	assert.True(t, ready)
}

func TestDeploymentIsReady(t *testing.T) {
	deployment := &appsv1api.Deployment{
		Status: appsv1api.DeploymentStatus{
			Conditions: []appsv1api.DeploymentCondition{
				{
					Type:               appsv1api.DeploymentAvailable,
					Status:             corev1api.ConditionTrue,
					LastTransitionTime: metav1.NewTime(time.Now().Add(-15 * time.Second)),
				},
			},
		},
	}
	obj, err := runtime.DefaultUnstructuredConverter.ToUnstructured(deployment)
	require.NoError(t, err)

	dc := &test.FakeDynamicClient{}
	dc.On("Get", mock.Anything, mock.Anything).Return(&unstructured.Unstructured{Object: obj}, nil)

	factory := &test.FakeDynamicFactory{}
	factory.On("ClientForGroupVersionResource", mock.Anything, mock.Anything, mock.Anything).Return(dc, nil)

	ready, err := DeploymentIsReady(factory, "velero")
	require.NoError(t, err)
	assert.True(t, ready)
}

func TestNodeAgentIsReady(t *testing.T) {
	daemonset := &appsv1api.DaemonSet{
		Status: appsv1api.DaemonSetStatus{
			NumberAvailable:        1,
			DesiredNumberScheduled: 1,
		},
	}
	obj, err := runtime.DefaultUnstructuredConverter.ToUnstructured(daemonset)
	require.NoError(t, err)

	dc := &test.FakeDynamicClient{}
	dc.On("Get", mock.Anything, mock.Anything).Return(&unstructured.Unstructured{Object: obj}, nil)

	factory := &test.FakeDynamicFactory{}
	factory.On("ClientForGroupVersionResource", mock.Anything, mock.Anything, mock.Anything).Return(dc, nil)

	ready, err := NodeAgentIsReady(factory, "velero")
	require.NoError(t, err)
	assert.True(t, ready)
}

func TestNodeAgentWindowsIsReady(t *testing.T) {
	daemonset := &appsv1api.DaemonSet{
		Status: appsv1api.DaemonSetStatus{
			NumberAvailable:        0,
			DesiredNumberScheduled: 0,
		},
	}
	obj, err := runtime.DefaultUnstructuredConverter.ToUnstructured(daemonset)
	require.NoError(t, err)

	dc := &test.FakeDynamicClient{}
	dc.On("Get", mock.Anything, mock.Anything).Return(&unstructured.Unstructured{Object: obj}, nil)

	factory := &test.FakeDynamicFactory{}
	factory.On("ClientForGroupVersionResource", mock.Anything, mock.Anything, mock.Anything).Return(dc, nil)

	ready, err := NodeAgentWindowsIsReady(factory, "velero")
	require.NoError(t, err)
	assert.True(t, ready)
}

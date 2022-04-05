package kubernetes

import (
	"github.com/stretchr/testify/suite"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/components/pods"
	corev1 "k8s.io/api/core/v1"
	"testing"
)

func TestPodVolumesSuite(t *testing.T) {
	suite.Run(t, &PodVolumesTestSuite{})
}

type PodVolumesTestSuite struct {
	suite.Suite
	volumeName      string
	volumeMountPath string
	volumeSubPath   string
	volumeReadOnly  bool
	base            pods.VolumeBase
}

func (suite *PodVolumesTestSuite) SetupSuite() {
	suite.volumeName = "test"
	suite.volumeMountPath = "/tmp/test"
	suite.volumeSubPath = "subpath/test"
	suite.volumeReadOnly = true

	suite.base = pods.VolumeBase{
		Name:      suite.volumeName,
		MountPath: suite.volumeMountPath,
		SubPath:   suite.volumeSubPath,
		ReadOnly:  suite.volumeReadOnly,
	}
}

func (suite *PodVolumesTestSuite) TestParseVolumeMount() {
	expected := corev1.VolumeMount{
		Name:      suite.volumeName,
		ReadOnly:  suite.volumeReadOnly,
		MountPath: suite.volumeMountPath,
		SubPath:   suite.volumeSubPath,
	}

	out := ParseVolumeMount(suite.base)

	suite.Equal(expected, out)
}

func (suite *PodVolumesTestSuite) TestParseVolumeHostPath() {
	in := pods.VolumeHostPath{
		VolumeBase:   suite.base,
		HostPath:     "/tmp/test",
		HostPathType: pods.HostPathUnset,
	}

	expectedType := corev1.HostPathType(in.HostPathType)
	expected := corev1.Volume{
		Name: suite.volumeName,
		VolumeSource: corev1.VolumeSource{
			HostPath: &corev1.HostPathVolumeSource{
				Path: in.HostPath,
				Type: &expectedType,
			},
		},
	}

	out := ParseVolume(in)

	suite.Equal(expected.Name, out.Name)
	suite.Equal(expected.VolumeSource.HostPath, out.VolumeSource.HostPath)
}

func (suite *PodVolumesTestSuite) TestParseVolumeConfiguration() {
	in := pods.VolumeConfiguration{
		VolumeBase:        suite.base,
		ConfigurationName: "test-config",
		Items: map[string]string{
			"1": "1",
			"2": "2",
			"3": "3",
		},
	}

	items := make([]corev1.KeyToPath, 0, len(in.Items))
	for key, path := range in.Items {
		items = append(
			items, corev1.KeyToPath{
				Key:  key,
				Path: path,
			},
		)
	}

	expected := corev1.Volume{
		Name: suite.volumeName,
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: in.ConfigurationName,
				},
				Items: items,
			},
		},
	}

	out := ParseVolume(in)

	suite.Equal(expected.Name, out.Name)
	suite.Equal(
		expected.VolumeSource.ConfigMap.LocalObjectReference,
		out.VolumeSource.ConfigMap.LocalObjectReference,
	)
	suite.Equal(
		len(expected.VolumeSource.ConfigMap.Items),
		len(out.VolumeSource.ConfigMap.Items),
	)
}

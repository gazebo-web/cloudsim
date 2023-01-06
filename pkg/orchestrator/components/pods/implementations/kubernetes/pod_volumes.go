package kubernetes

import (
	"github.com/gazebo-web/cloudsim/pkg/orchestrator/components/pods"
	corev1 "k8s.io/api/core/v1"
)

// ParseVolume parses a generic pods.Volume and returns a Kubernetes corev1.Volume instance.
func ParseVolume(volume pods.Volume) corev1.Volume {
	kv := corev1.Volume{
		Name: volume.Base().Name,
	}

	// Fill in VolumeSource based on the specific type of generic Volume
	switch v := volume.(type) {
	case pods.VolumeHostPath:
		hostPathType := corev1.HostPathType(v.HostPathType)
		kv.HostPath = &corev1.HostPathVolumeSource{
			Path: v.HostPath,
			Type: &hostPathType,
		}

	case pods.VolumeConfiguration:
		items := make([]corev1.KeyToPath, 0, len(v.Items))
		for key, path := range v.Items {
			items = append(items, corev1.KeyToPath{
				Key:  key,
				Path: path,
			})
		}

		kv.ConfigMap = &corev1.ConfigMapVolumeSource{
			LocalObjectReference: corev1.LocalObjectReference{
				Name: v.ConfigurationName,
			},
			Items: items,
		}
	default:
		panic("kubernetes volume type not implemented")
	}

	return kv
}

// ParseVolumeMount parses a generic pods.Volume and returns a Kubernetes corev1.VolumeMount instance.
func ParseVolumeMount(volume pods.Volume) corev1.VolumeMount {
	base := volume.Base()
	return corev1.VolumeMount{
		Name:      base.Name,
		ReadOnly:  base.ReadOnly,
		MountPath: base.MountPath,
		SubPath:   base.SubPath,
	}
}

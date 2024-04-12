package pg

import (
	v1 "k8s.io/api/apps/v1"
	v13 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"

	postgresv1alpha1 "github.com/vygos/postgres-operator/api/v1alpha1"
)

func CreateStatefulSet(postgres postgresv1alpha1.PostgresSQL, namespace string) *v1.StatefulSet {
	replicas := int32(1)
	return &v1.StatefulSet{
		ObjectMeta: v12.ObjectMeta{
			Name:      postgres.Name,
			Namespace: namespace,
		},
		Spec: v1.StatefulSetSpec{
			Replicas:    &replicas,
			ServiceName: "postgres",
			Selector:    postgres.Spec.Selector,
			Template: v13.PodTemplateSpec{
				ObjectMeta: v12.ObjectMeta{
					Labels: postgres.Spec.Selector.MatchLabels,
				},
				Spec: v13.PodSpec{
					Containers: []v13.Container{
						{
							ImagePullPolicy: v13.PullIfNotPresent,
							Name:            "postgres",
							Ports: []v13.ContainerPort{
								{
									ContainerPort: int32(5432),
								},
							},
							Image: "postgres",
							VolumeMounts: []v13.VolumeMount{
								{
									MountPath: "var/lib/postgresql/data",
									Name:      "postgres-data",
								},
							},
							Env: postgres.Spec.Env,
						},
					},
				},
			},
			VolumeClaimTemplates: []v13.PersistentVolumeClaim{
				{
					ObjectMeta: v12.ObjectMeta{
						Name:   "postgres-data",
						Labels: postgres.Spec.Selector.MatchLabels,
					},
					Spec: v13.PersistentVolumeClaimSpec{
						Selector:    postgres.Spec.Selector,
						AccessModes: []v13.PersistentVolumeAccessMode{v13.ReadWriteOnce},
						Resources: v13.ResourceRequirements{
							Requests: v13.ResourceList{
								"storage": resource.MustParse("1Gi"),
							},
						},
					},
				},
			},
		},
	}
}

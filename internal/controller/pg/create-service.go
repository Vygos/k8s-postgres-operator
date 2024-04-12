package pg

import (
	"fmt"

	v1 "k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/vygos/postgres-operator/api/v1alpha1"
)

const SvcSufix = "service"

func CreateService(postgres v1alpha1.PostgresSQL, namespace string) *v1.Service {
	return &v1.Service{
		ObjectMeta: v12.ObjectMeta{
			Name:      fmt.Sprintf("%s-%s", postgres.Name, SvcSufix),
			Namespace: namespace,
		},
		Spec: v1.ServiceSpec{
			Selector:  postgres.Spec.Selector.MatchLabels,
			ClusterIP: "None",
			Ports: []v1.ServicePort{
				{
					Port:       5432,
					TargetPort: intstr.FromInt(5432),
				},
			},
		},
	}
}

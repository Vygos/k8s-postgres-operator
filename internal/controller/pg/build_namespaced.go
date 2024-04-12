package pg

import (
	"fmt"

	"k8s.io/apimachinery/pkg/types"
)

func BuildNamespacedSvcName(
	namespace string,
	objName string,
) types.NamespacedName {

	return types.NamespacedName{
		Namespace: namespace,
		Name:      fmt.Sprintf("%s-%s", objName, SvcSufix),
	}
}

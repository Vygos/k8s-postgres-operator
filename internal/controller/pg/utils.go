package pg

import (
	"fmt"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
)

func CreateIfNotFound(execute func() error, err error) (ctrl.Result, error) {
	if errors.IsNotFound(err) {
		err := execute()
		if err != nil {
			return ctrl.Result{}, err
		}

		return ctrl.Result{Requeue: true}, nil
	}

	return ctrl.Result{}, err
}

func BuildNamespacedSvcName(
	namespace string,
	objName string,
) types.NamespacedName {

	return types.NamespacedName{
		Namespace: namespace,
		Name:      fmt.Sprintf("%s-%s", objName, SvcSufix),
	}
}

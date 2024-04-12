package pg

import (
	"k8s.io/apimachinery/pkg/api/errors"
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

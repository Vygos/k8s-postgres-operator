package glog

import (
	"context"
	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

func Log(ctx context.Context) logr.Logger {
	return log.FromContext(ctx)
}

package validation

import (
	"context"
	rbac "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/apiserver/pkg/storage/names"
)

type Validator[T any] func(T) field.ErrorList

type Strategy[T any] struct {
	runtime.ObjectTyper
	names.NameGenerator
	NamespacedScoped bool
	Validator        func(x *T) field.ErrorList
}

func ValidateClusterRoleNew(role *rbac.ClusterRole) field.ErrorList {
	return ValidateClusterRole(role, ClusterRoleValidationOptions{
		AllowInvalidLabelValueInSelector: false,
	})
}

func (s Strategy[T]) NamespaceScoped() bool {
	return s.NamespacedScoped
}

func (s Strategy[T]) WarningsOnCreate(ctx context.Context, obj runtime.Object) []string {
	return nil
}

func (s Strategy[T]) Canonicalize(obj runtime.Object) {}

func (s Strategy[T]) PrepareForCreate(ctx context.Context, obj runtime.Object) {}

func (s Strategy[T]) Validate(ctx context.Context, obj runtime.Object) field.ErrorList {
	u := obj.(*unstructured.Unstructured).Object
	var v T
	_ = runtime.DefaultUnstructuredConverter.FromUnstructuredWithValidation(u, &v, false)
	return s.Validator(&v)
}

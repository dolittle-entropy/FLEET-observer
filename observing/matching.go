package observing

import (
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
)

func createEnvironmentLabelSelectorFor(object metaV1.ObjectMeta) (labels.Selector, error) {
	tenant, err := labels.NewRequirement("tenant", selection.Equals, []string{object.GetLabels()["tenant"]})
	if err != nil {
		return nil, err
	}
	application, err := labels.NewRequirement("application", selection.Equals, []string{object.GetLabels()["application"]})
	if err != nil {
		return nil, err
	}
	environment, err := labels.NewRequirement("environment", selection.Equals, []string{object.GetLabels()["environment"]})
	if err != nil {
		return nil, err
	}

	return labels.NewSelector().Add(*tenant, *application, *environment), nil
}

func createMicroserviceLabelSelectorFor(object metaV1.ObjectMeta) (labels.Selector, error) {
	microservice, err := labels.NewRequirement("microservice", selection.Equals, []string{object.GetLabels()["microservice"]})
	if err != nil {
		return nil, err
	}

	environmentSelector, err := createEnvironmentLabelSelectorFor(object)
	if err != nil {
		return nil, err
	}

	return environmentSelector.Add(*microservice), nil
}

func environmentEquals(left, right metaV1.ObjectMeta) bool {
	if left.GetAnnotations()["dolittle.io/tenant-id"] != right.GetAnnotations()["dolittle.io/tenant-id"] {
		return false
	}
	if left.GetAnnotations()["dolittle.io/application-id"] != right.GetAnnotations()["dolittle.io/application-id"] {
		return false
	}
	if left.GetLabels()["tenant"] != right.GetLabels()["tenant"] {
		return false
	}
	if left.GetLabels()["application"] != right.GetLabels()["application"] {
		return false
	}
	if left.GetLabels()["environment"] != right.GetLabels()["environment"] {
		return false
	}
	return true
}

func microserviceEquals(left, right metaV1.ObjectMeta) bool {
	if !environmentEquals(left, right) {
		return false
	}
	if left.GetAnnotations()["dolittle.io/microservice-id"] != right.GetAnnotations()["dolittle.io/microservice-id"] {
		return false
	}
	if left.GetLabels()["microservice"] != right.GetLabels()["microservice"] {
		return false
	}
	return true
}

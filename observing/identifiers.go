package observing

import metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"

func GetMicroserviceIdentifiers(meta metaV1.ObjectMeta) (tenantID, applicationID, environmentName, microserviceID string, ok bool) {
	tenantID, ok = meta.GetAnnotations()["dolittle.io/tenant-id"]
	if !ok {
		return
	}

	applicationID, ok = meta.GetAnnotations()["dolittle.io/application-id"]
	if !ok {
		return
	}

	environmentName, ok = meta.GetLabels()["environment"]
	if !ok {
		return
	}

	microserviceID, ok = meta.GetAnnotations()["dolittle.io/microservice-id"]
	if !ok {
		return
	}

	return
}

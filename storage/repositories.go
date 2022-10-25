/*
 * Copyright (c) Dolittle. All rights reserved.
 * Licensed under the MIT license. See LICENSE file in the project root for full license information.
 */

package storage

type Repositories struct {
	Nodes          Nodes
	Customers      Customers
	Applications   Applications
	Environments   Environments
	Artifacts      Artifacts
	Runtimes       Runtimes
	Deployments    Deployments
	Configurations Configurations
	Events         Events
}

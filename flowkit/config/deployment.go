/*
 * Flow CLI
 *
 * Copyright 2019 Dapper Labs, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package config

import (
	"fmt"

	"github.com/onflow/cadence"
	"golang.org/x/exp/slices"
)

// ContractDeployment defines the deployment of the contract with possible args.
type ContractDeployment struct {
	Name string
	Args []cadence.Value
}

// Deployment defines the configuration for a contract deployment.
type Deployment struct {
	Network   string               // network name to deploy to
	Account   string               // account name to which to deploy to
	Contracts []ContractDeployment // contracts to deploy
}

// AddContract to deployment list on the account name and network name.
func (d *Deployment) AddContract(contract ContractDeployment) {
	for _, c := range d.Contracts {
		if contract.Name == c.Name {
			return // don't allow adding duplicates
		}
	}

	d.Contracts = append(d.Contracts, contract)
}

// RemoveContract removes a specific contract by name from an existing deployment identified by account name and network name.
func (d *Deployment) RemoveContract(contractName string) {
	for i, contract := range d.Contracts {
		if contract.Name == contractName {
			d.Contracts = slices.Delete(d.Contracts, i, i+1)
		}
	}
}

type Deployments []Deployment

// ByNetwork get all deployments by network.
func (d *Deployments) ByNetwork(network string) Deployments {
	var deployments Deployments

	for _, deploy := range *d {
		if deploy.Network == network {
			deployments = append(deployments, deploy)
		}
	}

	return deployments
}

// ByAccountAndNetwork get deploy by account and network.
func (d *Deployments) ByAccountAndNetwork(account string, network string) *Deployment {
	for i, deploy := range *d {
		if deploy.Network == network && deploy.Account == account {
			return &(*d)[i]
		}
	}

	return nil
}

// AddOrUpdate add new or update if already present.
func (d *Deployments) AddOrUpdate(deployment Deployment) {
	for i, existingDeployment := range *d {
		if existingDeployment.Account == deployment.Account &&
			existingDeployment.Network == deployment.Network {
			(*d)[i] = deployment
			return
		}
	}

	*d = append(*d, deployment)
}

// Remove removes deployment by account and network.
func (d *Deployments) Remove(account string, network string) error {
	deployment := d.ByAccountAndNetwork(account, network)
	if deployment == nil {
		return fmt.Errorf(
			"deployment for account %s on network %s does not exist in configuration",
			account,
			network,
		)
	}

	for i, deployment := range *d {
		if deployment.Network == network && deployment.Account == account {
			*d = slices.Delete(*d, i, i+1)
		}
	}

	return nil
}

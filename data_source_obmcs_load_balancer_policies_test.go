// Copyright (c) 2017, Oracle and/or its affiliates. All rights reserved.

package main

import (
	"testing"

	"github.com/MustWin/baremetal-sdk-go"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func TestLoadBalancerPoliciesDatasource(t *testing.T) {
	client := GetTestProvider()
	providers := map[string]terraform.ResourceProvider{
		"baremetal": Provider(func(d *schema.ResourceData) (interface{}, error) {
			return client, nil
		}),
	}
	resourceName := "data.baremetal_load_balancer_policies.t"
	config := `
data "baremetal_load_balancer_policies" "t" {
  compartment_id = "${var.compartment_id}"
}
`
	config += testProviderConfig()

	compartmentID := "${var.compartment_id}"
	list := &baremetal.ListLoadBalancerPolicies{
		LoadBalancerPolicies: []baremetal.LoadBalancerPolicy{
			{Name: "stub_name1"},
			{Name: "stub_name2"},
		},
	}
	client.On(
		"ListLoadBalancerPolicies",
		compartmentID,
		(*baremetal.ListLoadBalancerPolicyOptions)(nil),
	).Return(list, nil)

	resource.UnitTest(t, resource.TestCase{
		PreventPostDestroyRefresh: true,
		Providers:                 providers,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "compartment_id", compartmentID),
					resource.TestCheckResourceAttr(resourceName, "policies.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "policies.0.name", "stub_name1"),
					resource.TestCheckResourceAttr(resourceName, "policies.1.name", "stub_name2"),
				),
			},
		},
	})
}

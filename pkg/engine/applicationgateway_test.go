// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package engine

import (
	"testing"

	"github.com/Azure/aks-engine-azurestack/pkg/api"
	"github.com/Azure/aks-engine-azurestack/pkg/helpers"

	"github.com/google/go-cmp/cmp"

	"github.com/Azure/aks-engine-azurestack/pkg/api/common"
	"github.com/Azure/azure-sdk-for-go/profiles/2020-09-01/network/mgmt/network"
)

func TestCreateApplicationGateway(t *testing.T) {
	cs := &api.ContainerService{
		Properties: &api.Properties{
			OrchestratorProfile: &api.OrchestratorProfile{
				KubernetesConfig: &api.KubernetesConfig{},
			},
		},
	}
	actual := createApplicationGateway(cs.Properties)

	expected := ApplicationGatewayARM{
		ARMResource: ARMResource{
			APIVersion: "[variables('apiVersionNetwork')]",
			DependsOn: []string{
				"[concat('Microsoft.Network/publicIPAddresses/', variables('appGwPublicIPAddressName'))]",
				"[concat('Microsoft.Network/virtualNetworks/', variables('virtualNetworkName'))]",
			},
		},
		ApplicationGateway: network.ApplicationGateway{
			Location: helpers.PointerToString("[variables('location')]"),
			Name:     helpers.PointerToString("[variables('appGwName')]"),
			ApplicationGatewayPropertiesFormat: &network.ApplicationGatewayPropertiesFormat{
				Sku: &network.ApplicationGatewaySku{
					Name:     network.ApplicationGatewaySkuName("[parameters('appGwSku')]"),
					Tier:     network.ApplicationGatewayTier("[parameters('appGwSku')]"),
					Capacity: helpers.PointerToInt32(2),
				},
				GatewayIPConfigurations: &[]network.ApplicationGatewayIPConfiguration{
					{
						Name: helpers.PointerToString("gatewayIP"),
						ApplicationGatewayIPConfigurationPropertiesFormat: &network.ApplicationGatewayIPConfigurationPropertiesFormat{
							Subnet: &network.SubResource{
								ID: helpers.PointerToString("[concat(variables('vnetID'),'/subnets/',variables('appGwSubnetName'))]"),
							},
						},
					},
				},
				FrontendIPConfigurations: &[]network.ApplicationGatewayFrontendIPConfiguration{
					{
						Name: helpers.PointerToString("frontendIP"),
						ApplicationGatewayFrontendIPConfigurationPropertiesFormat: &network.ApplicationGatewayFrontendIPConfigurationPropertiesFormat{
							PublicIPAddress: &network.SubResource{
								ID: helpers.PointerToString("[resourceId('Microsoft.Network/publicIpAddresses',variables('appGwPublicIPAddressName'))]"),
							},
						},
					},
				},
				FrontendPorts: &[]network.ApplicationGatewayFrontendPort{
					{
						Name: helpers.PointerToString("httpPort"),
						ApplicationGatewayFrontendPortPropertiesFormat: &network.ApplicationGatewayFrontendPortPropertiesFormat{
							Port: helpers.PointerToInt32(80),
						},
					},
				},
				BackendAddressPools: &[]network.ApplicationGatewayBackendAddressPool{
					{
						Name: helpers.PointerToString("pool"),
						ApplicationGatewayBackendAddressPoolPropertiesFormat: &network.ApplicationGatewayBackendAddressPoolPropertiesFormat{
							BackendAddresses: &[]network.ApplicationGatewayBackendAddress{},
						},
					},
				},
				HTTPListeners: &[]network.ApplicationGatewayHTTPListener{
					{
						Name: helpers.PointerToString("httpListener"),
						ApplicationGatewayHTTPListenerPropertiesFormat: &network.ApplicationGatewayHTTPListenerPropertiesFormat{
							Protocol: network.HTTP,
							FrontendPort: &network.SubResource{
								ID: helpers.PointerToString("[concat(variables('appGwId'), '/frontendPorts/httpPort')]"),
							},
							FrontendIPConfiguration: &network.SubResource{
								ID: helpers.PointerToString("[concat(variables('appGwId'), '/frontendIPConfigurations/frontendIP')]"),
							},
						},
					},
				},
				BackendHTTPSettingsCollection: &[]network.ApplicationGatewayBackendHTTPSettings{
					{
						Name: helpers.PointerToString("setting"),
						ApplicationGatewayBackendHTTPSettingsPropertiesFormat: &network.ApplicationGatewayBackendHTTPSettingsPropertiesFormat{
							Port:     helpers.PointerToInt32(80),
							Protocol: network.HTTP,
						},
					},
				},
				RequestRoutingRules: &[]network.ApplicationGatewayRequestRoutingRule{
					{
						Name: helpers.PointerToString("rule"),
						ApplicationGatewayRequestRoutingRulePropertiesFormat: &network.ApplicationGatewayRequestRoutingRulePropertiesFormat{
							HTTPListener: &network.SubResource{
								ID: helpers.PointerToString("[concat(variables('appGwId'), '/httpListeners/httpListener')]"),
							},
							BackendAddressPool: &network.SubResource{
								ID: helpers.PointerToString("[concat(variables('appGwId'), '/backendAddressPools/pool')]"),
							},
							BackendHTTPSettings: &network.SubResource{
								ID: helpers.PointerToString("[concat(variables('appGwId'), '/backendHttpSettingsCollection/setting')]"),
							},
						},
					},
				},
			},
			Type: helpers.PointerToString("Microsoft.Network/applicationGateways"),
		},
	}

	diff := cmp.Diff(actual, expected)

	if diff != "" {
		t.Errorf("unexpected error while comparing application gateways: %s", diff)
	}

}

func TestCreateApplicationGatewayWAF(t *testing.T) {
	cs := &api.ContainerService{
		Properties: &api.Properties{
			OrchestratorProfile: &api.OrchestratorProfile{
				KubernetesConfig: &api.KubernetesConfig{
					Addons: []api.KubernetesAddon{
						{
							Name:    common.AppGwIngressAddonName,
							Enabled: helpers.PointerToBool(true),
							Config: map[string]string{
								"appgw-sku": "WAF_v2",
							},
						},
					},
				},
			},
		},
	}
	actual := createApplicationGateway(cs.Properties)

	expected := ApplicationGatewayARM{
		ARMResource: ARMResource{
			APIVersion: "[variables('apiVersionNetwork')]",
			DependsOn: []string{
				"[concat('Microsoft.Network/publicIPAddresses/', variables('appGwPublicIPAddressName'))]",
				"[concat('Microsoft.Network/virtualNetworks/', variables('virtualNetworkName'))]",
			},
		},
		ApplicationGateway: network.ApplicationGateway{
			Location: helpers.PointerToString("[variables('location')]"),
			Name:     helpers.PointerToString("[variables('appGwName')]"),
			ApplicationGatewayPropertiesFormat: &network.ApplicationGatewayPropertiesFormat{
				Sku: &network.ApplicationGatewaySku{
					Name:     network.ApplicationGatewaySkuName("[parameters('appGwSku')]"),
					Tier:     network.ApplicationGatewayTier("[parameters('appGwSku')]"),
					Capacity: helpers.PointerToInt32(2),
				},
				GatewayIPConfigurations: &[]network.ApplicationGatewayIPConfiguration{
					{
						Name: helpers.PointerToString("gatewayIP"),
						ApplicationGatewayIPConfigurationPropertiesFormat: &network.ApplicationGatewayIPConfigurationPropertiesFormat{
							Subnet: &network.SubResource{
								ID: helpers.PointerToString("[concat(variables('vnetID'),'/subnets/',variables('appGwSubnetName'))]"),
							},
						},
					},
				},
				FrontendIPConfigurations: &[]network.ApplicationGatewayFrontendIPConfiguration{
					{
						Name: helpers.PointerToString("frontendIP"),
						ApplicationGatewayFrontendIPConfigurationPropertiesFormat: &network.ApplicationGatewayFrontendIPConfigurationPropertiesFormat{
							PublicIPAddress: &network.SubResource{
								ID: helpers.PointerToString("[resourceId('Microsoft.Network/publicIpAddresses',variables('appGwPublicIPAddressName'))]"),
							},
						},
					},
				},
				FrontendPorts: &[]network.ApplicationGatewayFrontendPort{
					{
						Name: helpers.PointerToString("httpPort"),
						ApplicationGatewayFrontendPortPropertiesFormat: &network.ApplicationGatewayFrontendPortPropertiesFormat{
							Port: helpers.PointerToInt32(80),
						},
					},
				},
				BackendAddressPools: &[]network.ApplicationGatewayBackendAddressPool{
					{
						Name: helpers.PointerToString("pool"),
						ApplicationGatewayBackendAddressPoolPropertiesFormat: &network.ApplicationGatewayBackendAddressPoolPropertiesFormat{
							BackendAddresses: &[]network.ApplicationGatewayBackendAddress{},
						},
					},
				},
				HTTPListeners: &[]network.ApplicationGatewayHTTPListener{
					{
						Name: helpers.PointerToString("httpListener"),
						ApplicationGatewayHTTPListenerPropertiesFormat: &network.ApplicationGatewayHTTPListenerPropertiesFormat{
							Protocol: network.HTTP,
							FrontendPort: &network.SubResource{
								ID: helpers.PointerToString("[concat(variables('appGwId'), '/frontendPorts/httpPort')]"),
							},
							FrontendIPConfiguration: &network.SubResource{
								ID: helpers.PointerToString("[concat(variables('appGwId'), '/frontendIPConfigurations/frontendIP')]"),
							},
						},
					},
				},
				BackendHTTPSettingsCollection: &[]network.ApplicationGatewayBackendHTTPSettings{
					{
						Name: helpers.PointerToString("setting"),
						ApplicationGatewayBackendHTTPSettingsPropertiesFormat: &network.ApplicationGatewayBackendHTTPSettingsPropertiesFormat{
							Port:     helpers.PointerToInt32(80),
							Protocol: network.HTTP,
						},
					},
				},
				RequestRoutingRules: &[]network.ApplicationGatewayRequestRoutingRule{
					{
						Name: helpers.PointerToString("rule"),
						ApplicationGatewayRequestRoutingRulePropertiesFormat: &network.ApplicationGatewayRequestRoutingRulePropertiesFormat{
							HTTPListener: &network.SubResource{
								ID: helpers.PointerToString("[concat(variables('appGwId'), '/httpListeners/httpListener')]"),
							},
							BackendAddressPool: &network.SubResource{
								ID: helpers.PointerToString("[concat(variables('appGwId'), '/backendAddressPools/pool')]"),
							},
							BackendHTTPSettings: &network.SubResource{
								ID: helpers.PointerToString("[concat(variables('appGwId'), '/backendHttpSettingsCollection/setting')]"),
							},
						},
					},
				},
				WebApplicationFirewallConfiguration: &network.ApplicationGatewayWebApplicationFirewallConfiguration{
					Enabled:      helpers.PointerToBool(true),
					FirewallMode: network.Detection,
				},
			},
			Type: helpers.PointerToString("Microsoft.Network/applicationGateways"),
		},
	}

	diff := cmp.Diff(actual, expected)

	if diff != "" {
		t.Errorf("unexpected error while comparing application gateways: %s", diff)
	}

}

func TestCreateApplicationGatewayPrivateIP(t *testing.T) {
	cs := &api.ContainerService{
		Properties: &api.Properties{
			OrchestratorProfile: &api.OrchestratorProfile{
				KubernetesConfig: &api.KubernetesConfig{
					Addons: []api.KubernetesAddon{
						{
							Name:    common.AppGwIngressAddonName,
							Enabled: helpers.PointerToBool(true),
							Config: map[string]string{
								"appgw-private-ip": "10.0.0.1",
							},
						},
					},
				},
			},
		},
	}
	actual := createApplicationGateway(cs.Properties)

	expected := ApplicationGatewayARM{
		ARMResource: ARMResource{
			APIVersion: "[variables('apiVersionNetwork')]",
			DependsOn: []string{
				"[concat('Microsoft.Network/publicIPAddresses/', variables('appGwPublicIPAddressName'))]",
				"[concat('Microsoft.Network/virtualNetworks/', variables('virtualNetworkName'))]",
			},
		},
		ApplicationGateway: network.ApplicationGateway{
			Location: helpers.PointerToString("[variables('location')]"),
			Name:     helpers.PointerToString("[variables('appGwName')]"),
			ApplicationGatewayPropertiesFormat: &network.ApplicationGatewayPropertiesFormat{
				Sku: &network.ApplicationGatewaySku{
					Name:     network.ApplicationGatewaySkuName("[parameters('appGwSku')]"),
					Tier:     network.ApplicationGatewayTier("[parameters('appGwSku')]"),
					Capacity: helpers.PointerToInt32(2),
				},
				GatewayIPConfigurations: &[]network.ApplicationGatewayIPConfiguration{
					{
						Name: helpers.PointerToString("gatewayIP"),
						ApplicationGatewayIPConfigurationPropertiesFormat: &network.ApplicationGatewayIPConfigurationPropertiesFormat{
							Subnet: &network.SubResource{
								ID: helpers.PointerToString("[concat(variables('vnetID'),'/subnets/',variables('appGwSubnetName'))]"),
							},
						},
					},
				},
				FrontendIPConfigurations: &[]network.ApplicationGatewayFrontendIPConfiguration{
					{
						Name: helpers.PointerToString("frontendIP"),
						ApplicationGatewayFrontendIPConfigurationPropertiesFormat: &network.ApplicationGatewayFrontendIPConfigurationPropertiesFormat{
							PublicIPAddress: &network.SubResource{
								ID: helpers.PointerToString("[resourceId('Microsoft.Network/publicIpAddresses',variables('appGwPublicIPAddressName'))]"),
							},
						},
					},
					{
						Name: helpers.PointerToString("privateIp"),
						ApplicationGatewayFrontendIPConfigurationPropertiesFormat: &network.ApplicationGatewayFrontendIPConfigurationPropertiesFormat{
							PrivateIPAddress: helpers.PointerToString("10.0.0.1"),
						},
					},
				},
				FrontendPorts: &[]network.ApplicationGatewayFrontendPort{
					{
						Name: helpers.PointerToString("httpPort"),
						ApplicationGatewayFrontendPortPropertiesFormat: &network.ApplicationGatewayFrontendPortPropertiesFormat{
							Port: helpers.PointerToInt32(80),
						},
					},
				},
				BackendAddressPools: &[]network.ApplicationGatewayBackendAddressPool{
					{
						Name: helpers.PointerToString("pool"),
						ApplicationGatewayBackendAddressPoolPropertiesFormat: &network.ApplicationGatewayBackendAddressPoolPropertiesFormat{
							BackendAddresses: &[]network.ApplicationGatewayBackendAddress{},
						},
					},
				},
				HTTPListeners: &[]network.ApplicationGatewayHTTPListener{
					{
						Name: helpers.PointerToString("httpListener"),
						ApplicationGatewayHTTPListenerPropertiesFormat: &network.ApplicationGatewayHTTPListenerPropertiesFormat{
							Protocol: network.HTTP,
							FrontendPort: &network.SubResource{
								ID: helpers.PointerToString("[concat(variables('appGwId'), '/frontendPorts/httpPort')]"),
							},
							FrontendIPConfiguration: &network.SubResource{
								ID: helpers.PointerToString("[concat(variables('appGwId'), '/frontendIPConfigurations/frontendIP')]"),
							},
						},
					},
				},
				BackendHTTPSettingsCollection: &[]network.ApplicationGatewayBackendHTTPSettings{
					{
						Name: helpers.PointerToString("setting"),
						ApplicationGatewayBackendHTTPSettingsPropertiesFormat: &network.ApplicationGatewayBackendHTTPSettingsPropertiesFormat{
							Port:     helpers.PointerToInt32(80),
							Protocol: network.HTTP,
						},
					},
				},
				RequestRoutingRules: &[]network.ApplicationGatewayRequestRoutingRule{
					{
						Name: helpers.PointerToString("rule"),
						ApplicationGatewayRequestRoutingRulePropertiesFormat: &network.ApplicationGatewayRequestRoutingRulePropertiesFormat{
							HTTPListener: &network.SubResource{
								ID: helpers.PointerToString("[concat(variables('appGwId'), '/httpListeners/httpListener')]"),
							},
							BackendAddressPool: &network.SubResource{
								ID: helpers.PointerToString("[concat(variables('appGwId'), '/backendAddressPools/pool')]"),
							},
							BackendHTTPSettings: &network.SubResource{
								ID: helpers.PointerToString("[concat(variables('appGwId'), '/backendHttpSettingsCollection/setting')]"),
							},
						},
					},
				},
			},
			Type: helpers.PointerToString("Microsoft.Network/applicationGateways"),
		},
	}

	diff := cmp.Diff(actual, expected)

	if diff != "" {
		t.Errorf("unexpected error while comparing application gateways: %s", diff)
	}

}

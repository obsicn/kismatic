resource "azurerm_resource_group" "kismatic" {
  name     = "${var.cluster_name}"
  location = "${var.location}"
}

resource "azurerm_virtual_network" "kismatic" {
  name                = "${var.cluster_name}"
  address_space       = ["10.0.0.0/16"]
  location            = "${azurerm_resource_group.kismatic.location}"
  resource_group_name = "${azurerm_resource_group.ket.name}"
}

# Required for cloud provider integration
resource "azurerm_route_table" "ket" {
  name                = "kubernetes"
  location            = "${azurerm_resource_group.kismatic.location}"
  resource_group_name = "${azurerm_resource_group.kismatic.name}"
}

resource "azurerm_subnet" "kubenodes" {
  name                      = "kubenodes"
  resource_group_name       = "${azurerm_resource_group.ket.name}"
  virtual_network_name      = "${azurerm_virtual_network.kubenet.name}"
  address_prefix            = "${var.nodes_subnet}"
  route_table_id            = "${azurerm_route_table.ket.id}"
}
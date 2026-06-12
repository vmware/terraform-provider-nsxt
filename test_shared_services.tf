# HCL script to test the fix for issue #2104: include_shared_services and built_in_only in nsxt_policy_services data source

# 1. Create a Project
resource "nsxt_policy_project" "project1" {
  display_name = "project1"
  description  = "Terraform provisioned Project"
}

# 2. Create a Service in the default space (not in the project)
resource "nsxt_policy_service" "service_default" {
  display_name = "service_default_for_share"
  description  = "Service in default space to be shared"
  l4_port_set_entry {
    display_name      = "default_l4_entry"
    description       = "L4 Port Set Entry"
    protocol          = "TCP"
    destination_ports = ["80"]
  }
}

# 3. Share the Service with project1
resource "nsxt_policy_share" "share_to_project1" {
  display_name = "share_to_project1"
  description  = "Share container for project1"
  shared_with  = [nsxt_policy_project.project1.path]
}

resource "nsxt_policy_shared_resource" "shared_service" {
  display_name = "shared_service_resource"
  share_path   = nsxt_policy_share.share_to_project1.path
  resource_object {
    resource_path    = nsxt_policy_service.service_default.path
    include_children = false
  }
}

# 4. Read services inside project1 including shared services
data "nsxt_policy_services" "project_services_with_shared" {
  context {
    project_id = nsxt_policy_project.project1.id
  }
  include_shared_services = true
  depends_on              = [nsxt_policy_shared_resource.shared_service]
}

# 5. Read services inside project1 excluding shared services (default behavior)
data "nsxt_policy_services" "project_services_without_shared" {
  context {
    project_id = nsxt_policy_project.project1.id
  }
  include_shared_services = false
  depends_on              = [nsxt_policy_shared_resource.shared_service]
}

# 6. Read only built-in services from the default project
data "nsxt_policy_services" "default_built_in" {
  built_in_only = true
}

# Outputs to verify the results
output "project_services_with_shared_items" {
  value = data.nsxt_policy_services.project_services_with_shared.items
}

output "project_services_without_shared_items" {
  value = data.nsxt_policy_services.project_services_without_shared.items
}

output "default_built_in_items" {
  value = data.nsxt_policy_services.default_built_in.items
}

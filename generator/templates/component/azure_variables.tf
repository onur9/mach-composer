{% set azure = definition.Azure %}
### azure related
azure_short_name              = "{{ azure.ShortName }}"
azure_name_prefix             = local.name_prefix
azure_subscription_id         = local.subscription_id
azure_tenant_id               = local.tenant_id
azure_region                  = local.region
azure_service_object_ids      = local.service_object_ids
azure_resource_group          = {
  name     = local.resource_group_name
  location = local.resource_group_location
}
{% if azure.ServicePlan -%}
azure_app_service_plan        = {
  id                  = azurerm_app_service_plan.{{ azure.ServicePlan|service_plan_resource_name }}.id
  name                = azurerm_app_service_plan.{{ azure.ServicePlan|service_plan_resource_name }}.name
  resource_group_name = azurerm_app_service_plan.{{ azure.ServicePlan|service_plan_resource_name }}.resource_group_name
}
{% endif %}
{% if site.Azure.AlertGroup %}
azure_monitor_action_group_id = azurerm_monitor_action_group.alert_action_group.id
{% endif %}
{% for component_endpoint, site_endpoint in component.Endpoints -%}
azure_endpoint_{{ component_endpoint|slugify }} = {
  url = local.endpoint_url_{{ site_endpoint|slugify }}
  frontdoor_id = azurerm_frontdoor.app-service.header_frontdoor_id
}
{% endfor %}

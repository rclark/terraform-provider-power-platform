---
page_title: "Provider: Power Platform"
description: |-
  The Power Platform Provider allows managing environments and other resources within [Power Platform](https://powerplatform.microsoft.com/)
---

# Power Platform Provider

The Power Platform provider allows managing environments and other resources within [Power Platform](https://powerplatform.microsoft.com/).

!> This code is experimental and provided solely for evaluation purposes. It is **NOT** intended for production use and may contain bugs, incomplete features, or other issues. Use at your own risk, as it may undergo significant changes without notice, and no guarantees or support are provided. By using this code, you acknowledge and agree to these conditions. Consult the documentation or contact the maintainer if you have questions or concerns.

## Requirements

This provider requires **Terraform >= 0.12**.  For more information on provider installation and constraining provider versions, see the [Provider Requirements documentation](https://developer.hashicorp.com/terraform/language/providers/requirements).

## Installation

To use this provider, add the following to your Terraform configuration:

```terraform
terraform {
  required_providers {
    powerplatform = {
      source  = "microsoft/power-platform"
      version = "~> 1.0" # Replace with the latest version
    }
  }
}
```

See the official Terraform documentation for more information about [requiring providers](https://developer.hashicorp.com/terraform/language/providers/requirements).

## Authenticating to Power Platform

Terraform supports a number of different methods for authenticating to Power Platform.

* [Authenticating to Power Platform using the Azure CLI](#authenticating-to-power-platform-using-the-azure-cli)
* [Authenticating to Power Platform using a Service Principal with OIDC](#authenticating-to-power-platform-using-a-service-principal-with-oidc)
* [Authenticating to Power Platform using a Service Principal and a Client Secret](#authenticating-to-power-platform-using-a-service-principal-and-a-client-secret)

We recommend using either a Service Principal when running Terraform non-interactively (such as when running Terraform in a CI server) - and authenticating using the Azure CLI when running Terraform locally.

Important Notes about Authenticating using the Azure CLI:

* Terraform only supports authenticating using the az CLI (and this must be available on your PATH) - authenticating using the older azure CLI or PowerShell Cmdlets are not supported.
* Authenticating via the Azure CLI is only supported when using a User Account. If you're using a Service Principal (for example via az login --service-principal) you should instead authenticate via the Service Principal directly (either using a Client Secret or OIDC).

### Authenticating to Power Platform using the Azure CLI

The Power Platform provider can use the [Azure CLI](https://learn.microsoft.com/en-us/cli/azure/) to authenticate to Power Platform services. If you have the Azure CLI installed, you can use it to log in to your Microsoft Entra Id account and the Power Platform provider will use the credentials from the Azure CLI.

#### Prerequisites

1. [Install the Azure CLI](https://docs.microsoft.com/en-us/cli/azure/install-azure-cli)
1. [Create an app registration for the Power Platform Terraform Provider](guides/app_registration.md)
1. Login using the scope as the "expose API" you configured when creating the app registration

```bash
az login --allow-no-subscriptions --scope api://powerplatform_provider_terraform/.default
```

Configure the provider to use the Azure CLI with the following code:

```terraform
provider "powerplatform" {
  use_cli = true
}
```

### Authenticating to Power Platform using a Service Principal with OIDC

The Power Platform provider can use a Service Principal with OpenID Connect (OIDC) to authenticate to Power Platform services. By using [Microsoft Entra's workload identity federation](https://learn.microsoft.com/en-us/entra/workload-id/workload-identity-federation) your CI/CD pipelines in GitHub or Azure DevOps can access Power Platform resources without needing to manage secrets.

1. [Create an app registration for the Power Platform Terraform Provider](guides/app_registration.md)
1. [Register your app registration with Power Platform](https://learn.microsoft.com/en-us/power-platform/admin/powerplatform-api-create-service-principal#registering-an-admin-management-application)
1. [Create a trust relationship between your CI/CD pipeline and the app registration](https://learn.microsoft.com/en-us/entra/workload-id/workload-identity-federation-create-trust?pivots=identity-wif-apps-methods-azp)
1. Configure the provider to use OIDC with the following code:

```terraform
provider "powerplatform" {
  use_oidc = true
}
```

Additional Resources about OIDC:
* [OpenID Connect authentication with Microsoft Entra ID](https://learn.microsoft.com/en-us/entra/architecture/auth-oidc)
* [Configuring OpenID Connect for GitHub and Microsoft Entra ID](https://docs.github.com/en/actions/deployment/security-hardening-your-deployments/configuring-openid-connect-in-azure)

### Authenticating to Power Platform using a Service Principal and a Client Secret

The Power Platform provider can use a Service Principal with Client Secret to authenticate to Power Platform services.

1. [Create an app registration for the Power Platform Terraform Provider](guides/app_registration.md)
1. [Register your app registration with Power Platform](https://learn.microsoft.com/en-us/power-platform/admin/powerplatform-api-create-service-principal#registering-an-admin-management-application)
1. Configure the provider to use a Service Principal with a Client Secret with either environment variables or using Terraform variables

#### Using Environment Variables

We recomend using Environment Variables to pass the credentials to the provider.

| Name | Description | Default Value |
|------|-------------|---------------|
| `POWER_PLATFORM_CLIENT_ID` | The service principal client id | |
| `POWER_PLATFORM_CLIENT_SECRET` | The service principal secret | |
| `POWER_PLATFORM_TENANT_ID` | The guid of the tenant | |

-> Variables passed into the provider will override the environment variables.

#### Using Terraform Variables

Alternatively, you can configure the provider using variables in your Terraform configuration which can be passed in via [command line parameters](https://developer.hashicorp.com/terraform/language/values/variables#variables-on-the-command-line), [a `*.tfvars` file](https://developer.hashicorp.com/terraform/language/values/variables#variable-definitions-tfvars-files), or [environment variables](https://developer.hashicorp.com/terraform/language/values/variables#environment-variables).  If you choose to use variables, please be sure to [protect sensitive input variables](https://developer.hashicorp.com/terraform/tutorials/configuration-language/sensitive-variables) so that you do not expose your credentials in your Terraform configuration. 

```terraform
provider "powerplatform" {
  # Use a service principal to authenticate with the Power Platform service
  client_id     = var.client_id
  client_secret = var.client_secret
  tenant_id     = var.tenant_id
}
```

## Additional configuration

In addition to the authentication options, the following options are also supported in the provider block:

| Name | Description | Default Value |
|------|-------------|---------------|
| `telemetry_optout` | Opting out of telemetry will remove the hostheader from the requests made to the Power Platform service.  There is no other telemetry data collected by the provider.  This may affect the ability to identify and troubleshoot issues with the provider. | `false` |

## Resources and Data Sources

Use the navigation to the left to read about the available resources and data sources.

!> By calling `terraform destroy` all the resources, that you've created, will be deleted permamently deleted. Please be careful with this command when working with production environments. You can use [prevent-destory](https://developer.hashicorp.com/terraform/language/meta-arguments/lifecycle#prevent_destroy) lifecycle argument in your resources to prevent accidental deletion.  

## Examples

More detailed examples can be found in the [Power Platform Terraform Quickstarts](https://github.com/microsoft/power-platform-terraform-quickstarts) repo.  This repo contains a number of examples for using the Power Platform provider to manage environments and other resources within Power Platform along with Azure and Entra.

## Contributing

Contributions to this provider are always welcome! Please see the [Contribution Guidelines](https://github.com/microsoft/terraform-provider-power-platform/)

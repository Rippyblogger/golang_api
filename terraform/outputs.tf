output "oidc_provider_arn" {
  value = data.terraform_remote_state.infra.outputs.oidc_provider_arn
}

output "cluster_name" {
  value =  data.terraform_remote_state.infra.outputs.cluster_name
}

output "cluster_endpoint" {
  value = data.terraform_remote_state.infra.outputs.cluster_endpoint
}

output "cluster_certificate_authority_data" {
  value = data.terraform_remote_state.infra.outputs.cluster_ca_certificate
}

output "oidc_provider_url" {
  value = data.terraform_remote_state.infra.outputs.oidc_provider_url
}
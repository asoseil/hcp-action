# HCP Vault Secrets Reader

This GitHub Action authenticates with HashiCorp Vault Secrets using client credentials and retrieves secrets from multiple applications based on a provided input object.

## Features
- **Authentication**: Logs into HCP Vault using the provided client ID and client secret.
- **Secrets**: Fetches secrets from specified apps and keys based on a JSON object.
- **Export**: Exports the retrieved secrets as outputs or environment variables.

## Usage

```yaml
uses: magicmotorsport/hcp-vault-secrets-action@v1
with:
  client-id: 'client-id-key'
  client-secret: 'client-secret-key'
  stream: 'env'
  secrets: 'secret'
```

## To do

- [ ] Option to stream on out, env or both
- [ ] Build bin on push and commit
- [ ] Use bin as entrypoint

## Refs

- Docker container action: https://docs.github.com/en/actions/sharing-automations/creating-actions/creating-a-docker-container-action
- HCP auth login: https://developer.hashicorp.com/hcp/docs/cli/commands/auth/login
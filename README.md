# Terraform Provider Manidae

Manidae is a Terraform provider implemented with the [Terraform Plugin Framework](https://github.com/hashicorp/terraform-plugin-framework).

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.24

## Build

```shell
make install
```

## Development

```shell
make fmt
make lint
make test
```

To regenerate docs (and format Terraform examples):

```shell
make generate
```

Acceptance tests (creates real resources):

```shell
make testacc
```

## Data Source: `manidae_parameter`

`data "manidae_parameter"` resolves a value from an environment variable derived from `name`, falling back to `default`.

Example:

```hcl
data "manidae_parameter" "root_volume_size_gb" {
  name    = "root_volume_size_gb"
  type    = "number"
  default = 30

  validation {
    min = 20
  }
}
```

Enum-style options (for `type = "string"`):

```hcl
data "manidae_parameter" "instance_type" {
  name    = "instance_type"
  default = "SA2.MEDIUM8"

  option { value = "SA2.MEDIUM2" }
  option { value = "SA2.MEDIUM4" }
  option { value = "SA2.MEDIUM8" }
}
```

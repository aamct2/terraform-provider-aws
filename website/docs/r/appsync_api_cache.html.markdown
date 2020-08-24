---
subcategory: "AppSync"
layout: "aws"
page_title: "AWS: aws_appsync_api_cache"
description: |-
  Provides an AppSync API Cache.
---

# Resource: aws_appsync_api_cache

Provides an AppSync API Cache.

## Example Usage

```hcl
resource "aws_appsync_graphql_api" "example" {
  authentication_type = "API_KEY"
  name                = "example"
}

resource "aws_appsync_api_cache" "example" {
  api_caching_behavior = "FULL_REQUEST_CACHING"
  api_id               = aws_appsync_graphql_api.example.id
  ttl                  = 60
  type                 = "SMALL"
}
```

## Argument Reference

The following arguments are supported:

- `api_caching_behavior` - (Required) Caching behavior.
- `api_id` - (Required) The ID of the associated AppSync API
- `at_rest_encryption_enabled` - (Optional) At rest encryption flag for cache. This setting cannot be updated after creation.
- `transit_encryption_enabled` - (Optional) Transit encryption flag when connecting to cache. This setting cannot be updated after creation.
- `ttl` - (Required) TTL in seconds for cache entries. Valid values are between 1 and 3600 seconds.
- `type` - (Required) The cache instance type.

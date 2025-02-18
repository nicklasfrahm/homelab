# API

This document describes the API of my configuration management tool.

## Endpoints

The API is inspired by Kubernetes and is available at [`https://cloud.nicklasfrahm.dev`](https://cloud.nicklasfrahm.dev/).

### `GET /v1beta1/machines`

Returns a list of all machines.

```json
{
  "kind": "MachineList",
  "apiVersion": "cloud.nicklasfrahm.dev/v1beta1",
  "metadata": {},
  "items": [
    {
      "kind": "Machine",
      "apiVersion": "cloud.nicklasfrahm.dev/v1beta1",
      "metadata": {
        "name": "ant",
        "creationTimestamp": null
      },
      "spec": {
        "hardware": {
          "vendor": "FriendlyElec",
          "model": "NanoPiR5S"
        },
        "interfaces": [
          {
            "mac": "32:de:fa:97:71:4f"
          }
        ]
      },
      "status": {}
    }
  ]
}
```

### `GET /v1beta1/machines/{name}`

Returns the configuration of a single machine.

```json
{
  "kind": "Machine",
  "apiVersion": "cloud.nicklasfrahm.dev/v1beta1",
  "metadata": {
    "name": "ant",
    "creationTimestamp": null
  },
  "spec": {
    "hardware": {
      "vendor": "FriendlyElec",
      "model": "NanoPiR5S"
    },
    "interfaces": [
      {
        "mac": "32:de:fa:97:71:4f"
      }
    ]
  },
  "status": {}
}
```

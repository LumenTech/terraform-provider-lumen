| Page_Title      | Description                                 |
|-----------------|---------------------------------------------|
| Resource_Order_Delete_Bare_Metal_Instance  | Details on bare metal instance delete |

## Introduction
This document provides details on resource order to delete Lumen bare metal instance(s). The API details are provided in Ref [[1]](#1).

## Working
Deleting a bare metal instance with terraform doesn't need any terraform utility files. Initializing terraform in an existing workspace (with tfstate files) enables terraform to understand what all resources can be deleted. Here's the CLIs:
```shell
$ terraform init
$ terraform show / state
$ terraform destroy --auto-approve
```

## References
<a id="1">[1]</a> Lumen Developer API doc: https://developer.lumen.com/apis/edge-bare-metal#api-reference_edge-bare-metal-api_instances_instances-id_delete

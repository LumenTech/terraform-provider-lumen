| Page_Title      | Description                                 |
|-----------------|---------------------------------------------|
| Resource_Order_Delete_Network_Instance  | Details on deleting network instance |

## Introduction
This document provides details on resource order to delete Lumen network instance(s). The API details are provided in Ref [[1]](#1).

## Working
To delete a network instance with terraform doesn't need any terraform utility files. Initializing terraform in an existing workspace (with tfstate files) enables terraform to understand what all resources can be deleted. Here's the CLIs:
```shell
$ terraform init
$ terraform show / state
$ terraform destroy --auto-approve
```

## References
<a id="1">[1]</a> http://apidocs.edge.lumen.com/#delete-an-instance

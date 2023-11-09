## Lumen User Registration
As a first step to use the Lumen terraform provider, users need to register themselves in Lumen Developer Center. The url for Lumen Developer Center is mentioned in Ref [[1]](#1). Once the registration is successful, users can then proceed to provision Lumen resources using APIs. Users will need their username, password and account number.

## Deprecated - Lumen API Key
NOTE: Once you have migrated to the new version of Lumen Edge Bare Metal API, the below access-token and refresh-token will no longer be required.

To generate API keys users need to pass on their username and password in terraform.tfvars obtained from Lumen Developer Center, for generating oauth token for Lumen API endpoint authentication. This can be done with an API request as mentioned in Ref [[2]](#2). Also, users need to grab the access-token and refresh-token under `morph-api` from Lumen Edge Orchestrator as mentioned in Ref [[2]](#2) Steps-[3-4] for executing API calls to Lumen Edge endpoint.

## References
<a id="1">[1]</a> Lumen Developer Center: https://developer.lumen.com/apis/edge-bare-metal#getting-started

<a id="2">[2]</a> Lumen API doc: https://developer.lumen.com/apis/edge-bare-metal#authentication

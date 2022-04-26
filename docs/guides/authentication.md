## Lumen-Technologies Authentication
Users need to authenticate with Lumen API marketplace to be able to build bare-metal instances as mentioned in Ref [[1]](#1).

## Lumen API key
Users need to pass on their username and password in terraform.tfvars obtained from Lumen API Marketplace Ref [[2]](#2), for generating and refreshing access token for authentication with Lumen API endpoint Ref [[3]](#3). Also, users should provide `morph-api` token for executing API calls to Lumen Technologies endpoint.

## References
<a id="1">[1]</a> Lumen API Marketplace doc: https://apimarketplace.lumen.com/api/edge-compute?tab=getting-started

<a id="2">[2]</a> Lumen API Marketplace (for registering new users) "Getting started": https://apimarketplace.lumen.com/api/edge-compute?tab=getting-started

<a id="3">[3]</a> Lumen API doc: https://apimarketplace.lumen.com/api/edge-compute?tab=authentication-authorization

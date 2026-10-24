# Rhyolite

Rhyolite is a backend for server deployment.

Each folder with a deployable item(s) is called a "Rock".

Each rock has a `.rhyolite.yaml` which defines which commands can be executed, and their tokens.

It takes post requests to /v1/{path}/{token}, where path is the folder where the config for the Rock is.  
The token specifies the command to run, as well as authorizing it, as defined in the `.rhyolite.yaml`

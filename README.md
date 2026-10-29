# Rhyolite

Rhyolite is a backend for server deployment.

Each folder with a deployable item(s) is called a "Rock".

Each rock has a `.rhyolite.yaml` which defines which commands can be executed, and their tokens.

It takes post requests to /v1/{path}/{command}/{token}, where path is the folder where the config for the Rock is.  
The token is unique per command and is what authorized it.

The main configuration file is located at `/etc/rhyolite/config.yaml`

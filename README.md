# Rhyolite

Rhyolite is a backend for server deployment.

It takes post requests to /v1/{path}, and a header "X-Deploy-Token".

If the token matches that in $RHYOLITE_DIR_ROOT/rhyolite.yaml, then run the command in the yaml file.

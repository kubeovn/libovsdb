package serverdb

// server_model is a database model for the special _Server database that all
// ovsdb instances export. It reports back status of the server process itself.

//go:generate go tool github.com/ovn-kubernetes/libovsdb/cmd/modelgen --extended -p serverdb -o . _server.ovsschema

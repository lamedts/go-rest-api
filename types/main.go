package types

const (
	ServerAPI     string = "API server"
	ServerGraphql string = "Graphql server"
)

type Server interface {
	Serve() bool
}

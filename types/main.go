package types

const (
	ServerAPI     string = "API server"
	ServerGraphql string = "Graphql server"
)

type Server interface {
	Serve() bool
	// Stop() bool
	// GetPurchaseRequestChannel() chan PurchaseRequest
	// GetPaymentResponseChannel() chan []byte
	// GetPurchaseMetricChannel() chan<- MetricRecord
	//Send(data []byte) bool
}

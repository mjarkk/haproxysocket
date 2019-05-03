package haproxysocket

// HaproxyInstace this will contain all haproxy unix socket settings
type HaproxyInstace struct {
	Network string
	Address string
}

// New creates a new haproxyInstace instace
func New(network, address string) *HaproxyInstace {
	toReturn := HaproxyInstace{
		Network: network,
		Address: address,
	}
	return &toReturn
}

package osipsclasses

// NetDescription function to get the description for NET type
func NetDescription(record string) Metric {
	switch record {
	case "waiting_udp":
		return (Metric{"Number of bytes waiting to be consumed on UDP interfaces that OpenSIPS is listening on.", "gauge"})
	case "waiting_tcp":
		return (Metric{"Number of bytes waiting to be consumed on TCP interface that OpenSIPS is listening on.", "gauge"})
	case "waiting_tls":
		return (Metric{"Number of bytes waiting to be consumed on an interface that OpenSIPS is listening on.", "gauge"})
	default:
		return (Metric{"Unrecognized metric.", "gauge"})
	}
}


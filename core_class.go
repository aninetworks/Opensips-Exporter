package osipsclasses

// CoreDescription function to get the description for CORE type
func CoreDescription(record string) Metric {
	switch record {
	case "rcv_requests":
		return (Metric{"Total number of received requests by OpenSIPS.", "counter"})
	case "rcv_replies":
		return (Metric{"Total number of received replies by OpenSIPS.", "counter"})
	case "fwd_requests":
		return (Metric{"Number of requests by OpenSIPS.", "counter"})
	case "fwd_replies":
		return (Metric{"Number of received replies by OpenSIPS.", "counter"})
	case "drop_requests":
		return (Metric{"Number of requests by OpenSIPS.", "counter"})
	case "drop_replies":
		return (Metric{"Number of received replies by OpenSIPS.", "counter"})
	case "err_requests":
		return (Metric{"Number of requests by OpenSIPS.", "counter"})
	case "err_replies":
		return (Metric{"Number of received replies by OpenSIPS.", "counter"})
	case "bad_URIs_rcvd":
		return (Metric{"Number of URIs that OpenSIPS failed to parse.", "counter"})
	case "unsupported_methods":
		return (Metric{"Number of non-standard methods encountered by OpenSIPS while parsing SIP methods.", "counter"})
	case "bad_msg_hdr":
		return (Metric{"Number of SIP headers that OpenSIPS failed to parse.", "counter"})
	case "timestamp":
		return (Metric{"Number of seconds elapsed from OpenSIPS starting.", "counter"})
	default:
		return (Metric{"Unrecognized metric.", "gauge"})
	}

}


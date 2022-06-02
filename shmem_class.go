package osipsclasses

// ShmemDescription function to get the description for SHMEM type
func ShmemDescription(record string) Metric {
	switch record {
	case "total_size":
		return (Metric{"Total size of shared memory available to OpenSIPS processes.", "gauge"})
	case "used_size":
		return (Metric{"Amount of shared memory requested and used by OpenSIPS processes.", "gauge"})
	case "real_used_size":
		return (Metric{"Amount of shared memory requested by OpenSIPS processes + malloc overhead.", "gauge"})
	case "max_used_size":
		return (Metric{"Maximum amount of shared memory ever used by OpenSIPS processes.", "gauge"})
	case "free_size":
		return (Metric{"Free memory available. Computed as total_size - real_used_size.", "gauge"})
	case "fragments":
		return (Metric{"Total number of fragments in the shared memory.", "gauge"})
	default:
		return (Metric{"Unrecognized metric.", "gauge"})
	}

}


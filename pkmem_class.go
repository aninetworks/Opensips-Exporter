package osipsclasses

import "strings"

// PkmemDescription function to get the description for PKMEM type
func PkmemDescription(record string) Metric {
	desc := "Unrecognized metric."
	val := "gauge"
	if strings.Contains(record, "total_size") {
		desc = "Total size of private memory available to the OpenSIPS process."
	} else if strings.Contains(record, "used_size") {
		desc = "Amount of private memory requested and used by the OpenSIPS process."
	} else if strings.Contains(record, "real_used_size") {
		desc = "Amount of private memory requested by the OpenSIPS process, including allocator-specific metadata."
	} else if strings.Contains(record, "max_used_size") {
		desc = "The maximum amount of private memory ever used by the OpenSIPS process."
	} else if strings.Contains(record, "free_size") {
		desc = "Free private memory available for the OpenSIPS process. Computed as total_size - real_used_size."
	} else if strings.Contains(record, "fragments") {
		desc = "Currently available number of free fragments in the private memory for OpenSIPS process."
	}
	return (Metric{desc, val})
}


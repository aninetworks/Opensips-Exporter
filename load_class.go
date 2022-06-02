package osipsclasses

import "strings"

// LoadDescription function to get the description for LOAD type
func LoadDescription(record string) Metric {
	switch record {
	case "load":
		return (Metric{"The realtime load of entire OpenSIPS - this counts all the core processes of OpenSIPS; the additional processes requested by modules are not counted in this load.", "gauge"})
	case "load1m":
		return (Metric{"The last minute average load of core OpenSIPS (covering only core or SIP processes)", "gauge"})
	case "load10m":
		return (Metric{"The last 10 minute average load of core OpenSIPS (covering only core or SIP processes)", "gauge"})
	case "load-all":
		return (Metric{"The realtime load of entire OpenSIPS, counting both core and module processes.", "gauge"})
	case "load1m-all":
		return (Metric{"The last minute average load of entire OpenSIPS (covering all processes).", "gauge"})
	case "load10m-all":
		return (Metric{"The last 10 minute average load of entire OpenSIPS (covering all processes).", "gauge"})
	default:
		desc := "Unrecognized metric."
		if strings.Contains(record, "load-proc-id") || strings.Contains(record, "load1m-proc-id") || strings.Contains(record, "load10m-proc-id") {
			desc = "The realtime load of the process ID."
		}
		return (Metric{desc, "gauge"})
	}
}


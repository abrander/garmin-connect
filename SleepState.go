package connect

// SleepState is used to describe the state of sleep with a device capable
// of measuring sleep health.
type SleepState int

// Known sleep states in Garmin Connect.
const (
	SleepStateUnknown SleepState = -1
	SleepStateDeep    SleepState = 0
	SleepStateLight   SleepState = 1
	SleepStateREM     SleepState = 2
	SleepStateAwake   SleepState = 3
)

// UnmarshalJSON implements json.Unmarshaler.
func (s *SleepState) UnmarshalJSON(value []byte) error {
	// Garmin abuses floats to transfers enums. We ignore the value, and
	// simply compares them as strings.
	switch string(value) {
	case "0.0":
		*s = SleepStateDeep
	case "1.0":
		*s = SleepStateLight
	case "2.0":
		*s = SleepStateREM
	case "3.0":
		*s = SleepStateAwake
	default:
		*s = SleepStateUnknown
	}

	return nil
}

// Sleep implements fmt.Stringer.
func (s SleepState) String() string {
	m := map[SleepState]string{
		SleepStateUnknown: "Unknown",
		SleepStateDeep:    "Deep",
		SleepStateLight:   "Light",
		SleepStateREM:     "REM",
		SleepStateAwake:   "Awake",
	}

	str, found := m[s]
	if !found {
		str = m[SleepStateUnknown]
	}

	return str
}

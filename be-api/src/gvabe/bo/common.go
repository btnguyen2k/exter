package bo

import (
	"github.com/btnguyen2k/henge"
)

const (
	// SerKeyAttrs is a key used by BO's ToMap(), MarshalJSON() and UnmarshalJSON() functions to store BO's custom attributes.
	SerKeyAttrs = "_attrs"

	// SerKeyFields is a key used by BO's ToMap, MarshalJSON() and UnmarshalJSON() functions to store BO's top-level custom fields.
	SerKeyFields = "_fields"
)

var (
	// UboTimeLayout controls how henge.UniversalBo will format create/update timestamp.
	UboTimeLayout = henge.DefaultTimeLayout

	// UboTimestampRounding controls how henge.UniversalBo will round create/update timestamp.
	UboTimestampRounding = henge.DefaultTimestampRoundingSetting
)
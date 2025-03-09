package entities

// Zone defines if the land use zoning for a tile
type Zone string

const (
	ResidentialZone Zone = "zone-residential"
	RetailZone      Zone = "zone-retail"
	NoZone          Zone = ""
)

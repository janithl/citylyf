package entities

// LandUse defines what a tile is zoned for
type LandUse string

const (
	ResidentialUse LandUse = "residential"
	RetailUse      LandUse = "retail"
	ReserveUse     LandUse = "reserve"
	TransportUse   LandUse = "transport"
	NoUse          LandUse = ""
)

// LandStatus defines at which stage of development a tile is at
type LandStatus string

const (
	DevelopedStatus   LandStatus = "developed"
	DevelopingStatus  LandStatus = "developing"
	UndevelopedStatus LandStatus = ""
)

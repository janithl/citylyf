package entities

// LandSlope defines the slope of the land
type LandSlope string

const (
	Flat            LandSlope = "flat"
	Top             LandSlope = "top"
	Bottom          LandSlope = "bottom"
	Left            LandSlope = "left"
	Right           LandSlope = "right"
	TopLeft         LandSlope = "top-left"
	BottomLeft      LandSlope = "bottom-left"
	TopRight        LandSlope = "top-right"
	BottomRight     LandSlope = "bottom-right"
	TopLeftRight    LandSlope = "top-left-right"
	BottomLeftRight LandSlope = "bottom-left-right"
	TopBottomRight  LandSlope = "top-bottom-right"
	TopBottomLeft   LandSlope = "top-bottom-left"
)

// LandUse defines what a tile is zoned for
type LandUse string

const (
	ResidentialUse LandUse = "residential"
	RetailUse      LandUse = "retail"
	ReserveUse     LandUse = "reserve"
	TransportUse   LandUse = "transport"
	AgricultureUse LandUse = "agriculture"
	NoUse          LandUse = ""
)

// LandStatus defines at which stage of development a tile is at
type LandStatus string

const (
	DevelopedStatus   LandStatus = "developed"
	DevelopingStatus  LandStatus = "developing"
	UndevelopedStatus LandStatus = ""
)

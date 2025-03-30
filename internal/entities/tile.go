package entities

const TileSize = 15.0 // Tile size in metres

type Tile struct {
	Elevation    int
	LandSlope    LandSlope
	Intersection IntersectionType
	LandUse      LandUse
	LandStatus   LandStatus
}

func (t *Tile) IsBuildable() bool { // buildable on sealevel if flat land
	return (t.Elevation > Sim.Geography.SeaLevel ||
		(t.Elevation == Sim.Geography.SeaLevel && t.LandSlope == Flat)) &&
		t.Elevation < Sim.Geography.HillLevel && t.LandUse != ReserveUse
}

func (t *Tile) IsBuilt() bool {
	return t.LandUse == ReserveUse || t.LandStatus != UndevelopedStatus
}

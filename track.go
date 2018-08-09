package main

type Track struct {
	ID            int      `json:"id"`
	Start         Location `json:"startLocation"`
	End           Location `json:"endLocation"`
	totalDistance Distance
	// TODO: Expand
}

type InternalTrack struct {
	ID       int
	Start    Location
	End      Location
	Progress []Progress
	// TODO: Expand
}

func (t InternalTrack) toTrack() Track {
	var distance, err = t.Distance()
	if err != nil {
		distance = Distance{0, "m"}
	}
	return Track{ID: t.ID, Start: t.Start, End: t.End, totalDistance: distance}
}

func (t InternalTrack) Distance() (Distance, error) {
	// TODO: Calculate Distance between longitude and latitude
	return Distance{}, nil
}

package venue

// Fields returns the set of fields that have been retrieved for this venue.
func (v *Venue) Fields() Fields {
	return v.fields
}

// ID is the unique identifier of the Venue.
func (v *Venue) ID() ID {
	if v == nil {
		return 0
	}
	if v.fields&FID == 0 {
		panic("Venue.ID called without having fetched FID")
	}
	return v.id
}

// Name is the name of the Venue.
func (v *Venue) Name() string {
	if v.fields&FName == 0 {
		panic("Venue.Name called without having fetched FName")
	}
	return v.name
}

// URL is the Google Maps URL for a map to the Venue.
func (v *Venue) URL() string {
	if v.fields&FURL == 0 {
		panic("Venue.URL called without having fetched FURL")
	}
	return v.url
}

// Flags is the set of flags for the Venue.
func (v *Venue) Flags() Flag {
	if v.fields&FFlags == 0 {
		panic("Venue.Flags called without having fetched FFlags")
	}
	return v.flags
}

package class

// Fields returns the set of fields that have been retrieved for this venue.
func (c *Class) Fields() Fields {
	return c.fields
}

// ID is the unique identifier of the Class.
func (c *Class) ID() ID {
	if c == nil {
		return 0
	}
	if c.fields&FID == 0 {
		panic("Class.ID called without having fetched FID")
	}
	return c.id
}

// Type is the type of class being taught (i.e., the curriculum).
func (c *Class) Type() Type {
	if c.fields&FType == 0 {
		panic("Class.Type called without having fetched FType")
	}
	return c.ctype
}

// Start is the date of the first session of the class, in YYYY-MM-DD format.
// It is intended to identify the class.
func (c *Class) Start() string {
	if c.fields&FStart == 0 {
		panic("Class.Start called without having fetched FStart")
	}
	return c.start
}

// EnDesc is English description of the class instance's date(s), time(s),
// location(s), and language if appropriate.  It is HTML text.
func (c *Class) EnDesc() string {
	if c.fields&FEnDesc == 0 {
		panic("Class.EnDesc called without having fetched FEnDesc")
	}
	return c.enDesc
}

// EsDesc is Spanish description of the class instance's date(s), time(s),
// location(s), and language if appropriate.  It is HTML text.
func (c *Class) EsDesc() string {
	if c.fields&FEsDesc == 0 {
		panic("Class.EsDesc called without having fetched FEsDesc")
	}
	return c.esDesc
}

// Limit is the limit on how many people can register for the class, with zero
// meaning no limit.
func (c *Class) Limit() uint {
	if c.fields&FLimit == 0 {
		panic("Class.Limit called without having fetched FLimit")
	}
	return c.limit
}

// Referrals is a slice, indexed by Referral, of the number of people who
// indicated they learned of the class through that referral method.
func (c *Class) Referrals() []uint {
	if c.fields&FLimit == 0 {
		panic("Class.Referrals called without having fetched FReferrals")
	}
	return c.referrals
}

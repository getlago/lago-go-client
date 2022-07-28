package lago

type ChargeModel string

const (
	Standard   ChargeModel = "standard"
	Graduated  ChargeModel = "graduated"
	Package    ChargeModel = "package"
	Percentage ChargeModel = "percentage"
)

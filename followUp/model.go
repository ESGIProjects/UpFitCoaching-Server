package followUp

import "server/user"

type Appraisal struct {
	Id					int64		`json:"id"`
	User				user.Info	`json:"user"`
	Date				string		`json:"date"`
	Goal				string		`json:"goal"`
	SessionsByWeek		*int		`json:"sessionsByWeek"`
	Contraindication	string		`json:"contraindication"`
	SportAntecedents	string		`json:"sportAntecedents"`
	HelpNeeded			*int		`json:"helpNeeded"`
	HasNutritionist		*int		`json:"hasNutritionist"`
	Comments			string		`json:"comments"`
}

type Measurements struct {
	Id					int64		`json:"id"`
	User				user.Info	`json:"user"`
	Date				string		`json:"date"`
	Weight				*int		`json:"weight"`
	Height				*int		`json:"height"`
	HipCircumference	*int		`json:"hipCircumference"`
	WaistCircumference	*int		`json:"waistCircumference"`
	ThighCircumference	*int		`json:"thighCircumference"`
	ArmCircumference	*int		`json:"armCircumference"`
}

type Test struct {
	Id							int64		`json:"id"`
	User						user.Info	`json:"user"`
	Date						string		`json:"date"`
	WarmUp						*float64	`json:"warmUp"`
	StartSpeed					*float64	`json:"startSpeed"`
	Increase					*float64	`json:"increase"`
	Frequency					*float64	`json:"frequency"`
	KneeFlexibility				*int		`json:"kneeFlexibility"`
	ShinFlexibility				*int		`json:"shinFlexibility"`
	HitFootFlexibility			*int		`json:"hitFootFlexibility"`
	ClosedFistGroundFlexibility	*int		`json:"closedFistGroundFlexibility"`
	HandFlatGroundFlexibility	*int		`json:"handFlatGroundFlexibility"`
}
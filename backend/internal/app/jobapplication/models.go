package jobapplication

import "gorm.io/gorm"


type JobApplication struct {
	gorm.Model
	ID uint `json:"id" gorm:"primaryKey"`
	Title string `json:"title"`
	Description *string `json:"description,omitempty"`
	Employment Employment `json:"employment" gorm:"embedded"`
	Location *string `json:"location,omitempty"`
}

type Employment struct {
	Type EmploymentType `json:"type"`
	Duration *string `json:"duration,omitempty"`
	WorkLocation WorkLocation `json:"workLocation" gorm:"column:work_location"`
	SalaryRange *SalaryRange `json:"salaryRange,omitempty" gorm:"embedded;embeddedPrefix:salary_range_"`
}

type SalaryRange struct {
	From int `json:"from"`
	To int `json:"to"`
	Currency string `json:"currency"`
}

// EmploymentType defines employment status
// @enum EmploymentType
// @description Type of employment: FullTime, PartTime, or Contract
type EmploymentType string

const (
	FullTime EmploymentType = "FullTime"
	PartTime EmploymentType = "PartTime"
	Contract EmploymentType = "Contract"
)

// WorkLocation defines working location/type
// @enum WorkLocation
// @description Type of work environment/location: Onsite, Hybrid, Remote
type WorkLocation string

const (
	Onsite WorkLocation = "Onsite"
	Hybrid WorkLocation = "Hybrid"
	Remote WorkLocation = "Remote"
)

package jobapplication

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type JobApplication struct {
	gorm.Model
	UserID          uint              `json:"-" gorm:"index;uniqueIndex:uniq_user_extid_src"`
	Company         string            `json:"company"`
	Title           string            `json:"title"`
	Description     *string           `json:"description,omitempty"`
	Status          ApplicationStatus `json:"status" gorm:"type:VARCHAR(24);index"`
	Source          ApplicationSource `json:"source" gorm:"type:VARCHAR(24);index;uniqueIndex:uniq_user_extid_src"`
	AppliedAt       *time.Time        `json:"appliedAt,omitempty" gorm:"index"`
	NextFollowUpAt  *time.Time        `json:"nextFollowUpAt,omitempty" gorm:"index"`
	LastContactAt   *time.Time        `json:"lastContactAt,omitempty" gorm:"index"`
	PostingURL      *string           `json:"postingUrl,omitempty"`
	CompanyURL      *string           `json:"companyUrl,omitempty"`
	ContactName     *string           `json:"contactName,omitempty"`
	ContactEmail    *string           `json:"contactEmail,omitempty"`
	ExternalJobID   *string           `json:"externalJobId,omitempty" gorm:"index;uniqueIndex:uniq_user_extid_src"`
	Employment      Employment        `json:"employment" gorm:"embedded;embeddedPrefix:employment_"`
	Location        *string           `json:"location,omitempty"`
	Tags            datatypes.JSON    `json:"tags,omitempty" gorm:"type:jsonb;default:'[]'"`
}

type Employment struct {
	Type          EmploymentType `json:"type"`
	Duration      *string        `json:"duration,omitempty"`
	WorkLocation  WorkLocation   `json:"workLocation" gorm:"column:work_location"`
	Seniority     *string        `json:"seniority,omitempty"`
	HoursPerWeek  *int           `json:"hoursPerWeek,omitempty"`
	SalaryRange   *SalaryRange   `json:"salaryRange,omitempty" gorm:"embedded;embeddedPrefix:salary_range_"`
}

type SalaryRange struct {
	From       int          `json:"from"`
	To         int          `json:"to"`
	Currency   string       `json:"currency"`
	Period     SalaryPeriod `json:"period"`
	Negotiable bool         `json:"negotiable"`
}

type EmploymentType string
const (
	FullTime EmploymentType = "FullTime"
	PartTime EmploymentType = "PartTime"
	Contract EmploymentType = "Contract"
)

type WorkLocation string
const (
	Onsite WorkLocation = "Onsite"
	Hybrid WorkLocation = "Hybrid"
	Remote WorkLocation = "Remote"
)

type ApplicationStatus string
const (
	StatusApplied    ApplicationStatus = "Applied"
	StatusScreening  ApplicationStatus = "Screening"
	StatusInterview  ApplicationStatus = "Interview"
	StatusOffer      ApplicationStatus = "Offer"
	StatusRejected   ApplicationStatus = "Rejected"
	StatusHired      ApplicationStatus = "Hired"
	StatusWithdrawn  ApplicationStatus = "Withdrawn"
)

type ApplicationSource string
const (
	SourceReferral    ApplicationSource = "Referral"
	SourceCompanySite ApplicationSource = "CompanySite"
	SourceJobBoard    ApplicationSource = "JobBoard"
	SourceLinkedIn    ApplicationSource = "LinkedIn"
	SourceIndeed      ApplicationSource = "Indeed"
	SourceAgency      ApplicationSource = "Agency"
	SourceOther       ApplicationSource = "Other"
)

type SalaryPeriod string
const (
	PerYear  SalaryPeriod = "Year"
	PerMonth SalaryPeriod = "Month"
	PerWeek  SalaryPeriod = "Week"
	PerDay   SalaryPeriod = "Day"
	PerHour  SalaryPeriod = "Hour"
)

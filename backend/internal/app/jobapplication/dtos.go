package jobapplication

import "time"

type BaseJobApplicationDto struct {
	Company        string        `json:"company"`
	Title          string        `json:"title"`
	Description    *string       `json:"description,omitempty"`
	Status         string        `json:"status"`
	Source         string        `json:"source"`
	AppliedAt      *time.Time    `json:"appliedAt,omitempty"`
	NextFollowUpAt *time.Time    `json:"nextFollowUpAt,omitempty"`
	LastContactAt  *time.Time    `json:"lastContactAt,omitempty"`
	PostingURL     *string       `json:"postingUrl,omitempty"`
	CompanyURL     *string       `json:"companyUrl,omitempty"`
	ContactName    *string       `json:"contactName,omitempty"`
	ContactEmail   *string       `json:"contactEmail,omitempty"`
	ExternalJobID  *string       `json:"externalJobId,omitempty"`
	Employment     EmploymentDto `json:"employment"`
	Location       *string       `json:"location,omitempty"`
	Tags           []string      `json:"tags,omitempty"`
}

type JobApplicationCreateDto struct {
	BaseJobApplicationDto
}

type JobApplicationPublicDto struct {
	ID        uint                   `json:"id"`
	Created   string                 `json:"created"`
	BaseJobApplicationDto
}

type JobApplicationPatchDto struct {
	Company        *string              `json:"company,omitempty"`
	Title          *string              `json:"title,omitempty"`
	Description    *string              `json:"description,omitempty"`
	Status         *string              `json:"status,omitempty"`
	Source         *string              `json:"source,omitempty"`
	AppliedAt      *time.Time           `json:"appliedAt,omitempty"`
	NextFollowUpAt *time.Time           `json:"nextFollowUpAt,omitempty"`
	LastContactAt  *time.Time           `json:"lastContactAt,omitempty"`
	PostingURL     *string              `json:"postingUrl,omitempty"`
	CompanyURL     *string              `json:"companyUrl,omitempty"`
	ContactName    *string              `json:"contactName,omitempty"`
	ContactEmail   *string              `json:"contactEmail,omitempty"`
	ExternalJobID  *string              `json:"externalJobId,omitempty"`
	Employment     *EmploymentPatchDto  `json:"employment,omitempty"`
	Location       *string              `json:"location,omitempty"`
	Tags           *[]string            `json:"tags,omitempty"`
}

type EmploymentDto struct {
	Type         string          `json:"type"`
	Duration     *string         `json:"duration,omitempty"`
	WorkLocation string          `json:"workLocation"`
	Seniority    *string         `json:"seniority,omitempty"`
	HoursPerWeek *int            `json:"hoursPerWeek,omitempty"`
	SalaryRange  *SalaryRangeDto `json:"salaryRange,omitempty"`
}

type EmploymentPatchDto struct {
	Type         *string               `json:"type,omitempty"`
	Duration     *string               `json:"duration,omitempty"`
	WorkLocation *string               `json:"workLocation,omitempty"`
	Seniority    *string               `json:"seniority,omitempty"`
	HoursPerWeek *int                  `json:"hoursPerWeek,omitempty"`
	SalaryRange  *SalaryRangePatchDto  `json:"salaryRange,omitempty"`
}

type SalaryRangeDto struct {
	From       int    `json:"from"`
	To         int    `json:"to"`
	Currency   string `json:"currency"`
	Period     string `json:"period"`
	Negotiable bool   `json:"negotiable"`
}

type SalaryRangePatchDto struct {
	From       *int    `json:"from,omitempty"`
	To         *int    `json:"to,omitempty"`
	Currency   *string `json:"currency,omitempty"`
	Period     *string `json:"period,omitempty"`
	Negotiable *bool   `json:"negotiable,omitempty"`
}

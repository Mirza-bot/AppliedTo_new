package jobapplicationdtos

type BaseJobApplicationDto struct {
	Title string `json:"title"`
	Description *string `json:"description,omitempty"`
	Employment EmploymentDto `json:"employment" gorm:"embedded"`
	Location *string `json:"location,omitempty"`
}

type JobApplicationCreateDto struct {
	BaseJobApplicationDto
}

type JobApplicationPublicDto struct {
	ID uint `json:"id"`
	Created string `json:"created"`
	BaseJobApplicationDto
}

type JobApplicationPatchDto struct {
	Title *string `json:"title"`
	Description *string `json:"description,omitempty"`
	Employment *EmploymentPatchDto `json:"employment" gorm:"embedded"`
	Location *string `json:"location,omitempty"`
}

type EmploymentDto struct {
	Type string `json:"type"`
	Duration *string `json:"duration,omitempty"`
	WorkLocation string `json:"workLocation" gorm:"column:work_location"`
	SalaryRange *SalaryRangeDto `json:"salaryRange,omitempty" gorm:"column:salary_range"`
}

type EmploymentPatchDto struct {
	Type *string `json:"type"`
	Duration *string `json:"duration,omitempty"`
	WorkLocation *string `json:"workLocation" gorm:"column:work_location"`
	SalaryRange *SalaryRangePatchDto `json:"salaryRange,omitempty" gorm:"column:salary_range"`
}

type SalaryRangeDto struct {
	From int `json:"from"`
	To int `json:"to"`
	Currency string `json:"currency"`
}

type SalaryRangePatchDto struct {
	From *int `json:"from"`
	To *int `json:"to"`
	Currency *string `json:"currency"`
}



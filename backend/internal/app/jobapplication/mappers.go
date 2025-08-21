package jobapplication

import (
	"appliedTo/internal/platform/patch"
	"appliedTo/internal/utils"
	"log"
	"time"
)

// ---- helpers: SalaryRange ----

func toModelSalaryRange(in *SalaryRangeDto) *SalaryRange {
	if in == nil {
		return nil
	}
	return &SalaryRange{
		From:       in.From,
		To:         in.To,
		Currency:   in.Currency,
		Period:     SalaryPeriod(in.Period),
		Negotiable: in.Negotiable,
	}
}

func toDtoSalaryRange(in *SalaryRange) *SalaryRangeDto {
	if in == nil {
		return nil
	}
	return &SalaryRangeDto{
		From:       in.From,
		To:         in.To,
		Currency:   in.Currency,
		Period:     string(in.Period),
		Negotiable: in.Negotiable,
	}
}

// ---- helpers: Employment ----

func mapEmploymentDtoToModel(e EmploymentDto) Employment {
	return Employment{
		Type:         EmploymentType(e.Type),
		Duration:     e.Duration,
		WorkLocation: WorkLocation(e.WorkLocation),
		Seniority:    e.Seniority,
		HoursPerWeek: e.HoursPerWeek,
		SalaryRange:  toModelSalaryRange(e.SalaryRange),
	}
}

func toDtoEmployment(in Employment) EmploymentDto {
	return EmploymentDto{
		Type:         string(in.Type),
		Duration:     in.Duration,
		WorkLocation: string(in.WorkLocation),
		Seniority:    in.Seniority,
		HoursPerWeek: in.HoursPerWeek,
		SalaryRange:  toDtoSalaryRange(in.SalaryRange),
	}
}

func patchEmploymentModel(e *Employment, dto EmploymentPatchDto) {
	if dto.Type != nil {
		e.Type = EmploymentType(*dto.Type)
	}
	patch.PatchRef(&e.Duration, dto.Duration)

	if dto.WorkLocation != nil {
		e.WorkLocation = WorkLocation(*dto.WorkLocation)
	}

	patch.PatchRef(&e.Seniority, dto.Seniority)
	patch.PatchRef(&e.HoursPerWeek, dto.HoursPerWeek)

	if dto.SalaryRange != nil {
		if e.SalaryRange == nil {
			e.SalaryRange = &SalaryRange{}
		}
		patchSalaryRange(e.SalaryRange, *dto.SalaryRange)
	}
}

func patchSalaryRange(sr *SalaryRange, dto SalaryRangePatchDto) {
	patch.Patch(&sr.From, dto.From)
	patch.Patch(&sr.To, dto.To)
	patch.Patch(&sr.Currency, dto.Currency)
	if dto.Period != nil {
		sr.Period = SalaryPeriod(*dto.Period)
	}
	patch.Patch(&sr.Negotiable, dto.Negotiable)
}

// --- INPUT MAPPERS ---

func OverwriteModel(m *JobApplication, dto JobApplicationCreateDto) {
	m.Company        = dto.Company
	m.Title          = dto.Title
	m.Description    = dto.Description
	m.Status         = ApplicationStatus(dto.Status)
	m.Source         = ApplicationSource(dto.Source)
	m.AppliedAt      = dto.AppliedAt
	m.NextFollowUpAt = dto.NextFollowUpAt
	m.LastContactAt  = dto.LastContactAt
	m.PostingURL     = dto.PostingURL
	m.CompanyURL     = dto.CompanyURL
	m.ContactName    = dto.ContactName
	m.ContactEmail   = dto.ContactEmail
	m.ExternalJobID  = dto.ExternalJobID
	m.Employment     = mapEmploymentDtoToModel(dto.Employment)
	m.Location       = dto.Location
	m.Tags           = utils.ToJSONTags(dto.Tags)
}

func CreateModel(dto JobApplicationCreateDto) JobApplication {
	var m JobApplication
	OverwriteModel(&m, dto)
	return m
}

func PatchModel(m *JobApplication, dto JobApplicationPatchDto) {
	patch.Patch(&m.Company, dto.Company)
	patch.Patch(&m.Title, dto.Title)
	patch.PatchRef(&m.Description, dto.Description)

	if dto.Status != nil {
		m.Status = ApplicationStatus(*dto.Status)
	}
	if dto.Source != nil {
		m.Source = ApplicationSource(*dto.Source)
	}

	patch.PatchRef(&m.AppliedAt, dto.AppliedAt)
	patch.PatchRef(&m.NextFollowUpAt, dto.NextFollowUpAt)
	patch.PatchRef(&m.LastContactAt, dto.LastContactAt)

	patch.PatchRef(&m.PostingURL, dto.PostingURL)
	patch.PatchRef(&m.CompanyURL, dto.CompanyURL)
	patch.PatchRef(&m.ContactName, dto.ContactName)
	patch.PatchRef(&m.ContactEmail, dto.ContactEmail)
	patch.PatchRef(&m.ExternalJobID, dto.ExternalJobID)

	if dto.Employment != nil {
		patchEmploymentModel(&m.Employment, *dto.Employment)
	}

	patch.PatchRef(&m.Location, dto.Location)

	if dto.Tags != nil {
		m.Tags = utils.ToJSONTags(*dto.Tags)
	}
}

// --- OUTPUT MAPPER ---

func MapModelToPublicDto(m JobApplication) JobApplicationPublicDto {
	tags := []string{}
	if t, err := utils.FromJSONTags(m.Tags); err == nil {
		tags = t
	} else { log.Printf("invalid tags JSON for job_application id=%d: %v", m.ID, err) }
	return JobApplicationPublicDto{
		ID:      m.ID,
		Created: m.CreatedAt.UTC().Format(time.RFC3339),
		BaseJobApplicationDto: BaseJobApplicationDto{
			Company:        m.Company,
			Title:          m.Title,
			Description:    m.Description,
			Status:         string(m.Status),
			Source:         string(m.Source),
			AppliedAt:      m.AppliedAt,
			NextFollowUpAt: m.NextFollowUpAt,
			LastContactAt:  m.LastContactAt,
			PostingURL:     m.PostingURL,
			CompanyURL:     m.CompanyURL,
			ContactName:    m.ContactName,
			ContactEmail:   m.ContactEmail,
			ExternalJobID:  m.ExternalJobID,
			Employment:     toDtoEmployment(m.Employment),
			Location:       m.Location,
			Tags:           tags,
		},
	}
}

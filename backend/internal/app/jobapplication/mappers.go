package jobapplication

import (
	"appliedTo/internal/platform/patch"
)

// --- PRIVATE HELPER ---

func mapEmploymentDtoToModel(e EmploymentDto) Employment {
	return Employment{
		Type:         EmploymentType(e.Type),
		Duration:     e.Duration,
		WorkLocation: WorkLocation(e.WorkLocation),
		SalaryRange:  (*SalaryRange)(e.SalaryRange),
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
}

// --- INPUT MAPPERS ---

func CreateModel(dto JobApplicationCreateDto) JobApplication {
	return JobApplication{
		Title:       dto.Title,
		Description: dto.Description,
		Location:    dto.Location,
		Employment:  mapEmploymentDtoToModel(dto.Employment),
	}
}

func PatchModel(m *JobApplication, dto JobApplicationPatchDto) {
	patch.Patch(&m.Title, dto.Title)
	patch.PatchRef(&m.Description, dto.Description)
	patch.PatchRef(&m.Location, dto.Location)

	if dto.Employment != nil {
		patchEmploymentModel(&m.Employment, *dto.Employment)
	}
}

// --- OUTPUT MAPPER ---

func MapModelToPublicDto(m JobApplication) JobApplicationPublicDto {
	return JobApplicationPublicDto{
		ID: m.ID,
		BaseJobApplicationDto: BaseJobApplicationDto{
			Title:       m.Title,
			Description: m.Description,
			Location:    m.Location,
			Employment: EmploymentDto{
				Type:         string(m.Employment.Type),
				Duration:     m.Employment.Duration,
				WorkLocation: string(m.Employment.WorkLocation),
				SalaryRange:  (*SalaryRangeDto)(m.Employment.SalaryRange),
			},
		},
	}
}

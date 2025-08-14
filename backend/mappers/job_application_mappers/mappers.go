package mappers

import (
	"appliedTo/dtos/job_application_dtos"
	"appliedTo/models"
	"appliedTo/utils"
)

// --- PRIVATE HELPER ---

func mapEmploymentDtoToModel(e jobapplicationdtos.EmploymentDto) models.Employment {
	return models.Employment{
		Type:         models.EmploymentType(e.Type),
		Duration:     e.Duration,
		WorkLocation: models.WorkLocation(e.WorkLocation),
		SalaryRange:  (*models.SalaryRange)(e.SalaryRange),
	}
}

func patchEmploymentModel(e *models.Employment, dto jobapplicationdtos.EmploymentPatchDto) {
	if dto.Type != nil {
		e.Type = models.EmploymentType(*dto.Type)
	}
	utils.PatchRef(&e.Duration, dto.Duration)

	if dto.WorkLocation != nil {
		e.WorkLocation = models.WorkLocation(*dto.WorkLocation)
	}

	if dto.SalaryRange != nil {
		if e.SalaryRange == nil {
			e.SalaryRange = &models.SalaryRange{}
		}
		patchSalaryRange(e.SalaryRange, *dto.SalaryRange)
	}
}

func patchSalaryRange(sr *models.SalaryRange, dto jobapplicationdtos.SalaryRangePatchDto) {
	utils.Patch(&sr.From, dto.From)
	utils.Patch(&sr.To, dto.To)
	utils.Patch(&sr.Currency, dto.Currency)
}

// --- INPUT MAPPERS ---

func CreateModel(dto jobapplicationdtos.JobApplicationCreateDto) models.JobApplication {
	return models.JobApplication{
		Title:       dto.Title,
		Description: dto.Description,
		Location:    dto.Location,
		Employment:  mapEmploymentDtoToModel(dto.Employment),
	}
}

func PatchModel(m *models.JobApplication, dto jobapplicationdtos.JobApplicationPatchDto) {
	utils.Patch(&m.Title, dto.Title)
	utils.PatchRef(&m.Description, dto.Description)
	utils.PatchRef(&m.Location, dto.Location)

	if dto.Employment != nil {
		patchEmploymentModel(&m.Employment, *dto.Employment)
	}
}

// --- OUTPUT MAPPER ---

func MapModelToPublicDto(m models.JobApplication) jobapplicationdtos.JobApplicationPublicDto {
	return jobapplicationdtos.JobApplicationPublicDto{
		ID: m.ID,
		BaseJobApplicationDto: jobapplicationdtos.BaseJobApplicationDto{
			Title:       m.Title,
			Description: m.Description,
			Location:    m.Location,
			Employment: jobapplicationdtos.EmploymentDto{
				Type:         string(m.Employment.Type),
				Duration:     m.Employment.Duration,
				WorkLocation: string(m.Employment.WorkLocation),
				SalaryRange:  (*jobapplicationdtos.SalaryRangeDto)(m.Employment.SalaryRange),
			},
		},
	}
}

package controllers

import (
	"net/http"
	"strconv"

	"appliedTo/dtos/job_application_dtos"
	mappers "appliedTo/mappers/job_application_mappers"
	"appliedTo/models"
	"appliedTo/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// @Summary Create a new job application
// @Description Creates a new job application with the provided title.
// @Tags jobApplication
// @Accept  json
// @Produce  json
// @Param   jobApplication  body  jobapplicationdtos.JobApplicationCreateDto  true  "JobApplication data"
// @Success 200 {object} map[string]interface{} "Job Application created successfully"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Could not create job application"
// @Router /job_application [post]
func CreateJobApplication(c *gin.Context) {
	var application jobapplicationdtos.JobApplicationCreateDto

	if err := c.ShouldBindJSON(&application); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	if !utils.ValidateRequiredFields(c, []utils.RequiredField{
		{Value: application.BaseJobApplicationDto.Title, Name: "title"},
		{Value: application.BaseJobApplicationDto.Employment.Type, Name: "employment Type"},
		{Value: application.BaseJobApplicationDto.Employment.WorkLocation, Name: "work location"},
	}) {
		return
	}

	jobApplication := mappers.CreateModel(application)

	if err := db.Create(&jobApplication).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create job application"})
		return
	}

	response := mappers.MapModelToPublicDto(jobApplication)

	c.JSON(http.StatusOK, gin.H{"message": "Job application created successfully", "job_application": response})
}

// @Summary Get a job application by ID
// @Description Get detailed information about a job application
// @Tags jobApplication
// @Accept  json
// @Produce  json
// @Param   id  path  int  true  "JobApplication ID"
// @Success 200 {object} jobapplicationdtos.JobApplicationPublicDto "Successfully retrieved job application"
// @Failure 400 "Invalid ID format"
// @Failure 404 "Application not found"
// @Failure 404 "Database query failed"
// @Router /job_application/{id} [get]
func GetJobApplication(c *gin.Context) {
	var applicationId = c.Param("id")
	id, err := strconv.Atoi(applicationId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var application models.JobApplication

	if err := db.First(&application, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Application not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Databse query failed"})
		return
	}

	response := mappers.MapModelToPublicDto(application)
	c.JSON(http.StatusOK, gin.H{"job_application": response})
}

// @Summary Patch a job application
// @Description Partially update a job application. Only fields provided in the body will be modified.
// @Tags jobApplication
// @Accept  json
// @Produce  json
// @Param   id             path  int                                         true  "JobApplication ID"
// @Param   jobApplication body  jobapplicationdtos.JobApplicationPatchDto   true  "Fields to patch"
// @Success 200 {object} jobapplicationdtos.JobApplicationPublicDto "Updated job application"
// @Failure 400 {object} map[string]string "Invalid ID or payload"
// @Failure 404 "Not found"
// @Failure 500 "Update failed"
// @Router /job_application/{id} [patch]
func PatchJobApplication(c *gin.Context) {
	id64, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	id := uint(id64)

	var patch jobapplicationdtos.JobApplicationPatchDto
	if err := c.ShouldBindJSON(&patch); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	var m models.JobApplication
	if err := db.First(&m, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	mappers.PatchModel(&m, patch)

	if err := db.Save(&m).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "update failed"})
		return
	}

	resp := mappers.MapModelToPublicDto(m)
	c.JSON(http.StatusOK, gin.H{"message": "updated", "job_application": resp})
}

// @Summary Update a job application (full replace)
// @Description Replace a job application with the provided data. All fields should be supplied.
// @Tags jobApplication
// @Accept  json
// @Produce  json
// @Param   id              path  int                                         true  "JobApplication ID"
// @Param   jobApplication  body  jobapplicationdtos.JobApplicationCreateDto  true  "JobApplication data"
// @Success 200 {object} jobapplicationdtos.JobApplicationPublicDto "Job application successfully updated."
// @Failure 400 "Invalid ID format"
// @Failure 404 "Application not found"
// @Failure 500 "Database query failed"
// @Router /job_application/{id} [put]
func UpdateJobApplication(c *gin.Context) {
	var application jobapplicationdtos.JobApplicationCreateDto

	if err := c.ShouldBindJSON(&application); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	applicationID := c.Param("id")
	id, err := strconv.Atoi(applicationID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var jobApplication models.JobApplication
	if err := db.First(&jobApplication, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Job application not found"})
		return
	}

	if !utils.ValidateRequiredFields(c, []utils.RequiredField{
		{Value: application.BaseJobApplicationDto.Title, Name: "title"},
		{Value: application.BaseJobApplicationDto.Employment.Type, Name: "employment Type"},
		{Value: application.BaseJobApplicationDto.Employment.WorkLocation, Name: "work location"},
	}) {
		return
	}

	jobApplication = mappers.CreateModel(application)


	if err := db.Model(&jobApplication).Updates(jobApplication).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update job application"})
		return
	}

	response := mappers.MapModelToPublicDto(jobApplication)

	c.JSON(http.StatusOK, gin.H{"message": "Job application updated successfully", "job_application": response})
}

// @Summary Delete a job application.
// @Description Remove the job application from the database by providing the job-application-ID.
// @Tags jobApplication
// @Accept  json
// @Produce  json
// @Param   id  path  int  true  "job application ID"
// @Success 200 "Job application deleted."
// @Failure 400 "Invalid ID format"
// @Failure 500 "Could not delete job application"
// @Router /job_application/{id} [delete]
func DeleteJobApplication(c *gin.Context) {
	applicationID := c.Param("id")
	id, err := strconv.Atoi(applicationID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	if err := db.Delete(&models.JobApplication{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete job application"})
	}

	c.JSON(http.StatusOK, gin.H{"message": "Job application deleted successfully"})
}


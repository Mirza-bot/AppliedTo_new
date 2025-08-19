package jobapplicationapi

import (
	"appliedTo/internal/app/jobapplication"
	"appliedTo/internal/platform/http/middleware"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handlers struct {
	Svc *jobapplication.Service
}

func NewHandlers(s *jobapplication.Service) *Handlers { return &Handlers{Svc: s} }

// @Summary Create a new job application
// @Description Creates a new job application with the provided title.
// @Tags jobApplication
// @Accept  json
// @Produce  json
// @Param   jobApplication  body  jobapplication.JobApplicationCreateDto  true  "JobApplication data"
// @Success 200 {object} map[string]interface{} "Job application created successfully"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Could not create job application"
// @Router /job_application [post]
func (h *Handlers) CreateJobApplication(c *gin.Context) {
	var in jobapplication.JobApplicationCreateDto
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	out, err := h.Svc.Create(c.Request.Context(), in)
	if err != nil {
		// validation.Required(...) returns a descriptive error string
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Job application created successfully", "job_application": out})
}

// @Summary Get a job application by ID
// @Description Get detailed information about a job application
// @Tags jobApplication
// @Accept  json
// @Produce  json
// @Param   id  path  int  true  "JobApplication ID"
// @Success 200 {object} jobapplication.JobApplicationPublicDto "Successfully retrieved job application"
// @Failure 404 {object} map[string]string "Application not found"
// @Failure 500 {object} map[string]string "Database query failed"
// @Router /job_application/{id} [get]
func (h *Handlers) GetJobApplication(c *gin.Context) {
	id := c.GetUint(middleware.CtxKeyJobApplicationID)

	out, err := h.Svc.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Application not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database query failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"job_application": out})
}

// @Summary Patch a job application
// @Description Partially update a job application. Only fields provided in the body will be modified.
// @Tags jobApplication
// @Accept  json
// @Produce  json
// @Param   id             path  int                                       true  "JobApplication ID"
// @Param   jobApplication body  jobapplication.JobApplicationPatchDto true  "Fields to patch"
// @Success 200 {object} jobapplication.JobApplicationPublicDto "Updated job application"
// @Failure 400 {object} map[string]string "Invalid payload"
// @Failure 404 {object} map[string]string "Not found"
// @Failure 500 {object} map[string]string "Update failed"
// @Router /job_application/{id} [patch]
func (h *Handlers) PatchJobApplication(c *gin.Context) {
	id := c.GetUint(middleware.CtxKeyJobApplicationID)

	var patch jobapplication.JobApplicationPatchDto
	if err := c.ShouldBindJSON(&patch); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	out, err := h.Svc.Patch(c.Request.Context(), id, patch)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Update failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Updated", "job_application": out})
}

// @Summary Update a job application (full replace)
// @Description Replace a job application with the provided data. All fields should be supplied.
// @Tags jobApplication
// @Accept  json
// @Produce  json
// @Param   id              path  int                                       true  "JobApplication ID"
// @Param   jobApplication  body  jobapplication.JobApplicationCreateDto true  "JobApplication data"
// @Success 200 {object} jobapplication.JobApplicationPublicDto "Job application successfully updated."
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 404 {object} map[string]string "Application not found"
// @Failure 500 {object} map[string]string "Database query failed"
// @Router /job_application/{id} [put]
func (h *Handlers) UpdateJobApplication(c *gin.Context) {
	id := c.GetUint(middleware.CtxKeyJobApplicationID)

	var in jobapplication.JobApplicationCreateDto
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	out, err := h.Svc.Update(c.Request.Context(), id, in)
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "Application not found"})
		default:
			// could be validation error or DB error
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Job application updated successfully", "job_application": out})
}

// @Summary Delete a job application.
// @Description Remove the job application from the database by providing the job-application-ID.
// @Tags jobApplication
// @Accept  json
// @Produce  json
// @Param   id  path  int  true  "Job application ID"
// @Success 200 {object} map[string]string "Job application deleted."
// @Failure 404 {object} map[string]string "Not found"
// @Failure 500 {object} map[string]string "Could not delete job application"
// @Router /job_application/{id} [delete]
func (h *Handlers) DeleteJobApplication(c *gin.Context) {
	id := c.GetUint(middleware.CtxKeyJobApplicationID)

	if err := h.Svc.Delete(c.Request.Context(), id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete job application"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Job application deleted successfully"})
}


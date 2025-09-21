package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// Assignment represents a bus-staff assignment
type Assignment struct {
	ID        int        `json:"id" db:"id"`
	BusID     int        `json:"bus_id" db:"bus_id"`
	StaffID   int        `json:"staff_id" db:"staff_id"`
	Role      string     `json:"role" db:"role"` // driver, conductor
	StartDate time.Time  `json:"start_date" db:"start_date"`
	EndDate   *time.Time `json:"end_date,omitempty" db:"end_date"`
	Status    string     `json:"status" db:"status"` // active, completed, cancelled
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
}

// AssignmentWithDetails includes bus and staff information
type AssignmentWithDetails struct {
	Assignment
	BusPlateNumber string `json:"bus_plate_number,omitempty"`
	BusModel       string `json:"bus_model,omitempty"`
	StaffName      string `json:"staff_name,omitempty"`
	StaffPosition  string `json:"staff_position,omitempty"`
}

// Request structs
type CreateAssignmentRequest struct {
	BusID     int    `json:"bus_id" binding:"required"`
	StaffID   int    `json:"staff_id" binding:"required"`
	Role      string `json:"role" binding:"required"`
	StartDate string `json:"start_date" binding:"required"` // YYYY-MM-DD format
	EndDate   string `json:"end_date,omitempty"`
}

// Mock data for demonstration (would come from other services in production)
var mockBuses = map[int]map[string]string{
	1: {"plate_number": "ABC-1234", "model": "Toyota Coaster"},
	2: {"plate_number": "XYZ-5678", "model": "Isuzu NPR"},
}

var mockStaff = map[int]map[string]string{
	1: {"name": "John Driver", "position": "driver"},
	2: {"name": "Jane Conductor", "position": "conductor"},
}

func handleCreateAssignment(c *gin.Context) {
	var req CreateAssignmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse start date
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_date format. Use YYYY-MM-DD"})
		return
	}

	// Parse end date if provided
	var endDate *time.Time
	if req.EndDate != "" {
		ed, err := time.Parse("2006-01-02", req.EndDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end_date format. Use YYYY-MM-DD"})
			return
		}
		endDate = &ed
	}

	// Validate role
	if req.Role != "driver" && req.Role != "conductor" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Role must be 'driver' or 'conductor'"})
		return
	}

	assignment := Assignment{
		BusID:     req.BusID,
		StaffID:   req.StaffID,
		Role:      req.Role,
		StartDate: startDate,
		EndDate:   endDate,
		Status:    "active",
	}

	if err := CreateAssignment(&assignment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create assignment"})
		return
	}

	c.JSON(http.StatusCreated, assignment)
}

func handleGetAssignments(c *gin.Context) {
	assignments, err := GetAllAssignments()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve assignments"})
		return
	}

	assignmentList := make([]AssignmentWithDetails, 0, len(assignments))
	for _, assignment := range assignments {
		details := AssignmentWithDetails{
			Assignment: assignment,
		}

		// Add bus details if available
		if bus, exists := mockBuses[assignment.BusID]; exists {
			details.BusPlateNumber = bus["plate_number"]
			details.BusModel = bus["model"]
		}

		// Add staff details if available
		if staff, exists := mockStaff[assignment.StaffID]; exists {
			details.StaffName = staff["name"]
			details.StaffPosition = staff["position"]
		}

		assignmentList = append(assignmentList, details)
	}

	c.JSON(http.StatusOK, gin.H{"assignments": assignmentList, "count": len(assignmentList)})
}

func handleGetAssignment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid assignment ID"})
		return
	}

	assignment, err := GetAssignmentByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	if assignment == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Assignment not found"})
		return
	}

	c.JSON(http.StatusOK, assignment)
}

func handleUpdateAssignment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid assignment ID"})
		return
	}

	// Check if assignment exists
	existingAssignment, err := GetAssignmentByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	if existingAssignment == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Assignment not found"})
		return
	}

	var req CreateAssignmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse start date
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_date format. Use YYYY-MM-DD"})
		return
	}

	// Parse end date if provided
	var endDate *time.Time
	if req.EndDate != "" {
		ed, err := time.Parse("2006-01-02", req.EndDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end_date format. Use YYYY-MM-DD"})
			return
		}
		endDate = &ed
	}

	// Update assignment fields
	existingAssignment.BusID = req.BusID
	existingAssignment.StaffID = req.StaffID
	existingAssignment.Role = req.Role
	existingAssignment.StartDate = startDate
	existingAssignment.EndDate = endDate

	if err := UpdateAssignment(existingAssignment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update assignment"})
		return
	}

	c.JSON(http.StatusOK, existingAssignment)
}

func handleDeleteAssignment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid assignment ID"})
		return
	}

	// Check if assignment exists
	existingAssignment, err := GetAssignmentByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	if existingAssignment == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Assignment not found"})
		return
	}

	if err := DeleteAssignment(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete assignment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Assignment deleted successfully"})
}

func handleGetStaffForBus(c *gin.Context) {
	busIDStr := c.Param("busId")
	busID, err := strconv.Atoi(busIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid bus ID"})
		return
	}

	assignments, err := GetAssignmentsByBusID(busID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve assignments"})
		return
	}

	busAssignments := make([]AssignmentWithDetails, 0)
	for _, assignment := range assignments {
		if assignment.Status == "active" {
			details := AssignmentWithDetails{
				Assignment: assignment,
			}

			// Add staff details if available
			if staff, exists := mockStaff[assignment.StaffID]; exists {
				details.StaffName = staff["name"]
				details.StaffPosition = staff["position"]
			}

			busAssignments = append(busAssignments, details)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"bus_id":      busID,
		"assignments": busAssignments,
		"count":       len(busAssignments),
	})
}

func handleGetAssignmentsForStaff(c *gin.Context) {
	staffIDStr := c.Param("staffId")
	staffID, err := strconv.Atoi(staffIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid staff ID"})
		return
	}

	assignments, err := GetAssignmentsByStaffID(staffID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve assignments"})
		return
	}

	staffAssignments := make([]AssignmentWithDetails, 0)
	for _, assignment := range assignments {
		details := AssignmentWithDetails{
			Assignment: assignment,
		}

		// Add bus details if available
		if bus, exists := mockBuses[assignment.BusID]; exists {
			details.BusPlateNumber = bus["plate_number"]
			details.BusModel = bus["model"]
		}

		staffAssignments = append(staffAssignments, details)
	}

	c.JSON(http.StatusOK, gin.H{
		"staff_id":    staffID,
		"assignments": staffAssignments,
		"count":       len(staffAssignments),
	})
}

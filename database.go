package main

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var db *pgxpool.Pool

// InitDB initializes the database connection pool
func InitDB() error {
	var err error
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	// Create connection pool
	db, err = pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		return err
	}

	// Test the connection
	if err := db.Ping(context.Background()); err != nil {
		return err
	}

	log.Println("Connected to Supabase database")

	// Create tables if they don't exist
	if err := createTables(); err != nil {
		return err
	}

	return nil
}

// CloseDB closes the database connection pool
func CloseDB() {
	if db != nil {
		db.Close()
	}
}

// createTables creates the assignments table if it doesn't exist
func createTables() error {
	query := `
	CREATE TABLE IF NOT EXISTS assignments (
		id SERIAL PRIMARY KEY,
		bus_id INTEGER NOT NULL,
		staff_id INTEGER NOT NULL,
		role VARCHAR(20) NOT NULL CHECK (role IN ('driver', 'conductor')),
		start_date DATE NOT NULL,
		end_date DATE,
		status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'completed', 'cancelled')),
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		UNIQUE(bus_id, staff_id, role, start_date)
	);

	-- Create indexes for better performance
	CREATE INDEX IF NOT EXISTS idx_assignments_bus_id ON assignments(bus_id);
	CREATE INDEX IF NOT EXISTS idx_assignments_staff_id ON assignments(staff_id);
	CREATE INDEX IF NOT EXISTS idx_assignments_status ON assignments(status);
	CREATE INDEX IF NOT EXISTS idx_assignments_start_date ON assignments(start_date);
	`

	_, err := db.Exec(context.Background(), query)
	if err != nil {
		log.Printf("Error creating assignments table: %v", err)
		return err
	}

	log.Println("Assignments table created successfully")
	return nil
}

// Assignment database operations

// CreateAssignment inserts a new assignment into the database
func CreateAssignment(assignment *Assignment) error {
	query := `
		INSERT INTO assignments (bus_id, staff_id, role, start_date, end_date, status)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`

	err := db.QueryRow(context.Background(), query, assignment.BusID, assignment.StaffID,
		assignment.Role, assignment.StartDate, assignment.EndDate, assignment.Status).
		Scan(&assignment.ID, &assignment.CreatedAt, &assignment.UpdatedAt)

	return err
}

// GetAssignmentByID retrieves an assignment by ID
func GetAssignmentByID(id int) (*Assignment, error) {
	assignment := &Assignment{}
	query := `
		SELECT id, bus_id, staff_id, role, start_date, end_date, status, created_at, updated_at
		FROM assignments
		WHERE id = $1
	`

	err := db.QueryRow(context.Background(), query, id).
		Scan(&assignment.ID, &assignment.BusID, &assignment.StaffID, &assignment.Role,
			&assignment.StartDate, &assignment.EndDate, &assignment.Status,
			&assignment.CreatedAt, &assignment.UpdatedAt)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // Assignment not found
		}
		return nil, err
	}

	return assignment, nil
}

// GetAllAssignments retrieves all assignments from the database
func GetAllAssignments() ([]Assignment, error) {
	var assignments []Assignment
	query := `
		SELECT id, bus_id, staff_id, role, start_date, end_date, status, created_at, updated_at
		FROM assignments
		ORDER BY created_at DESC
	`

	rows, err := db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var assignment Assignment
		err := rows.Scan(&assignment.ID, &assignment.BusID, &assignment.StaffID, &assignment.Role,
			&assignment.StartDate, &assignment.EndDate, &assignment.Status,
			&assignment.CreatedAt, &assignment.UpdatedAt)
		if err != nil {
			return nil, err
		}
		assignments = append(assignments, assignment)
	}

	return assignments, nil
}

// GetAssignmentsByBusID retrieves all assignments for a specific bus
func GetAssignmentsByBusID(busID int) ([]Assignment, error) {
	var assignments []Assignment
	query := `
		SELECT id, bus_id, staff_id, role, start_date, end_date, status, created_at, updated_at
		FROM assignments
		WHERE bus_id = $1
		ORDER BY created_at DESC
	`

	rows, err := db.Query(context.Background(), query, busID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var assignment Assignment
		err := rows.Scan(&assignment.ID, &assignment.BusID, &assignment.StaffID, &assignment.Role,
			&assignment.StartDate, &assignment.EndDate, &assignment.Status,
			&assignment.CreatedAt, &assignment.UpdatedAt)
		if err != nil {
			return nil, err
		}
		assignments = append(assignments, assignment)
	}

	return assignments, nil
}

// GetAssignmentsByStaffID retrieves all assignments for a specific staff member
func GetAssignmentsByStaffID(staffID int) ([]Assignment, error) {
	var assignments []Assignment
	query := `
		SELECT id, bus_id, staff_id, role, start_date, end_date, status, created_at, updated_at
		FROM assignments
		WHERE staff_id = $1
		ORDER BY created_at DESC
	`

	rows, err := db.Query(context.Background(), query, staffID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var assignment Assignment
		err := rows.Scan(&assignment.ID, &assignment.BusID, &assignment.StaffID, &assignment.Role,
			&assignment.StartDate, &assignment.EndDate, &assignment.Status,
			&assignment.CreatedAt, &assignment.UpdatedAt)
		if err != nil {
			return nil, err
		}
		assignments = append(assignments, assignment)
	}

	return assignments, nil
}

// UpdateAssignment updates an existing assignment
func UpdateAssignment(assignment *Assignment) error {
	query := `
		UPDATE assignments
		SET bus_id = $1, staff_id = $2, role = $3, start_date = $4, end_date = $5, status = $6, updated_at = CURRENT_TIMESTAMP
		WHERE id = $7
		RETURNING updated_at
	`

	err := db.QueryRow(context.Background(), query, assignment.BusID, assignment.StaffID,
		assignment.Role, assignment.StartDate, assignment.EndDate, assignment.Status, assignment.ID).
		Scan(&assignment.UpdatedAt)

	return err
}

// DeleteAssignment deletes an assignment by ID
func DeleteAssignment(id int) error {
	query := `DELETE FROM assignments WHERE id = $1`
	_, err := db.Exec(context.Background(), query, id)
	return err
}

# Bus Staff Assignment Service

Assignment management microservice for assigning bus staff to buses in the Choreo platform testing project.

## Features

- Create staff-to-bus assignments
- Manage assignment lifecycle (active, completed, cancelled)
- Query assignments by bus or staff member
- Role-based assignments (driver, conductor)
- Date-based assignment periods

## API Endpoints

### Health Check

- `GET /health` - Service health check

### Assignment Management

- `POST /api/assignments` - Create new assignment
- `GET /api/assignments` - List all assignments
- `GET /api/assignments/:id` - Get specific assignment
- `PUT /api/assignments/:id` - Update assignment
- `DELETE /api/assignments/:id` - Delete assignment

### Query Operations

- `GET /api/assignments/bus/:busId` - Get all staff assigned to a specific bus
- `GET /api/assignments/staff/:staffId` - Get all bus assignments for a specific staff member

## Request/Response Examples

### Create Assignment

```bash
POST /api/assignments
Content-Type: application/json

{
  "bus_id": 1,
  "staff_id": 1,
  "role": "driver",
  "start_date": "2025-09-21",
  "end_date": "2025-12-31"
}
```

Response:

```json
{
  "id": 1,
  "bus_id": 1,
  "staff_id": 1,
  "role": "driver",
  "start_date": "2025-09-21T00:00:00Z",
  "end_date": "2025-12-31T00:00:00Z",
  "status": "active",
  "created_at": "2025-09-21T13:30:00Z",
  "updated_at": "2025-09-21T13:30:00Z"
}
```

### Get Staff for Bus

```bash
GET /api/assignments/bus/1
```

Response:

```json
{
  "bus_id": 1,
  "assignments": [
    {
      "id": 1,
      "bus_id": 1,
      "staff_id": 1,
      "role": "driver",
      "start_date": "2025-09-21T00:00:00Z",
      "end_date": "2025-12-31T00:00:00Z",
      "status": "active",
      "staff_name": "John Driver",
      "staff_position": "driver",
      "created_at": "2025-09-21T13:30:00Z",
      "updated_at": "2025-09-21T13:30:00Z"
    }
  ],
  "count": 1
}
```

### Get Assignments for Staff

```bash
GET /api/assignments/staff/1
```

Response:

```json
{
  "staff_id": 1,
  "assignments": [
    {
      "id": 1,
      "bus_id": 1,
      "staff_id": 1,
      "role": "driver",
      "start_date": "2025-09-21T00:00:00Z",
      "end_date": "2025-12-31T00:00:00Z",
      "status": "active",
      "bus_plate_number": "ABC-1234",
      "bus_model": "Toyota Coaster",
      "created_at": "2025-09-21T13:30:00Z",
      "updated_at": "2025-09-21T13:30:00Z"
    }
  ],
  "count": 1
}
```

## Running the Service

```bash
# Install dependencies
go mod tidy

# Run the service
go run .

# Or build and run
go build -o assignment-service
./assignment-service
```

## Environment Variables

- `PORT` - Server port (default: 8082)
- `GIN_MODE` - Gin framework mode (debug/release)
- `DB_HOST` - Database host
- `DB_PORT` - Database port
- `DB_USER` - Database user
- `DB_PASSWORD` - Database password
- `DB_NAME` - Database name
- `AUTH_SERVICE_URL` - Auth service URL for validation
- `BUS_MANAGEMENT_SERVICE_URL` - Bus management service URL

## Docker

```bash
# Build image
docker build -t assignment-service .

# Run container
docker run -p 8082:8082 assignment-service
```

## Data Models

### Assignment

- `id` - Unique identifier
- `bus_id` - Reference to bus (from bus-management service)
- `staff_id` - Reference to staff member (from bus-management service)
- `role` - Assignment role (driver, conductor)
- `start_date` - Assignment start date
- `end_date` - Assignment end date (optional)
- `status` - Assignment status (active, completed, cancelled)
- `created_at` - Creation timestamp
- `updated_at` - Last update timestamp

## Business Rules

- Each assignment must have a valid bus_id and staff_id
- Role can be either "driver" or "conductor"
- Start date is required, end date is optional
- Multiple staff can be assigned to the same bus with different roles
- Staff can have multiple assignments over time
- Only one active assignment per staff member per bus at a time (in production, add validation)

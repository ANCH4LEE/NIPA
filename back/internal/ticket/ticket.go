package ticket

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

type Ticket struct {
	ID          int       `json:"ticket_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	ContactInfo string    `json:"contact"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Ticketdatabase interface {
	GetAllTicket(ctx context.Context, status string, orderBy string) ([]Ticket, error)
	CreateTicket(ctx context.Context, title, description, contact string) error
	UpdateTicket(ctx context.Context, id int, status string) error
	Close() error
	Ping() error
	Reconnect(connStr string) error
}

type PostgresDatabase struct {
	db *sql.DB
}

func NewPostgresDatabase(connStr string) (*PostgresDatabase, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		fmt.Printf("Database ping failed: %v\n", err)
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	fmt.Println("Database connected successfully")
	return &PostgresDatabase{db: db}, nil
}

func (pdb *PostgresDatabase) GetAllTicket(ctx context.Context, status string, orderBy string) ([]Ticket, error) {
	if pdb.db == nil {
		return nil, fmt.Errorf("database connection is not initialized")
	}
	query := `SELECT ticket_id, title, description, contact, status, created_at, updated_at FROM tickets`
	var args []interface{}
	var conditions []string

	if status != "" {
		conditions = append(conditions, "status = $1")
		args = append(args, status)
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	validOrders := map[string]bool{"updated_at DESC": true, "updated_at ASC": true, "status ASC": true, "status DESC": true}
	if !validOrders[orderBy] {
		orderBy = "updated_at DESC" // ให้ DESC เป็นค่าเริ่มต้น
	}

	query += " ORDER BY " + orderBy

	rows, err := pdb.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query error: %v", err)
	}
	defer rows.Close()

	var tickets []Ticket
	for rows.Next() {
		var t Ticket
		if err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.ContactInfo, &t.Status, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan error: %v", err)
		}
		tickets = append(tickets, t)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %v", err)
	}

	return tickets, nil
}

func (pdb *PostgresDatabase) CreateTicket(ctx context.Context, title, description, contact string) error {
	query := `
        INSERT INTO tickets (title, description, contact, status, created_at, updated_at) 
        vALUES ($1, $2, $3, 'Pending', NOW(), NOW())`

	_, err := pdb.db.ExecContext(ctx, query, title, description, contact)
	if err != nil {
		return fmt.Errorf("insert error: %v", err)
	}
	return nil
}

func (pdb *PostgresDatabase) UpdateTicket(ctx context.Context, id int, status string) error {
	// ตรวจสอบค่า status
	validStatuses := map[string]bool{"Pending": true, "Accepted": true, "Resolved": true, "Rejected": true}
	if !validStatuses[status] {
		return fmt.Errorf("invalid status value: %s", status)
	}

	var ticket Ticket
	query := `
		UPDATE tickets 
		SET status = $1, updated_at = NOW() 
		WHERE ticket_id = $2 
		RETURNING ticket_id, status, created_at, updated_at`

	err := pdb.db.QueryRowContext(ctx, query, status, id).
		Scan(&ticket.ID, &ticket.Status, &ticket.CreatedAt, &ticket.UpdatedAt)

	if err == sql.ErrNoRows {
		return errors.New("not found")
	} else if err != nil {
		return fmt.Errorf("update error: %v", err)
	}

	return nil
}

func (pdb *PostgresDatabase) Close() error {
	return pdb.db.Close()
}

func (pdb *PostgresDatabase) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return pdb.db.PingContext(ctx)
}

func (pdb *PostgresDatabase) Reconnect(connStr string) error {
	if pdb.db != nil {
		pdb.db.Close()
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("failed to ping database: %v", err)
	}

	pdb.db = db
	return nil
}

type TicketRepo struct {
	db Ticketdatabase
}

func NewTicketRepo(db Ticketdatabase) *TicketRepo {
	if db == nil {
		fmt.Println("Error: Ticketdatabase is nil! Check database initialization.")
	}
	return &TicketRepo{db: db}
}

func (tk *TicketRepo) GetAllTicket(ctx context.Context, status string, orderBy string) ([]Ticket, error) {
	if tk.db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}
	return tk.db.GetAllTicket(ctx, status, orderBy)
}

func (tk *TicketRepo) CreateTicket(ctx context.Context, title, description, contact string) error {
	return tk.db.CreateTicket(ctx, title, description, contact)
}

func (tk *TicketRepo) UpdateTicket(ctx context.Context, id int, status string) error {
	return tk.db.UpdateTicket(ctx, id, status)
}

func (tk *TicketRepo) Close() error {
	return tk.db.Close()
}

func (tk *TicketRepo) Ping() error {
	if tk.db == nil {
		return fmt.Errorf("database connection is not initialized")
	}
	return tk.db.Ping()
}

func (tk *TicketRepo) Reconnect(connStr string) error {
	return tk.db.Reconnect(connStr)
}

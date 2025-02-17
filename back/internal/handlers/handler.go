package handlers

import (
	"back/internal/ticket"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TicketHandlers struct {
	repo *ticket.TicketRepo
}

func NewTicketHandlers(repo *ticket.TicketRepo) *TicketHandlers {
	return &TicketHandlers{repo: repo}
}

func (h *TicketHandlers) GetTickets(c *gin.Context) {
	status := c.DefaultQuery("status", "")
	orderBy := c.DefaultQuery("orderBy", "updated_at DESC")

	tickets, err := h.repo.GetAllTicket(c.Request.Context(), status, orderBy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tickets)
}

func (h *TicketHandlers) CreateTicket(c *gin.Context) {
	var newTicket struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		ContactInfo string `json:"contact"`
	}

	if err := c.ShouldBindJSON(&newTicket); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
		return
	}

	err := h.repo.CreateTicket(c.Request.Context(), newTicket.Title, newTicket.Description, newTicket.ContactInfo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Ticket created successfully"})
}

func (h *TicketHandlers) UpdateTicketStatus(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ticket ID"})
		return
	}

	var req struct {
		Status string `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	if req.Status == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Status is required"})
		return
	}

	// ตรวจสอบค่า status ว่าถูกต้องหรือไม่
	validStatuses := map[string]bool{"Pending": true, "Accepted": true, "Resolved": true, "Rejected": true}
	if !validStatuses[req.Status] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status value"})
		return
	}

	err = h.repo.UpdateTicket(c.Request.Context(), id, req.Status)
	if err != nil {
		if err.Error() == "not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Ticket not found"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Ticket status updated successfully"})
}

func (h *TicketHandlers) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "healthy"})
}

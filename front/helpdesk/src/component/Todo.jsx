import React, { useState, useEffect } from "react";
import { Container, Row, Col, Card, Button, Modal, Form, Alert } from "react-bootstrap";
import axios from "axios";
import "./Change.css";
import "bootstrap-icons/font/bootstrap-icons.css";

const TodoApp = () => {
  const [showModal, setShowModal] = useState(false);
  const [newTask, setNewTask] = useState({ title: "", description: "", contact: "", status: "Pending" });
  const [tasks, setTasks] = useState([]);
  const [filter, setFilter] = useState("All Tickets");
  const [error, setError] = useState("");
  const [showStatusConfirmModal, setShowStatusConfirmModal] = useState(false);
  const [selectedTask, setSelectedTask] = useState(null);
  const [newStatus, setNewStatus] = useState("");

  useEffect(() => {
    const fetchTasks = async () => {
      try {
        const response = await axios.get('http://localhost:8080/api/v1/tickets');
        setTasks(response.data);
      } catch (err) {
        console.error(err);
      }
    };

    fetchTasks();
  }, []);

  const handleShow = () => { setError(""); setShowModal(true); };
  const handleClose = () => setShowModal(false);
  const handleChange = (e) => {
    setNewTask({ ...newTask, [e.target.name]: e.target.value });
  };

  const handleSubmit = async () => {
    if (!newTask.title || !newTask.description || !newTask.contact) {
      setError("Please complete all required fields.");
      return;
    }

    try {
      await axios.post('http://localhost:8080/api/v1/tickets', newTask);
      const response = await axios.get('http://localhost:8080/api/v1/tickets');
      setTasks(response.data);
      setShowModal(false);
      setNewTask({ title: "", description: "", contact: "", status: "Pending" });
      setError("");
    } catch (err) {
      console.error(err);
    }
  };

  const handleStatusChange = async (task, status) => {
    setSelectedTask(task);
    setNewStatus(status);
    setShowStatusConfirmModal(true);
  };

  const handleConfirmStatusChange = async () => {
    if (!selectedTask) return;

    const updatedTask = { ...selectedTask, status: newStatus, updatedAt: new Date().toISOString() };

    try {
      await axios.put(`http://localhost:8080/api/v1/tickets/${selectedTask.ticket_id}`, updatedTask);
      setTasks((prevTasks) =>
        prevTasks.map((t) =>
          t.ticket_id === selectedTask.ticket_id
            ? { ...t, status: newStatus, updatedAt: updatedTask.updatedAt }
            : t
        )
      );
      setShowStatusConfirmModal(false);
    } catch (err) {
      console.error("Error updating status:", err);
    }
  };

  const filteredTasks = tasks.filter(task => {
    if (filter === "All Tickets") return true;
    return task.status === filter;
  });

  return (
    <Container fluid className="bg-white text-dark min-vh-100">
      <Row>
        {/**/ }
        <Col md={2} className="sidebar bg-light p-3 border-end">
          <ul className="list-unstyled">
            <li className={`p-2 text-secondary ${filter === "All Tickets" ? "active" : ""}`} onClick={() => setFilter("All Tickets")}>All Tickets</li>
            <li className={`p-2 text-secondary ${filter === "Pending" ? "active" : ""}`} onClick={() => setFilter("Pending")}>Pending</li>
            <li className={`p-2 text-secondary ${filter === "Accepted" ? "active" : ""}`} onClick={() => setFilter("Accepted")}>Accepted</li>
            <li className={`p-2 text-secondary ${filter === "Resolved" ? "active" : ""}`} onClick={() => setFilter("Resolved")}>Resolved</li>
            <li className={`p-2 text-secondary ${filter === "Rejected" ? "active" : ""}`} onClick={() => setFilter("Rejected")}>Rejected</li>
          </ul>
        </Col>
        <Col md={10} className="p-4">
          <Row className="align-items-center mb-3">
            <Col><h3>{filter}</h3></Col>
            <Col className="text-end">
              <Button variant="outline-primary" onClick={handleShow}>+ Add New Problem</Button>
            </Col>
          </Row>

          <Row>
            {filteredTasks.map((task, index) => (
              <Col md={4} key={index} className="mb-3">
                <Card className="bg-white text-dark shadow-sm border">
                  <Card.Body>
                    <Card.Title>{task.title}</Card.Title>
                    <Card.Text>{task.description}</Card.Text>
                    <Card.Text><small>{task.contact}</small></Card.Text>
                    <div className="select-container">
                    <Form.Control 
                      as="select"
                      value={task.status} style={{ width: '100%', maxWidth: '120px' }}
                      onChange={async (e) => {
                        const newStatus = e.target.value;
                        handleStatusChange(task, newStatus);
                      }}
                      className={`status-${task.status ? task.status.toLowerCase() : ''}`}
                    >
                      <option value="Pending" disabled={task.status !== "Pending"}>Pending</option>
                      <option value="Accepted" disabled={task.status !== "Pending" && task.status !== "Accepted"}>Accepted</option>
                      <option value="Resolved" disabled={task.status === "Rejected" || task.status === "Resolved"}>Resolved</option>
                      <option value="Rejected" disabled={task.status === "Rejected"}>Rejected</option>
                    </Form.Control>
                    </div>
                  </Card.Body>
                </Card>
              </Col>
            ))}
          </Row>
        </Col>
      </Row>

      {/* Add new Ticket */}
      <Modal show={showModal} onHide={handleClose}>
        <Modal.Header closeButton>
          <Modal.Title>Create a Problem</Modal.Title>
        </Modal.Header>
        <Modal.Body>
          {error && <Alert variant="danger">{error}</Alert>}
          <Form>
            <Form.Group className="mb-3">
              <Form.Label>Title</Form.Label>
              <Form.Control type="text" placeholder="Enter problem title" name="title" value={newTask.title} onChange={handleChange} />
            </Form.Group>
            <Form.Group className="mb-3">
              <Form.Label>Description</Form.Label>
              <Form.Control as="textarea" rows={3} placeholder="Enter problem description" name="description" value={newTask.description} onChange={handleChange} />
            </Form.Group>
            <Form.Group className="mb-3">
              <Form.Label>Contact</Form.Label>
              <Form.Control type="text" placeholder="Enter email" name="contact" value={newTask.contact} onChange={handleChange} />
            </Form.Group>
          </Form>
        </Modal.Body>
        <Modal.Footer>
          <Button variant="secondary" onClick={handleClose}>Close</Button>
          <Button variant="primary" onClick={handleSubmit}>Create New Ticket</Button>
        </Modal.Footer>
      </Modal>

      {/* ask Status Change */}
      <Modal show={showStatusConfirmModal} onHide={() => setShowStatusConfirmModal(false)}>
        <Modal.Header closeButton>
          <Modal.Title>Confirm Status Change</Modal.Title>
        </Modal.Header>
        <Modal.Body>
          <p>Are you sure you want to change the status to <strong>{newStatus}</strong>?</p>
        </Modal.Body>
        <Modal.Footer>
          <Button variant="secondary" onClick={() => setShowStatusConfirmModal(false)}>Cancel</Button>
          <Button variant="primary" onClick={handleConfirmStatusChange}>Confirm</Button>
        </Modal.Footer>
      </Modal>
    </Container>
  );
};

export default TodoApp;

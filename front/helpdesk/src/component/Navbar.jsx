import React from "react";
import { Navbar, Container } from "react-bootstrap";
import "bootstrap/dist/css/bootstrap.min.css";

function TopNav(){
    return (
    <Navbar expand="lg" style={{ backgroundColor: "#5758a6" }} variant="light">
      <Container>
        <Navbar.Brand href="#" style={{ color: "white", pointerEvents: "none", fontWeight: "bold" }}
        >Helpdesk Support Ticket Management</Navbar.Brand>
      </Container>
    </Navbar>
  );
}; 

export default TopNav;

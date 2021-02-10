import React, { useState } from 'react';
import {
  Col,
  Collapse,
  Container,
  Navbar,
  NavbarToggler,
  NavbarBrand,
  Nav,
  NavItem,
  NavLink,
  Row,
} from 'reactstrap';
import { Link } from 'react-router-dom';

const TopNavbar = (props) => {
  const [collapsed, setCollapsed] = useState(true);

  const toggleNavbar = () => setCollapsed(!collapsed);

  return (
    <div>
      <Container>
        <Row>
          <Col>
            <Navbar color="faded" light>
              <NavbarBrand tag={Link} to="/">Cerebrum</NavbarBrand>
              <NavbarToggler onClick={toggleNavbar} className="mr-2" />
              <Collapse isOpen={!collapsed} navbar>
                <Nav navbar>
                  <NavItem>
                    <NavLink tag={Link} to="/" activeClassName="active">
                      Home
                    </NavLink>
                  </NavItem>
                  <NavItem>
                    <NavLink
                      tag={Link}
                      to="/login"
                      activeClassName="active"
                    >
                      Login
                    </NavLink>
                  </NavItem>
                </Nav>
              </Collapse>
            </Navbar>
          </Col>
        </Row>
      </Container>
    </div>
  );
};

export default TopNavbar;

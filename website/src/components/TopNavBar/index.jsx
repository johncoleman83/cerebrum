import React, { useState } from 'react';
import {
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
import PropTypes from 'prop-types';

const TopNavBar = ({ activeLink }) => {
  const [collapsed, setCollapsed] = useState(true);

  const toggleNavbar = () => setCollapsed(!collapsed);

  const activeClass = (link) => {
    return link === activeLink ? 'active' : '';
  };

  return (
    <div>
      <Container>
        <Row>
          <Navbar color="faded" light>
            <NavbarBrand tag={Link} to="/">Cerebrum</NavbarBrand>
            <NavbarToggler onClick={toggleNavbar} className="mr-2" />
            <Collapse isOpen={!collapsed} navbar>
              <Nav navbar>
                <NavItem className={activeClass('/')}>
                  <NavLink tag={Link} to="/">
                    Home
                  </NavLink>
                </NavItem>
                <NavItem className={activeClass('/profile')}>
                  <NavLink tag={Link} to="/profile">
                    Profile
                  </NavLink>
                </NavItem>
                <NavItem className={activeClass('/login')}>
                  <NavLink tag={Link} to="/login">
                    Login
                  </NavLink>
                </NavItem>
              </Nav>
            </Collapse>
          </Navbar>
        </Row>
      </Container>
    </div>
  );
};

TopNavBar.propTypes = {
  activeLink: PropTypes.string.isRequired,
};
export default TopNavBar;

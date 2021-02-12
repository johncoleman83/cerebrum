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
    link === activeLink ? 'active' : '';
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
                <NavItem>
                  <NavLink tag={Link} to="/" className={activeClass('/')}>
                    Home
                  </NavLink>
                </NavItem>
                <NavItem>
                  <NavLink
                    tag={Link}
                    to="/profile"
                    className={activeClass('/profile')}
                  >
                    Profile
                  </NavLink>
                  <NavLink
                    tag={Link}
                    to="/login"
                    className={activeClass('/login')}
                  >
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

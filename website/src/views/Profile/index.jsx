import React, { Component } from 'react';
import { connect } from 'react-redux';
import {
  Card,
  CardTitle,
  Container,
} from 'reactstrap';
import { Helmet } from 'react-helmet';
import PropTypes from 'prop-types';
import {
  user,
  isUserValid,
} from 'src/features/current-user/selectors';
import TopNavbar from 'src/components/TopNavBar';

class Profile extends Component {
  constructor(props) {
    super(props);
    this.state = {
      loaded: false,
    };
  }

  async componentDidMount() {
    this.setState({ loaded: true });
  }

  render() {
    console.info('Profile render()');
    console.info('this.props');
    console.info(this.props);

    return (
      <React.Fragment>
        <TopNavbar activeLink='/profile'/>

        <Helmet>
          <title>Home</title>
        </Helmet>

        <Container className="pb-4 h-100 d-flex flex-column">

          <h1 className="h1 font-weight-normal">Logged In!</h1>

          {this.state.loaded && (
            <Card body className="flex-grow-0 mt-1">
              <CardTitle tag="h2">Hello World!</CardTitle>

              <p>Congratulations on logging in.</p>

              <CardTitle tag="h2" className="mt-2">Your Data</CardTitle>

              information gathered from API request to /me:
              <ul style={{ listStyle: 'none' }}>
                <li>{`User id: ${this.props.currentUser.user.id}`}</li>
                <li>
                  {`User username: ${this.props.currentUser.user.username}`}
                </li>
                <li>{`User Email: ${this.props.currentUser.user.email}`}</li>
              </ul>
            </Card>
          )}
        </Container>
      </React.Fragment>
    );
  }
}

Profile.propTypes = {
  currentUser: PropTypes.shape({
    user: PropTypes.oneOfType([
      PropTypes.shape({
        id: PropTypes.number.isRequired,
        accountId: PropTypes.number.isRequired,
        username: PropTypes.string.isRequired,
        email: PropTypes.string.isRequired,
      }),
      PropTypes.shape({}),
    ]).isRequired,
    isUserValid: PropTypes.bool.isRequired,
    isFetching: PropTypes.bool,
    error: PropTypes.oneOfType([PropTypes.object]),
  }).isRequired,
};

const mapStateToProps = (state) => ({
  currentUser: {
    user: user(state),
    isUserValid: isUserValid(state),
    isFetching: state.currentUser.isFetching,
    error: state.currentUser.error,
  },
});

export default connect(mapStateToProps)(Profile);

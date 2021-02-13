import React, { Component } from 'react';
import { connect } from 'react-redux';
import { fetchMeAction } from 'src/features/current-user/actions';
import { BrowserRouter, Switch, Route } from 'react-router-dom';
import {
  user,
  isUserValid,
} from 'src/features/current-user/selectors';
import {
  isAuthValid,
} from 'src/features/authentication/selectors';
import {
  redirectToLogin,
} from 'src/features/authentication/redirect/util';
import Login from 'src/views/Login';
import Home from 'src/views/Home';
import Profile from 'src/views/Profile';
import PropTypes from 'prop-types';
import config from 'src/config/app';
import { isAuthenticated, isAuthError } from '../../util/general';

const getBaseName = () => config.URL_STRING;

const shouldLogout = (store) => {
  return (
    isAuthError(store.currentUser.error) ||
    !isAuthenticated(store)
  );
};

class App extends Component {
  constructor(props) {
    super(props);

    this.state = {
      loaded: false,
    };
  }

  async componentDidMount() {
    await this.props.fetchMeAction(this.props.authentication.authToken);

    if (shouldLogout(this.props) && window.location.pathname != '/login') {
      redirectToLogin();
      return;
    }
    this.setState({ loaded: true });
  }

  render() {
    console.info('calling App BrowserRouter this.props');
    console.info(this.props);
    return (
      <BrowserRouter basename={getBaseName()}>
        {this.state.loaded && (
          <Switch>
            <Route
              path="/login"
              component={Login} />
            <Route
              path='/profile'
              component={Profile} />
            <Route
              exact path='/'
              component={Home} />
          </Switch>
        )}
      </BrowserRouter>
    );
  }
}

App.propTypes = {
  fetchMeAction: PropTypes.func.isRequired,
  authentication: PropTypes.shape({
    authToken: PropTypes.string,
    isAuthValid: PropTypes.bool.isRequired,
    isFetching: PropTypes.bool,
    error: PropTypes.oneOfType([PropTypes.object]),
  }),
  currentUser: PropTypes.shape({
    user: PropTypes.oneOfType([
      PropTypes.shape({
        id: PropTypes.number.isRequired,
        accountId: PropTypes.number.isRequired,
        username: PropTypes.string.isRequired,
        email: PropTypes.string.isRequired,
      }),
      PropTypes.shape({}),
    ]),
    isUserValid: PropTypes.bool.isRequired,
    isFetching: PropTypes.bool,
    error: PropTypes.oneOfType([PropTypes.object]),
  }).isRequired,
};

const mapStateToProps = (state) => ({
  authentication: {
    authToken: state.authentication.authToken,
    isAuthValid: isAuthValid(state),
    error: state.authentication.error,
    isFetching: state.authentication.isFetching,
  },
  currentUser: {
    user: user(state),
    isUserValid: isUserValid(state),
    isFetching: state.currentUser.isFetching,
    error: state.currentUser.error,
  },
});
const mapActionsToProps = { fetchMeAction };
export default connect(mapStateToProps, mapActionsToProps)(App);

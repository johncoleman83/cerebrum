import React, { Component } from 'react';
import { connect } from 'react-redux';
import { BrowserRouter, Switch, Route } from 'react-router-dom';
import { fetchMeAction } from 'src/features/current-user/actions';
import {
  user,
  isUserValid,
} from 'src/features/current-user/selectors';
import {
  redirectToLogin,
} from 'src/features/authorization/redirect/util';
import PropTypes from 'prop-types';
import config from 'src/config/app';
import LogInPage from 'src/views/LogInPage';
import LoggedInPage from 'src/views/LoggedInPage';

const getBaseName = () => config.URL_STRING;

const isAuthError = (error) => error?.response?.status === 401 ?? false;

class App extends Component {
  constructor(props) {
    super(props);
    this.state = {
      loaded: false,
    };
  }

  async componentDidMount() {
    await this.props.fetchMeAction(this.props.authorization.authToken);

    // fetch error is 401 Unauthorized
    if (
      (
        this.props.currentUser.error &&
        isAuthError(this.props.currentUser.error)
      ) || !this.props.currentUser.isUserValid
    ) {
      if (window.location.pathname != '/login') {
        redirectToLogin();
        return;
      }
    }
    this.setState({ loaded: true });
  }

  render() {
    return (
      <BrowserRouter basename={getBaseName()}>
        {this.state.loaded && (
          <Switch>
            <Route exact path="/login" component={LogInPage} />
            <Route exact path="/" component={LoggedInPage} />
          </Switch>
        )}
      </BrowserRouter>
    );
  }
}

App.defaultProps = {
  error: null,
};

App.propTypes = {
  fetchMeAction: PropTypes.func.isRequired,
  authorization: PropTypes.shape({
    authToken: PropTypes.string,
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
  authorization: {
    authToken: state.authorization.authToken,
    error: state.authorization.error,
    isFetching: state.authorization.isFetching,
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

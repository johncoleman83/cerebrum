import React, { Component } from 'react';
import { connect } from 'react-redux';
import {
  Container,
  Breadcrumb,
  BreadcrumbItem,
  Button,
  Form,
  FormGroup,
  Label,
  Input,
} from 'reactstrap';
import { Redirect } from 'react-router';
import { Helmet } from 'react-helmet';
import PropTypes from 'prop-types';
import { loginAction, logoutAction } from 'src/features/authorization/actions';
import { fetchMeAction } from 'src/features/current-user/actions';
import {
  user,
  isUserValid,
} from 'src/features/current-user/selectors';

class LoginPage extends Component {
  constructor(props) {
    super(props);

    // this only logs out on the frontend
    // it should be more secure by actually
    // logging out in the backend
    this.props.logoutAction();

    this.handleSubmit = this.handleSubmit.bind(this);
    this.handleChange = this.handleChange.bind(this);

    this.state = {
      ...this.state,
      username: '',
      password: '',
      loaded: false,
    };
  }

  handleChange(e) {
    const target = e.target;
    this.setState({ [target.name]: target.value });
  }

  async handleSubmit(e) {
    e.preventDefault();
    await this.props.loginAction(this.state.username, this.state.password);
    const authToken = this.props.authorization.authToken;
    if (authToken.length > 0) {
      await this.props.fetchMeAction(authToken);
    }
  }

  async componentDidMount() {
    this.setState({ loaded: true });
  }

  render() {
    return (
      <React.Fragment>
        <Helmet>
          <title>Login Page</title>
        </Helmet>

        <Container className="pb-4 h-100 d-flex flex-column">
          <Breadcrumb>
            <BreadcrumbItem active>Login Page</BreadcrumbItem>
          </Breadcrumb>

          <h1 className="h1 font-weight-normal">Login</h1>
          {
            this.props.authorization.error && (
              <p>error =
                {
                  JSON.stringify(this.props.authorization.error.response.data)
                }
              </p>
            )
          }
          {
            this.state.loaded &&
            this.props.authorization.authToken.length > 0 &&
            this.props.currentUser.isUserValid &&
              <Redirect to="/" />
          }
          {
            this.state.loaded && (
              <Form inline>
                <FormGroup>
                  <Label for="username" hidden>username</Label>
                  <Input
                    type="text"
                    name="username"
                    id="username"
                    placeholder="username"
                    value={this.state.username}
                    onChange={(e) => this.handleChange(e)}
                  />
                </FormGroup>
                <FormGroup>
                  <Label
                    for="password"
                    hidden>Password
                  </Label>
                  <Input
                    type="password"
                    name="password"
                    id="password"
                    placeholder="Password"
                    value={this.state.password}
                    onChange={(e) => this.handleChange(e)}
                  />
                </FormGroup>
                <Button onClick={this.handleSubmit}>Submit</Button>
              </Form>
            )
          }
        </Container>
      </React.Fragment>
    );
  }
}

LoginPage.defaultProps = {
  error: null,
};

LoginPage.propTypes = {
  loginAction: PropTypes.func.isRequired,
  fetchMeAction: PropTypes.func.isRequired,
  logoutAction: PropTypes.func.isRequired,
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
    ]).isRequired,
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

const mapActionsToProps = { loginAction, fetchMeAction, logoutAction };
export default connect(mapStateToProps, mapActionsToProps)(LoginPage);

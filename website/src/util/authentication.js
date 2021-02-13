import React from 'react';
import { Redirect } from 'react-router-dom';

const isAuthError = (error) => {
  (error && error?.response?.status === 401) ?? false;
};

function isAuthenticated(props) {
  return (
    props.authentication?.isAuthValid &&
    props.currentUser?.isUserValid
  );
}

function shouldLogout(props) {
  return (
    isAuthError(props.currentUser.error) ||
    !isAuthenticated(props)
  );
};

function protectedComponent(Component, props) {
  return (
    isAuthenticated(props) ?
    <Component {...props} /> : <Redirect to='/login' />
  );
}

export {
  shouldLogout,
  protectedComponent,
};

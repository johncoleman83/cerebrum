import 'bootstrap/dist/css/bootstrap.min.css';
import 'react-app-polyfill/ie11';
import 'react-app-polyfill/stable';

import React from 'react';
import ReactDOM from 'react-dom';
import { Provider } from 'react-redux';
import axios from 'axios';
import cssVars from 'css-vars-ponyfill';
import createReduxStore from 'src/app/store';
import App from 'src/views/App';
import { setAxiosDefaults } from 'src/util/axios';
import * as serviceWorker from './serviceWorker';
setAxiosDefaults(axios);

// make css vars work in older browsers
cssVars({
  // set this to false to see the effects in Chrome
  onlyLegacy: true,
});

const store = createReduxStore();
ReactDOM.render(
    <React.StrictMode>
      <Provider store={store}>
        <App />
      </Provider>
    </React.StrictMode>,
    document.getElementById('root'),
);

// If you want your app to work offline and load faster, you can change
// unregister() to register() below. Note this comes with some pitfalls.
// Learn more about service workers: https://bit.ly/CRA-PWA
serviceWorker.register();

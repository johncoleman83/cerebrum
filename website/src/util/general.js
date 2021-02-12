import isPlainObject from 'lodash/isPlainObject';
import toPairs from 'lodash/toPairs';
import camelCase from 'lodash/camelCase';


const capitalize = (word) =>
  `${word.slice(0, 1).toUpperCase()}${word.slice(1).toLowerCase()}`;

const camelize = (text, separator = '_') => {
  const words = text.split(separator);
  return [words[0],
    words.slice(1).map((word) => capitalize(word)).join('')].join('');
};

export const isAuthError = (error) => {
  (error && error?.response?.status === 401) ?? false;
};


export function isAuthenticated(state) {
  return (
    state.authentication.isAuthValid &&
    state.currentUser.isUserValid
  );
}

export function camelizeKeys(object, keysToSkip = []) {
  return toPairs(object)
      .filter((pair) => !keysToSkip.includes(pair[0]))
      .reduce((acc, [key, value]) => {
        const newKey = camelCase(key);
        if (isPlainObject(value)) {
          return {
            ...acc,
            [newKey]: camelizeKeys(value),
          };
        }

        return {
          ...acc,
          [newKey]: value,
        };
      }, {});
}

// sets a defined list of properties from a source object to a target object
export const setNamedProperties = (
    targetObject, sourceObject, propertiesToSet, camelizePropertyNames = true,
) => {
  if (typeof targetObject !== 'object' || typeof sourceObject !== 'object') {
    return;
  }
  (propertiesToSet || []).forEach((propertyName) => {
    const name = camelizePropertyNames ? camelize(propertyName) : propertyName;
    targetObject[name] = sourceObject[propertyName];
  });
};

export const redirectTo = (href) => {
  window.location.assign(href);
};

export const reduceArrayToObject = (accumulator, element) => {
  if (!element) {
    throw new Error('Array element cannot be null.');
  }
  if (element && element.id == null) {
    throw new Error('Array elements must have an id property.');
  }

  accumulator[element.id] = element;
  return accumulator;
};

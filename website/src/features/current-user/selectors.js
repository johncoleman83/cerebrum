// General access
export const user = (state) => state.currentUser.user;

// Roles
export const userRole = (state) => user(state).role.name;

// Other
export const isUserValid = (state) => Boolean(user(state)?.id);

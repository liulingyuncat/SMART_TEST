import { configureStore } from '@reduxjs/toolkit';
import authReducer from './authSlice';
import projectReducer from './projectSlice';

const store = configureStore({
  reducer: {
    auth: authReducer,
    project: projectReducer,
  },
});

export default store;

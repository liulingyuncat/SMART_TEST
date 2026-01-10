import { createSlice } from '@reduxjs/toolkit';

const initialState = {
  currentProjectId: null,
  projects: [],
};

const projectSlice = createSlice({
  name: 'project',
  initialState,
  reducers: {
    setCurrentProject: (state, action) => {
      state.currentProjectId = action.payload;
    },
    setProjects: (state, action) => {
      state.projects = action.payload;
    },
    clearCurrentProject: (state) => {
      state.currentProjectId = null;
    },
    clearProjects: (state) => {
      state.projects = [];
      // 注意：不清除currentProjectId，以便在重新登录时可以恢复
    },
  },
});

export const { setCurrentProject, setProjects, clearCurrentProject, clearProjects } = projectSlice.actions;
export default projectSlice.reducer;
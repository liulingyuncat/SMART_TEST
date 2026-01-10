// Mock modules before imports
jest.mock('../../api/project');
jest.mock('../../i18n', () => ({}));

// Mock react-i18next
jest.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (key) => key,
  }),
}));

// Mock react-redux
jest.mock('react-redux', () => ({
  useSelector: jest.fn(),
}));

import React from 'react';
import { render, screen } from '@testing-library/react';
import { useSelector } from 'react-redux';
import ProjectCard from './ProjectCard';

const mockProject = {
  id: 1,
  name: 'Test Project',
  description: 'This is a test project description',
  created_at: '2025-01-01T00:00:00Z'
};

describe('ProjectCard Component', () => {
  beforeEach(() => {
    // Mock useSelector to return a default user
    useSelector.mockImplementation((selector) => {
      const mockState = {
        auth: {
          user: {
            id: 1,
            username: 'testuser',
            role: 'project_member',
          },
        },
      };
      return selector(mockState);
    });
  });

  afterEach(() => {
    jest.clearAllMocks();
  });

  test('renders project name', () => {
    render(<ProjectCard project={mockProject} />);

    expect(screen.getByText('Test Project')).toBeInTheDocument();
  });

  test('renders project description', () => {
    render(<ProjectCard project={mockProject} />);

    expect(screen.getByText('This is a test project description')).toBeInTheDocument();
  });

  test('renders active tag', () => {
    render(<ProjectCard project={mockProject} />);

    expect(screen.getByText('project.active')).toBeInTheDocument();
  });

  test('renders formatted creation date', () => {
    render(<ProjectCard project={mockProject} />);

    expect(screen.getByText('2025-01-01')).toBeInTheDocument();
  });

  test('shows no description when description is empty', () => {
    const projectWithoutDesc = { ...mockProject, description: '' };

    render(<ProjectCard project={projectWithoutDesc} />);

    expect(screen.getByText('project.noDescription')).toBeInTheDocument();
  });

  test('applies correct CSS classes', () => {
    const { container } = render(<ProjectCard project={mockProject} />);

    expect(container.firstChild).toHaveClass('project-card');
    expect(container.querySelector('.project-card-header')).toBeInTheDocument();
    expect(container.querySelector('.project-name')).toBeInTheDocument();
    expect(container.querySelector('.project-card-content')).toBeInTheDocument();
    expect(container.querySelector('.project-description')).toBeInTheDocument();
    expect(container.querySelector('.project-meta')).toBeInTheDocument();
    expect(container.querySelector('.project-date')).toBeInTheDocument();
  });

  test('has hoverable card', () => {
    const { container } = render(<ProjectCard project={mockProject} />);

    const card = container.querySelector('.ant-card-hoverable');
    expect(card).toBeInTheDocument();
  });

  test('truncates long project names with ellipsis', () => {
    const longNameProject = {
      ...mockProject,
      name: 'This is a very long project name that should be truncated with ellipsis'
    };

    render(<ProjectCard project={longNameProject} />);

    const title = screen.getByRole('heading', { level: 4 });
    expect(title).toHaveClass('ant-typography-ellipsis');
  });

  test('truncates long descriptions', () => {
    const longDescProject = {
      ...mockProject,
      description: 'This is a very long description that should be truncated after three rows of text content in the card component layout and display system.'
    };

    render(<ProjectCard project={longDescProject} />);

    const paragraph = screen.getByText(/This is a very long description/);
    expect(paragraph).toHaveClass('ant-typography-ellipsis');
  });

  test('displays calendar icon', () => {
    render(<ProjectCard project={mockProject} />);

    const calendarIcon = document.querySelector('.anticon-calendar');
    expect(calendarIcon).toBeInTheDocument();
  });

  test('validates required props', () => {
    const consoleSpy = jest.spyOn(console, 'error').mockImplementation(() => {});

    // Missing project prop should trigger PropTypes warning
    expect(() => {
      render(<ProjectCard />);
    }).toThrow();

    consoleSpy.mockRestore();
  });

  test('validates project prop structure', () => {
    const invalidProject = { name: 'Test' }; // Missing required fields

    expect(() => {
      render(<ProjectCard project={invalidProject} />);
    }).not.toThrow(); // PropTypes warnings don't throw in production
  });
});
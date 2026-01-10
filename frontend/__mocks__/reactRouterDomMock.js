const React = require('react');

const reactRouterDomMock = {
  BrowserRouter: ({ children }) => React.createElement('div', null, children),
  Routes: ({ children }) => React.createElement('div', null, children),
  Route: () => React.createElement('div'),
  Link: ({ children, to, ...props }) => React.createElement('a', { href: to, ...props }, children),
  useNavigate: () => jest.fn(),
  useParams: () => ({}),
  useLocation: () => ({ pathname: '/', search: '', hash: '', state: null }),
  Navigate: () => React.createElement('div'),
};

module.exports = reactRouterDomMock;
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import App from './App';

describe('App Component', () => {
  test('renders without crashing', () => {
    render(<App />);
    expect(screen.getByText('Vite + React')).toBeInTheDocument();
  });

  test('displays initial count of 0', () => {
    render(<App />);
    expect(screen.getByText('count is 0')).toBeInTheDocument();
  });

  test('increments count when button is clicked', async () => {
    const user = userEvent.setup();
    render(<App />);

    const button = screen.getByRole('button', { name: /count is 0/i });
    await user.click(button);

    expect(screen.getByText('count is 1')).toBeInTheDocument();
  });

  test('increments count multiple times', async () => {
    const user = userEvent.setup();
    render(<App />);

    const button = screen.getByRole('button', { name: /count is/i });

    await user.click(button);
    await user.click(button);
    await user.click(button);

    expect(screen.getByText('count is 3')).toBeInTheDocument();
  });

  test('renders Vite and React logos with correct links', () => {
    render(<App />);

    const viteLink = screen.getByRole('link', { name: /vite logo/i });
    const reactLink = screen.getByRole('link', { name: /react logo/i });

    expect(viteLink).toHaveAttribute('href', 'https://vite.dev');
    expect(reactLink).toHaveAttribute('href', 'https://react.dev');
  });

  // ✅ FIXED: supports text split by <code> tag
  test('displays HMR instruction text', () => {
    render(<App />);
    const code = screen.getByText('src/App.jsx');
    expect(code).toBeInTheDocument();

    const paragraph = code.closest('p');
    expect(paragraph).toHaveTextContent(/Edit\s+src\/App\.jsx\s+and save to test HMR/);
  });


  test('displays read the docs text', () => {
    render(<App />);
    expect(screen.getByText('Click on the Vite and React logos to learn more')).toBeInTheDocument();
  });

  test('button has correct initial text and updates on click', async () => {
    const user = userEvent.setup();
    render(<App />);

    const button = screen.getByRole('button');
    expect(button).toHaveTextContent('count is 0');

    await user.click(button);
    expect(button).toHaveTextContent('count is 1');
  });

  test('matches snapshot', () => {
    const { container } = render(<App />);
    expect(container.firstChild).toMatchSnapshot();
  });

  test('matches snapshot after state change', async () => {
    const user = userEvent.setup();
    const { container } = render(<App />);

    const button = screen.getByRole('button');
    await user.click(button);

    expect(container.firstChild).toMatchSnapshot();
  });
});

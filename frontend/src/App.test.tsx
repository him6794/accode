import React from 'react';
import { render, screen } from '@testing-library/react';
import App from './App';

test('renders legacy ACcode title', () => {
  render(<App />);
  const title = screen.getByText(/accode線上直譯器/i);
  expect(title).toBeInTheDocument();
});

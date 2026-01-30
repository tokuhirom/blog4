import { render } from 'preact';
import { App } from './app.jsx';

const mountEl = document.getElementById('login-app');
render(<App />, mountEl);

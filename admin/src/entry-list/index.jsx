import { render } from 'preact';
import { App } from './app.jsx';

const mountEl = document.getElementById('entry-list-app');
render(<App />, mountEl);

import { render } from 'preact';
import { App } from './app.jsx';

const dataEl = document.getElementById('entries-data');
const initData = JSON.parse(dataEl.textContent);

const mountEl = document.getElementById('entry-list-app');
render(<App initData={initData} />, mountEl);

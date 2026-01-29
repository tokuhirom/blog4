import { render } from 'preact';
import { App } from './app.jsx';

const dataEl = document.getElementById('entry-edit-data');
const initData = JSON.parse(dataEl.textContent);

const mountEl = document.getElementById('entry-edit-app');
render(<App initData={initData} />, mountEl);

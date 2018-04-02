import React from 'react';
import ReactDOM from 'react-dom';
import './index.css';
import MareyDiagram from './MareyDiagram';
import registerServiceWorker from './registerServiceWorker';

ReactDOM.render(<MareyDiagram />, document.getElementById('root'));
registerServiceWorker();

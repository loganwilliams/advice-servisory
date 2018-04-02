import React from 'react';
import ReactDOM from 'react-dom';
import './index.css';
import Diagrams from './Diagrams';
import registerServiceWorker from './registerServiceWorker';

ReactDOM.render(<Diagrams />, document.getElementById('root'));
registerServiceWorker();

import React, { Component } from 'react';
import { injectGlobal } from 'styled-components';
import AppPage from './components/layout';

import MaterialIconsEot from './assets/fonts/MaterialIcons-Regular.eot';
import MaterialIconsTtf from './assets/fonts/MaterialIcons-Regular.ttf';
import MaterialIconsWoff from './assets/fonts/MaterialIcons-Regular.woff';
import MaterialIconsWoff2 from './assets/fonts/MaterialIcons-Regular.woff2';
import MaterialIconsSvg from './assets/fonts/MaterialIcons-Regular.svg';

injectGlobal`
  @font-face {
  font-family: 'Material Icons';
  font-style: normal;
  font-weight: 400;
  src: url(${MaterialIconsEot}); /* For IE6-8 */
  src: local('Material Icons'),
       local('MaterialIcons-Regular'),
       url(${MaterialIconsWoff2}) format('woff2'),
       url(${MaterialIconsWoff}) format('woff'),
       url(${MaterialIconsTtf}) format('truetype');
       url(${MaterialIconsSvg}) format('svg');
}

.material-icons {
  font-family: 'Material Icons';
  font-weight: normal;
  font-style: normal;
  font-size: 24px;  /* Preferred icon size */
  display: inline-block;
  line-height: 1;
  text-transform: none;
  letter-spacing: normal;
  word-wrap: normal;
  white-space: nowrap;
  direction: ltr;

  /* Support for all WebKit browsers. */
  -webkit-font-smoothing: antialiased;
  /* Support for Safari and Chrome. */
  text-rendering: optimizeLegibility;

  /* Support for Firefox. */
  -moz-osx-font-smoothing: grayscale;

  /* Support for IE. */
  font-feature-settings: 'liga';
}
`;


class App extends Component {
  render() {
    return (
      <div className="App">
        <AppPage title="Artists" />
      </div>
    );
  }
}

export default App;

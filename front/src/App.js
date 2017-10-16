import React, { Component } from 'react';
import { injectGlobal, ThemeProvider } from 'styled-components';
import AppPage from './components/layout';

import MaterialIconsEot from './assets/fonts/MaterialIcons-Regular.eot';
import MaterialIconsTtf from './assets/fonts/MaterialIcons-Regular.ttf';
import MaterialIconsWoff from './assets/fonts/MaterialIcons-Regular.woff';
import MaterialIconsWoff2 from './assets/fonts/MaterialIcons-Regular.woff2';
import MaterialIconsSvg from './assets/fonts/MaterialIcons-Regular.svg';

import themeDefault from './themes/light';

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
`;

class App extends Component {
  render() {
    return (
      <ThemeProvider theme={themeDefault}>
        <div className="App">
          <AppPage title="Artists" />
        </div>
      </ThemeProvider>
    );
  }
}

export default App;

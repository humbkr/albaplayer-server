import React, { Component } from 'react';
import { ApolloProvider } from 'react-apollo';
import { injectGlobal, ThemeProvider } from 'styled-components';
import apolloClient from './graphql/apollo';
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
  constructor(props) {
    super(props);
    this.state = {
      isSidebarOpen: false,
    };

    this.handleToggleSidebar = this.handleToggleSidebar.bind(this);
  }

  handleToggleSidebar() {
    this.setState({
      isSidebarOpen: !this.state.isSidebarOpen,
    });
  }

  render() {
    return (
      <ThemeProvider theme={themeDefault}>
        <ApolloProvider client={apolloClient}>
          <div className="App">
            <AppPage title="Artists" />
          </div>
        </ApolloProvider>
      </ThemeProvider>
    );
  }
}

export default App;

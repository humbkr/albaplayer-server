import React, { Component } from 'react';

import { AppPageHeader } from './layout';

const homeViewTitle = 'Home';

class HomeView extends Component {
  constructor(props){
    super(props);
    document.title = homeViewTitle;
  }

  render() {
    return (
      <div>
        <AppPageHeader title={homeViewTitle} />
        <h1>Welcome to Alba Player</h1>
      </div>
    );
  }
}

export default HomeView;

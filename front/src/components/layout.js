import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { BrowserRouter as Router, Route, Link } from 'react-router-dom';

import { Icon, IconButton, Title, Header } from './commons/common';
import SidebarContainer from './sidebar';
import ArtistListView from "./artistListView";
import AlbumListView from "./albumListView";
import HomeView from "./homeView";

class AppPageHeader extends Component {
  render() {
    return (
      <Header>
        <Title>{this.props.title}</Title>
      </Header>
    );
  }
}
AppPageHeader.propTypes = {
  title: PropTypes.string,
};

class AppPage extends Component {
  render() {
    return (
      <Router>
        <div>
          <SidebarContainer />
          <Route exact path="/" component={HomeView}/>
          <Route path="/artists" component={ArtistListView}/>
          <Route path="/albums" component={AlbumListView}/>
        </div>
      </Router>
    );
  }
}
AppPage.propTypes = {
  title: PropTypes.string,
};

export default AppPage;

export {
  AppPageHeader,
};

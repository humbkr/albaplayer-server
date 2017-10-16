import React, { Component } from 'react';
import PropTypes from 'prop-types';

import { Icon, IconButton, Title, Header } from './commons/common';
import MainMenuContainer from './menu';
import ArtistListView from "./ArtistListView";

class AppPageHeader extends Component {
  render() {
    return (
      <Header>
        <IconButton onClick={this.props.menuButtonAction}><Icon>menu</Icon></IconButton>
        <Title>{this.props.title}</Title>
      </Header>
    );
  }
}
AppPageHeader.propTypes = {
  title: PropTypes.string,
  menuButtonAction: PropTypes.func.isRequired,
};

class AppPageContent extends Component {
  render() {
    return (
      <div className="appPageContent">
        {this.props.children}
      </div>
    );
  }
}

class AppPage extends Component {
  constructor(props) {
    super(props);
    this.state = {
      mainMenuIsOpen: false,
    };

    this.handleToggleMainMenu = this.handleToggleMainMenu.bind(this);
  }

  handleToggleMainMenu() {
    this.setState({
      mainMenuIsOpen: !this.state.mainMenuIsOpen,
    });
  }

  render() {
    return (
      <div>
        <MainMenuContainer isOpen={this.state.mainMenuIsOpen} closeButtonHandler={this.handleToggleMainMenu}/>
        <AppPageHeader title={this.props.title} menuButtonAction={this.handleToggleMainMenu} />
        <AppPageContent>
          <ArtistListView />
        </AppPageContent>
      </div>
    );
  }
}
AppPage.propTypes = {
  title: PropTypes.string,
};

export default AppPage;

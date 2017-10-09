import React, { Component } from 'react';
import styled from 'styled-components';
import PropTypes from 'prop-types';

import { Icon, IconButton } from './commons/common';
import ArtistListView from "./artistListView";

const Title = styled.h1`
  display: inline-block;
  vertical-align: top;
  font-size: 1.2em;
  font-weight: normal;
`;

const PageHeader = styled.header`
  display: inline-block;
  width: 100%;
  height: 50px;
  
  > * {
    position: relative;
    top: 50%;
    transform: translateY(-50%);
  }
  
  > ${Title} {
    margin-left: 7px;
  }
`;

class AppPageHeader extends Component {
  render() {
    return (
      <PageHeader>
        <IconButton><Icon>menu</Icon></IconButton>
        <Title>{this.props.title}</Title>
      </PageHeader>
    );
  }
}
AppPageHeader.propTypes = {
  title: PropTypes.string,
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
  render() {
    return (
      <div>
        <AppPageHeader title={this.props.title} />
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

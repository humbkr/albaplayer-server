import React, { Component } from 'react';
import styled from 'styled-components';
import onClickOutside from 'react-onclickoutside'

import { ListItem, List } from './commons/list';
import {IconButton, Icon, Header, Title } from "./commons/common";

const MainMenuLink = styled.a`
  width: 100%;
  height: ${props => props.theme.itemHeight};
  display: inline-block;
  color: ${props => (props.secondary) ? props.theme.sidebar.textSecondaryColor : props.theme.sidebar.textPrimaryColor};
  
  :hover {
    color: ${props => props.theme.sidebar.textPrimaryColorHover};
  }
  
  > span {
    padding-left: 10px;
  }

  > * {
    display: inline-block;
    vertical-align: top;
    position: relative;
    top: 50%;
    transform: translateY(-50%);
  }
`;

const MainMenuFooter = styled.footer`
  position: absolute;
  bottom: 0;
  padding-left: 10px;
  width: 100%;
  
  :hover {
    background-color: ${props => props.theme.highlight};
  }
  
  > MainMenuLink {
    color: ${props => props.theme.sidebar.textSecondaryColor};
  }
`;

const Sidebar = styled.div`
  position: fixed;
  top: 0;
  left: 0;
  height: 100%;
  width: ${props => props.isOpen ? '250px' : '0'};
  z-index: 1;
  overflow: hidden;
  transition: 0.2s;
  box-shadow: 5px 0 5px -2px rgba(0,0,0,.5);
  background-color: ${props => props.theme.sidebar.background};
  color: ${props => props.theme.sidebar.textPrimaryColor};
  
  button {
    color: ${props => props.theme.sidebar.textPrimaryColor};
  }
`;

class MainMenuContainer extends Component {
  // For onClickOutside package.
  handleClickOutside(event) {
    const isOpen = this.props.isOpen;
    if (isOpen) {
      this.props.closeButtonHandler();
    }
  }

  render() {
    const isOpen = this.props.isOpen;

    return (
      <Sidebar isOpen={isOpen}>
        <Header>
          <IconButton onClick={this.props.closeButtonHandler}><Icon>close</Icon></IconButton>
          <Title>Menu</Title>
        </Header>
        <List>
          <ListItem>
            <MainMenuLink href="#">
              <Icon>person</Icon>
              <span>Artists</span>
            </MainMenuLink>
          </ListItem>
          <ListItem>
            <MainMenuLink href="#">
              <Icon>album</Icon>
              <span>Albums</span>
            </MainMenuLink>
          </ListItem>
          <ListItem>
            <MainMenuLink href="#">
              <Icon>fingerprint</Icon>
              <span>Genres</span>
            </MainMenuLink>
          </ListItem>
          <ListItem>
            <MainMenuLink href="#">
              <Icon>playlist_play</Icon>
              <span>Playlists</span>
            </MainMenuLink>
          </ListItem>
        </List>
        <MainMenuFooter>
          <MainMenuLink secondary href="#">
            <Icon>settings</Icon>
            <span>Settings</span>
          </MainMenuLink>
        </MainMenuFooter>
      </Sidebar>
    );
  }
}

export default onClickOutside(MainMenuContainer);

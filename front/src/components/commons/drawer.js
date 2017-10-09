import React, { Component } from 'react';
import PropTypes from 'prop-types';
import styled, { css } from 'styled-components';
import onClickOutside from 'react-onclickoutside'

import {IconButton, Icon} from "./common";

const DrawerMenuWrapper = styled.div`
  display: flex;
  height: 50px;
`;

const DrawerMenuDecorated = styled.div`
  flex: 1;
  
  > * {
    position: relative;
    top: 50%;
    transform: translateY(-50%);
  }
`;

const DrawerMenu = styled.div`
  position: absolute;
  right: 0;
  width: 50px;
  
  ${props => props.open && css`
    display: flex;
    justify-content: flex-end;
    width: ${props => props.maxWidth ? props.maxWidth : '100%'};
    background-color: #cccccc;
  `}
`;

const DrawerMenuContent = styled.div`
  display: none;
  width: 0;
  
  ${props => props.open && css`
    display: block;
    flex: 1;
  `}
`;

class DrawerMenuDecorator extends Component {
  constructor(props) {
    super(props);
    this.state = {
      open: false,
    };

    this.handleButtonClick = this.handleButtonClick.bind(this);
  }

  handleButtonClick(event) {
    if (this.state.open && this.props.onClose) {
      // We are closing the menu.
      this.props.onClose(event);
    }

    this.setState({
      open: !this.state.open,
    });
  }

  // For onClickOutside package.
  handleClickOutside(event) {
    if (!this.props.persistant) {
      this.setState({
        open: false,
      });
    }
  }

  render() {
    const isOpen = this.state.open;

    return (
      <DrawerMenuWrapper>
        <DrawerMenuDecorated>
          {this.props.children}
        </DrawerMenuDecorated>
        <DrawerMenu
          open={isOpen}
          maxWidth={this.props.widthOpen}
        >
          <DrawerMenuContent open={isOpen}>
            {this.props.content}
          </DrawerMenuContent>
          <IconButton onClick={this.handleButtonClick}>
            {(isOpen) ? <Icon>close</Icon> : <Icon>{this.props.icon}</Icon>}
          </IconButton>
        </DrawerMenu>
      </DrawerMenuWrapper>
    );
  }
}
DrawerMenuDecorator.propTypes = {
  icon: PropTypes.string.isRequired,
  content: PropTypes.element.isRequired,
  widthOpen: PropTypes.string,
  persistant: PropTypes.bool,
  onClose: PropTypes.func,
};

export default onClickOutside(DrawerMenuDecorator);

import React, { Component } from 'react';
import styled from 'styled-components';

const ListItem = styled.li`
  width: 100%;
  height: ${props => props.theme.itemHeight};
  ${props => props.border ? 'border-bottom: 1px solid ' + props.theme.separatorColor : ''};
  padding-left: 15px;
  
  :hover {
    background-color: ${props => props.theme.highlight};
  }
  
  > * {
    position: relative;
    top: 50%;
    transform: translateY(-50%);
  }
`;

const List = styled.ul`
  list-style-type: none;
  ${props => props.border ? 'border-top: 1px solid ' + props.theme.separatorColor : ''};
`;

export {
  List,
  ListItem,
};

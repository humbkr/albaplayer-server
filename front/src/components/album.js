import React, { Component } from 'react';
import styled from 'styled-components';

import {IconButton, Icon} from "./commons/common";
import DrawerMenuDecorator from './commons/drawer';

const AlbumTeaserTitle = styled.h2`
  font-size: 1em;
  font-weight: normal;
`;

const AlbumSubInfo = styled.div`
  color: ${props => props.theme.textSecondaryColor};
  font-size: 0.8em;
  margin-top: 5px;
`;

const AlbumTeaserArtist = styled.span`
  font-style: italic;
`;

class AlbumTeaser extends Component {
  render() {
    const album = this.props.item;

    return (
      <div>
        <AlbumTeaserTitle>{album.title}</AlbumTeaserTitle>
        <AlbumSubInfo>
          <span>{album.year}</span>
        </AlbumSubInfo>
      </div>
    );
  }
}

class AlbumTeaserPlayable extends Component {
  render() {
    const album = this.props.item;
    const drawerContent = (
      <div>
        <IconButton><Icon>play_arrow</Icon></IconButton>
        <IconButton><Icon>playlist_add</Icon></IconButton>
      </div>
    );

    return (
      <DrawerMenuDecorator icon="more_vert" content={drawerContent} widthOpen="150px">
        <AlbumTeaser item={album} />
      </DrawerMenuDecorator>
    );
  }
}

export {
  AlbumTeaser,
  AlbumTeaserPlayable,
};

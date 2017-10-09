import React, { Component } from 'react';
import styled from 'styled-components';

import {IconButton, Icon} from "./commons/common";
import DrawerMenuDecorator from './commons/drawer';

const ArtistTeaserName = styled.h2`
  font-size: 1em;
  font-weight: normal;
`;

class ArtistTeaser extends Component {
  render() {
    const artist = this.props.item;

    return (
      <ArtistTeaserName>{artist.name}</ArtistTeaserName>
    );
  }
}

class ArtistTeaserPlayable extends Component {
  render() {
    const artist = this.props.item;
    const drawerContent = (
      <div>
        <IconButton><Icon>play_arrow</Icon></IconButton>
        <IconButton><Icon>playlist_add</Icon></IconButton>
      </div>
    );

    return (
      <DrawerMenuDecorator icon="more_vert" content={drawerContent} widthOpen="150px">
        <ArtistTeaser item={artist} />
      </DrawerMenuDecorator>
    );
  }
}

export {
  ArtistTeaser,
  ArtistTeaserPlayable,
};

import React, { Component } from 'react';
import { gql, graphql } from 'react-apollo';

import { Message, Loading } from './commons/common';
import LibraryListView from "./library/listView";
import { AlbumTeaserPlayable } from "./album";

class AlbumListView extends Component {
  constructor(props){
    super(props);
    document.title = 'Albums';
  }

  render() {
    if (this.props.data.loading) {
      return <Loading />;
    }
    if (this.props.data.error) {
      return <Message type="error">{this.props.data.error.message}</Message>;
    }

    // Copy objects into a new array so we can reorder them client-side.
    let items = [];
    this.props.data.albums.forEach((value) => {
      items.push(value);
    });

    const orderByOptions = [
      {value: 'title', label: 'title'},
      {value: 'year', label: 'year'},
      {value: 'artist.name', label: 'artist'},
    ];

    return (
      <LibraryListView
        itemDisplay={AlbumTeaserPlayable}
        items={items}
        orderOptions={orderByOptions}
        defaultOrder="title"
        searchProperty="title"
      />
    );
  }
}
const allAlbumsQuery = gql`
  query AllAlbumsQuery {
    albums {
      id
      title
      year
    }
  }
`;

export default graphql(allAlbumsQuery)(AlbumListView);

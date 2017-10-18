import React, { Component } from 'react';
import { gql, graphql } from 'react-apollo';

import { Message, Loading } from './commons/common';
import LibraryListView from "./library/listView";
import { ArtistTeaserPlayable } from "./artist";
import AppPage from "./layout";

class ArtistListView extends Component {
  constructor(props){
    super(props);
    document.title = 'Artists';
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
    this.props.data.artists.forEach((value) => {
      items.push(value);
    });

    const orderByOptions = [
      {value: 'name', label: 'name'},
      {value: 'id', label: 'id'},
    ];

    return (
      <LibraryListView
        itemDisplay={ArtistTeaserPlayable}
        items={items}
        orderOptions={orderByOptions}
        defaultOrder="name"
        searchProperty="name"
      />
    );
  }
}
const allArtistsQuery = gql`
  query AllArtistsQuery {
    artists {
      id
      name
    }
  }
`;

export default graphql(allArtistsQuery)(ArtistListView);

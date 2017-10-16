import React, { Component } from 'react';

import LibraryListView from "./library/listView";
import { ArtistTeaserPlayable } from "./artist";

class ArtistListView extends Component {
  render() {
    const items = getDataFromDataSource();
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
      />
    );
  }
}

/** ==============TEST=================================== */
const dataSource = [
  {id: 1, name: "Tool"},
  {id: 2, name: "Alice in Chains"},
  {id: 3, name: "Prodigy"},
  {id: 4, name: "Queens of the Stone Age"},
  {id: 5, name: "Arctic Monkeys"},
  {id: 6, name: "All Them Witches"},
  {id: 7, name: "Kool & the Gang"},
  {id: 8, name: "Marilyn Manson"},
];

function getDataFromDataSource() {
  return dataSource;
}
/** ===================================================== */

export default ArtistListView;

import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { List, ListItem } from '../commons/list';

class LibraryListContainer extends Component {
  // Sort function with nested objects capabilities.
  sortItems(property, items) {
    property = property.split('.');
    const len = property.length;

    items.sort(function (a, b) {
      let i = 0;
      while( i < len ) {
        a = a[property[i]];
        b = b[property[i]];
        i++;
      }

      return (a.toLowerCase() > b.toLowerCase()) ? 1 : ((b.toLowerCase() > a.toLowerCase()) ? -1 : 0);
    });

    return items;
  };

  render() {
    const sortProperty = this.props.orderBy;
    const searchFilter = this.props.searchFilter.toLowerCase();
    const searchProperty = this.props.searchProperty;

    let items = this.sortItems(sortProperty, this.props.items);
    // let items = this.props.items;

    const itemsList = items.map((item) => {
      // Filter items based on search term.
      if (item[searchProperty].toLowerCase().includes(searchFilter)) {

        // Return a list item using the required specific display.
        const Display = this.props.itemDisplay;
        return (
          <ListItem key={item.id} border>
            <Display item={item} />
          </ListItem>
        );
      }
    });

    return (
      <List border>{itemsList}</List>
    );
  }
}
LibraryListContainer.propTypes = {
  orderBy: PropTypes.string,
  searchFilter: PropTypes.string,
  searchProperty: PropTypes.string,
  items: PropTypes.array.isRequired,
  itemDisplay: PropTypes.func,
};

export default LibraryListContainer;

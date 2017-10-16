import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { List, ListItem } from '../commons/list';

class LibraryListContainer extends Component {
  render() {
    const sortProperty = this.props.orderBy;
    const searchFilter = this.props.searchFilter.toLowerCase();
    let items = this.props.items;

    // Sort items.
    items.sort((a,b) => {
      return (
        a[sortProperty] > b[sortProperty]) ? 1 : ((b[sortProperty] > a[sortProperty]) ? -1 : 0);
    });

    const itemsList = items.map((item) => {
      // Filter items based on search term.
      if (item.name.toLowerCase().includes(searchFilter)) {

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
  items: PropTypes.array.isRequired,
  itemDisplay: PropTypes.func,
};

export default LibraryListContainer;

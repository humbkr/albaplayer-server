import React, { Component } from 'react';
import PropTypes from 'prop-types';
import styled from 'styled-components';

const LibraryListItem = styled.li`
  width: 100%;
  height: 50px;
  border-bottom: 1px solid #cccccc;
  padding-left: 15px;
  
  :hover {
    background-color: #cccccc;
  }
  
  > * {
    position: relative;
    top: 50%;
    transform: translateY(-50%);
  }
`;

const LibraryList = styled.ul`
  list-style-type: none;
  border-top: 1px solid #cccccc;
`;

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
          <LibraryListItem key={item.id}>
            <Display item={item} />
          </LibraryListItem>
        );
      }
    });

    return (
      <LibraryList>{itemsList}</LibraryList>
    );
  }
}
LibraryListContainer.propTypes = {
  orderBy: PropTypes.string,
  searchFilter: PropTypes.string,
  items: PropTypes.array.isRequired,
  itemDisplay: PropTypes.Element,
};

export default LibraryListContainer;

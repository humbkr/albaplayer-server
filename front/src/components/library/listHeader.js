import React, { Component } from 'react';
import PropTypes from 'prop-types';
import styled from 'styled-components';

import { SelectContainer } from '../commons/common';
import DrawerMenuDecorator from '../commons/drawer';

const LibraryListSearch = styled.div`
  height: ${props => props.theme.itemHeight};
  padding: 8px;

  > input {
    height: 100%;
    width: 100%;
    font-size: 1em;
    padding-left: 5px;
  }
`;

const LibraryListHeader = styled.header`
  width: 100%;
  height: ${props => props.theme.itemHeight};
  padding-left: 15px;
`;

class LibraryListHeaderContainer extends Component {
  constructor(props) {
    super(props);

    this.focusInput = this.focusInput.bind(this);
  }

  focusInput() {
    this.drawerContent.refs
  }

  render() {
    const orderBy = this.props.orderBy;
    const searchValue = this.props.searchValue;
    const options = this.props.orderOptions;

    this.drawerContent = (
      <LibraryListSearch>
        <input
          onChange={this.props.handleChangeSearch}
          ref="searchInput"
          value={searchValue}
          type="text"
          placeholder="Search..."
        />
      </LibraryListSearch>
    );

    return (
      <LibraryListHeader>
        <DrawerMenuDecorator
          icon="search"
          content={this.drawerContent}
          onOpen={this.focusInput}
          onClose={this.props.handleCloseSearch}
          persistant
        >
          <SelectContainer
            label="order by:"
            options={options}
            value={orderBy}
            onChangeHandler={this.props.handleChangeOrderBy}
          />
        </DrawerMenuDecorator>
      </LibraryListHeader>
    );
  }
}
LibraryListHeaderContainer.propTypes = {
  orderBy: PropTypes.string,
  orderOptions: PropTypes.array.isRequired,
  searchValue: PropTypes.string,
  handleChangeOrderBy: PropTypes.func.isRequired,
  handleChangeSearch: PropTypes.func.isRequired,
  handleCloseSearch: PropTypes.func,
};

export default LibraryListHeaderContainer;

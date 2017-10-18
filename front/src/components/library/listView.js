import React, { Component } from 'react';
import PropTypes from 'prop-types';

import LibraryListHeaderContainer from './listHeader'
import LibraryListContainer from './list'

class LibraryListView extends Component {
  constructor(props) {
    super(props);

    const defaultOrder = this.props.defaultOrder;
    this.state = {
      orderBy: defaultOrder,
      search: '',
    };

    this.handleChangeOrderBy = this.handleChangeOrderBy.bind(this);
    this.handleChangeSearch = this.handleChangeSearch.bind(this);
    this.handleCloseSearch = this.handleCloseSearch.bind(this);
  }

  handleChangeOrderBy(event) {
    this.setState({
      orderBy: event.target.value,
    });
  }

  handleChangeSearch(event) {
    this.setState({
      search: event.target.value,
    });
  }

  handleCloseSearch() {
    this.setState({
      search: '',
    });
  }

  render() {
    return (
      <div>
        <LibraryListHeaderContainer
          orderBy={this.state.orderBy}
          orderOptions={this.props.orderOptions}
          searchValue={this.state.search}
          handleChangeOrderBy={this.handleChangeOrderBy}
          handleChangeSearch={this.handleChangeSearch}
          handleCloseSearch={this.handleCloseSearch}
        />
        <LibraryListContainer
          items={this.props.items}
          itemDisplay={this.props.itemDisplay}
          orderBy={this.state.orderBy}
          searchFilter={this.state.search}
          searchProperty={this.props.searchProperty}
        />
      </div>
    );
  }
}
LibraryListView.propTypes = {
  items: PropTypes.array.isRequired,
  itemDisplay: PropTypes.func,
  orderOptions: PropTypes.array.isRequired,
  defaultOrder: PropTypes.string,
  searchProperty: PropTypes.string,
};

export default LibraryListView;

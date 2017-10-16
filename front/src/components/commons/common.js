import React, { Component } from 'react';
import PropTypes from 'prop-types';
import styled from 'styled-components';

const Title = styled.h1`
  display: inline-block;
  vertical-align: top;
  font-size: 1.2em;
  font-weight: normal;
`;

const Header = styled.header`
  display: inline-block;
  width: 100%;
  height: ${props => props.theme.itemHeight};
  
  > * {
    position: relative;
    top: 50%;
    transform: translateY(-50%);
  }
  
  > ${Title} {
    margin-left: 7px;
  }
`;

const IconButton = styled.button`
  display: inline-block;
  vertical-align: top;
  width: ${props => props.theme.itemHeight};
  height: ${props => props.theme.itemHeight};
  border: none;
  background-color: transparent;
  padding: 8px;
  
  color: ${props => props.theme.textPrimaryColor};
  
  :hover {
    cursor: pointer;
  }
`;

const Icon = styled.i`
  font-family: 'Material Icons';
  font-weight: normal;
  font-style: normal;
  font-size: 24px;
  display: inline-block;
  line-height: 1;
  text-transform: none;
  letter-spacing: normal;
  word-wrap: normal;
  white-space: nowrap;
  direction: ltr;

  /* Support for all WebKit browsers. */
  -webkit-font-smoothing: antialiased;
  /* Support for Safari and Chrome. */
  text-rendering: optimizeLegibility;

  /* Support for Firefox. */
  -moz-osx-font-smoothing: grayscale;

  /* Support for IE. */
  font-feature-settings: 'liga';
`;

const Select = styled.select`
  border: none;
  background-color: transparent;
  font-weight: bold;
  font-size: 1em;
  text-align-last: center;
  height: 40px;
  
  :hover {
    cursor: pointer;
  }
`;

const SelectWrapper = styled.div`
  display: inline-block;
  flex: 1;
  vertical-align: top;
  width: 100%;
  height: ${props => props.theme.itemHeight};
  
  > label {
    display: inline-block;
    vertical-align: top;
  }
  
  > * {
    position: relative;
    top: 50%;
    transform: translateY(-50%);
  }
`;

class SelectContainer extends Component {
  render() {
    const defaultValue = this.props.value;
    const onChangeHandler = this.props.onChangeHandler;
    const optionsRaw = this.props.options;
    const options = optionsRaw.map((option) => {
      return (<option key={option.value} value={option.value}>{option.label}</option>);
    });

    return (
      <SelectWrapper>
        <label htmlFor="select">order by:</label>
        <Select id="select" value={defaultValue} onChange={onChangeHandler}>
          {options}
        </Select>
      </SelectWrapper>
    );
  };
}
SelectContainer.propTypes = {
  options: PropTypes.arrayOf(PropTypes.shape({
    value: PropTypes.string,
    label: PropTypes.string,
  })).isRequired,
  value: PropTypes.string,
  onChangeHandler: PropTypes.func.isRequired,
};

export {
  Icon,
  IconButton,
  SelectContainer,
  Title,
  Header,
};

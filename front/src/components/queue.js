import React, { Component } from 'react';

/**
 * Use this component for all different displays of quick menu.
 */
class QueueQuickMenu extends Component {
  render() {
    return (
      <div className="queueQuickMenu">
        <button className="addToQueue" type="button">Add to queue</button>
        <button className="playNow" type="button">Play now</button>
        <button className="closeMenu" type="button">Close</button>
      </div>
    );
  }
}

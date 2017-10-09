import React from 'react';
import { graphql, QueryRenderer } from 'react-relay';

import relay from '../relay';
import router from '../router';
import history from '../history';
import AppRenderer from './AppRenderer';

class App extends React.Component {
    state = {
        location: history.location,
        params: {},
        query: null,
        variables: {},
        component: null,
    };

    componentDidMount() {
        // Start watching for changes in the URL (window.location).
        this.unlisten = history.listen(this.resolveRoute);
        this.resolveRoute(history.location);
    }

    componentWillUnmount() {
        this.unlisten();
    }

    resolveRoute = location =>
        // Find the route that matches the provided URL path.
        router
            .resolve({ path: location.pathname })
            .then(route => this.setState({ ...route, location }));

    renderReadyState = ({ error, props, retry }) =>
        <AppRenderer
            error={error}
            data={props}
            retry={retry}
            query={this.state.query}
            location={this.state.location}
            params={this.state.params}
            component={this.state.component}
        />;

    render() {
        return (
            <QueryRenderer
                environment={relay}
                query={this.state.query}
                variables={this.state.variables}
                render={this.renderReadyState}
            />
        );
    }
}

export default App;

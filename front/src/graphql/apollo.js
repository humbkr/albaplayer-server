import { ApolloClient, createNetworkInterface } from 'react-apollo';

// By default, this client will send queries to the
//  `/graphql` endpoint on the same host
//const client = new ApolloClient();

const networkInterface = createNetworkInterface({
  uri: 'http://localhost:8888/graphql',
  opts: {
    headers: {

    }
  },
});

const apolloClient = new ApolloClient({
  networkInterface: networkInterface
});

export default apolloClient;

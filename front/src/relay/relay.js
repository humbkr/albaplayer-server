import { Environment, Network, RecordSource, Store } from 'relay-runtime';

// Data will be cached here.
const store = new Store(new RecordSource());

// Configure communication to the GraphQL server.
const network = Network.create((operation, variables) => {
  return fetch(
    'http://localhost:8888/graphql',
    {
      method: 'POST',
      headers: {
        // Add authentication and other headers here.
        'Accept': 'application/json',
        'content-type': 'application/json',
      },
      // credentials: 'include',
      body: JSON.stringify({
        // GraphQL text from input.
        query: operation.text,
        variables,
      }),
    },
  ).then(response => response.json());
});

const environment = new Environment({
  network: network,
  store: store,
});

export default environment;

import { createClient, cacheExchange, fetchExchange } from 'urql'

const API_URL = import.meta.env.VITE_API_URL || '/graphql'

export const graphqlClient = createClient({
  url: API_URL,
  exchanges: [cacheExchange, fetchExchange],
  fetchOptions: () => {
    const token = localStorage.getItem('access_token')
    return {
      headers: {
        authorization: token ? `Bearer ${token}` : '',
      },
    }
  },
})

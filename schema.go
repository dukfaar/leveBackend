package main

import (
	"github.com/dukfaar/goUtils/relay"
	"github.com/dukfaar/leveBackend/leve"
)

var Schema string = `
		schema {
			query: Query
			mutation: Mutation
		}

		type Query {
			leves(first: Int, last: Int, before: String, after: String): LeveConnection!
			leve(id: ID!): Leve!
		}

		type Mutation {
			importLeves(): Boolean
		}` +
	relay.PageInfoGraphQLString +
	leve.GraphQLType

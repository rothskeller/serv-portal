declare module 'smartystreets-javascript-sdk' {
    interface Metadata {
        latitude?: number
        longitude?: number
    }
    interface QueryResultItem {
        deliveryLine1?: string
        deliveryLine2?: string
        lastLine?: string
        metadata?: Metadata
    }
    type QueryResult = QueryResultItem[]
    export namespace usStreet {
        export class Lookup {
            street: string
        }
    }
    export type USStreetClient = {
        send: (lookup: usStreet.Lookup) => Promise<{ lookups: Array<{ result: QueryResult }> }>
    }
    export namespace core {
        export class SharedCredentials {
            constructor(key: string)
        }
        export const buildClient: {
            usStreet: (credentials: SharedCredentials) => USStreetClient
        }
    }
}
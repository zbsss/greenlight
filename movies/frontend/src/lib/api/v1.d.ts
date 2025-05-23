/**
 * This file was auto-generated by openapi-typescript.
 * Do not make direct changes to the file.
 */

export interface paths {
    "/v1/movies": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /** List all movies */
        get: {
            parameters: {
                query?: never;
                header?: never;
                path?: never;
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description List of movies */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            movies?: components["schemas"]["Movie"][];
                        };
                    };
                };
            };
        };
        put?: never;
        /** Create a new movie */
        post: {
            parameters: {
                query?: never;
                header?: never;
                path?: never;
                cookie?: never;
            };
            requestBody: {
                content: {
                    "application/json": components["schemas"]["CreateMovieRequest"];
                };
            };
            responses: {
                /** @description Movie created successfully */
                201: {
                    headers: {
                        /** @description URL of the newly created movie */
                        Location?: string;
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            movie?: components["schemas"]["Movie"];
                        };
                    };
                };
                /** @description Bad request */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            error?: string;
                        };
                    };
                };
            };
        };
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/movies/{id}": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /** Get a movie by ID */
        get: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    id: number;
                };
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description Movie found */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            movie?: components["schemas"]["Movie"];
                        };
                    };
                };
                /** @description Movie not found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content?: never;
                };
            };
        };
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        /** Update a movie */
        patch: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    id: number;
                };
                cookie?: never;
            };
            requestBody: {
                content: {
                    "application/json": components["schemas"]["UpdateMovieRequest"];
                };
            };
            responses: {
                /** @description Movie updated successfully */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            movie?: components["schemas"]["Movie"];
                        };
                    };
                };
                /** @description Bad request */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            error?: string;
                        };
                    };
                };
                /** @description Movie not found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content?: never;
                };
            };
        };
        trace?: never;
    };
}
export type webhooks = Record<string, never>;
export interface components {
    schemas: {
        Movie: {
            /** Format: int64 */
            id: number;
            /** Format: int32 */
            version: number;
            title: string;
            /** Format: int32 */
            year: number;
            /** @description Runtime in minutes, formatted as "X min" */
            runtime: string;
            genres: string[];
        };
        CreateMovieRequest: {
            title: string;
            /** Format: int32 */
            year: number;
            /** Format: int32 */
            runtimeMin: number;
            genres: string[];
        };
        UpdateMovieRequest: {
            title?: string;
            /** Format: int32 */
            year?: number;
            /** Format: int32 */
            runtimeMin?: number;
            genres?: string[];
        };
    };
    responses: never;
    parameters: never;
    requestBodies: never;
    headers: never;
    pathItems: never;
}
export type $defs = Record<string, never>;
export type operations = Record<string, never>;

import Cookies from 'universal-cookie';

const BACKEND_URL = import.meta.env.VITE_JURY_URL;

/**
 * Represents the standardized structure of a response from the API helpers.
 * @template T The type of the data returned in the response.
 */
export interface FetchResponse<T> {
    /** The HTTP status code of the response. */
    status: number;
    /** An error message if the request failed, otherwise an empty string. */
    error: string;
    /** The parsed JSON data from the response, or null if an error occurred. */
    data: T | null;
}

/**
 * Represents a standard success response with a message.
 */
export interface OkResponse {
    message: string;
}

/**
 * Represents a judge in the system.
 */
export interface Judge {
    /** The unique identifier of the judge. */
    id: string;
    /** The name of the judge. */
    name: string;
    /** The unique login code for the judge. */
    code: string;
}

export async function getRequest<T>(path: string, auth: string): Promise<FetchResponse<T>> {
    try {
        const options: RequestInit = {
            method: 'GET',
            headers: createHeaders(auth, true),
        };
        const response = await fetch(`${BACKEND_URL}${path}`, options);

        try {
            const data = await response.json();
            return { status: response.status, error: data.error ? data.error : '', data };
        } catch (jsonError) {
            console.error('JSON parsing error:', jsonError);
            return { status: response.status, error: 'Error parsing response', data: null };
        }
        // eslint-disable-next-line
    } catch (error: any) {
        console.error(error);
        return { status: 404, error: 'Network connection issue', data: null };
    }
}

export async function postRequest<T>(
    path: string,
    auth: string,
    // eslint-disable-next-line
    body: any,
    form?: boolean // Whether the body is a form data; otherwise it's JSON
): Promise<FetchResponse<T>> {
    try {
        const options: RequestInit = {
            method: 'POST',
            headers: createHeaders(auth, !form),
            body: !body ? null : form ? body : JSON.stringify(body)
        };
        const response = await fetch(`${BACKEND_URL}${path}`, options);

        try {
            const data = await response.json();
            return { status: response.status, error: data.error ? data.error : '', data };
        } catch (jsonError) {
            console.error('JSON parsing error:', jsonError);
            return { status: response.status, error: 'Error parsing response', data: null };
        }
        // eslint-disable-next-line
    } catch (error: any) {
        console.error(error);
        return { status: 404, error: 'Network connection issue', data: null };
    }
}

export async function putRequest<T>(
    path: string,
    auth: string,
    // eslint-disable-next-line
    body: any
): Promise<FetchResponse<T>> {
    try {
        const options: RequestInit = {
            method: 'PUT',
            headers: createHeaders(auth, true),
            body: body ? JSON.stringify(body) : null,
        };
        const response = await fetch(`${BACKEND_URL}${path}`, options);

        try {
            const data = await response.json();
            return { status: response.status, error: data.error ? data.error : '', data };
        } catch (jsonError) {
            console.error('JSON parsing error:', jsonError);
            return { status: response.status, error: 'Error parsing response', data: null };
        }
        // eslint-disable-next-line
    } catch (error: any) {
        console.error(error);
        return { status: 404, error: 'Network connection issue', data: null };
    }
}

export async function deleteRequest(
    path: string,
    auth: string
): Promise<FetchResponse<OkResponse>> {
    try {
        const options: RequestInit = {
            method: 'DELETE',
            headers: createHeaders(auth, true),
        };
        const response = await fetch(`${BACKEND_URL}${path}`, options);

        try {
            const data = await response.json();
            return { status: response.status, error: data.error ? data.error : '', data };
        } catch (jsonError) {
            console.error('JSON parsing error:', jsonError);
            return { status: response.status, error: 'Error parsing response', data: null };
        }
        // eslint-disable-next-line
    } catch (error: any) {
        console.error(error);
        return { status: 404, error: 'Network connection issue', data: null };
    }
}

export function createHeaders(auth: string, json: boolean): Headers {
    const headers = new Headers();
    if (json) {
        headers.append('Content-Type', 'application/json');
    }

    if (auth === 'admin') {
        const cookies = new Cookies();
        const pass = cookies.get('admin-pass');
        const basicAuth = btoa(`admin:${pass}`);
        headers.append('Authorization', `Basic ${basicAuth}`);
    } else if (auth === 'judge') {
        const cookies = new Cookies();
        const token = cookies.get('token');
        headers.append('Authorization', `Bearer ${token}`);
    }

    return headers;
}
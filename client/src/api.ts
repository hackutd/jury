import Cookies from 'universal-cookie';

const BACKEND_URL = import.meta.env.VITE_JURY_URL;

export async function getRequest<T>(path: string, auth: string): Promise<FetchResponse<T>> {
    try {
        const options: RequestInit = {
            method: 'GET',
            headers: createHeaders(auth, true),
        };
        const response = await fetch(`${BACKEND_URL}${path}`, options);
        let data;
        try {
            data = await response.json();
        } catch (jsonError) {
            console.error('JSON parsing error:', jsonError);
            return { status: response.status, error: 'Error parsing response', data: null };
        }
        return { status: response.status, error: data.error ? data.error : '', data };
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
    body: any
): Promise<FetchResponse<T>> {
    try {
        const options: RequestInit = {
            method: 'POST',
            headers: createHeaders(auth, true),
            body: body ? JSON.stringify(body) : null,
        };
        const response = await fetch(`${BACKEND_URL}${path}`, options);
        let data;
        try {
            data = await response.json();
        } catch (jsonError) {
            console.error('JSON parsing error:', jsonError);
            return { status: response.status, error: 'Error parsing response', data: null };
        }
        return { status: response.status, error: data.error ? data.error : '', data };
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
        let data;
        try {
            data = await response.json();
        } catch (jsonError) {
            console.error('JSON parsing error:', jsonError);
            return { status: response.status, error: 'Error parsing response', data: null };
        }
        return { status: response.status, error: data.error ? data.error : '', data };
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
        let data;
        try {
            data = await response.json();
        } catch (jsonError) {
            console.error('JSON parsing error:', jsonError);
            return { status: response.status, error: 'Error parsing response', data: null };
        }
        return { status: response.status, error: data.error ? data.error : '', data };
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

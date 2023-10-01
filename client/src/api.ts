import Cookies from 'universal-cookie';

const BACKEND_URL = process.env.REACT_APP_JURY_URL;

interface DataFetchResponse<T> {
    status: number;
    error: string;
    data: T | null;
}

export async function getRequest<T>(path: string, auth: string): Promise<DataFetchResponse<T>> {
    try {
        const options: RequestInit = {
            method: 'GET',
            headers: createHeaders(auth),
        };
        const response = await fetch(`${BACKEND_URL}${path}`, options);
        const data = await response.json();
        return { status: response.status, error: "", data };
    } catch (error: any) {
        console.error(error);
        return { status: 404, error: error, data: null };
    }
}

export async function postRequest<T>(
    path: string,
    auth: string,
    body: any
): Promise<DataFetchResponse<T>> {
    try {
        console.log(path)
        const options: RequestInit = {
            method: 'POST',
            headers: createHeaders(auth),
            body: body ? JSON.stringify(body) : null,
        };
        const response = await fetch(`${BACKEND_URL}${path}`, options);
        const data = await response.json();
        return { status: response.status, error: "", data };
    } catch (error: any) {
        console.error(error);
        return { status: 404, error: error, data: null };
    }
}

function createHeaders(auth: string): Headers {
    const headers = new Headers();
    headers.append('Content-Type', 'application/json');

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

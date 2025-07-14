---
sidebar_position: 4
title: Calling Backend API
description: How to call the backend API from the frontend
---

# Calling Backend API

To call the backend API, you should use the functions defined in `client/src/api.ts`. The helper methods defined in this function let you make GET, POST, PUT, and DELETE requests to the backend.

## Authentication

To authenticate a request, you will either need an admin password (`auth="admin"`) or a judge token (`auth="judge"`).

On all judge pages gated behind the judge login, you should check if the judge has logged in. The judge `token` cookie will automatically be passed into request as long as the `auth` parameter is set to `judge`. Here is an example from `live.tsx`, where we check for logged in on page load in a useEffect function:

```jsx
    useEffect(() => {
        async function fetchData() {
            // Check to see if the user is logged in
            const loggedInRes = await postRequest<OkResponse>('/judge/auth', 'judge', null);
            if (loggedInRes.data?.ok !== 1 && loggedInRes.status === 401) {
                console.error(`Judge is not logged in!`);
                navigate('/judge/login');
                return;
            }
            if (loggedInRes.status !== 200) {
                errorAlert(loggedInRes);
                return;
            }
            // ... other code
        }

        fetchData();
    }, []);
```

For admin pages gated behind the admin login, you should check if the admin has logged in. The admin `admin-pass` token will automatically passed into the request as long as the `auth` parameter is set to `admin`. The auth check is the same as above, except using `/admin/auth` endpoint and `'auth'` as the auth parameter.

The `createHeaders` function in `client/src/api.ts` shows that admins use Basic Authentication with judges use Bearer Token Authentication with the API:

```js
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
```

## Request Body

When using the POST or PUT requests, you can pass a request body into the function:

```js
export async function postRequest<T>(
    path: string,
    auth: string,
    body: any,
    form?: boolean // Whether the body is a form data; otherwise it's JSON
): Promise<FetchResponse<T>> {
    // ... details hidden
}
```

This body should be a POJO (plain-old JS object) as it will be cast to a JSON string. However, if you are passing FormData in, make sure to set the `form` parameter to true.

## Response Type

All API call helper functions are generic functions which returns a `FetchResponse<T>`. The FetchResponse object looks like the following:

```js
interface FetchResponse<T> {
    status: number;
    error: string;
    data: T | null;
}
```

Generally, you should check the response as following (given that `res` is the return from an API call):

```js
if (res.status !== 200) {
    errorAlert(res);
    return;
}
```

The `errorAlert` function is a wrapper for error logging and alerting. This is standardized across the app and should be used everywhere to report an error from the backend.

To retrieve the data from the response, use the `data` field. Make sure to cast it to the type you define. Here is an example from `client/src/store.tsx`:

```js
const flagsRes = await getRequest<Flag[]>('/admin/flags', 'admin');
if (flagsRes.status !== 200) {
    errorAlert(flagsRes);
    return;
}
set({ flags: flagsRes.data as Flag[] });
```

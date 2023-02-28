import React from 'react';
import ReactDOM from 'react-dom/client';
import './index.css';
import Home from './Home';
import { createBrowserRouter, RouterProvider } from 'react-router-dom';
import JudgeLogin from './judge/login';
import Judge from './judge';
import AdminLogin from './admin/login';
import Admin from './admin';

const root = ReactDOM.createRoot(document.getElementById('root') as HTMLElement);

const router = createBrowserRouter([
    {
        path: '/',
        element: <Home />,
    },
    {
        path: '/judge/login',
        element: <JudgeLogin />,
    },
    {
        path: '/judge',
        element: <Judge />,
    },
    {
        path: '/admin/login',
        element: <AdminLogin />,
    },
    {
        path: '/admin',
        element: <Admin />,
    },
]);

root.render(
    <React.StrictMode>
        <RouterProvider router={router} />
    </React.StrictMode>
);

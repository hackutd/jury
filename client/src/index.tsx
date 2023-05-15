import React from 'react';
import ReactDOM from 'react-dom/client';
import './index.css';
import Home from './Home';
import HomeHub from './components/home/HomeHub';
import { createBrowserRouter, RouterProvider } from 'react-router-dom';
import JudgeLogin from './judge/login';
import Judge from './judge';
import AdminLogin from './admin/login';
import Admin from './admin';
import AddProjects from './admin/AddProjects';
import AddJudges from './admin/AddJudges';

const root = ReactDOM.createRoot(document.getElementById('root') as HTMLElement);

const router = createBrowserRouter([
    {
        path: '/',
        element: process.env.REACT_APP_HUB ? <HomeHub /> : <Home />,
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
    {
        path: '/admin/add-projects',
        element: <AddProjects />,
    },
    {
        path: '/admin/add-judges',
        element: <AddJudges />,
    },
]);

root.render(
    <React.StrictMode>
        <RouterProvider router={router} />
    </React.StrictMode>
);

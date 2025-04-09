import React from 'react';
import ReactDOM from 'react-dom/client';
import { createBrowserRouter, RouterProvider } from 'react-router-dom';
import Home from './pages/Home';
import HomeHub from './components/home/HomeHub';
import JudgeLogin from './pages/judge/login';
import Judge from './pages/judge';
import AdminLogin from './pages/admin/login';
import Admin from './pages/admin';
import AddProjects from './pages/admin/AddProjects';
import AddJudges from './pages/admin/AddJudges';
import JudgeWelcome from './pages/judge/welcome';
import JudgeLive from './pages/judge/live';
import Project from './pages/judge/project';

import AdminSettings from './pages/admin/settings';
import Expo from './pages/Expo';

import './index.css';
import AdminLog from './pages/admin/log';
import QrCode from './pages/admin/QrCode';
import AddSelf from './pages/AddSelf';
import Info from './pages/Info';

const root = ReactDOM.createRoot(document.getElementById('root') as HTMLElement);

const router = createBrowserRouter([
    {
        path: '/',
        element: import.meta.env.VITE_HUB ? <HomeHub /> : <Home />,
    },
    {
        path: '/judge/login',
        element: <JudgeLogin />,
    },
    {
        path: '/judge/welcome',
        element: <JudgeWelcome />,
    },
    {
        path: '/judge/live',
        element: <JudgeLive />,
    },
    {
        path: '/judge',
        element: <Judge />,
    },
    {
        path: '/judge/project/:id',
        element: <Project />,
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
        path: '/expo',
        element: <Expo />,
    },
    {
        path: '/expo/:track',
        element: <Expo />,
    },
    {
        path: '/admin/add-projects',
        element: <AddProjects />,
    },
    {
        path: '/admin/add-judges',
        element: <AddJudges />,
    },
    {
        path: '/admin/settings',
        element: <AdminSettings />,
    },
    {
        path: '/admin/log',
        element: <AdminLog />,
    },
    {
        path: '/admin/qr',
        element: <QrCode />,
    },
    {
        path: '/add-self',
        element: <AddSelf />,
    },
    {
        path: '/add-self/done',
        element: (
            <Info>
                Thanks for adding yourself to the system! You should recieve an email soon with your
                judging login info. You can safely close this tab.
            </Info>
        ),
    },
]);

root.render(
    <React.StrictMode>
        <RouterProvider router={router} />
    </React.StrictMode>
);

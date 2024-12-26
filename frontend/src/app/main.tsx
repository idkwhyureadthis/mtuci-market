import { createRoot } from 'react-dom/client';
import { RouterProvider } from 'react-router-dom';
import { router } from './routes';
import './index.css';
import React from 'react';

createRoot(document.getElementById('root')!).render(
    <RouterProvider router={router} />
);
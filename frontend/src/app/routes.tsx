import { createBrowserRouter } from 'react-router-dom';
import { MainPage } from './pages/MainPage';
import { SignInPage } from './pages/SignInPage';
import { SignUpPage } from './pages/SignUpPage';
import { CreateProductPage } from './pages/CreateProductPage';
import { ModeratorPage } from './pages/ModeratorPage';
import UserProfile from './pages/UserProfile';
import { CreateModerator } from './pages/CreateModerator';

export const router = createBrowserRouter([
  {
    path: '/',
    element: <MainPage />,
  },
  {
    path: '/sign-in',
    element: <SignInPage />,
  },
  {
    path: '/sign-up',
    element: <SignUpPage />,
  },
  {
    path: '/create',
    element: <CreateProductPage />,
  },
  {
    path: '/moderator',
    element: <ModeratorPage />,
  },
  {
    path: '/profile',
    element: <UserProfile />
  },
  {
    path: "create_moderator",
    element: <CreateModerator />
  }
]);
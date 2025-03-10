import { ClientDetails } from 'components/ClientDetails/ClientDetails';
import { ClientsList } from 'components/ClientsList/ClientsList';
import { HomePage } from 'pages/HomePage/HomePage';
import { NotFoundPage } from 'pages/NotFoundPage/NotFoundPage';
import { UsersPage } from 'pages/UsersPage/UsersPage';
import { RoutesType } from 'routing/routesTypes';

export const baseRoutes: RoutesType[] = [
    // {
    //     path: '/',
    //     component: HomePage
    // },
    {
        path: '/devices',
        component: ClientsList
    },
    {
        path: '/devices/:id',
        component: ClientDetails
    },
    {
        path: '/404',
        component: NotFoundPage
    },
    {
        path: '*',
        component: NotFoundPage
    }
];

export const privateRoutes: RoutesType[] = [
    {
        path: '/users',
        component: UsersPage
    },
    {
        path: '*',
        component: NotFoundPage
    }
];

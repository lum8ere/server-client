import { ClientDetails } from 'components/ClientDetails/ClientDetails';
import { HomePage } from 'pages/HomePage/HomePage';
import { NotFoundPage } from 'pages/NotFoundPage/NotFoundPage';
import { RoutesType } from 'routing/routesTypes';

export const baseRoutes: RoutesType[] = [
    {
        path: '/',
        component: HomePage
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

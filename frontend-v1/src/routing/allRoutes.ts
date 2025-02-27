import { ClientsList } from 'components/ClientsList/ClientsList';
import { HomePage } from 'pages/HomePage/HomePage';
import { NotFoundPage } from 'pages/NotFoundPage/NotFoundPage';
import { RoutesType } from 'routing/routesTypes';

export const baseRoutes: RoutesType[] = [
    {
        path: '/',
        component: HomePage
    },
    {
        path: '/client/:id',
        component: ClientsList
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

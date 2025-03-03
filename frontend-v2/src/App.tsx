import { Navigate, Route, Routes } from 'react-router-dom';
import { DefaultLayout } from 'modules/CommonPage/DefaultLayout';
import { allRoutes as routes } from 'routing/routes';

export const App: React.FC = () => {
    return (
        <DefaultLayout>
            <Routes>
                {routes.map((route) => {
                    return (
                        <Route key={route.path} path={route.path} element={<route.component />} />
                    );
                })}

                <Route path="*" element={<Navigate to="/404" />} />
            </Routes>
        </DefaultLayout>
    );
};

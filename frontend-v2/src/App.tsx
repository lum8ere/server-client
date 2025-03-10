import { Navigate, Route, Routes } from 'react-router-dom';
import { DefaultLayout } from 'modules/CommonPage/DefaultLayout';
import { allRoutes as routes } from 'routing/routes';
import { LoginPage } from 'pages/LoginPage/LoginPage';
import { RegisterPage } from 'pages/RegisterPage/RegisterPage';

export const App: React.FC = () => {
    return (
        <DefaultLayout>
            <Routes>
                <Route path="/login" element={<LoginPage />} />
                <Route path="/register" element={<RegisterPage />} />

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

import { Navigate, Route, Routes } from 'react-router-dom';
import { DefaultLayout } from 'modules/CommonPage/DefaultLayout';
import { allRoutes as routes, privateRoutes } from 'routing/routes';
import { PrivateRoute } from 'components/PrivateRoute/PrivateRoute';
import { LoginPage } from 'pages/LoginPage/LoginPage';

export const App: React.FC = () => {
    return (
        <Routes>
            {/* Публичные маршруты */}
            <Route path="/login" element={<LoginPage />} />
            {/* Все маршруты, требующие авторизации, оборачиваются в PrivateRoute */}
            <Route element={<PrivateRoute />}>
                {routes.map((route) => (
                    <Route
                        key={route.path}
                        path={route.path}
                        element={
                            <DefaultLayout>
                                <route.component />
                            </DefaultLayout>
                        }
                    />
                ))}
            </Route>

            {/* Защищенный маршрут, доступный только администратору */}
            <Route element={<PrivateRoute requiredRole="admin" />}>
                {privateRoutes.map((route) => (
                    <Route
                        key={route.path}
                        path={route.path}
                        element={
                            <DefaultLayout>
                                <route.component />
                            </DefaultLayout>
                        }
                    />
                ))}
            </Route>
            {/* Любой несуществующий URL перенаправляем на страницу /login */}
            <Route path="*" element={<Navigate to="/login" replace />} />
        </Routes>
    );
};

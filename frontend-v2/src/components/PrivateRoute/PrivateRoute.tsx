import React from 'react';
import { Navigate, Outlet } from 'react-router-dom';
import { useSelector } from 'react-redux';
import { RootState } from 'store';

const PrivateRoute: React.FC = () => {
    const token = useSelector((state: RootState) => state.auth.token);

    // Если токен отсутствует, перенаправляем на /login
    if (!token) {
        return <Navigate to="/login" replace />;
    }
    // Если токен есть, разрешаем доступ к дочерним маршрутам
    return <Outlet />;
};

export default PrivateRoute;

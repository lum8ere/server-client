import React from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import ClientList from 'components/ClientList';
import ClientDashboard from 'components/ClientDashboard';

export const App: React.FC = () => {
    return (
        <Router>
            <Routes>
                <Route path="/" element={<ClientList />} />
                <Route path="/client/:id" element={<ClientDashboard />} />
            </Routes>
        </Router>
    );
};

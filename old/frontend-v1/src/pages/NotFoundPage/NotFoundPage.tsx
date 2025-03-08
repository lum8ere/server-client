import React from 'react';
import { useNavigate } from 'react-router-dom';
import { Button, Result } from 'antd';

export const NotFoundPage: React.FC = () => {
    const navigate = useNavigate();

    return (
        <Result
            status="404"
            title="404"
            subTitle="К сожалению, страница, которую вы посетили, не существует."
            extra={
                <Button
                    type="primary"
                    style={{ backgroundColor: '#8e3131' }}
                    onClick={() => navigate('/')}
                >
                    Back
                </Button>
            }
        />
    );
};

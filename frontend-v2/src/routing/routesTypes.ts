import { FC, MemoExoticComponent, ReactNode } from 'react';

export interface RoutesType {
    path: string;
    component: FC<any> | MemoExoticComponent<any>;
    protected?: boolean;
}
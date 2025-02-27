module.exports = {
    root: true,
    parser: '@typescript-eslint/parser',
    plugins: ['react', 'react-hooks', 'import', 'prettier', '@typescript-eslint'],
    extends: [
        'plugin:import/errors',
        'plugin:import/warnings',
        'plugin:import/typescript',
        'plugin:@typescript-eslint/recommended',
        'plugin:prettier/recommended'
    ],
    env: {
        browser: true,
        es2021: true
    },
    rules: {
        camelcase: 'off',
        'linebreak-style': 'off', // Пока отключим строгие окончания строк ("unix" / "windows"). Разработка идет на разных плаформах.
        'react/jsx-pascal-case': 'off', // Так как в нашем проекте используется неймин вида Lorem_Ipsum у компонентов (напр. Model_TablePage), данное правило отключаем, чтобы не было ошибок
        'react/react-in-jsx-scope': 'off', // Правило, которое просит/не просит импортировать React в каждый файл (...from React 17 onwards, you no longer need to import the React object to use JSX. This is due to the new JSX Transform introduced in React 17.)
        'react/jsx-filename-extension': ['error', { extensions: ['.tsx', '.jsx'] }],
        'react/prop-types': 'off', // У нас связка React + TS. PropsTypes не актуальны
        'react/require-default-props': 'off', // У нас связка React + TS. PropsTypes не актуальны
        'react-hooks/exhaustive-deps': 'warn', // По умолчанию в проекте 'error',
        'react/no-unused-prop-types': 'off', // Выключаем проверку на неиспользуемые propTypes, так как ими не пользуемся
        // 'prettier/prettier': 'error',
        'import/extensions': [
            'error',
            'ignorePackages',
            {
                js: 'never',
                jsx: 'never',
                ts: 'never',
                tsx: 'never'
            }
        ],
        'import/no-extraneous-dependencies': [
            'error',
            {
                devDependencies: ['**/*.test.ts', '**/*.test.tsx']
            }
        ],
        'no-unused-vars': 'off',
        '@typescript-eslint/no-explicit-any': 'off',
        '@typescript-eslint/no-unused-vars': 'off',
        'import/prefer-default-export': 'off',
        'import/no-default-export': 'off',
        'import/no-unresolved': 'off', // Отключаем ругань на "глобальные импорты"
        'no-underscore-dangle': 'off', // Предупреждения на использование переменных вида _state. Updated: отключено
        'no-console': ['warn', { allow: ['info'] }],
        'no-plusplus': 'off', // Allow ++ operations
        'no-shadow': 'off', // Выключил, иначе ругается на enum'ы
        'max-len': [
            'warn',
            {
                ignoreComments: true,
                ignoreUrls: true,
                ignoreStrings: true,
                ignoreRegExpLiterals: true,
                code: 100
            }
        ], // Настройки максимально допустимой длины строки
        'no-continue': 'off', // Разрешили оператор "continue"
        'jsx-a11y/click-events-have-key-events': 'off', // onClick в статичных тэгах
        'jsx-a11y/no-static-element-interactions': 'off', // onClick в статичных тэгах
        'jsx-a11y/no-noninteractive-element-interactions': 'off', // onClick в статичных тэгах
        'max-classes-per-file': 'off', // Более одного класса в файле
        'prefer-destructuring': 'off', // Отключили требование деструктуризации - т.к. она менее читаемаая
        'import/order': 'error', // Порядок сортировка импортов
        'import/no-cycle': 'off', // иногда ругается необоснованно (разобраться)
        'func-call-spacing': 'off', // conflicts with TypeScript
        'no-spaced-func': 'off', // deprecated
        'brace-style': 'error', // форматирование переносов скобок по-умолчанию
        'arrow-body-style': 'off', // Правило, которое просит/не просит ставить фигурные скобки в стрелочных функциях
        'implicit-arrow-linebreak': 'off', // Форматирование стрелочной функции, связанное с max-len правилом
        'no-multiple-empty-lines': 'off', // Максимальное количество пустых строк в коде
        '@typescript-eslint/ban-types': [
            'warn',
            {
                types: {
                    Function: false,
                    String: {
                        message: 'Use string instead',
                        fixWith: 'string'
                    }
                },
                extendDefaults: true
            }
        ],
        '@typescript-eslint/ban-ts-comment': [
            'error',
            {
                'ts-expect-error': 'allow-with-description',
                'ts-ignore': 'allow-with-description',
                'ts-nocheck': 'allow-with-description',
                'ts-check': 'allow-with-description'
            }
        ],
        'no-nested-ternary': 'off',
        'no-restricted-syntax': [
            'error',
            'WithStatement',
            "BinaryExpression[operator='in']",
            'ForInStatement',
            'DoWhileStatement',
            'EmptyStatement'
        ], // что не используем, от чего отказываемся и что не безопасно использовать
        'no-await-in-loop': 'off',
        'class-methods-use-this': 'off',
        'prettier/prettier': 'error',
        '@typescript-eslint/no-empty-function': ['error', { allow: ['arrowFunctions'] }]
    }
};

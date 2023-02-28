db.judges.insert([
    {
        token: '',
        code: '610699',
        name: 'Michael Zhao',
        email: 'michaelzhao314@gmail.com',
        active: true,
        last_activity: {
            $date: {
                $numberLong: '0',
            },
        },
        read_welcome: true,
        notes: 'catboy | actually a furry',
        next: null,
        prev: null,
        alpha: 10,
        beta: 1,
    },
    {
        token: '',
        code: '610699',
        code: '348646',
        name: 'John Smith',
        email: 'alecbrightman@gmail.com',
        active: true,
        last_activity: {
            $date: {
                $numberLong: '0',
            },
        },
        read_welcome: false,
        notes: 'this is a fake account uwu',
        next: null,
        prev: null,
        alpha: 10,
        beta: 1,
    },
]);

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

db.projects.insert([
    {
        name: 'Arke',
        location: 1,
        description:
            'This is a fancy boat! Arke is a hackathon project that aims to revolutionize the way people interact with augmented reality. The project seeks to develop an AR platform that enables users to create and experience immersive AR environments without the need for any technical expertise.\nArke will provide an easy-to-use drag-and-drop interface that allows users to create and customize 3D AR objects, scenes, and animations. The platform will also provide a library of pre-built objects and environments that users can use as a starting point for their AR creations.\nThe goal of Arke is to democratize AR and make it accessible to everyone, regardless of their technical background. The platform will enable users to create AR experiences for a range of applications, including education, entertainment, marketing, and more.\nThroughout the hackathon, participants will have the opportunity to work on various aspects of the Arke platform, including its user interface, 3D object creation and animation, AR scene management, and more. They will also have access to mentorship and resources to help them bring their ideas to life.\nAt the end of the hackathon, participants will present their AR creations to a panel of judges, who will evaluate them based on creativity, technical complexity, and overall impact. The winning team will receive a prize, as well as recognition for their contribution to the future of AR.',
        try_link: 'https://github.com/rust-lang/rust-analyzer/',
        video_link: null,
        challenge_list: [],
        seen: 0,
        votes: 0,
        mu: 0,
        sigma_sq: 1,
        active: true,
        prioritized: false,
        last_activity: {
            $date: {
                $numberLong: '0',
            },
        },
    },
    {
        name: 'EyeAlert',
        location: 2,
        description:
            "EyeAlert is a hackathon project that focuses on improving workplace safety and productivity by leveraging computer vision technology. The project aims to address the issue of employees suffering from eye strain and fatigue due to prolonged computer usage.\nThe EyeAlert system uses a camera to continuously monitor the user's eye movements and facial expressions. It then analyzes this data to detect signs of eye fatigue and provides timely alerts to the user. These alerts could include reminders to take breaks, suggestions to adjust the lighting, or recommendations to change their computer's display settings.\nThe EyeAlert system could also provide insights to employers, helping them identify which employees are at a higher risk of eye strain and fatigue, and take proactive measures to prevent these issues from occurring. This could include training sessions on proper ergonomic practices, adjusting work schedules to allow for more frequent breaks, or providing eye-care resources to employees.\nOverall, the EyeAlert project has the potential to improve employee well-being, boost productivity, and reduce workplace injuries caused by eye strain and fatigue.",
        try_link: 'https://github.com/adityarathod/eyealert',
        video_link: null,
        challenge_list: ["Children's Health", 'McKesson', 'Lockheed Martin'],
        seen: 0,
        votes: 0,
        mu: 0,
        sigma_sq: 1,
        active: true,
        prioritized: false,
        last_activity: {
            $date: {
                $numberLong: '0',
            },
        },
    },
    {
        name: 'Nub',
        location: 3,
        description:
            'Nub is a hackathon project that aims to provide a simple and intuitive solution for individuals who struggle with organizing their daily tasks and schedules. This project is designed to help users easily prioritize their tasks, manage their time more effectively, and ultimately increase productivity.\nThe Nub platform will have a clean and user-friendly interface, allowing users to quickly input their tasks and schedule them for the appropriate time. The platform will also provide helpful reminders and notifications to ensure that users stay on track and meet their deadlines.\nOne of the key features of Nub is its ability to learn from user behavior and preferences. Over time, the platform will become more personalized, providing tailored suggestions and recommendations to help users optimize their workflow and improve their efficiency.\nIn addition to its core task management functionality, Nub will also offer a range of supplementary features, such as integrations with popular productivity tools and resources, gamification elements to encourage users to stay on track, and a community forum where users can share tips and best practices.',
        try_link: 'https://github.com/google/rustcxx',
        video_link: null,
        challenge_list: ["Children's Health"],
        seen: 0,
        votes: 0,
        mu: 0,
        sigma_sq: 1,
        active: true,
        prioritized: false,
        last_activity: {
            $date: {
                $numberLong: '0',
            },
        },
    },
]);

db.options.insert([
    {
        next_table_num: 4,
    },
]);

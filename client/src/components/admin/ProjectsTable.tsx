import { useEffect, useState } from 'react';

const ProjectsTable = () => {
    const [projects, setProjects] = useState([]);

    // Fetch projects list
    const fetchProjects = async () => {
        const fetchedProjects = await fetch(`${process.env.REACT_APP_JURY_URL}/project/list`, {
            method: 'GET',
            headers: { 'Content-Type': 'application/json' },
            credentials: 'include',
        }).then((data) => data.json());
        // TODO: Add error handling
        setProjects(fetchedProjects);
    };

    // Convert millisecond time to "x secs/mins/hours ago"
    const timeSince = (date: number) => {
        // eslint-disable-next-line eqeqeq
        if (date == 0) {
            return 'never';
        }
        const seconds = Math.floor((new Date().getTime() - date) / 1000);
        if (seconds < 0) {
            return Math.abs(seconds) + ' seconds in the future?!?!';
        }
        let interval = seconds / 31536000;
        if (interval > 1) {
            return Math.floor(interval) + ' years ago';
        }
        interval = seconds / 2592000;
        if (interval > 1) {
            return Math.floor(interval) + ' months ago';
        }
        interval = seconds / 86400;
        if (interval > 1) {
            return Math.floor(interval) + ' days ago';
        }
        interval = seconds / 3600;
        if (interval > 1) {
            return Math.floor(interval) + ' hours ago';
        }
        interval = seconds / 60;
        if (interval > 1) {
            return Math.floor(interval) + ' minutes ago';
        }
        return 'just now';
    };

    useEffect(() => {
        fetchProjects();
    }, []);

    return (
        <div className="w-full px-8">
            <table className="table-auto w-full text-lg">
                <tbody>
                    <tr>
                        <th></th>
                        <th className="text-left py-1">Name</th>
                        <th className="text-center">Table Number</th>
                        <th className="text-center">Mu</th>
                        <th className="text-center">Sigma^2</th>
                        <th className="text-center">Votes</th>
                        <th className="text-center">Seen</th>
                        <th className="text-center">Updated</th>
                        <th className="text-right">Actions</th>
                    </tr>
                    {projects.map((project: any, idx) => (
                        <tr key={idx} className="border-t-2 border-backgroundDark">
                            <td>
                                <input type="checkbox"></input>
                            </td>
                            <td>{project.name}</td>
                            <td className="text-center py-1">Table {project.location}</td>
                            <td className="text-center">{project.mu}</td>
                            <td className="text-center">{project.sigma_sq}</td>
                            <td className="text-center">{project.votes}</td>
                            <td className="text-center">{project.seen}</td>
                            {/* TODO: What the fuck is this; just change the datatype to a long pls */}
                            <td className="text-center">
                                {timeSince(project.last_activity.$date.$numberLong)}
                            </td>
                            <td className="text-right">:</td>
                        </tr>
                    ))}
                </tbody>
            </table>
        </div>
    );
};

export default ProjectsTable;

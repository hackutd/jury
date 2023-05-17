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

    useEffect(() => {
        fetchProjects();
    }, []);

    return (
        <div className="w-full">
            <table>
                <tr>
                    <th></th>
                    <th>Name</th>
                    <th>Table Number</th>
                    <th>Mu</th>
                    <th>Sigma^2</th>
                    <th>Votes</th>
                    <th>Seen</th>
                    <th>Updated</th>
                    <th>Actions</th>
                </tr>
                {projects.map((project: any, idx) => (
                    <tr key={idx}>
                        <td>
                            <input type="checkbox"></input>
                        </td>
                        <td>{project.name}</td>
                        <td>Table {project.location}</td>
                        <td>{project.mu}</td>
                        <td>{project.sigma_sq}</td>
                        <td>{project.votes}</td>
                        <td>{project.seen}</td>
                        {/* TODO: What the fuck is this; just change the datatype to a long pls */}
                        <td>{project.last_activity.$date.$numberLong}</td>
                    </tr>
                ))}
            </table>
        </div>
    );
};

export default ProjectsTable;

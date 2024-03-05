import { useEffect, useState } from 'react';
import Container from '../components/Container';
import { getRequest } from '../api';
import { errorAlert } from '../util';

const Expo = () => {
    const [rawProjects, setRawProjects] = useState<PublicProject[]>([]);
    const [projects, setProjects] = useState<PublicProject[]>([]);
    const [nameSort, setNameSort] = useState(false);

    // Fetch public project list from DB
    useEffect(() => {
        async function fetchProjects() {
            const res = await getRequest<PublicProject[]>('/project/list/public', '');
            if (res.status !== 200) {
                errorAlert(res);
                return;
            }
            setRawProjects(res.data as PublicProject[]);
        }

        fetchProjects();
    }, []);

    // On load or when sort changes, sort by name/table #
    useEffect(() => {
        if (!rawProjects) return;

        const sortedProjects = [...rawProjects];
        if (nameSort) {
            sortedProjects.sort((a, b) => a.name.localeCompare(b.name));
        } else {
            sortedProjects.sort((a, b) => a.location - b.location);
        }

        setProjects(sortedProjects);
    }, [rawProjects, nameSort]);

    return (
        <Container noCenter>
            <h1 className="mt-4 text-4xl text-center font-bold">Project Expo</h1>
            <h2 className="text-2xl text-center font-bold text-primary">
                {import.meta.env.VITE_JURY_NAME}
            </h2>
            <table className="mb-4">
                <thead>
                    <tr>
                        <th
                            onClick={() => setNameSort(true)}
                            className={
                                'px-4 py-2 cursor-pointer text-left ' + (nameSort && 'underline')
                            }
                        >
                            Name
                        </th>
                        <th
                            onClick={() => setNameSort(false)}
                            className={
                                'px-4 py-2 cursor-pointer text-left ' + (!nameSort && 'underline')
                            }
                        >
                            Table
                        </th>
                    </tr>
                </thead>
                <tbody>
                    {projects.map((project, idx) => (
                        <tr key={idx}>
                            <td className="px-4 py-2">
                                <a href={project.url} target="_blank" rel="noopener noreferrer">
                                    {project.name}
                                </a>
                            </td>
                            <td className="px-4 py-2">{project.location}</td>
                        </tr>
                    ))}
                </tbody>
            </table>
        </Container>
    );
};

export default Expo;

import { useEffect, useState } from 'react';
import Container from '../components/Container';
import { getRequest } from '../api';
import { errorAlert } from '../util';
import { useNavigate, useParams } from 'react-router-dom';
import Dropdown from '../components/Dropdown';

const Expo = () => {
    const { track: trackParam } = useParams();
    const [rawProjects, setRawProjects] = useState<PublicProject[]>([]);
    const [projects, setProjects] = useState<PublicProject[]>([]);
    const [nameSort, setNameSort] = useState(false);
    const [track, setTrack] = useState('');
    const [challenges, setChallenges] = useState<string[]>([]);
    const navigate = useNavigate();

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

        const filteredProjects =
            track === ''
                ? sortedProjects
                : sortedProjects.filter((p) => p.challenge_list.indexOf(track) !== -1);

        setProjects(filteredProjects);
    }, [rawProjects, nameSort, track]);

    useEffect(() => {
        if (!trackParam) return;

        console.log(trackParam);

        setTrack(trackParam);
    }, [trackParam]);

    useEffect(() => {
        async function fetchChallenges() {
            const res = await getRequest<string[]>('/challenges', '');
            if (res.status !== 200) {
                errorAlert(res);
                return;
            }

            setChallenges(['', ...(res.data as string[])]);
        }

        fetchChallenges();
    }, []);

    return (
        <Container noCenter>
            <h1 className="mt-4 text-4xl text-center font-bold">Project Expo</h1>
            <h2 className="text-2xl text-center font-bold text-primary">
                {import.meta.env.VITE_JURY_NAME}
            </h2>
            <Dropdown
                selected={track}
                setSelected={setTrack}
                onChange={(t) => {
                    navigate('/expo/' + t.replace(/\s/g, '%20'));
                }}
                options={challenges ?? []}
                className="my-2"
            />
            <h3 className="text-center font-bold text-lighter text-2xl">Count: {projects.length}</h3>
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

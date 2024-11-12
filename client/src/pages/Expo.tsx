import { useEffect, useState } from 'react';
import Container from '../components/Container';
import { getRequest } from '../api';
import { errorAlert } from '../util';
import { useNavigate, useParams, useSearchParams } from 'react-router-dom';
import Dropdown from '../components/Dropdown';
import Button from '../components/Button';
import { Helmet } from 'react-helmet';

const Expo = () => {
    const { track: trackParam } = useParams();
    const [rawProjects, setRawProjects] = useState<PublicProject[]>([]);
    const [projects, setProjects] = useState<PublicProject[]>([]);
    const [nameSort, setNameSort] = useState(false);
    const [track, setTrack] = useState('');
    const [challenges, setChallenges] = useState<string[]>([]);
    const [groupNames, setGroupNames] = useState<string[]>([]);
    const [searchParams, _] = useSearchParams();
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

        async function fetchData() {
            const res = await getRequest<string[]>('/challenges', '');
            if (res.status !== 200) {
                errorAlert(res);
                return;
            }

            const groupRes = await getRequest<string[]>('/group-names', '');
            if (groupRes.status !== 200) {
                errorAlert(groupRes);
                return;
            }
            setGroupNames(groupRes.data as string[]);

            setChallenges(['', ...(res.data as string[])]);
        }

        fetchProjects();
        fetchData();
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

    // Print view
    if (searchParams.get('print') == 'true') {
        return (
            <div>
                <h1 className="text-3xl">{track} Projects</h1>
                <h3 className="font-bold text-lighter text-xl">Count: {projects.length}</h3>
                <table className="table-fixed">
                    <thead className="text-left">
                        <tr>
                            <th className="pr-2">Table</th>
                            <th className="pr-2">Group</th>
                            <th className="pr-2">Name</th>
                        </tr>
                    </thead>
                    <tbody>
                        {projects.map((project, idx) => (
                            <tr key={idx}>
                                <td className="pr-2">{project.location}</td>
                                <td className="pr-2">{groupNames[project.group]}</td>
                                <td className="pr-2">{project.name}</td>
                            </tr>
                        ))}
                    </tbody>
                </table>
            </div>
        );
    }

    return (
        <>
            <Helmet>
                <title>Project Expo | {import.meta.env.VITE_JURY_NAME} | Jury</title>
                <meta
                    name="description"
                    content={`Project expo with table numbers for ${
                        import.meta.env.VITE_JURY_NAME
                    }`}
                />
            </Helmet>
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
                <h3 className="text-center font-bold text-lighter text-2xl">
                    Count: {projects.length}
                </h3>
                <div className="flex items-center flex-col">
                    <Button
                        type="text"
                        href={window.location.href + '?print=true'}
                        className="text-md w-auto py-0 hover:underline"
                    >
                        Print this page
                    </Button>
                </div>
                <table className="mb-4">
                    <thead>
                        <tr>
                            <th
                                onClick={() => setNameSort(true)}
                                className={
                                    'px-4 py-2 cursor-pointer text-left ' +
                                    (nameSort && 'underline')
                                }
                            >
                                Name
                            </th>
                            <th
                                onClick={() => setNameSort(false)}
                                className={
                                    'px-4 py-2 cursor-pointer text-left ' +
                                    (!nameSort && 'underline')
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
                                    <a
                                        href={project.url}
                                        target="_blank"
                                        rel="noopener noreferrer"
                                        className="hover:underline"
                                    >
                                        {project.name}
                                    </a>
                                </td>
                                <td className="px-4 py-2">{project.location}</td>
                            </tr>
                        ))}
                    </tbody>
                </table>
            </Container>
        </>
    );
};

export default Expo;

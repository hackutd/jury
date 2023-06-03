import { useEffect, useState } from 'react';
import ProjectRow from './ProjectRow';

enum SortField {
    Name,
    TableNumber,
    Mu,
    Sigma,
    Votes,
    Seen,
    Updated,
    None,
}

interface SortState {
    field: SortField;
    ascending: boolean;
}

const ProjectsTable = () => {
    const [projects, setProjects] = useState<Project[]>([]);
    const [rawProjects, setRawProjects] = useState<Project[]>([]);
    const [checked, setChecked] = useState<boolean[]>([]);
    const [sortState, setSortState] = useState<SortState>({
        field: SortField.None,
        ascending: true,
    });

    // Fetch projects list
    const fetchProjects = async () => {
        const res = await fetch(`${process.env.REACT_APP_JURY_URL}/project/list`, {
            method: 'GET',
            headers: { 'Content-Type': 'application/json' },
            credentials: 'include',
        });

        if (!res.ok) {
            alert('Error fetching projects');
            return;
        }

        const fetchedProjects = await res.json();
        setProjects(fetchedProjects);
        setRawProjects(fetchedProjects);
        setChecked(Array(fetchedProjects.length).fill(false));
    };

    const handleCheckedChange = (e: React.ChangeEvent<HTMLInputElement>, i: number) => {
        setChecked({
            ...checked,
            [i]: e.target.checked,
        });
    };

    const sort = (field: SortField) => {
        let asc = 1;
        if (sortState.field === field) {
            // If sorted on same field and descending, reset sort state
            if (!sortState.ascending) {
                setSortState({
                    field: SortField.None,
                    ascending: true,
                });
                setProjects(rawProjects);
                return;
            }

            // Otherwise, sort descending
            setSortState({
                field,
                ascending: false,
            });
            asc = -1;
        } else {
            // If in different sorted state, sort ascending on new field
            setSortState({
                field,
                ascending: true,
            });
        }

        let sortFunc = (a: Project, b: Project) => 0;
        switch (field) {
            case SortField.Name:
                sortFunc = (a, b) => a.name.localeCompare(b.name) * asc;
                break;
            case SortField.TableNumber:
                sortFunc = (a, b) => (a.location - b.location) * asc;
                break;
            case SortField.Mu:
                sortFunc = (a, b) => (a.mu - b.mu) * asc;
                break;
            case SortField.Sigma:
                sortFunc = (a, b) => (a.sigma_sq - b.sigma_sq) * asc;
                break;
            case SortField.Votes:
                sortFunc = (a, b) => (a.votes - b.votes) * asc;
                break;
            case SortField.Seen:
                sortFunc = (a, b) => (a.seen - b.seen) * asc;
                break;
            case SortField.Updated:
                sortFunc = (a, b) =>
                    (a.last_activity.$date.$numberLong - b.last_activity.$date.$numberLong) * asc;
                break;
        }
        setProjects([...projects].sort(sortFunc));
    };

    // Onload, fetch projects
    useEffect(() => {
        fetchProjects();
    }, []);

    const arrow = () => (sortState.ascending ? '▲' : '▼');

    return (
        <div className="w-full px-8">
            <table className="table-fixed w-full text-lg">
                <tbody>
                    <tr>
                        <th className="w-12"></th>
                        <th
                            className={
                                'text-left py-1 cursor-pointer hover:text-primary duration-100' +
                                (sortState.field === SortField.Name
                                    ? ' text-primary'
                                    : ' text-black')
                            }
                            onClick={() => sort(SortField.Name)}
                        >
                            Name {sortState.field === SortField.Name && arrow()}
                        </th>
                        <th
                            className={
                                'text-center cursor-pointer hover:text-primary duration-100' +
                                (sortState.field === SortField.TableNumber
                                    ? ' text-primary'
                                    : ' text-black')
                            }
                            onClick={() => sort(SortField.TableNumber)}
                        >
                            Table Number {sortState.field === SortField.TableNumber && arrow()}
                        </th>
                        <th
                            className={
                                'text-center cursor-pointer hover:text-primary duration-100' +
                                (sortState.field === SortField.Mu ? ' text-primary' : ' text-black')
                            }
                            onClick={() => sort(SortField.Mu)}
                        >
                            Mu {sortState.field === SortField.Mu && arrow()}
                        </th>
                        <th
                            className={
                                'text-center cursor-pointer hover:text-primary duration-100' +
                                (sortState.field === SortField.Sigma
                                    ? ' text-primary'
                                    : ' text-black')
                            }
                            onClick={() => sort(SortField.Sigma)}
                        >
                            Sigma^2 {sortState.field === SortField.Sigma && arrow()}
                        </th>
                        <th
                            className={
                                'text-center cursor-pointer hover:text-primary duration-100' +
                                (sortState.field === SortField.Votes
                                    ? ' text-primary'
                                    : ' text-black')
                            }
                            onClick={() => sort(SortField.Votes)}
                        >
                            Votes {sortState.field === SortField.Votes && arrow()}
                        </th>
                        <th
                            className={
                                'text-center cursor-pointer hover:text-primary duration-100' +
                                (sortState.field === SortField.Seen
                                    ? ' text-primary'
                                    : ' text-black')
                            }
                            onClick={() => sort(SortField.Seen)}
                        >
                            Seen {sortState.field === SortField.Seen && arrow()}
                        </th>
                        <th
                            className={
                                'text-center cursor-pointer hover:text-primary duration-100' +
                                (sortState.field === SortField.Updated
                                    ? ' text-primary'
                                    : ' text-black')
                            }
                            onClick={() => sort(SortField.Updated)}
                        >
                            Updated {sortState.field === SortField.Updated && arrow()}
                        </th>
                        <th className="text-right w-24">Actions</th>
                    </tr>
                    {projects.map((project: Project, idx) => (
                        <ProjectRow
                            key={idx}
                            idx={idx}
                            project={project}
                            checked={checked[idx]}
                            handleCheckedChange={handleCheckedChange}
                        />
                    ))}
                </tbody>
            </table>
        </div>
    );
};

export default ProjectsTable;

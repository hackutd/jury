import { useEffect, useState } from 'react';
import { errorAlert, timeSince } from '../../../util';
import DeletePopup from './DeletePopup';
import EditProjectPopup from './EditProjectPopup';
import { useAdminStore, useAdminTableStore, useFlagsStore, useOptionsStore } from '../../../store';
import { putRequest } from '../../../api';
import FlagsPopup from '../FlagsPopup';
import { twMerge } from 'tailwind-merge';
import ActionsDropdown from '../../ActionsDropdown';
import MoveGroupPopup from './MoveGroupPopup';
import MovePopup from './MovePopup';

interface ProjectRowProps {
    project: Project;
    idx: number;
}

const ProjectRow = ({ project, idx }: ProjectRowProps) => {
    const [popup, setPopup] = useState(false);
    const [flagPopup, setFlagPopup] = useState(false);
    const [flags, setFlags] = useState<Flag[]>([]);
    const [flagPopupProjectId, setFlagPopupProjectId] = useState('');
    const [editPopup, setEditPopup] = useState(false);
    const [deletePopup, setDeletePopup] = useState(false);
    const [moveGroupPopup, setMoveGroupPopup] = useState(false);
    const [movePopup, setMovePopup] = useState(false);
    const fetchProjects = useAdminStore((state) => state.fetchProjects);
    const options = useOptionsStore((state) => state.options);
    const track = useOptionsStore((state) => state.selectedTrack);
    const allFlags = useFlagsStore((state) => state.flags);
    const setSelected = useAdminTableStore((state) => state.setSelected);
    const selected = useAdminTableStore((state) => state.selected);

    const handleCheckedChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        // Change the selected state of the project
        const newSelected = selected.slice();
        newSelected[idx] = e.target.checked;
        setSelected(newSelected);
    };

    const hideProject = async () => {
        const res = await putRequest<OkResponse>(`/project/hide/${project.id}`, 'admin', {
            hide: project.active,
        });
        if (res.status !== 200) {
            errorAlert(res);
            return;
        }

        alert(`Project ${project.active ? 'hidden' : 'un-hidden'} successfully!`);
        fetchProjects();
    };

    const prioritizeProject = async () => {
        const res = await putRequest<OkResponse>(`/project/prioritize/${project.id}`, 'admin', {
            prioritize: !project.prioritized,
        });
        if (res.status !== 200) {
            errorAlert(res);
            return;
        }

        alert(`Project ${project.prioritized ? 'un-prioritized' : 'prioritized'} successfully!`);
        fetchProjects();
    };

    useEffect(() => {
        const flags = allFlags.filter(
            (flag) => flag.project_id === project.id && flag.reason !== 'busy'
        );
        setFlags(flags);

        if (flags.length > 0) {
            setFlagPopupProjectId(flags[0].project_id);
        } else {
            setFlagPopupProjectId('null');
        }
    }, [allFlags, project]);

    let stars = project.stars;
    let seen = project.seen;
    if (options.judge_tracks && track !== '') {
        stars = project.track_stars[track] || 0;
        seen = project.track_seen[track] || 0;
    }

    return (
        <>
            <tr
                key={idx}
                className={twMerge(
                    'border-t-2 border-backgroundDark duration-150 bg-background',
                    flags.length >= 1 && 'bg-error/30',
                    project.prioritized && 'bg-gold/30',
                    !project.active && 'bg-lightest',
                    selected && selected[idx] && 'bg-primary/20'
                )}
            >
                <td className="px-2">
                    <input
                        type="checkbox"
                        checked={selected && selected[idx]}
                        onChange={(e) => {
                            handleCheckedChange(e);
                        }}
                        className="cursor-pointer hover:text-primary duration-100"
                    ></input>
                </td>
                <td>
                    <a href={project.url} className="hover:underline" target="_blank">
                        {project.name}
                    </a>
                </td>
                <td className="flex justify-center">
                    <button
                        onClick={() => {
                            setFlagPopup(true);
                        }}
                    >
                        {flags.length >= 1 && (
                            <div className="flex items-center gap-1">
                                <svg
                                    xmlns="http://www.w3.org/2000/svg"
                                    viewBox="0 -960 1000 1000"
                                    className="h-[28px] w-[28px] fill-error mt-1 cursor-pointer hover:scale-110"
                                >
                                    <path d="M200-120v-680h360l16 80h224v400H520l-16-80H280v280h-80Zm300-440Zm86 160h134v-240H510l-16-80H280v240h290l16 80Z" />
                                </svg>
                                {flags.length > 1 && (
                                    <h1 className="mt-1 text-error">({flags.length})</h1>
                                )}
                            </div>
                        )}
                    </button>
                </td>
                <td className="text-center py-1">Table {project.location}</td>
                {options.multi_group && track === '' && (
                    <td className="text-center">{project.group}</td>
                )}
                {track === '' && <td className="text-center">{project.score}</td>}
                <td className="text-center">{stars}</td>
                <td className="text-center">{seen}</td>
                <td className="text-center">{timeSince(project.last_activity)}</td>
                <td className="text-right align-center">
                    <ActionsDropdown
                        open={popup}
                        setOpen={setPopup}
                        actions={[
                            'Edit',
                            project.active ? 'Hide' : 'Unhide',
                            project.prioritized ? 'Unprioritize' : 'Prioritize',
                            'Move Table',
                            'Move Groups',
                            'Delete',
                        ]}
                        actionFunctions={[
                            setEditPopup.bind(null, true),
                            hideProject,
                            prioritizeProject,
                            setMovePopup.bind(null, true),
                            setMoveGroupPopup.bind(null, true),
                            setDeletePopup.bind(null, true),
                        ]}
                        redIndices={[5]}
                    />
                    <div
                        className="cursor-pointer hover:text-primary duration-150 mr-2"
                        onClick={() => {
                            setPopup(!popup);
                        }}
                    >
                        ···
                    </div>
                </td>
            </tr>
            <FlagsPopup
                enabled={flagPopup}
                setEnabled={setFlagPopup}
                projectID={flagPopupProjectId}
            />
            <DeletePopup enabled={deletePopup} setEnabled={setDeletePopup} element={project} />
            <EditProjectPopup enabled={editPopup} setEnabled={setEditPopup} project={project} />
            <MovePopup enabled={movePopup} setEnabled={setMovePopup} item={project} />
            <MoveGroupPopup enabled={moveGroupPopup} setEnabled={setMoveGroupPopup} item={project} isProject />
        </>
    );
};

export default ProjectRow;

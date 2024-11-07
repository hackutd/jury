import { useEffect, useRef, useState } from 'react';
import { errorAlert, fixIfFloatDigits, timeSince } from '../../../util';
import DeletePopup from './DeletePopup';
import EditProjectPopup from './EditProjectPopup';
import { useAdminStore, useOptionsStore } from '../../../store';
import { postRequest } from '../../../api';
import FlagsPopup from '../FlagsPopup';
import { twMerge } from 'tailwind-merge';

interface ProjectRowProps {
    project: Project;
    idx: number;
    flags: Flag[];
    checked: boolean;
    handleCheckedChange: (e: React.ChangeEvent<HTMLInputElement>, idx: number) => void;
}

const ProjectRow = ({ project, idx, flags, checked, handleCheckedChange }: ProjectRowProps) => {
    const [popup, setPopup] = useState(false);
    const [flagPopup, setFlagPopup] = useState(false);
    const [flagCount, setflagCount] = useState(0);
    const [flagPopupProjectId, setflagPopupProjectId] = useState('');
    const [editPopup, setEditPopup] = useState(false);
    const [deletePopup, setDeletePopup] = useState(false);
    const ref = useRef<HTMLDivElement>(null);
    const fetchProjects = useAdminStore((state) => state.fetchProjects);
    const options = useOptionsStore((state) => state.options);
    const track = useOptionsStore((state) => state.selectedTrack);

    useEffect(() => {
        function closeClick(event: MouseEvent) {
            if (ref && ref.current && !ref.current.contains(event.target as Node)) {
                setPopup(false);
            }
        }

        // Bind the event listener
        document.addEventListener('mousedown', closeClick);
        return () => {
            // Unbind the event listener on clean up
            document.removeEventListener('mousedown', closeClick);
        };
    }, [ref]);

    const doAction = (action: 'edit' | 'prioritize' | 'hide' | 'delete') => {
        switch (action) {
            case 'edit':
                // Open edit popup
                setEditPopup(true);
                break;
            case 'hide':
                // Hide
                hideProject();
                break;
            case 'delete':
                // Open delete popup
                setDeletePopup(true);
                break;
            case 'prioritize':
                prioritizeProject();
                break;
            default:
                console.error('Invalid action');
                alert('Invalid action');
                break;
        }

        setPopup(false);
    };

    const hideProject = async () => {
        const res = await postRequest<OkResponse>(
            project.active ? '/project/hide' : '/project/unhide',
            'admin',
            { id: project.id }
        );
        if (res.status === 200) {
            alert(`Project ${project.active ? 'hidden' : 'un-hidden'} successfully!`);
            fetchProjects();
        } else {
            errorAlert(res);
        }
    };

    const prioritizeProject = async () => {
        const res = await postRequest<OkResponse>(
            project.prioritized ? '/project/unprioritize' : '/project/prioritize',
            'admin',
            { id: project.id }
        );
        if (res.status === 200) {
            alert(
                `Project ${project.prioritized ? 'un-prioritized' : 'prioritized'} successfully!`
            );
            fetchProjects();
        } else {
            errorAlert(res);
        }
    };

    useEffect(() => {
        if (flags.length > 0) {
            setflagCount(flagCount);
            setflagPopupProjectId(flags[0].project_id);
        } else {
            setflagPopupProjectId('');
        }
    }, [flags, project]);

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
                    checked && 'bg-primary/20'
                )}
            >
                <td className="px-2">
                    <input
                        type="checkbox"
                        checked={checked}
                        onChange={(e) => {
                            handleCheckedChange(e, idx);
                        }}
                        className="cursor-pointer hover:text-primary duration-100"
                    ></input>
                </td>
                <td>{project.name}</td>
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
                <td className="text-center py-1">
                    Table {project.location} {checked}
                </td>
                {options.multi_group && track === '' && (
                    <td className="text-center">{project.group}</td>
                )}
                {track === '' && <td className="text-center">{project.score}</td>}
                <td className="text-center">{stars}</td>
                <td className="text-center">{seen}</td>
                <td className="text-center">{timeSince(project.last_activity)}</td>
                <td className="text-right align-center">
                    {popup && (
                        <div
                            className="absolute flex flex-col bg-background rounded-md border-lightest border-2 font-normal text-sm"
                            ref={ref}
                        >
                            <div
                                className="py-1 pl-4 pr-2 cursor-pointer hover:bg-primary/20 duration-150"
                                onClick={() => doAction('edit')}
                            >
                                Edit
                            </div>
                            <div
                                className="py-1 pl-4 pr-2 cursor-pointer hover:bg-primary/20 duration-150"
                                onClick={() => doAction('hide')}
                            >
                                {project.active ? 'Hide' : 'Un-hide'}
                            </div>
                            <div
                                className="py-1 pl-4 pr-2 cursor-pointer hover:bg-primary/20 duration-150"
                                onClick={() => doAction('prioritize')}
                            >
                                {project.prioritized ? 'Unprioritize' : 'Prioritize'}
                            </div>

                            <div
                                className="py-1 pl-4 pr-2 cursor-pointer hover:bg-primary/20 duration-150 text-error"
                                onClick={() => doAction('delete')}
                            >
                                Delete
                            </div>
                        </div>
                    )}
                    <div
                        className="cursor-pointer hover:text-primary duration-150 mr-2"
                        onClick={() => {
                            setPopup(!popup);
                        }}
                    >
                        ...
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
        </>
    );
};

export default ProjectRow;

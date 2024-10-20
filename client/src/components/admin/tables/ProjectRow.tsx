import { useEffect, useRef, useState } from 'react';
import { errorAlert, fixIfFloatDigits, timeSince } from '../../../util';
import DeletePopup from './DeletePopup';
import EditProjectPopup from './EditProjectPopup';
import useAdminStore from '../../../store';
import { postRequest } from '../../../api';
import { getRequest } from '../../../api';
import FlagsPopup from '../FlagsPopup';

interface ProjectRowProps {
    project: Project;
    idx: number;
    flags: Flag[];
    checked: boolean;
    handleCheckedChange: (e: React.ChangeEvent<HTMLInputElement>, idx: number) => void;
}

const ProjectRow = ({ project, idx, flags, checked, handleCheckedChange }: ProjectRowProps) => {
    const [popup, setPopup] = useState(false);
    const [showFlags, setShowFlags] = useState(false);
    const [matchingFlagCount, setMatchingFlagCount] = useState(0);
    const [matchingFlag, setMatchingFlag] = useState<Flag | undefined>(undefined);
    const [projectIDState, setProjectIDState] = useState('');
    const [editPopup, setEditPopup] = useState(false);
    const [deletePopup, setDeletePopup] = useState(false);
    const ref = useRef<HTMLDivElement>(null);
    const fetchProjects = useAdminStore((state) => state.fetchProjects);

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

    useEffect(() => {
        // Filter flags to find all matching ones
        const matchingFlags = flags.filter((flag) => flag.project_id === project.id);
        // Set matching flag count
        const matchingFlagCount = matchingFlags.length;
        // If there is at least one matching flag, set the first one as the matching flag
        if (matchingFlagCount > 0) {
            setMatchingFlag(matchingFlags[0]); // Update matchingFlag state to the first match
            setMatchingFlagCount(matchingFlagCount); // Update matchingFlagCount state
            setProjectIDState(matchingFlags[0].project_id); // Update projectIDState based on first matching flag
        } else {
            setMatchingFlag(undefined); // No match found
            setProjectIDState(''); // Clear project ID state
        }
    }, [flags, project]); // Recalculate when flags or project changes

    const rowClassNameFlag = [
        'border-t-2 border-backgroundDark duration-150',
        matchingFlag ? 'bg-error bg-opacity-30' : '',
        checked ? 'bg-primary/20' : !project.active ? 'bg-lightest' : 'bg-background',
    ]
        .filter(Boolean)
        .join(' ');
    return (
        <>
            <tr key={idx} className={rowClassNameFlag}>
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
                            setShowFlags(true);
                        }}
                    >
                        {matchingFlag ? (
                            <div className="flex items-center gap-1">
                                <svg
                                    xmlns="http://www.w3.org/2000/svg"
                                    viewBox="0 -960 1000 1000"
                                    className="h-[28px] w-[28px] fill-error mt-1 cursor-pointer hover:scale-110"
                                >
                                    <path d="M200-120v-680h360l16 80h224v400H520l-16-80H280v280h-80Zm300-440Zm86 160h134v-240H510l-16-80H280v240h290l16 80Z" />
                                </svg>
                                {matchingFlagCount > 1 ? (
                                    <h1 className="mt-1 text-error">({matchingFlagCount})</h1>
                                ) : (
                                    ''
                                )}
                            </div>
                        ) : (
                            ''
                        )}
                    </button>
                </td>
                <td className="text-center py-1">
                    Table {project.location} {checked}
                </td>
                <td className="text-center">{project.score}</td>
                <td className="text-center">{project.seen}</td>
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
                                className="py-1 pl-4 pr-2 cursor-pointer hover:bg-primary/20 duration-150 text-error"
                                onClick={() => doAction('delete')}
                            >
                                Delete
                            </div>
                        </div>
                    )}
                    <div
                        className="cursor-pointer hover:text-primary duration-150"
                        onClick={() => {
                            setPopup(!popup);
                        }}
                    >
                        ...
                    </div>
                </td>
            </tr>
            {showFlags && <FlagsPopup close={setShowFlags} projectID={projectIDState} />}
            {deletePopup && <DeletePopup element={project} close={setDeletePopup} />}
            {editPopup && <EditProjectPopup project={project} close={setEditPopup} />}
        </>
    );
};

export default ProjectRow;

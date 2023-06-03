import { useEffect, useRef, useState } from 'react';
import { timeSince } from '../../util';
import DeletePopup from './DeletePopup';
import EditProjectPopup from './EditProjectPopup';

interface ProjectRowProps {
    project: Project;
    idx: number;
    checked: boolean;
    handleCheckedChange: (e: React.ChangeEvent<HTMLInputElement>, idx: number) => void;
}

const ProjectRow = ({ project, idx, checked, handleCheckedChange }: ProjectRowProps) => {
    const [popup, setPopup] = useState(false);
    const [editPopup, setEditPopup] = useState(false);
    const [deletePopup, setDeletePopup] = useState(false);
    const ref = useRef<HTMLDivElement>(null);

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
            case 'prioritize':
                // Prioritize
                break;
            case 'hide':
                // Hide
                break;
            case 'delete':
                // Open delete popup
                setDeletePopup(true);
                break;
        }

        setPopup(false);
    };

    return (
        <>
            <tr
                key={idx}
                className={
                    'border-t-2 border-backgroundDark duration-150 ' +
                    (checked ? 'bg-primary/20' : 'bg-background')
                }
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
                <td className="text-center py-1">
                    Table {project.location} {checked}
                </td>
                <td className="text-center">{project.mu}</td>
                <td className="text-center">{project.sigma_sq}</td>
                <td className="text-center">{project.votes}</td>
                <td className="text-center">{project.seen}</td>
                {/* TODO: What the fuck is this; just change the datatype to a long pls */}
                <td className="text-center">
                    {timeSince(project.last_activity.$date.$numberLong)}
                </td>
                <td className="text-right font-bold flex align-center justify-end">
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
                                onClick={() => doAction('prioritize')}
                            >
                                Prioritize
                            </div>
                            <div
                                className="py-1 pl-4 pr-2 cursor-pointer hover:bg-primary/20 duration-150"
                                onClick={() => doAction('hide')}
                            >
                                Hide
                            </div>
                            <div
                                className="py-1 pl-4 pr-2 cursor-pointer hover:bg-primary/20 duration-150 text-error"
                                onClick={() => doAction('delete')}
                            >
                                Delete
                            </div>
                        </div>
                    )}
                    <span
                        className="cursor-pointer px-1 hover:text-primary duration-150"
                        onClick={() => {
                            setPopup(!popup);
                        }}
                    >
                        ...
                    </span>
                </td>
            </tr>
            {deletePopup && <DeletePopup element={project} close={setDeletePopup} />}
            {editPopup && <EditProjectPopup project={project} close={setEditPopup} />}
        </>
    );
};

export default ProjectRow;

import { useEffect, useRef, useState } from 'react';
import { errorAlert, fixIfFloatDigits, timeSince } from '../../../util';
import DeletePopup from './DeletePopup';
import EditJudgePopup from './EditJudgePopup';
import { postRequest } from '../../../api';
import useAdminStore from '../../../store';
import { twMerge } from 'tailwind-merge';

interface JudgeRowProps {
    judge: Judge;
    idx: number;
    checked: boolean;
    handleCheckedChange: (e: React.ChangeEvent<HTMLInputElement>, idx: number) => void;
}

const JudgeRow = ({ judge, idx, checked, handleCheckedChange }: JudgeRowProps) => {
    const [popup, setPopup] = useState(false);
    const [editPopup, setEditPopup] = useState(false);
    const [deletePopup, setDeletePopup] = useState(false);
    const ref = useRef<HTMLDivElement>(null);
    const fetchJudges = useAdminStore((state) => state.fetchJudges);

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
                hideJudge();
                break;
            case 'delete':
                // Open delete popup
                setDeletePopup(true);
                break;
        }

        setPopup(false);
    };

    const hideJudge = async () => {
        const res = await postRequest<OkResponse>(judge.active ? '/judge/hide' : '/judge/unhide', 'admin', { id: judge.id });
        if (res.status === 200) {
            alert(`Judge ${judge.active ? 'hidden' : 'un-hidden'} successfully!`);
            fetchJudges();
        } else {
            errorAlert(res.status);
        }
    };

    return (
        <>
            <tr
                key={idx}
                className={twMerge(
                    'border-t-2 border-backgroundDark duration-150',
                    checked ? 'bg-primary/20' : !judge.active ? 'bg-lightest' : 'bg-background'
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
                <td>{judge.name}</td>
                <td className="text-center">{judge.votes}</td>
                <td className="text-center">{fixIfFloatDigits(judge.alpha, 5)}</td>
                <td className="text-center">{fixIfFloatDigits(judge.beta, 5)}</td>
                <td className="text-center">{timeSince(judge.last_activity)}</td>
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
                                onClick={() => doAction('hide')}
                            >
                                {judge.active ? 'Hide' : 'Un-hide'}
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
            {deletePopup && <DeletePopup element={judge} close={setDeletePopup} />}
            {editPopup && <EditJudgePopup judge={judge} close={setEditPopup} />}
        </>
    );
};

export default JudgeRow;

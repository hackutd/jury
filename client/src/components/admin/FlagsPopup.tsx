import { useEffect, useState } from 'react';
import { errorAlert, timeSince } from '../../util';
import Popup from '../Popup';
import { useFlagsStore } from '../../store';
import { deleteRequest } from '../../api';

// Map reasons to text to display
const REASON_MAP = {
    busy: 'Busy',
    absent: 'Absent',
    'cannot-demo': "Can't Demo",
    'too-complex': 'Too Complex',
    offensive: 'Offensive',
    'hidden-absent': 'Hidden bc Absent',
} as any;

interface FlagsPopupProps {
    /* State variable for open/closed */
    enabled: boolean;

    /* Function to modify the popup state variable */
    setEnabled: React.Dispatch<React.SetStateAction<boolean>>;

    /* Project ID to filter flags */
    projectID?: string;
}

const FlagsPopup = (props: FlagsPopupProps) => {
    const flags = useFlagsStore((state) => state.flags);
    const [sortedFlags, setSortedFlags] = useState<Flag[]>([]);
    const [sortMethod, setSortMethod] = useState(3);
    const [displayedFlags, setDisplayedFlags] = useState<Flag[]>([]);
    const fetchFlags = useFlagsStore((state) => state.fetchFlags);

    useEffect(() => {
        // Filter out busy flags and sort by time
        const filteredFlags = sortedFlags
            .filter((flag) => !flag.reason.includes('busy'))
            .toSorted((a, b) => b.time - a.time);

        // If a projectID is provided, only display flags for that project
        if (props.projectID) {
            setDisplayedFlags(filteredFlags.filter((flag) => flag.project_id === props.projectID));
            return;
        }

        setDisplayedFlags(filteredFlags);
    }, [props.projectID, sortedFlags]);

    useEffect(() => {
        switch (sortMethod) {
            case 0:
                // Sort by project_name
                setSortedFlags(
                    flags.toSorted((a, b) => a.project_name.localeCompare(b.project_name))
                );
                break;
            case 1:
                // Sort by judge_name
                setSortedFlags(flags.toSorted((a, b) => a.judge_name.localeCompare(b.judge_name)));
                break;
            case 2:
                // Sort by reason
                setSortedFlags(flags.toSorted((a, b) => a.reason.localeCompare(b.reason)));
                break;
            case 3:
                // Sort by time
                setSortedFlags(flags.toSorted((a, b) => a.time - b.time));
                break;
        }
    }, [sortMethod, flags]);

    const clearFlag = async (id: string) => {
        const res = await deleteRequest(`/admin/flag/${id}`, 'admin');
        if (res.status !== 200) {
            errorAlert(res);
            return;
        }

        fetchFlags();
    };

    return (
        <Popup enabled={props.enabled} setEnabled={props.setEnabled} className="overflow-hidden">
            <h1 className="text-5xl text-center font-bold mb-6">Flags</h1>
            <div className="overflow-x-scroll">
                <div className="h-[50vh] overflow-y-auto sm:w-auto w-[200%]">
                    <div className="flex flex-row sm:items-center gap-2 text-xl border-b-2 border-backgroundDark py-2">
                        <h2
                            className={
                                'basis-2/6 text-left font-bold cursor-pointer ' +
                                (sortMethod === 0 ? 'text-primary' : 'text-light')
                            }
                            onClick={() => setSortMethod(0)}
                        >
                            Name
                        </h2>
                        <h2
                            className={
                                'basis-1/6 text-center font-bold cursor-pointer ' +
                                (sortMethod === 1 ? 'text-primary' : 'text-light')
                            }
                            onClick={() => setSortMethod(1)}
                        >
                            Judge
                        </h2>
                        <h2
                            className={
                                'basis-1/6 text-center font-bold cursor-pointer ' +
                                (sortMethod === 2 ? 'text-primary' : 'text-light')
                            }
                            onClick={() => setSortMethod(2)}
                        >
                            Reason
                        </h2>
                        <h2
                            className={
                                'basis-1/6 text-center font-bold cursor-pointer ' +
                                (sortMethod === 3 ? 'text-primary' : 'text-light')
                            }
                            onClick={() => setSortMethod(3)}
                        >
                            Time
                        </h2>
                        <h2
                            className={
                                'basis-1/6 text-right font-bold cursor-pointer ' +
                                (sortMethod === 3 ? 'text-primary' : 'text-light')
                            }
                        >
                            Resolve
                        </h2>
                    </div>
                    {displayedFlags.map((flag) => (
                        <div
                            key={`${flag.id}`}
                            className="flex flex-row items-center gap-2 text-xl border-b-2 border-backgroundDark py-1"
                        >
                            <h2 className="basis-2/6 text-left text-lg">
                                {`[${flag.project_location}] ${flag.project_name}`}
                            </h2>
                            <h2 className="basis-1/6 text-center text-lg">{flag.judge_name}</h2>
                            <h2 className="basis-1/6 text-center text-lg">
                                {REASON_MAP[flag.reason] || flag.reason}
                            </h2>
                            <h2 className="basis-1/6 text-center text-lg">
                                {timeSince(flag.time)}
                            </h2>
                            <button
                                className="basis-1/6 text-right cursor-pointer text-error hover:text-errorDark duration-150"
                                onClick={() => clearFlag(flag.id)}
                            >
                                resolve
                            </button>
                        </div>
                    ))}
                </div>
            </div>
        </Popup>
    );
};

export default FlagsPopup;

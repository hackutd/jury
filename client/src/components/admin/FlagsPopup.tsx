import { useEffect, useState } from 'react';
import { getRequest } from '../../api';
import { errorAlert, timeSince } from '../../util';

interface FlagsPopupProps {
    /* Function to modify the popup state variable */
    close: React.Dispatch<React.SetStateAction<boolean>>;
    projectID?: string;
}

const FlagsPopup = ({ close, projectID }: FlagsPopupProps) => {
    const [flags, setFlags] = useState<Flag[]>([]);
    const [sortedFlags, setSortedFlags] = useState<Flag[]>([]);
    const [sortMethod, setSortMethod] = useState(3);
    const [displayedFlags, setDisplayedFlags] = useState<Flag[]>([]);

    useEffect(() => {
        // Filter out busy flags and sort by time
        const filteredFlags = sortedFlags
            .filter((flag) => !flag.reason.includes('busy'))
            .toSorted((a, b) => b.time - a.time);

        // If a projectID is provided, only display flags for that project
        if (projectID) {
            setDisplayedFlags(filteredFlags.filter((flag) => flag.project_id === projectID));
            return;
        }

        setDisplayedFlags(filteredFlags);
    }, [projectID, sortedFlags]);

    useEffect(() => {
        async function getFlags() {
            const res = await getRequest<Flag[]>('/admin/flags', 'admin');
            if (res.status !== 200) {
                errorAlert(res);
            }
            setFlags(res.data as Flag[]);
        }

        getFlags();
    }, []);

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

    return (
        <>
            <div
                className="fixed left-0 top-0 z-20 w-screen h-screen bg-black/30"
                onClick={() => close(false)}
            ></div>
            <div className="bg-background fixed z-30 left-1/2 top-1/2 translate-x-[-50%] translate-y-[-50%] py-6 px-10 w-2/3">
                <h1 className="text-5xl text-center font-bold mb-6">Flags</h1>
                <div className="h-[50vh] overflow-y-auto">
                    <div className="flex flex-row items-center text-xl border-b-2 border-backgroundDark py-2">
                        <h2
                            className={
                                'basis-2/5 text-left font-bold cursor-pointer ' +
                                (sortMethod === 0 ? 'text-primary' : 'text-light')
                            }
                            onClick={() => setSortMethod(0)}
                        >
                            Name
                        </h2>
                        <h2
                            className={
                                'basis-1/5 text-left font-bold cursor-pointer ' +
                                (sortMethod === 1 ? 'text-primary' : 'text-light')
                            }
                            onClick={() => setSortMethod(1)}
                        >
                            Judge
                        </h2>
                        <h2
                            className={
                                'basis-1/5 text-left font-bold cursor-pointer ' +
                                (sortMethod === 2 ? 'text-primary' : 'text-light')
                            }
                            onClick={() => setSortMethod(2)}
                        >
                            Reason
                        </h2>
                        <h2
                            className={
                                'basis-1/5 text-left font-bold cursor-pointer ' +
                                (sortMethod === 3 ? 'text-primary' : 'text-light')
                            }
                            onClick={() => setSortMethod(3)}
                        >
                            Time
                        </h2>
                    </div>
                    {displayedFlags.map((flag) => (
                        <div
                            key={`${flag.id}`}
                            className="flex flex-row items-center text-xl border-b-2 border-backgroundDark py-1"
                        >
                            <h2 className="basis-2/5 text-left text-lg">
                                {`[${flag.project_location}] ${flag.project_name}`}
                            </h2>
                            <h2 className="basis-1/5 text-left text-lg">{flag.judge_name}</h2>
                            <h2 className="basis-1/5 text-left text-lg">{flag.reason}</h2>
                            <h2 className="basis-1/5 text-left text-lg">{timeSince(flag.time)}</h2>
                        </div>
                    ))}
                </div>
            </div>
        </>
    );
};

export default FlagsPopup;

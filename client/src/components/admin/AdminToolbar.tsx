import { useEffect, useState } from 'react';
import Button from '../Button';
import FlagsPopup from './FlagsPopup';
import Dropdown from '../Dropdown';
import { useAdminStore, useAdminTableStore, useOptionsStore } from '../../store';
import ActionsDropdown from '../ActionsDropdown';
import { postRequest } from '../../api';
import { errorAlert } from '../../util';
import MovePopup from './tables/MovePopup';

const AdminToolbar = (props: { showProjects: boolean; lastUpdate: Date }) => {
    const [showFlags, setShowFlags] = useState(false);
    const options = useOptionsStore((state) => state.options);
    const selectedTrack = useOptionsStore((state) => state.selectedTrack);
    const setSelectedTrack = useOptionsStore((state) => state.setSelectedTrack);
    const fetchStats = useAdminStore((state) => state.fetchStats);
    const selected = useAdminTableStore((state) => state.selected);
    const setSelected = useAdminTableStore((state) => state.setSelected);
    const [batchPopup, setBatchPopup] = useState(false);
    const [selectedIds, setSelectedIds] = useState<string[]>([]);
    const projects = useAdminTableStore((state) => state.projects);
    const judges = useAdminTableStore((state) => state.judges);
    const [movePopup, setMovePopup] = useState(false);
    const fetchJudges = useAdminStore((state) => state.fetchJudges);
    const fetchProjects = useAdminStore((state) => state.fetchProjects);

    // Convert date to string
    const dateToString = (date: Date) => {
        return date.toLocaleDateString('en-US', {
            year: 'numeric',
            month: 'long',
            day: 'numeric',
            hour: 'numeric',
            minute: 'numeric',
            second: 'numeric',
            timeZoneName: 'short',
        });
    };

    // Fetch stats when track changes
    useEffect(() => {
        fetchStats();
    }, [selectedTrack]);

    // When selected changes, update selectedIdxs
    useEffect(() => {
        if (!selected) return;

        const ref = props.showProjects ? projects : judges;
        console.log(ref);

        const newSelectedIds: string[] = [];
        selected.forEach((isSelected, idx) => {
            if (isSelected) newSelectedIds.push(ref[idx].id);
        });
        setSelectedIds(newSelectedIds);

        if (newSelectedIds.length === 0) {
            setBatchPopup(false);
        }
    }, [selected]);

    const hide = async (h: boolean) => {
        if (!selected) return;

        const res = await postRequest<OkResponse>(
            props.showProjects ? '/project/hide' : '/judge/hide',
            'admin',
            {
                items: selectedIds,
                hide: h,
            }
        );
        if (res.status !== 200) {
            errorAlert(res);
            return;
        }

        alert(
            `Successfully ${h ? 'hidden' : 'unhidden'} ${selectedIds.length} ${
                props.showProjects ? 'projects' : 'judges'
            }!`
        );

        reset();
    };

    const prioritize = async (p: boolean) => {
        if (!selected) return;

        const res = await postRequest<OkResponse>('/project/prioritize', 'admin', {
            items: selectedIds,
            prioritize: p,
        });
        if (res.status !== 200) {
            errorAlert(res);
            return;
        }

        alert(
            `Successfully ${p ? 'prioritized' : 'unprioritized'} ${selectedIds.length} projects!`
        );

        reset();
    };

    const move = async (group: number) => {
        if (!selected) return;

        const res = await postRequest<OkResponse>('/judge/move', 'admin', {
            items: selectedIds,
            group,
        });
        if (res.status !== 200) {
            errorAlert(res);
            return;
        }

        alert(`Successfully moved ${selectedIds.length} judges to group ${group}!`);

        reset();
    };

    const reset = () => {
        props.showProjects ? fetchProjects() : fetchJudges();
        setSelected(new Array(selected.length).fill(false));
        setBatchPopup(false);
    };

    return (
        <>
            <div className="flex flex-row px-8 py-4">
                <Button
                    type="outline"
                    bold
                    small
                    className="py-2 px-4 rounded-md"
                    href={props.showProjects ? '/admin/add-projects' : '/admin/add-judges'}
                >
                    Add {props.showProjects ? 'Projects' : 'Judges'}
                </Button>
                {props.showProjects && (
                    <Button
                        type="outline"
                        bold
                        small
                        className="py-2 px-4 rounded-md ml-4"
                        onClick={() => {
                            setShowFlags(true);
                        }}
                    >
                        See All Flags
                    </Button>
                )}
                {options && options.judge_tracks && (
                    <Dropdown
                        options={['Main Judging', ...options.tracks]}
                        selected={selectedTrack === '' ? 'Main Judging' : selectedTrack}
                        setSelected={setSelectedTrack}
                        className="ml-4"
                    />
                )}
                {selectedIds.length > 0 && (
                    <div className="flex flex-row relative">
                        <Button
                            type="outline"
                            bold
                            small
                            className="px-4 rounded-md ml-4"
                            onClick={setBatchPopup.bind(null, true)}
                        >
                            Batch Ops
                        </Button>
                        <ActionsDropdown
                            actions={
                                props.showProjects
                                    ? ['Hide', 'Unhide', 'Prioritize', 'Unprioritize']
                                    : ['Hide', 'Unhide', 'Move']
                            }
                            actionFunctions={
                                props.showProjects
                                    ? [
                                          hide.bind(null, true),
                                          hide.bind(null, false),
                                          prioritize.bind(null, true),
                                          prioritize.bind(null, false),
                                      ]
                                    : [
                                          hide.bind(null, true),
                                          hide.bind(null, false),
                                          setMovePopup.bind(null, true),
                                      ]
                            }
                            open={batchPopup}
                            setOpen={setBatchPopup}
                            large
                            className="left-4 top-12"
                        />
                        <span className="text-xl text-lighter flex flex-row items-center">
                            <span className="font-bold ml-4">Selected:</span>&nbsp;
                            {selectedIds.length}
                        </span>
                    </div>
                )}

                <p className="grow self-end text-right text-lighter">
                    Last Update: {dateToString(props.lastUpdate)}
                </p>
            </div>
            <FlagsPopup enabled={showFlags} setEnabled={setShowFlags} />
            <MovePopup open={movePopup} setOpen={setMovePopup} judges={selectedIds} />
        </>
    );
};

export default AdminToolbar;

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
    const options = useOptionsStore((state) => state.options);
    const selectedTrack = useOptionsStore((state) => state.selectedTrack);
    const setSelectedTrack = useOptionsStore((state) => state.setSelectedTrack);
    const fetchStats = useAdminStore((state) => state.fetchStats);
    const fetchJudges = useAdminStore((state) => state.fetchJudges);
    const fetchProjects = useAdminStore((state) => state.fetchProjects);
    const selected = useAdminTableStore((state) => state.selected);
    const setSelected = useAdminTableStore((state) => state.setSelected);
    const projects = useAdminTableStore((state) => state.projects);
    const judges = useAdminTableStore((state) => state.judges);
    const [showFlags, setShowFlags] = useState(false);
    const [batchPopup, setBatchPopup] = useState(false);
    const [movePopup, setMovePopup] = useState(false);
    const [projectMovePopup, setProjectMovePopup] = useState(false);
    const [selectedIds, setSelectedIds] = useState<string[]>([]);

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

    const reset = () => {
        props.showProjects ? fetchProjects() : fetchJudges();
        setSelected(new Array(selected.length).fill(false));
        setBatchPopup(false);
    };

    return (
        <>
            <div className="flex xl:flex-row flex-col md:px-8 md:py-4 gap-4">
                <div className="flex md:flex-row flex-col md:items-start items-center gap-4">
                    <div className="flex flex-row gap-4 h-12">
                        <Button
                            type="outline"
                            bold
                            href={props.showProjects ? '/admin/add-projects' : '/admin/add-judges'}
                            className="text-md md:text-xl px-2 md:px-6"
                        >
                            Add {props.showProjects ? 'Projects' : 'Judges'}
                        </Button>
                        {props.showProjects && (
                            <Button
                                type="outline"
                                bold
                                onClick={setShowFlags.bind(null, true)}
                                className="text-md md:text-xl px-2 md:px-6"
                            >
                                See All Flags
                            </Button>
                        )}
                    </div>
                    {options && options.judge_tracks && (
                        <Dropdown
                            options={['Main Judging', ...options.tracks]}
                            selected={selectedTrack === '' ? 'Main Judging' : selectedTrack}
                            setSelected={setSelectedTrack}
                            large
                            className="text-center text-md md:text-xl h-12"
                        />
                    )}
                </div>
                {selectedIds.length > 0 && (
                    <div className="flex flex-row relative h-12">
                        <Button type="outline" bold onClick={setBatchPopup.bind(null, true)}>
                            Batch Ops
                        </Button>
                        <ActionsDropdown
                            actions={
                                props.showProjects
                                    ? [
                                          'Hide',
                                          'Unhide',
                                          'Prioritize',
                                          'Unprioritize',
                                          'Move Groups',
                                      ]
                                    : ['Hide', 'Unhide', 'Move']
                            }
                            actionFunctions={
                                props.showProjects
                                    ? [
                                          hide.bind(null, true),
                                          hide.bind(null, false),
                                          prioritize.bind(null, true),
                                          prioritize.bind(null, false),
                                          setProjectMovePopup.bind(null, true),
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

                <p className="grow self-end text-center w-full md:w-auto md:text-right text-lighter text-sm md:text-md">
                    Last Update: {dateToString(props.lastUpdate)}
                </p>
            </div>
            <FlagsPopup enabled={showFlags} setEnabled={setShowFlags} />
            <MovePopup enabled={movePopup} setEnabled={setMovePopup} items={selectedIds} />
            <MovePopup
                enabled={projectMovePopup}
                setEnabled={setProjectMovePopup}
                items={selectedIds}
                isProject
            />
        </>
    );
};

export default AdminToolbar;

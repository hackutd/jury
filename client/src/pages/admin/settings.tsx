import { useEffect, useState } from 'react';
import { createHeaders, getRequest, postRequest } from '../../api';
import Button from '../../components/Button';
import JuryHeader from '../../components/JuryHeader';
import { errorAlert } from '../../util';
import ConfirmPopup from '../../components/ConfirmPopup';
import Loading from '../../components/Loading';
import Checkbox from '../../components/Checkbox';
import TextInput from '../../components/TextInput';
import SelectionButton from '../../components/SelectionButton';
import { useClockStore, useOptionsStore } from '../../store';
import ChallengeBlock from '../../components/admin/ChallengeBlock';
import Card from '../../components/Card';
import ToTopButton from '../../components/ToTopButton';
import { useNavigate } from 'react-router-dom';
import { twMerge } from 'tailwind-merge';

// Text components
const Section = ({ children: c }: { children: string }) => (
    <h2 id={c.toLowerCase().replace(/ /g, '-')} className="text-3xl md:text-4xl text-primary">
        {c}
    </h2>
);
const SubSection = ({ children: c }: { children: React.ReactNode }) => (
    <h3 className="text-lg md:text-xl mt-4 mb-1 font-bold">{c}</h3>
);
const Description = ({ children: c }: { children: React.ReactNode }) => (
    <p className="text-sm md:text-md text-light mb-[0.125rem]">{c}</p>
);

// Custom button to use on settings page
const SettingsButton = ({
    children,
    onClick,
    disabled,
    className,
    type = 'primary',
}: {
    children: React.ReactNode;
    onClick: () => void;
    disabled?: boolean;
    className?: string;
    type?: 'primary' | 'error' | 'gold' | 'outline';
}) => (
    <Button
        type={type}
        onClick={onClick}
        disabled={disabled}
        className={twMerge('my-2 h-max text-lg md:text-2xl py-[2px] md:py-1', className)}
    >
        {children}
    </Button>
);

// Custom button at the top to nav to sections
const NavButton = ({ children }: { children: string }) => (
    <Button
        type="primary"
        onClick={() => {
            window.location.href = `#${children.toLowerCase().replace(/ /g, '-')}`;
        }}
        className="text-md px-2 md:px-6 py-[2px] md:py-1 md:text-xl"
    >
        {children}
    </Button>
);

const FieldButton = ({ children }: { children: React.ReactNode }) => {
    return <div className="flex flex-col md:flex-row items-start md:items-end">{children}</div>;
};

const AdminSettings = () => {
    const [reassignPopup, setReassignPopup] = useState(false);
    const [judgeReassignPopup, setJudgeReassignPopup] = useState(false);
    const [clockResetPopup, setClockResetPopup] = useState(false);
    const [judgingTimer, setJudgingTimer] = useState('');
    const [minViews, setMinViews] = useState(3);
    const [syncClock, setSyncClock] = useState(false);
    const [loading, setLoading] = useState(true);
    const [multiGroup, setMultiGroup] = useState(false);
    const [numGroups, setNumGroups] = useState(3);
    const [switchingMode, setSwitchingMode] = useState('auto');
    const [autoSwitchProp, setAutoSwitchProp] = useState(0.1);
    const [groupSizes, setGroupSizes] = useState('30, 30');
    const [judgeTracks, setJudgeTracks] = useState(false);
    const [tracks, setTracks] = useState<string>('');
    const [trackViews, setTrackViews] = useState('');
    const [groupNames, setGroupNames] = useState('');
    const [ignoreTracks, setIgnoreTracks] = useState<string>('');
    const [maxReqPerMin, setMaxReqPerMin] = useState(100);
    const [blockReqs, setBlockReqs] = useState(false);
    const [resetPopup, setResetPopup] = useState<ResetPopup>({
        open: false,
        text: '',
        handler: null,
    });
    const [running, setRunning] = useState(false);
    const fetchOptions = useOptionsStore((state) => state.fetchOptions);
    const options = useOptionsStore((state) => state.options);
    const fetchClock = useClockStore((state) => state.fetchClock);
    const clock = useClockStore((state) => state.clock);
    const navigate = useNavigate();

    async function getOptions() {
        await fetchOptions();

        // Calculate judging timer MM:SS
        const timer = options.judging_timer;
        if (timer) {
            const minutes = Math.floor(timer / 60);
            const seconds = timer % 60;
            const timerStr = `${minutes}:${seconds < 10 ? '0' : ''}${seconds}`;
            setJudgingTimer(timerStr);
        }

        // Set min views
        setMinViews(options.min_views);

        // Set sync clock
        setSyncClock(options.clock_sync);

        // Set group options
        setMultiGroup(options.multi_group);
        setNumGroups(options.num_groups);
        setGroupSizes(options.group_sizes.join(', '));
        setSwitchingMode(options.switching_mode);
        setAutoSwitchProp(options.auto_switch_prop);
        setJudgeTracks(options.judge_tracks);
        setTracks(options.tracks.join(', '));
        setTrackViews(options.track_views.join(', '));
        setGroupNames(options.group_names.join(', '));
        setIgnoreTracks(options.ignore_tracks.join(', '));
        setBlockReqs(options.block_reqs);
        setMaxReqPerMin(options.max_req_per_min);

        // Get active clock status
        await fetchClock();
        setRunning(clock.running);

        setLoading(false);
    }

    // Get the previous options on load
    useEffect(() => {
        async function checkAuthAndGetOptions() {
            const loggedInRes = await postRequest<OkResponse>('/admin/auth', 'admin', null);
            if (loggedInRes.status === 401) {
                console.error(`Admin is not logged in!`);
                navigate('/admin/login');
                return;
            }
            if (loggedInRes.status === 200) {
                getOptions();
                return;
            }

            errorAlert(loggedInRes);
        }

        checkAuthAndGetOptions();
    }, []);

    const updateBlockReqs = async () => {
        const res = await postRequest<OkResponse>('/admin/block-reqs', 'admin', {
            block_reqs: !blockReqs,
        });
        if (res.status !== 200 || res.data?.ok !== 1) {
            errorAlert(res);
            return;
        }

        setBlockReqs(!blockReqs);
        alert(`Judge login is now ${blockReqs ? 'enabled' : 'BLOCKED'}`);
    };

    const updateMaxReqPerMin = async () => {
        if (isNaN(maxReqPerMin) || maxReqPerMin < 0) {
            alert('Max requests per minute should be a positive integer!');
            return;
        }

        const res = await postRequest<OkResponse>('/admin/max-reqs', 'admin', {
            max_req_per_min: maxReqPerMin,
        });
        if (res.status !== 200 || res.data?.ok !== 1) {
            errorAlert(res);
            return;
        }

        alert('Max requests per minute updated!');
    };

    const reassignTables = async () => {
        const res = await postRequest<OkResponse>('/project/reassign', 'admin', null);
        if (res.status !== 200 || res.data?.ok !== 1) {
            errorAlert(res);
            return;
        }
        alert('Table numbers reassigned!');
        setReassignPopup(false);
    };

    const reassignJudges = async () => {
        const res = await postRequest<OkResponse>('/judge/reassign', 'admin', null);
        if (res.status !== 200 || res.data?.ok !== 1) {
            errorAlert(res);
            return;
        }
        alert('Judges reassigned!');
        setJudgeReassignPopup(false);
    };

    const updateTimer = async () => {
        // Convert judging timer to time
        const [minutes, seconds] = judgingTimer.split(':');
        const timer = judgingTimer === '' ? 0 : parseInt(minutes) * 60 + parseInt(seconds);

        // Check to make sure timer is valid and positive
        if (isNaN(timer)) {
            alert('Invalid timer format!');
            return;
        }
        if (timer < 0) {
            alert('Timer must be a positive number!');
            return;
        }

        // Update the timer
        const res = await postRequest<OkResponse>('/admin/options', 'admin', {
            judging_timer: timer,
        });
        if (res.status !== 200 || res.data?.ok !== 1) {
            errorAlert(res);
            return;
        }

        alert('Timer updated!');
        getOptions();
    };

    const updateMinViews = async () => {
        // Convert minViews to integer
        if (isNaN(minViews) || minViews < 0) {
            alert('Minimum views should be a positive integer!');
            return;
        }

        // Update min views
        const res = await postRequest<OkResponse>('/admin/options', 'admin', {
            min_views: minViews,
        });
        if (res.status !== 200 || res.data?.ok !== 1) {
            errorAlert(res);
            return;
        }

        alert('Min views updated!');
        getOptions();
    };

    const toggleJudgeTracks = async () => {
        const res = await postRequest<OkResponse>('/admin/options', 'admin', {
            judge_tracks: !judgeTracks,
        });
        if (res.status !== 200 || res.data?.ok !== 1) {
            errorAlert(res);
            return;
        }

        setJudgeTracks(!judgeTracks);
        alert('Track judging toggled!');
    };

    const updateTracks = async () => {
        // Split tracks by comma and remove empty strings
        const filteredTracks = tracks
            .split(',')
            .map((track) => track.trim())
            .filter((track) => track !== '');

        const res = await postRequest<OkResponse>('/admin/tracks', 'admin', {
            tracks: filteredTracks,
        });
        if (res.status !== 200 || !res.data) {
            errorAlert(res);
            return;
        }

        alert('Tracks updated!');
        getOptions();
    };

    const updateTrackViews = async () => {
        const views = trackViews.split(',').map((v) => parseInt(v.trim()));
        const ts = tracks
            .split(',')
            .map((t) => t.trim())
            .filter((t) => t !== '');
        if (views.some((v) => isNaN(v) || v < 1)) {
            alert('Track views should be positive integers!');
            return;
        }

        // Make sure number of views is equal to num of tracks
        if (views.length !== ts.length) {
            alert(
                'Number of track views should be the same as number of tracks (you need a min number of views for each track)!'
            );
            return;
        }

        const res = await postRequest<OkResponse>('/admin/track-views', 'admin', {
            track_views: views,
        });
        if (res.status !== 200 || res.data?.ok !== 1) {
            errorAlert(res);
            return;
        }

        alert('Track views updated!');
    };

    const toggleMultiGroup = async () => {
        const res = await postRequest<OkResponse>('/admin/options', 'admin', {
            multi_group: !multiGroup,
        });
        if (res.status !== 200 || res.data?.ok !== 1) {
            errorAlert(res);
            return;
        }

        setMultiGroup(!multiGroup);
        fetchOptions();
        alert('Multi-group judging toggled!');
    };

    const updateNumGroups = async () => {
        if (isNaN(numGroups) || numGroups < 1) {
            alert('Number of groups should be an integer greater than 0!');
            return;
        }

        const res = await postRequest<OkResponse>('/admin/num-groups', 'admin', {
            num_groups: numGroups,
        });
        if (res.status !== 200 || res.data?.ok !== 1) {
            errorAlert(res);
            return;
        }

        // Refreah options
        await getOptions();

        alert('Number of groups updated!');
    };

    const updateSwitchingMode = async (newMode: string) => {
        const res = await postRequest<OkResponse>('/admin/options', 'admin', {
            switching_mode: newMode,
        });
        if (res.status !== 200 || res.data?.ok !== 1) {
            errorAlert(res);
            return;
        }

        fetchOptions();
        alert('Switching method updated!');
    };

    const updateAutoSwitchProp = async () => {
        if (isNaN(autoSwitchProp) || autoSwitchProp < 0 || autoSwitchProp > 1) {
            alert('Auto switch proportion should be a decimal between 0 and 1!');
            return;
        }

        const res = await postRequest<OkResponse>('/admin/options', 'admin', {
            auto_switch_prop: autoSwitchProp,
        });
        if (res.status !== 200 || res.data?.ok !== 1) {
            errorAlert(res);
            return;
        }

        alert('Auto switch proportion updated!');
    };

    const updateGroupSizes = async () => {
        const sizes = groupSizes.split(',').map((size) => parseInt(size.trim()));
        if (sizes.some((size) => isNaN(size) || size < 1)) {
            alert('Group sizes should be positive integers!');
            return;
        }

        // Make sure number of counts is numGroups
        if (sizes.length !== numGroups) {
            alert(
                'Number of sizes should be the same as number of groups (you need a size for each group)!'
            );
            return;
        }

        const res = await postRequest<OkResponse>('/admin/group-sizes', 'admin', {
            group_sizes: sizes,
        });
        if (res.status !== 200 || res.data?.ok !== 1) {
            errorAlert(res);
            return;
        }

        alert('Group sizes updated!');
    };

    const updateGroupNames = async () => {
        const names = groupNames.split(',').map((name) => name.trim());
        if (names.some((name) => name === '')) {
            alert('Group names should not be empty!');
            return;
        }

        // Make sure number of names is numGroups
        if (names.length !== numGroups) {
            alert(
                'Number of names should be the same as number of groups (you need a name for each group)!'
            );
            return;
        }

        const res = await postRequest<OkResponse>('/admin/options', 'admin', {
            group_names: names,
        });
        if (res.status !== 200 || res.data?.ok !== 1) {
            errorAlert(res);
            return;
        }

        alert('Group names updated!');
    };

    const updateIgnoreTracks = async () => {
        const tracks = ignoreTracks.split(',').map((track) => track.trim());
        if (tracks.some((track) => track === '')) {
            alert('Track names should not be empty!');
            return;
        }

        const res = await postRequest<OkResponse>('/admin/options', 'admin', {
            ignore_tracks: tracks,
        });
        if (res.status !== 200 || res.data?.ok !== 1) {
            errorAlert(res);
            return;
        }

        alert('Ignored tracks updated!');
    };

    const resetClock = async () => {
        const res = await postRequest<OkResponse>('/admin/clock/reset', 'admin', null);
        if (res.status !== 200 || res.data?.ok !== 1) {
            errorAlert(res);
            return;
        }

        alert('Clock reset!');
        setClockResetPopup(false);
    };

    const toggleSyncClock = async () => {
        const res = await postRequest<OkResponse>('/admin/options', 'admin', {
            clock_sync: !syncClock,
        });
        if (res.status !== 200 || res.data?.ok !== 1) {
            errorAlert(res);
            return;
        }

        setSyncClock(!syncClock);
    };

    const backupClock = async () => {
        const res = await postRequest<OkResponse>('/admin/clock/backup', 'admin', null);
        if (res.status !== 200 || res.data?.ok !== 1) {
            errorAlert(res);
            return;
        }

        alert('Clock data backed up!');
    };

    const reset = async (type: string) => {
        const res = await postRequest<OkResponse>('/admin/reset', 'admin', { type });
        if (res.status !== 200 || res.data?.ok !== 1) {
            errorAlert(res);
            return;
        }

        alert(`Reset ${type}`);
        setResetPopup({ open: false, text: '', handler: null });
    };

    const exportCsv = async (type: string) => {
        const res = await fetch(`${import.meta.env.VITE_JURY_URL}/admin/export/${type}`, {
            method: 'GET',
            headers: createHeaders('admin', false),
        });

        if (res.status !== 200) {
            const error = 'Error exporting data: ' + res.statusText;
            alert(error);
            console.error(error);
            return;
        }

        saveToFile((await res.blob()) as Blob, type, 'csv');
    };

    const exportByChallenge = async () => {
        const res = await fetch(`${import.meta.env.VITE_JURY_URL}/admin/export/challenges`, {
            method: 'GET',
            headers: createHeaders('admin', false),
        });

        if (res.status !== 200) {
            const error = 'Error exporting data: ' + res.statusText;
            alert(error);
            console.error(error);
            return;
        }

        saveToFile((await res.blob()) as Blob, 'challenge-projects', 'zip');
    };

    const saveToFile = (blob: Blob, name: string, ext: string) => {
        const url = window.URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.setAttribute('download', `${name}.${ext}`);
        document.body.appendChild(a);
        a.click();
        a.remove();
    };

    return (
        <>
            <JuryHeader withBack withLogout isAdmin />
            <div className="flex flex-col items-start justify-center w-full px-8 py-4 md:px-16 md:py-8">
                <h1 className="text-4xl font-bold">Settings</h1>

                <div className="my-4 flex flex-row flex-wrap gap-1 md:gap-2">
                    <NavButton>Judge Login</NavButton>
                    <NavButton>Judging Parameters</NavButton>
                    <NavButton>Judging Clock and Timer</NavButton>
                    <NavButton>Track Judging</NavButton>
                    <NavButton>Multi-Group Judging</NavButton>
                    <NavButton>Export Data</NavButton>
                    <NavButton>Reset Database</NavButton>
                </div>

                {running && (
                    <p className="p-2 text-center border-gold border-2 bg-gold/20 rounded-md">
                        Judging has started. Some settings have been disabled to ensure consistency
                        with judging assignments and rank aggregations. Please pause judging to
                        change the disabled settings.
                    </p>
                )}

                <Card>
                    <Section>Judge Login</Section>

                    <SubSection>Disable Logins</SubSection>
                    <Description>
                        Disable ALL judge login endpoints. Judges who are currently logged in can
                        still judge, but no new judges can log in. This is useful if you want to
                        prevent brute force enumeration attacks.
                    </Description>
                    <div className="flex flex-col md:flex-row items-center md:gap-4">
                        <SettingsButton
                            onClick={updateBlockReqs}
                            type={blockReqs ? 'primary' : 'error'}
                        >
                            {blockReqs ? 'Enable Logins' : 'Disable Logins'}
                        </SettingsButton>
                        <p
                            className={twMerge(
                                'text-xl',
                                blockReqs ? 'text-error' : 'text-primary'
                            )}
                        >
                            {blockReqs ? 'Logins currently BLOCKED' : 'Logins currently ENABLED'}
                        </p>
                    </div>

                    <SubSection>Max Logins Per Minute</SubSection>
                    <Description>
                        Set the maximum number of requests that can be made to the login endpoint
                        per minute. This is also useful for preventing brute force attacks.
                    </Description>
                    <FieldButton>
                        <TextInput
                            text={maxReqPerMin}
                            setText={setMaxReqPerMin}
                            label="Max Logins Per Minute"
                            placeholder="Enter an integer..."
                            large
                            className="my-2 mr-4"
                            number
                        />
                        <SettingsButton onClick={updateMaxReqPerMin}>
                            Update Max Requests
                        </SettingsButton>
                    </FieldButton>
                </Card>

                <Card>
                    <Section>Judging Parameters</Section>

                    <SubSection>Reassign Project Numbers</SubSection>
                    <Description>
                        Reassign all project numbers to the projects. This will keep the relative
                        order but reassign the project numbers starting from the first project. If
                        groups are enabled, projects will be assigned round-robin into groups.
                        Otherwise, projects will be assigned in order.
                    </Description>
                    <FieldButton>
                        <SettingsButton
                            onClick={() => setReassignPopup(true)}
                            disabled={running}
                            type="gold"
                        >
                            Reassign
                        </SettingsButton>
                    </FieldButton>

                    <SubSection>Reassign Judge Groups</SubSection>
                    <Description>
                        Reassigns all judges to groups. The number of judges in each group will be
                        kept as even as possible.
                    </Description>
                    <SettingsButton onClick={() => setJudgeReassignPopup(true)} type="gold">
                        Reassign
                    </SettingsButton>

                    <SubSection>Set Minimum Project Views</SubSection>
                    <Description>
                        Set the minimum amount of times that a project should be seen during
                        judging. This will ensure all projects get seen at LEAST this many times
                        before switching over to the optimal method of assigning projects. Set to 0
                        to ignore this condition (recommended: 3-5).
                    </Description>
                    <FieldButton>
                        <TextInput
                            text={minViews}
                            setText={setMinViews}
                            label="Minimum Views"
                            placeholder="Enter an integer..."
                            large
                            className="my-2 mr-4"
                            number
                        />
                        <SettingsButton onClick={updateMinViews}>Update Min Views</SettingsButton>
                    </FieldButton>

                    <SubSection>Ignore Tracks</SubSection>
                    <Description>
                        Set a list of tracks to ignore when uploading projects. This will be applied
                        when uploading projects -- projects that are ignored will NOT be added to
                        Jury at all. This is most effective when using Devpost CSV upload.
                    </Description>
                    <ChallengeBlock />
                    <TextInput
                        text={ignoreTracks}
                        setText={setIgnoreTracks}
                        label="Tracks to Ignore"
                        placeholder="Track 1, Track 2, ..."
                        large
                        full
                        className="my-2"
                    />
                    <SettingsButton onClick={updateIgnoreTracks} type="gold">
                        Update Ignore Tracks
                    </SettingsButton>
                </Card>

                <Card>
                    <Section>Judging Clock and Timer</Section>

                    <SubSection>Reset Main Clock</SubSection>
                    <Description>Reset the clock back to 00:00:00</Description>
                    <SettingsButton
                        onClick={() => setClockResetPopup(true)}
                        disabled={running}
                        type="gold"
                    >
                        Reset
                    </SettingsButton>

                    <SubSection>Sync Clock with Database Automatically</SubSection>
                    <Description>
                        Backup clock data to database automatically -- note that this functionality
                        won't work properly if you are using the database with multiple instances of
                        Jury. This will save the clock state if the Jury backend server crashes for
                        any reason. Automatic means every 10 minutes, the clock state will be saved
                        to the database, as well as whenever the clock is paused/unpaused/reset.
                    </Description>
                    <Checkbox checked={syncClock} onChange={toggleSyncClock}>
                        Sync automatically
                    </Checkbox>

                    <SubSection>Backup Clock</SubSection>
                    <Description>Save clock state to database right now.</Description>
                    <SettingsButton onClick={backupClock}>Backup</SettingsButton>

                    <SubSection>Set Judging Timer</SubSection>
                    <Description>
                        Set how long judges have to view each project. This will reflect on the
                        timer that shows on the judging page. Leave this field blank (or 0) if you
                        do not wish to have a timer for each judge.
                    </Description>
                    <FieldButton>
                        <TextInput
                            text={judgingTimer}
                            setText={setJudgingTimer}
                            label="Judging Timer (MM:SS)"
                            placeholder="MM:SS"
                            large
                            className="my-2 mr-4"
                        />
                        <SettingsButton onClick={updateTimer}>Update Timer</SettingsButton>
                    </FieldButton>
                </Card>

                <Card>
                    <Section>Track Judging</Section>

                    <SubSection>Enable Track Judging</SubSection>
                    <Description>
                        Enable judging for individual tracks. This allows judges to be assigned to a
                        specific track and ONLY judge projects that have submitted to that track.
                        This does not affect the main general judging, and all projects will still
                        be considered for the main judging. However, if you have a different track
                        (eg. design), this feature will allow that track's judges to only see
                        projects that have submitted to that track.
                    </Description>
                    <Checkbox checked={judgeTracks} onChange={toggleJudgeTracks} disabled={running}>
                        Enable Track Judging
                    </Checkbox>

                    {judgeTracks && (
                        <>
                            <SubSection>Set Tracks</SubSection>
                            <Description>
                                Set the tracks that will be judged. Note this should match the name
                                under the 'Opt-In Prizes' category. Only the tracks listed here will
                                be judged!
                            </Description>
                            <ChallengeBlock />
                            <TextInput
                                text={tracks}
                                setText={setTracks}
                                label="List of Tracks"
                                placeholder="Track 1, Track 2, ..."
                                large
                                full
                                className="my-2"
                            />
                            <SettingsButton disabled={running} onClick={updateTracks}>
                                Update Tracks
                            </SettingsButton>

                            <SubSection>Set Track Min Views</SubSection>
                            <Description>
                                Set the number of views a track project should get across all
                                judges. For example, if this is set to 3, then a project will be
                                seen exactly 3 times by that track's judges. We suggest this to be
                                set to around 2-4, depending on the size of the track.
                            </Description>
                            <TextInput
                                label="Track min views"
                                placeholder="3, 3, 2, ..."
                                text={trackViews}
                                setText={setTrackViews}
                                full
                                large
                                className="my-2 mr-4"
                            />
                            <SettingsButton onClick={updateTrackViews}>
                                Update Track Min Views
                            </SettingsButton>
                        </>
                    )}
                </Card>
                <Card>
                    <Section>Multi-group Judging</Section>

                    <SubSection>Enable Multiple Groups</SubSection>
                    <Description>
                        Enable multiple groups for judging. This will allow judges to be assigned to
                        different groups and score projects within their group only, switching after
                        either a certain number of projects or manually by admins. Note that
                        toggling this option{' '}
                        <span className="font-bold">will reassign project numbers</span>.
                    </Description>
                    <Checkbox checked={multiGroup} onChange={toggleMultiGroup} disabled={running}>
                        Enable Multi-Group Judging
                    </Checkbox>

                    {multiGroup && (
                        <>
                            <SubSection>Number of Groups</SubSection>
                            <Description>
                                Set the number of groups judges will be split into.
                            </Description>
                            <FieldButton>
                                <TextInput
                                    text={numGroups}
                                    setText={setNumGroups}
                                    label="Number of Groups"
                                    placeholder="Enter an integer..."
                                    large
                                    number
                                    className="my-2 mr-4"
                                />
                                <SettingsButton onClick={updateNumGroups} disabled={running}>
                                    Update Number of Groups
                                </SettingsButton>
                            </FieldButton>

                            <SubSection>Group Sizes</SubSection>
                            <Description>
                                Set how many projects each group will get. Please separate each
                                count with a comma. This list should be exactly the same length as
                                the number of groups (eg. 3 groups: "30, 30, 40"). Note that the
                                last group will be used as overflow if all groups fill up.
                            </Description>
                            <TextInput
                                label="Group Sizes"
                                placeholder="30, 30, 30, ..."
                                text={groupSizes}
                                setText={setGroupSizes}
                                full
                                large
                                className="my-2 mr-4"
                            />
                            <SettingsButton onClick={updateGroupSizes} disabled={running}>Update Sizes</SettingsButton>

                            <SubSection>Switching Mode</SubSection>
                            <Description>
                                Choose how judges will switch between projects. If set to "auto",
                                judges will switch after a certain number/proportion of projects. If
                                set to "manual", judges will switch after an admin manually presses
                                a button to switch judges' groups.
                            </Description>
                            <SelectionButton
                                selected={switchingMode}
                                setSelected={setSwitchingMode}
                                onClick={updateSwitchingMode}
                                options={['auto', 'manual']}
                                className="my-2"
                            />

                            {switchingMode === 'auto' && (
                                <>
                                    <SubSection>Auto Switch Proportions</SubSection>
                                    <Description>
                                        Set the proportion of projects judges will view before
                                        switching groups. This should be a decimal between 0 and 1.
                                        All numbers will be rounded up, with a minimum of 1 view per
                                        group. Note that too low of values may cause judges to
                                        experience a marathon, especially if your rooms/groups are
                                        far apart. However, too high of values may bias the
                                        aggregated ranking results if judges do not visit enough
                                        different groups.
                                    </Description>
                                    <FieldButton>
                                        <TextInput
                                            text={autoSwitchProp}
                                            setText={setAutoSwitchProp}
                                            label="Auto Switch Proportions"
                                            placeholder="Enter a decimal..."
                                            large
                                            number
                                            className="my-2 mr-4"
                                        />
                                        <SettingsButton onClick={updateAutoSwitchProp}>
                                            Update Proportion
                                        </SettingsButton>
                                    </FieldButton>
                                </>
                            )}

                            <SubSection>Group Names</SubSection>
                            <Description>
                                Set the names of the groups. This will be displayed to judges and
                                admins to help them keep track of which group they are in. This is
                                especially useful if you have multiple rooms or groups of judges.
                            </Description>
                            <TextInput
                                label="Group Names"
                                placeholder="Group 1, Group 2, Group 3, ..."
                                text={groupNames}
                                setText={setGroupNames}
                                full
                                large
                                className="my-2"
                            />
                            <SettingsButton onClick={updateGroupNames}>
                                Update Group Names
                            </SettingsButton>
                        </>
                    )}
                </Card>

                <Card>
                    <Section>Export Data</Section>

                    <SubSection>Export Collection</SubSection>
                    <Description>
                        Export each collection individually as a CSV download.
                    </Description>
                    <div className="flex flex-col lg:flex-row mt-2 lg:mt-0 gap-4">
                        <SettingsButton
                            onClick={() => exportCsv('judges')}
                            className="my-0 lg:my-2"
                        >
                            Export Judges
                        </SettingsButton>
                        <SettingsButton
                            onClick={() => exportCsv('projects')}
                            className="my-0 lg:my-2"
                        >
                            Export Projects
                        </SettingsButton>
                        <SettingsButton onClick={exportByChallenge} className="my-0 lg:my-2">
                            Export by Challenges
                        </SettingsButton>
                        <SettingsButton
                            onClick={() => exportCsv('rankings')}
                            className="my-0 lg:my-2"
                        >
                            Export Rankings
                        </SettingsButton>
                    </div>
                </Card>

                <Card>
                    <Section>Reset Data</Section>

                    <SubSection>Clear Rankings and Stars</SubSection>
                    <Description>
                        Removes ALL rankings that judges have made. This will not change anything
                        else but remove all judges' rankings and stars for projects.
                    </Description>
                    <SettingsButton
                        onClick={() =>
                            setResetPopup({
                                open: true,
                                text: 'Are you sure you want to delete all ranking data? This includes rankings and stars of ALL judges!',
                                handler: reset.bind(this, 'rankings'),
                            })
                        }
                        type="error"
                        disabled={running}
                    >
                        Delete Rankings
                    </SettingsButton>

                    <SubSection>Clear Judging Data</SubSection>
                    <Description>
                        Remove all active judging data. This will keep judges, projects, and
                        configurations; essentially it will revert Jury to a state before judging
                        started. This is useful if you accidentally started judging and want to
                        revert all judging without removing projects or judges.
                    </Description>
                    <SettingsButton
                        onClick={() =>
                            setResetPopup({
                                open: true,
                                text: 'Are you sure you want to delete all judging data? This will wipe all judging activity but retain projects/judges/settings.',
                                handler: reset.bind(this, 'judging-data'),
                            })
                        }
                        type="error"
                        disabled={running}
                    >
                        Delete Judging Data
                    </SettingsButton>

                    <SubSection>Delete all Projects</SubSection>
                    <Description>
                        This will delete all projects in the database, along with any judging and
                        flag data.
                    </Description>
                    <SettingsButton
                        onClick={() =>
                            setResetPopup({
                                open: true,
                                text: 'Are you sure you want to delete all projects? This will wipe everything except for judges and settings.',
                                handler: reset.bind(this, 'projects'),
                            })
                        }
                        type="error"
                        disabled={running}
                    >
                        Delete Projects
                    </SettingsButton>

                    <SubSection>Delete all Judges</SubSection>
                    <Description>
                        This will delete all judges in the database, along with any judging and flag
                        data.
                    </Description>
                    <SettingsButton
                        onClick={() =>
                            setResetPopup({
                                open: true,
                                text: 'Are you sure you want to delete all judges? This will wipe everything except for projects and settings.',
                                handler: reset.bind(this, 'judges'),
                            })
                        }
                        type="error"
                        disabled={running}
                    >
                        Delete Judges
                    </SettingsButton>

                    <SubSection>THIS WILL DELETE THE ENTIRE DATABASE</SubSection>
                    <Description>
                        Mostly used for testing purposes/before the event if you want to reset
                        everything bc something got messed up. Do NOT use this during judging (duh).
                    </Description>
                    <SettingsButton
                        onClick={() => {
                            setResetPopup({
                                open: true,
                                text: 'THIS WILL ACTUALLY DELETE ALL DATA!!!!! YOU NEED TO BE ABSOLUTELY SURE YOU WANT TO DO THIS. THIS IS YOUR LAST WARNING! Deleting the entire database will not only remove judges/projects/judging stats but also remove all configuration.',
                                handler: reset.bind(this, 'all'),
                            });
                        }}
                        type="error"
                        disabled={running}
                    >
                        Drop Database
                    </SettingsButton>
                </Card>
            </div>
            <ConfirmPopup
                enabled={reassignPopup}
                setEnabled={setReassignPopup}
                onSubmit={reassignTables}
                submitText="Reassign"
                title="Heads Up!"
                red
            >
                Are you sure you want to reassign project numbers? This should NOT be done DURING
                judging; only beforehand!!
            </ConfirmPopup>
            <ConfirmPopup
                enabled={judgeReassignPopup}
                setEnabled={setJudgeReassignPopup}
                onSubmit={reassignJudges}
                submitText="Reassign"
                title="Heads Up!"
                red
            >
                Are you sure you want to judge group numbers? This should NOT be done DURING
                judging; only beforehand!!
            </ConfirmPopup>
            <ConfirmPopup
                enabled={clockResetPopup}
                setEnabled={setClockResetPopup}
                onSubmit={resetClock}
                submitText="Reset"
                title="Heads Up!"
                red
            >
                Are you sure you want to reset the main clock? This will reset the main clock to
                00:00:00.
            </ConfirmPopup>
            <ConfirmPopup
                enabled={resetPopup.open}
                setEnabled={((e: boolean) => setResetPopup({ ...resetPopup, open: e })) as any}
                onSubmit={resetPopup.handler as any}
                submitText="RESET DATA"
                title="Heads Up!"
                red
            >
                {resetPopup.text}
            </ConfirmPopup>
            <Loading disabled={!loading} />
            <ToTopButton />
        </>
    );
};

export default AdminSettings;

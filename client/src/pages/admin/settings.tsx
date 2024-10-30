import { useEffect, useState } from 'react';
import { createHeaders, getRequest, postRequest } from '../../api';
import Button from '../../components/Button';
import JuryHeader from '../../components/JuryHeader';
import { errorAlert } from '../../util';
import ConfirmPopup from '../../components/ConfirmPopup';
import Loading from '../../components/Loading';
import Checkbox from '../../components/Checkbox';
import RawTextInput from '../../components/RawTextInput';
import SelectionButton from '../../components/SelectionButton';
import { twMerge } from 'tailwind-merge';
import { useOptionsStore } from '../../store';

// Text components
const Section = ({ children: c }: { children: React.ReactNode }) => (
    <h2 className="text-4xl mt-8 text-primary">{c}</h2>
);
const SubSection = ({ children: c }: { children: React.ReactNode }) => (
    <h3 className="text-xl mt-4 mb-1 font-bold">{c}</h3>
);
const Description = ({ children: c }: { children: React.ReactNode }) => (
    <p className="text-light mb-[0.125rem]">{c}</p>
);

// Custom button to use on settings page
const SettingsButton = ({
    children,
    onClick,
    className,
    type = 'primary',
}: {
    children: React.ReactNode;
    onClick: () => void;
    className?: string;
    type?: 'primary' | 'error' | 'gold';
}) => (
    <Button type={type} onClick={onClick} small className={'my-2 ' + className}>
        {children}
    </Button>
);

const AdminSettings = () => {
    const [reassignPopup, setReassignPopup] = useState(false);
    const [judgeReassignPopup, setJudgeReassignPopup] = useState(false);
    const [clockResetPopup, setClockResetPopup] = useState(false);
    const [dropPopup, setDropPopup] = useState(false);
    const [judgingTimer, setJudgingTimer] = useState('');
    const [minViews, setMinViews] = useState(3);
    const [syncClock, setSyncClock] = useState(false);
    const [categories, setCategories] = useState('');
    const [loading, setLoading] = useState(true);
    const [multiGroup, setMultiGroup] = useState(false);
    const [numGroups, setNumGroups] = useState(3);
    const [switchingMode, setSwitchingMode] = useState('auto');
    const [autoSwitchMethod, setAutoSwitchMethod] = useState('count');
    const [autoSwitchCount, setAutoSwitchCount] = useState(3);
    const [autoSwitchProp, setAutoSwitchProp] = useState(0.1);
    const [groupSizes, setGroupSizes] = useState('30, 30');
    const fetchOptions = useOptionsStore((state) => state.fetchOptions);

    async function getOptions() {
        const res = await getRequest<Options>('/admin/options', 'admin');
        if (res.status !== 200) {
            errorAlert(res);
            return;
        }
        if (!res.data) {
            alert('error: could not get options data');
            return;
        }

        // Calculate judging timer MM:SS
        const timer = res.data.judging_timer;
        if (timer) {
            const minutes = Math.floor(timer / 60);
            const seconds = timer % 60;
            const timerStr = `${minutes}:${seconds < 10 ? '0' : ''}${seconds}`;
            setJudgingTimer(timerStr);
        }

        // Set categories
        const cats = res.data.categories.join(', ');
        setCategories(cats ?? '');

        // Set min views
        setMinViews(res.data.min_views);

        // Set sync clock
        setSyncClock(res.data.clock_sync);

        // Set group options
        setMultiGroup(res.data.multi_group);
        setNumGroups(res.data.num_groups);
        setGroupSizes(res.data.group_sizes.join(', '));
        setSwitchingMode(res.data.main_group.switching_mode);
        setAutoSwitchMethod(res.data.main_group.auto_switch_method);
        setAutoSwitchCount(res.data.main_group.auto_switch_count);
        setAutoSwitchProp(res.data.main_group.auto_switch_prop);

        setLoading(false);
    }

    // Get the previous options
    useEffect(() => {
        getOptions();
    }, []);

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
    }

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
        const res = await postRequest<OkResponse>('/admin/timer', 'admin', {
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
        const res = await postRequest<OkResponse>('/admin/min-views', 'admin', {
            min_views: minViews,
        });
        if (res.status !== 200 || res.data?.ok !== 1) {
            errorAlert(res);
            return;
        }

        alert('Min views updated!');
        getOptions();
    };

    const updateCategories = async () => {
        // Split categories by comma and remove empty strings
        const filteredCats = categories
            .split(',')
            .map((cat) => cat.trim())
            .filter((cat) => cat !== '');

        // Post the new categories
        const res = await postRequest<OkResponse>('/admin/categories', 'admin', {
            categories: filteredCats,
        });
        if (res.status !== 200 || res.data?.ok !== 1) {
            errorAlert(res);
            return;
        }

        alert('Categories updated!');
        getOptions();
    };

    const toggleMultiGroup = async () => {
        const res = await postRequest<OkResponse>('/admin/groups/toggle', 'admin', {
            multi_group: !multiGroup,
        });
        if (res.status !== 200 || res.data?.ok !== 1) {
            errorAlert(res);
            return;
        }

        setMultiGroup(!multiGroup);
        alert('Multi-group judging toggled!');
    };

    const updateNumGroups = async () => {
        if (isNaN(numGroups) || numGroups < 1) {
            alert('Number of groups should be an integer greater than 0!');
            return;
        }

        const res = await postRequest<OkResponse>('/admin/groups/num', 'admin', {
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
        const res = await postRequest<OkResponse>('/admin/groups/options', 'admin', {
            switching_mode: newMode,
        });
        if (res.status !== 200 || res.data?.ok !== 1) {
            errorAlert(res);
            return;
        }

        fetchOptions();
        alert('Switching method updated!');
    };

    const updateAutoSwitchMethod = async (newMethod: string) => {
        const res = await postRequest<OkResponse>('/admin/groups/options', 'admin', {
            auto_switch_method: newMethod,
        });
        if (res.status !== 200 || res.data?.ok !== 1) {
            errorAlert(res);
            return;
        }

        alert('Auto switch method updated!');
    };

    const updateAutoSwitchCount = async () => {
        if (isNaN(autoSwitchCount) || autoSwitchCount < 1) {
            alert('Auto switch count should be a positive integer!');
            return;
        }

        const res = await postRequest<OkResponse>('/admin/groups/options', 'admin', {
            auto_switch_count: autoSwitchCount,
        });
        if (res.status !== 200 || res.data?.ok !== 1) {
            errorAlert(res);
            return;
        }

        alert('Auto switch count updated!');
    };

    const updateAutoSwitchProp = async () => {
        if (isNaN(autoSwitchProp) || autoSwitchProp < 0 || autoSwitchProp > 1) {
            alert('Auto switch proportion should be a decimal between 0 and 1!');
            return;
        }

        const res = await postRequest<OkResponse>('/admin/groups/options', 'admin', {
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

        // Make sure number of counts is numGroups - 1
        if (sizes.length !== numGroups - 1) {
            alert('Number of sizes should be one less than the number of groups!');
            return;
        }

        const res = await postRequest<OkResponse>('/admin/groups/options', 'admin', {
            group_sizes: sizes,
        });
        if (res.status !== 200 || res.data?.ok !== 1) {
            errorAlert(res);
            return;
        }

        alert('Split counts updated!');
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
        const res = await postRequest<OkResponse>('/admin/clock/sync', 'admin', {
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

    const dropDatabase = async () => {
        const res = await postRequest<OkResponse>('/admin/reset', 'admin', null);
        if (res.status !== 200 || res.data?.ok !== 1) {
            errorAlert(res);
            return;
        }

        alert('Database reset!');
        setDropPopup(false);
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
                <Section>Judging Parameters</Section>

                <SubSection>Reassign Project Numbers</SubSection>
                <Description>
                    Reassign all project numbers to the projects. This will keep the relative order
                    but reassign the project numbers starting from the first project. If groups are
                    enabled, projects will be assigned round-robin into groups. Otherwise, projects
                    will be assigned in order.
                </Description>
                <SettingsButton
                    onClick={() => setReassignPopup(true)}
                    className="bg-gold text-black hover:bg-goldDark hover:text-black"
                >
                    Reassign
                </SettingsButton>

                <SubSection>Reassign Judge Groups</SubSection>
                <Description>
                    Reassigns all judges to groups. The number of judges in each group will be kept
                    as even as possible.
                </Description>
                <SettingsButton onClick={() => setJudgeReassignPopup(true)} type="gold">
                    Reassign
                </SettingsButton>

                <SubSection>Set Minimum Project Views</SubSection>
                <Description>
                    Set the minimum amount of times that a project should be seen during judging.
                    This will ensure all projects get seen at LEAST this many times before switching
                    over to the optimal method of assigning projects. Set to 0 to ignore this
                    condition (recommended: 3-5).
                </Description>
                <div className="flex flex-row">
                    <RawTextInput
                        name="min-views"
                        text={minViews}
                        setText={setMinViews}
                        placeholder="Enter an integer..."
                        large
                        className="my-2 mr-4"
                        number
                    />
                    <SettingsButton onClick={updateMinViews}>Update Min Views</SettingsButton>
                </div>

                <SubSection>Set Categories</SubSection>
                <Description>
                    Set the categories that the judges will be scoring each project on. Please
                    separate each category with a comma.
                </Description>
                <RawTextInput
                    name="categories"
                    placeholder="Cat 1, Cat 2, Cat 3, ..."
                    text={categories}
                    setText={setCategories}
                    full
                    large
                    className="my-2"
                />
                <SettingsButton onClick={updateCategories}>Update Categories</SettingsButton>

                <Section>Judging Clock and Timer</Section>

                <SubSection>Reset Main Clock</SubSection>
                <Description>Reset the clock back to 00:00:00</Description>
                <SettingsButton
                    onClick={() => setClockResetPopup(true)}
                    className="bg-gold text-black hover:bg-goldDark hover:text-black"
                >
                    Reset
                </SettingsButton>

                <SubSection>Sync Clock with Database Automatically</SubSection>
                <Description>
                    Backup clock data to database automatically -- note that this functionality
                    won't work properly if you are using the database with multiple instances of
                    Jury. This will save the clock state if the Jury backend server crashes for any
                    reason. Automatic means every 10 minutes, the clock state will be saved to the
                    database, as well as whenever the clock is paused/unpaused/reset.
                </Description>
                <Checkbox checked={syncClock} onChange={toggleSyncClock}>
                    Sync automatically
                </Checkbox>

                <SubSection>Backup Clock</SubSection>
                <Description>Save clock state to database right now.</Description>
                <SettingsButton onClick={backupClock}>Backup</SettingsButton>

                <SubSection>Set Judging Timer</SubSection>
                <Description>
                    Set how long judges have to view each project. This will reflect on the timer
                    that shows on the judging page. Leave this field blank (or 0) if you do not wish
                    to have a timer for each judge.
                </Description>
                <div className="flex flex-row">
                    <RawTextInput
                        name="judging-timer"
                        text={judgingTimer}
                        setText={setJudgingTimer}
                        placeholder="MM:SS"
                        large
                        className="my-2 mr-4"
                    />
                    <SettingsButton onClick={updateTimer}>Update Timer</SettingsButton>
                </div>

                <Section>Multi-Group Judging</Section>

                <SubSection>Enable Multiple Groups</SubSection>
                <Description>
                    Enable multiple groups for judging. This will allow judges to be assigned to
                    different groups and score projects within their group only, switching after
                    either a certain number of projects or manually by admins. Note that toggling
                    this option <span className="font-bold">will reassign project numbers</span>.
                </Description>
                <Checkbox checked={multiGroup} onChange={toggleMultiGroup}>
                    Enable Multi-Group Judging
                </Checkbox>

                {multiGroup && (
                    <>
                        <SubSection>Number of Groups</SubSection>
                        <Description>
                            Set the number of groups judges will be split into.
                        </Description>
                        <div className="flex flex-row">
                            <RawTextInput
                                name="num-groups"
                                text={numGroups}
                                setText={setNumGroups}
                                placeholder="Enter an integer..."
                                large
                                number
                                className="my-2 mr-4"
                            />
                            <SettingsButton onClick={updateNumGroups}>
                                Update Number of Groups
                            </SettingsButton>
                        </div>

                        <SubSection>Group Sizes</SubSection>
                        <Description>
                            Set how many projects each group will get. Please separate each count
                            with a comma. This list should contain one less than the total number of
                            groups (eg. 3 groups: "30, 30").
                        </Description>
                        <RawTextInput
                            name="group-sizes"
                            placeholder="30, 30, 30, ..."
                            text={groupSizes}
                            setText={setGroupSizes}
                            full
                            large
                            className="my-2 mr-4"
                        />
                        <SettingsButton onClick={updateGroupSizes}>Update Sizes</SettingsButton>

                        <SubSection>Switching Mode</SubSection>
                        <Description>
                            Choose how judges will switch between projects. If set to "auto", judges
                            will switch after a certain number/proportion of projects. If set to
                            "manual", judges will switch after an admin manually presses a button to
                            switch judges' groups.
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
                                <SubSection>Auto Switch Method</SubSection>
                                <Description>
                                    Choose when judges will automatically switch between projects.
                                    If set to "count", judges will switch after viewing a specific
                                    number of projects in the group. If set to "proportion", judges
                                    will switch after a certain proportion of projects in the group.
                                    Note that when using "count" switching, rooms must contain
                                    approximately the same number of projects.
                                </Description>
                                <SelectionButton
                                    selected={autoSwitchMethod}
                                    setSelected={setAutoSwitchMethod}
                                    onClick={updateAutoSwitchMethod}
                                    options={['count', 'proportion']}
                                    className="my-2"
                                />

                                {autoSwitchMethod === 'count' ? (
                                    <>
                                        <SubSection>Auto Switch Count</SubSection>
                                        <Description>
                                            Set how many projects judges will view before switching
                                            groups.
                                        </Description>
                                        <div className="flex flex-row">
                                            <RawTextInput
                                                name="auto-switch-count"
                                                text={autoSwitchCount}
                                                setText={setAutoSwitchCount}
                                                placeholder="Enter an integer..."
                                                large
                                                number
                                                className="my-2 mr-4"
                                            />
                                            <SettingsButton onClick={updateAutoSwitchCount}>
                                                Update Count
                                            </SettingsButton>
                                        </div>
                                    </>
                                ) : (
                                    <>
                                        <SubSection>Auto Switch Proportion</SubSection>
                                        <Description>
                                            Set the proportion of projects judges will view before
                                            switching groups. This should be a decimal between 0 and
                                            1. All numbers will be rounded up, with a minimum of 1
                                            view per group. Note that too low of values may cause
                                            judges to experience a marathon, especially if your
                                            rooms/groups are far apart. However, too high of values
                                            may bias the aggregated ranking results if judges do not
                                            visit enough different groups.
                                        </Description>
                                        <div className="flex flex-row">
                                            <RawTextInput
                                                name="auto-switch-prop"
                                                text={autoSwitchProp}
                                                setText={setAutoSwitchProp}
                                                placeholder="Enter a decimal..."
                                                large
                                                number
                                                className="my-2 mr-4"
                                            />
                                            <SettingsButton onClick={updateAutoSwitchProp}>
                                                Update Proportion
                                            </SettingsButton>
                                        </div>
                                    </>
                                )}
                            </>
                        )}
                    </>
                )}

                <Section>Export Data</Section>

                <SubSection>Export Collection</SubSection>
                <Description>Export each collection individually as a CSV download.</Description>
                <div className="flex">
                    <SettingsButton onClick={() => exportCsv('users')} className="mr-4">
                        Export Judges
                    </SettingsButton>
                    <SettingsButton onClick={() => exportCsv('projects')} className="mr-4">
                        Export Projects
                    </SettingsButton>
                    <SettingsButton onClick={exportByChallenge} className="mr-4">
                        Export by Challenges
                    </SettingsButton>
                    <SettingsButton onClick={() => exportCsv('rankings')}>
                        Export Rankings
                    </SettingsButton>
                </div>

                <Section>Reset Database</Section>

                <SubSection>THIS WILL DELETE THE ENTIRE DATABASE</SubSection>
                <Description>
                    Mostly used for testing purposes/before the event if you want to reset
                    everything bc something got messed up. Do NOT use this during judging (duh).
                </Description>
                <SettingsButton onClick={() => setDropPopup(true)} type="error">
                    Drop Database
                </SettingsButton>
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
                enabled={dropPopup}
                setEnabled={setDropPopup}
                onSubmit={dropDatabase}
                submitText="RESET DATA"
                title="Heads Up!"
                red
            >
                THIS WILL ACTUALLY DELETE ALL DATA!!!!! YOU NEED TO BE ABSOLUTELY SURE YOU WANT TO
                DO THIS. THIS IS YOUR LAST WARNING!
            </ConfirmPopup>
            <Loading disabled={!loading} />
        </>
    );
};

export default AdminSettings;

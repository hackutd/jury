import { useEffect, useState } from 'react';
import { createHeaders, getRequest, postRequest } from '../../api';
import Button from '../../components/Button';
import JuryHeader from '../../components/JuryHeader';
import { errorAlert } from '../../util';
import TextPopup from '../../components/TextPopup';
import Loading from '../../components/Loading';
import Checkbox from '../../components/Checkbox';
import RawTextInput from '../../components/RawTextInput';

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
    type?: 'primary' | 'error';
}) => (
    <Button type={type} onClick={onClick} small className={'my-2 ' + className}>
        {children}
    </Button>
);

const AdminSettings = () => {
    const [reassignPopup, setReassignPopup] = useState(false);
    const [clockResetPopup, setClockResetPopup] = useState(false);
    const [dropPopup, setDropPopup] = useState(false);
    const [judgingTimer, setJudgingTimer] = useState('');
    const [minViews, setMinViews] = useState('');
    const [syncClock, setSyncClock] = useState(false);
    const [categories, setCategories] = useState('');
    const [loading, setLoading] = useState(true);

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
        setMinViews(res.data.min_views.toString());

        // Set sync clock
        setSyncClock(res.data.clock_sync);

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
        const v = parseInt(minViews);
        if (isNaN(v)) {
            alert('Minimum views should be a positive integer!');
            return;
        }

        // Update min views
        const res = await postRequest<OkResponse>('/admin/min-views', 'admin', {
            min_views: v,
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
                <Section>Judging Settings</Section>

                <SubSection>Reset Main Clock</SubSection>
                <Description>Reset the clock back to 00:00:00</Description>
                <SettingsButton onClick={() => setClockResetPopup(true)}>Reset</SettingsButton>

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

                <Section>Judging Parameters</Section>

                <SubSection>Reassign Project Numbers</SubSection>
                <Description>
                    Reassign all project numbers to the projects. This will keep the relative order
                    but reassign the project numbers starting from the first project.
                </Description>
                <SettingsButton
                    onClick={() => setReassignPopup(true)}
                    className="bg-gold text-black hover:bg-goldDark hover:text-black"
                >
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
                    />
                    <SettingsButton onClick={updateMinViews}>Update Min Views</SettingsButton>
                </div>

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
            <TextPopup
                enabled={reassignPopup}
                setEnabled={setReassignPopup}
                onSubmit={reassignTables}
                submitText="Reassign"
                title="Heads Up!"
                red
            >
                Are you sure you want to reassign project numbers? This should NOT be done DURING
                judging; only beforehand!!
            </TextPopup>
            <TextPopup
                enabled={clockResetPopup}
                setEnabled={setClockResetPopup}
                onSubmit={resetClock}
                submitText="Reset"
                title="Heads Up!"
                red
            >
                Are you sure you want to reset the main clock? This will reset the clock to 00:00:00
            </TextPopup>
            <TextPopup
                enabled={dropPopup}
                setEnabled={setDropPopup}
                onSubmit={dropDatabase}
                submitText="RESET DATA"
                title="Heads Up!"
                red
            >
                THIS WILL ACTUALLY DELETE ALL DATA!!!!! YOU NEED TO BE ABSOLUTELY SURE YOU WANT TO
                DO THIS. THIS IS YOUR LAST WARNING!
            </TextPopup>
            <Loading disabled={!loading} />
        </>
    );
};

export default AdminSettings;

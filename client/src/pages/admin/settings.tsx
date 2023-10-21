import { useEffect, useState } from 'react';
import { createHeaders, getRequest, postRequest } from '../../api';
import Button from '../../components/Button';
import JuryHeader from '../../components/JuryHeader';
import { errorAlert } from '../../util';
import Popup from '../../components/Popup';
import Checkbox from '../../components/Checkbox';
import Loading from '../../components/Loading';

const AdminSettings = () => {
    const [reassignPopup, setReassignPopup] = useState(false);
    const [groupChecked, setGroupChecked] = useState(false);
    const [groups, setGroups] = useState('');
    const [loading, setLoading] = useState(true);

    async function getOptions() {
        const res = await getRequest<Options>('/admin/options', 'admin');
        if (res.status !== 200) {
            errorAlert(res);
            return;
        }
        setGroupChecked(res.data?.use_groups ?? false);

        // Create groupings
        const groupStr = res.data?.groups?.map((g) => `${g.start} ${g.end}`).join('\n');
        setGroups(groupStr || '');

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

    const updateGroupings = async () => {
        const res = await postRequest<OkResponse>('/admin/groups', 'admin', {
            use_groups: groupChecked,
            raw_groups: groups,
        });
        if (res.status !== 200 || res.data?.ok !== 1) {
            errorAlert(res);
            return;
        }

        alert('Groups updated!');
        getOptions();
    };

    const exportData = async (type: string) => {
        const res = await fetch(`${process.env.REACT_APP_JURY_URL}/admin/export/${type}`, {
            method: 'GET',
            headers: createHeaders('admin', false),
        });

        if (res.status !== 200) {
            alert('Error exporting data: ' + res.statusText);
            return;
        }
        
        const blob = await res.blob();
        const url = window.URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.setAttribute('download', `${type}.csv`);
        document.body.appendChild(a);
        a.click();
        a.remove();
    };

    return (
        <>
            <JuryHeader withBack withLogout isAdmin />
            <div className="flex flex-col items-start justify-center w-full px-8 py-4 md:px-16 md:py-8">
                <h1 className="text-4xl font-bold">Settings</h1>
                {/* TODO: Add other settings */}
                <h2 className="text-xl mt-4">Project Numbers</h2>
                <h3 className="text-lg font-bold">Reassign Project Numbers</h3>
                <p className="text-light">
                    Reassign all table numbers to the projects. This will keep the relative order
                    but reassign the table numbers starting from the first project.
                </p>
                <Button
                    type="primary"
                    onClick={() => {
                        setReassignPopup(true);
                    }}
                    className="mt-4 w-auto md:w-auto px-4 py-2 bg-gold text-black"
                >
                    Reassign
                </Button>
                <h3 className="text-lg font-bold mt-4">Table Groupings</h3>
                <p className="text-light">
                    Check this box to use table groupings. This will force judges to stay in a
                    grouping for 3 rounds before moving on. This ideally should decrease the
                    distance judges will have to walk, if groups are defined correctly. Note that
                    group sizes <span className="font-bold">must be greater than 3</span> otherwise
                    the groupings will be ignored.
                </p>
                <Checkbox
                    checked={groupChecked}
                    onChange={(c) => {
                        setGroupChecked(c);
                    }}
                >
                    Enable Table Groupings
                </Checkbox>
                <p className="text-light mt-4">
                    {
                        'List table groupings below. It should be in the form <tableStartNum> <tableEndNum>. Each group should be on its own line. No table numbers should overlap. If groups are defined, table numbers will be assigned according to the ranges defined here (ie. groups of 1-10, 101-110 will skip table numbers 11-100). If there are more table numbers than groups defined, the default behavior is to append incrementally to the last group.'
                    }
                </p>
                <textarea
                    className="w-full h-36 px-4 py-4 text-2xl border-lightest border-2 rounded-sm focus:border-primary focus:border-4 focus:outline-none"
                    placeholder="1 10"
                    onChange={(e) => {
                        setGroups(e.target.value);
                    }}
                    value={groups}
                ></textarea>
                <Button
                    type="primary"
                    onClick={updateGroupings}
                    className="mt-4 w-auto md:w-auto px-4 py-2"
                >
                    Update Groupings
                </Button>
                <h2 className="text-xl mt-4">Export Data</h2>
                <h3 className="text-lg font-bold">Export Collections</h3>
                <p className="text-light">Export each collection individually as a CSV download.</p>
                <div className="flex">
                    <Button
                        type="primary"
                        onClick={() => {
                            exportData('judges');
                        }}
                        className="mt-4 w-auto md:w-auto px-4 py-2 mr-4"
                    >
                        Export Judges
                    </Button>
                    <Button
                        type="primary"
                        onClick={() => {
                            exportData('projects');
                        }}
                        className="mt-4 w-auto md:w-auto px-4 py-2"
                    >
                        Export Projects
                    </Button>
                </div>
            </div>
            <Popup
                enabled={reassignPopup}
                setEnabled={setReassignPopup}
                onSubmit={reassignTables}
                submitText="Reassign"
                title="Heads Up!"
                red
            >
                Are you sure you want to reassign project numbers? This should NOT be done DURING
                judging; only beforehand!!
            </Popup>
            <Loading disabled={!loading} />
        </>
    );
};

export default AdminSettings;

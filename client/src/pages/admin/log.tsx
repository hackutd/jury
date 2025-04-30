import { useEffect, useState } from 'react';
import JuryHeader from '../../components/JuryHeader';
import Loading from '../../components/Loading';
import { errorAlert } from '../../util';
import { getRequest } from '../../api';
import ToTopButton from '../../components/ToTopButton';

const AdminLog = () => {
    const [log, setLog] = useState('');
    const [loaded, setLoaded] = useState(false);

    useEffect(() => {
        async function fetchLog() {
            const res = await getRequest<Log>('/admin/log', 'admin');
            if (res.status !== 200) {
                errorAlert(res);
                return;
            }

            setLog(res.data?.log || '');
            setLoaded(true);
        }

        fetchLog();
    }, []);

    return (
        <>
            <JuryHeader isAdmin withLogout withBack />
            <div className="ml-4 mt-4">
                <h1 className="text-2xl text-primary">Audit Log</h1>
                <p className="text-light">
                    Here you can see a log of everything that has happened in the system.
                </p>
            </div>
            <div className="w-auto m-4 p-4 border-2 border-lightest rounded-md overflow-x-auto">
                <pre className="w-full text-sm md:text-md">{log}</pre>
            </div>
            <Loading disabled={loaded} />
            <ToTopButton />
        </>
    );
};

export default AdminLog;

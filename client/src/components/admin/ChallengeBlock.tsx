import { useEffect, useState } from 'react';
import { getRequest } from '../../api';
import { errorAlert } from '../../util';

const ChallengeBlock = () => {
    const [challenges, setChallenges] = useState<string[]>([]);

    useEffect(() => {
        async function fetchData() {
            // Get challenges
            const challengesRes = await getRequest<string[]>('/challenges', '');
            if (challengesRes.status !== 200) {
                errorAlert(challengesRes);
                return;
            }
            setChallenges(challengesRes.data as string[]);
        }

        fetchData();
    }, []);

    if (!challenges) return;

    return (
        <div className="border-2 border-primary bg-primary/10 rounded-md p-2 my-2">
            <h1 className="text-xl font-bold">Challenge List</h1>
            <p className="text-light">{challenges.join(',')}</p>
        </div>
    );
};

export default ChallengeBlock;

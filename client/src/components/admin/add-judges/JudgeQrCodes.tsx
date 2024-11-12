import { useEffect, useState } from 'react';
import Button from '../../Button';
import { useOptionsStore } from '../../../store';
import Dropdown from '../../Dropdown';

const JudgeQrCodes = () => {
    const [track, setTrack] = useState('');
    const options = useOptionsStore((state) => state.options);
    const fetchOptions = useOptionsStore((state) => state.fetchOptions);

    useEffect(() => {
        fetchOptions();
    }, []);

    useEffect(() => {
        if (options.tracks.length > 0) {
            setTrack(options.tracks[0]);
        }
    }, [options]);

    return (
        <div className="w-full h-full border-lightest border-2 p-8 rounded-sm">
            <div className="flex flex-col items-start h-full">
                <h1 className="text-3xl mb-4">Add Judge by QR Code</h1>
                <div className="flex flex-row w-full">
                    <Button type="primary" href="/admin/qr" full flat className="mr-4 rounded-md">
                        General Judges
                    </Button>
                    {options && options.judge_tracks && (
                        <>
                            <Button
                                type="primary"
                                href={`/admin/qr?track=${track}`}
                                full
                                flat
                                className="rounded-md"
                            >
                                {track} Track Judges
                            </Button>
                            <Dropdown
                                options={[...options.tracks]}
                                selected={track}
                                setSelected={setTrack}
                                className="ml-4"
                            />
                        </>
                    )}
                </div>
            </div>
        </div>
    );
};

export default JudgeQrCodes;

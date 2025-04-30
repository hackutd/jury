import { useEffect, useState } from 'react';
import Button from '../../Button';
import { useOptionsStore } from '../../../store';
import Dropdown from '../../Dropdown';
import Card from '../../Card';

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
        <Card>
            <div className="flex flex-col items-start h-full">
                <h1 className="text-3xl mb-4">Add Judge by QR Code</h1>
                <div className="flex flex-col lg:flex-row w-full justify-center gap-4">
                    <Button type="primary" href="/admin/qr">
                        General Judges
                    </Button>
                    {options && options.judge_tracks && (
                        <div className='flex flex-col-reverse mt-6 md:mt-0 md:flex-row justify-center lg:justify-start gap-4'>
                            <Button
                                type="primary"
                                href={`/admin/qr?track=${track}`}
                                className="rounded-md"
                            >
                                {track} Track Judges
                            </Button>
                            <Dropdown
                                options={[...options.tracks]}
                                selected={track}
                                setSelected={setTrack}
                                large
                                className="bg-white text-center"
                            />
                        </div>
                    )}
                </div>
            </div>
        </Card>
    );
};

export default JudgeQrCodes;

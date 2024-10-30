import { useState } from 'react';
import Button from '../Button';
import FlagsPopup from './FlagsPopup';
import Dropdown from '../Dropdown';
import { useOptionsStore } from '../../store';

const AdminToolbar = (props: { showProjects: boolean; lastUpdate: Date }) => {
    const [showFlags, setShowFlags] = useState(false);
    const options = useOptionsStore((state) => state.options);
    const selectedTrack = useOptionsStore((state) => state.selectedTrack);
    const setSelectedTrack = useOptionsStore((state) => state.setSelectedTrack);

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
                        selected={selectedTrack}
                        setSelected={setSelectedTrack}
                        className="ml-4"
                    />
                )}

                <p className="grow self-end text-right text-lighter">
                    Last Update: {dateToString(props.lastUpdate)}
                </p>
            </div>
            <FlagsPopup enabled={showFlags} setEnabled={setShowFlags} />
        </>
    );
};

export default AdminToolbar;

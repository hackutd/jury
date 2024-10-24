import { useState } from 'react';
import Button from '../Button';
import FlagsPopup from './FlagsPopup';

const AdminToolbar = (props: { showProjects: boolean; lastUpdate: Date }) => {
    const [showFlags, setShowFlags] = useState(false);

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
                <div>
                    <Button
                        type="outline"
                        bold
                        full
                        className="py-2 px-4 rounded-md"
                        href={props.showProjects ? '/admin/add-projects' : '/admin/add-judges'}
                    >
                        Add {props.showProjects ? 'Projects' : 'Judges'}
                    </Button>
                </div>
                <div className="ml-4">
                    {props.showProjects && (
                        <Button
                            type="outline"
                            bold
                            full
                            className="py-2 px-4 rounded-md"
                            onClick={() => {
                                setShowFlags(true);
                            }}
                        >
                            See All Flags
                        </Button>
                    )}
                </div>
                <p className="grow self-end text-right text-lighter">
                    Last Update: {dateToString(props.lastUpdate)}
                </p>
            </div>
            {showFlags && <FlagsPopup close={setShowFlags} />}
        </>
    );
};

export default AdminToolbar;

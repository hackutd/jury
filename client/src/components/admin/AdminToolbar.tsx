import { useState } from 'react';
import Button from '../Button';
import FlagsPopup from './FlagsPopup';

const AdminToolbar = (props: { showProjects: boolean }) => {
    const [showFlags, setShowFlags] = useState(false);
    return (
        <div className="flex flex-row px-8 py-4">
            <div>
                <Button
                    type="outline"
                    square
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
                        square
                        bold
                        full
                        className="py-2 px-4 rounded-md"
                        onClick={() => {
                            setShowFlags(true);
                        }}
                    >
                        See Flags
                    </Button>
                )}
            </div>
            {showFlags && <FlagsPopup close={setShowFlags} />}
        </div>
    );
};

export default AdminToolbar;

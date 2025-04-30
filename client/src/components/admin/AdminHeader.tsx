import { useNavigate } from 'react-router-dom';
import { useState } from 'react';
import Button from '../Button';
import PauseButton from './PauseButton';
import ActionsPopup from './ActionsPopup';

const AdminHeader = () => {
    const navigate = useNavigate();
    const [actionsPopup, setActionsPopup] = useState(false);

    return (
        <>
            <div className="absolute top-4 w-full flex items-center z-0">
                <PauseButton />
                <Button
                    type="outline"
                    onClick={() => {
                        navigate('/admin/settings');
                    }}
                    bold
                    className="ml-4 py-2"
                >
                    Settings
                </Button>
                <div className="grow mr-36 flex flex-row justify-end items-center">
                    <Button
                        type="outline"
                        onClick={setActionsPopup.bind(null, true)}
                        bold
                        tooltip="Open actions menu"
                        className="ml-4 py-2"
                    >
                        Actions
                    </Button>
                </div>
            </div>
            <ActionsPopup enabled={actionsPopup} setEnabled={setActionsPopup} />
        </>
    );
};

export default AdminHeader;

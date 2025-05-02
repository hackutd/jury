import { useNavigate } from 'react-router-dom';
import { useState } from 'react';
import Button from '../Button';
import PauseButton from './PauseButton';
import ActionsPopup from './ActionsPopup';

const AdminHeaderMobile = () => {
    const navigate = useNavigate();
    const [actionsPopup, setActionsPopup] = useState(false);

    return (
        <>
            <div className="w-full flex flex-col md:flex-row lg:hidden items-center gap-2 mt-4">
                <PauseButton />
                <div className="w-full flex flex-row items-center md:justify-start justify-center py-1 md:py-2 gap-4 md:pl-4">
                    <Button
                        type="outline"
                        onClick={() => {
                            navigate('/admin/settings');
                        }}
                        bold
                    >
                        Settings
                    </Button>
                    <Button
                        type="outline"
                        onClick={setActionsPopup.bind(null, true)}
                        bold
                        tooltip="Open actions menu"
                    >
                        Actions
                    </Button>
                </div>
            </div>
            <ActionsPopup enabled={actionsPopup} setEnabled={setActionsPopup} />
        </>
    );
};

export default AdminHeaderMobile;

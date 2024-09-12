import Popup from '../../Popup';
import Ratings from '../Ratings';
import { useNavigate } from 'react-router-dom';

interface RatePopupProps {
    /* Function to modify the popup state variable */
    setEnabled: React.Dispatch<React.SetStateAction<boolean>>;

    /* Judge to vote on */
    judge: Judge;

    /* State variable for determining if popup is open */
    enabled: boolean;

    /* Callback function for flagging a project */
    callback: () => Promise<void>;
}

/**
 * Component to show when the user clicks the "Submit" button
 */
const RatePopup = (props: RatePopupProps) => {
    const navigate = useNavigate();

    if (!props.enabled) return null;

    const done = async () => {
        await props.callback();
        navigate('/judge');
    };

    return (
        <Popup enabled={props.enabled} setEnabled={props.setEnabled}>
            <h1 className="text-3xl font-bold text-primary">Judge Project</h1>
            <h2 className="text-xl font-bold">Please rate the current project</h2>
            <div className="flex flex-col items-center w-full my-4">
                <Ratings callback={done} />
            </div>
        </Popup>
    );
};

export default RatePopup;

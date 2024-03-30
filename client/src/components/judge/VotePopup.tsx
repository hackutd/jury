import Ratings from './Ratings';
import { useNavigate } from 'react-router-dom';

interface VotePopupProps {
    /* Function to modify the popup state variable */
    close: React.Dispatch<React.SetStateAction<boolean>>;

    /* Judge to vote on */
    judge: Judge;

    /* State variable for determining if popup is open */
    open: boolean;

    /* Callback function for flagging a project */
    callback: () => void;
}

const VotePopup = (props: VotePopupProps) => {
    const navigate = useNavigate();

    if (!props.open) return null;

    const done = () => {
        props.callback();
        navigate('/judge');
    };

    return (
        <>
            <div
                className="fixed left-0 top-0 z-20 w-screen h-screen bg-black/30"
                onClick={() => props.close(false)}
            ></div>
            <div className="bg-background fixed z-30 left-1/2 top-1/2 translate-x-[-50%] translate-y-[-50%] px-6 py-6 md:px-10 md:w-1/3 w-11/12 flex flex-col items-center">
                <h1 className="text-3xl font-bold text-primary">Judge Project</h1>
                <h2 className="text-xl font-bold">Please rate the current project</h2>
                <div className="flex flex-col items-center w-full my-4">
                    <Ratings callback={done} />
                </div>
            </div>
        </>
    );
};

export default VotePopup;

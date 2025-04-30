import Button from '../../Button';
import Popup from '../../Popup';
import Star from '../Star';
import TextArea from '../../TextArea';

interface FinishPopupProps {
    /* Function to modify the popup state variable */
    setEnabled: React.Dispatch<React.SetStateAction<boolean>>;

    /* Judge to vote on */
    judge: Judge;

    /* State variable for determining if popup is open */
    enabled: boolean;

    /* Callback function for flagging a project */
    callback: () => Promise<void>;

    // TODO: Export all this to a global store for the judge
    /* Starred status of project */
    starred: boolean;

    /* Setter function for starred status */
    setStarred: React.Dispatch<React.SetStateAction<boolean>>;

    /* Notes for project */
    notes: string;

    /* Setter function for notes */
    setNotes: React.Dispatch<React.SetStateAction<string>>;
}

/**
 * Component to show when the user clicks the "Submit" button
 */
const FinishPopup = (props: FinishPopupProps) => {
    if (!props.enabled) return null;

    const done = async () => {
        await props.callback();
    };

    return (
        <Popup enabled={props.enabled} setEnabled={props.setEnabled} className="text-center">
            <h1 className="text-3xl font-bold text-primary">Judge Project</h1>
            <h2 className="text-xl font-bold">Finish judging this project</h2>
            <div className="flex flex-row justify-center text-left mt-2">
                <Star active={props.starred} setActive={props.setStarred} className="mr-4" />
                <p className="text-light">
                    Star projects you think should win the top places in the hackathon.
                </p>
            </div>
            <h3 className="text-lighter text-sm text-left mt-2">Personal Notes</h3>
            <TextArea
                label="Type any personal comments here"
                value={props.notes}
                setValue={props.setNotes}
                className='mt-1'
            />
            <Button type="primary" onClick={done} className="mt-4">
                Submit
            </Button>
        </Popup>
    );
};

export default FinishPopup;

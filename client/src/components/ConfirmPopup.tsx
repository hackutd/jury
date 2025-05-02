import Button from './Button';
import Popup from './Popup';

interface PopupProps {
    /* State variable for open/closed */
    enabled: boolean;

    /* Function to modify the popup state variable */
    setEnabled: React.Dispatch<React.SetStateAction<boolean>>;

    /* Title Text */
    title: string;

    /* Submit Text */
    submitText: string;

    /* On submit function */
    onSubmit: () => void;

    /* React children, corresponds to the body content */
    children?: React.ReactNode;

    /* If true, submit button is red */
    red?: boolean;

    /* Disable the submit button */
    disabledSubmit?: boolean;
}

const ConfirmPopup = (props: PopupProps) => {
    if (!props.enabled) {
        return null;
    }

    return (
        <Popup enabled={props.enabled} setEnabled={props.setEnabled}>
            <h1 className="text-5xl font-bold mb-2 text-center">{props.title}</h1>
            <p className="text-xl">{props.children}</p>
            <div className="flex flex-row justify-around mt-4 space-x-4">
                <Button type="outline" full onClick={() => props.setEnabled(false)}>
                    Cancel
                </Button>
                <Button
                    type={props.red ? 'error' : 'primary'}
                    full
                    onClick={props.onSubmit}
                    disabled={props.disabledSubmit}
                >
                    {props.submitText}
                </Button>
            </div>
        </Popup>
    );
};

export default ConfirmPopup;

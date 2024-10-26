import Button from './Button';
import Popup from './Popup';

interface InfoPopupProps {
    /* State variable for open/closed */
    enabled: boolean;

    /* Function to modify the popup state variable */
    setEnabled: React.Dispatch<React.SetStateAction<boolean>>;

    /* Title Text */
    title: string;

    /* Submit Text */
    submitText: string;

    /* React children, corresponds to the body content */
    children?: React.ReactNode;

    /* If true, button is red */
    red?: boolean;
}

const InfoPopup = (props: InfoPopupProps) => {
    if (!props.enabled) {
        return null;
    }

    return (
        <Popup enabled={props.enabled} setEnabled={props.setEnabled}>
            <h1 className="text-5xl font-bold mb-2 text-center">{props.title}</h1>
            <p className="text-xl">{props.children}</p>
            <div className="flex justify-center">
                <Button
                    type={props.red ? 'error' : 'primary'}
                    onClick={() => props.setEnabled(false)}
                    flat
                    className="mt-4"
                >
                    {props.submitText}
                </Button>
            </div>
        </Popup>
    );
};

export default InfoPopup;

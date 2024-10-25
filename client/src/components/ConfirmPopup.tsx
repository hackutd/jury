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

    /* If true, button is red */
    red?: boolean;
}

const ConfirmPopup = (props: PopupProps) => {
    if (!props.enabled) {
        return null;
    }

    return (
        <Popup enabled={props.enabled} setEnabled={props.setEnabled}>
            <h1 className="text-5xl font-bold mb-2 text-center">{props.title}</h1>
            <p className="text-xl">{props.children}</p>
            <div className="flex flex-row justify-around">
                <button
                    className=" border-lightest border-2 rounded-full px-6 py-1 mt-4 w-2/5 font-bold text-2xl text-lighter hover:bg-lighter/30 duration-200"
                    onClick={() => props.setEnabled(false)}
                >
                    Cancel
                </button>
                <button
                    className={
                        'text-white rounded-full px-4 py-2 mt-4 w-2/5 font-bold text-2xl hover:brightness-110 duration-200 ' +
                        (props.red ? 'bg-error' : 'bg-primary')
                    }
                    onClick={props.onSubmit}
                >
                    {props.submitText}
                </button>
            </div>
        </Popup>
    );
};

export default ConfirmPopup;

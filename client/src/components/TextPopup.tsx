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

const TextPopup = (props: PopupProps) => {
    if (!props.enabled) {
        return null;
    }

    return (
        <>
            <div
                className="fixed left-0 top-0 z-20 w-screen h-screen bg-black/30"
                onClick={() => props.setEnabled(false)}
            ></div>
            <div className="bg-background fixed z-30 left-1/2 top-1/2 translate-x-[-50%] translate-y-[-50%] py-6 px-10 w-1/3">
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
            </div>
        </>
    );
};

export default TextPopup;

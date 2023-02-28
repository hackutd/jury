import { twMerge } from 'tailwind-merge';

const LoginBlock = (props: { disabled: boolean }) => {
    return (
        <div
            className={twMerge(
                'fixed top-0 left-0 w-screen h-screen bg-[#00000077] justify-center items-center',
                props.disabled ? 'hidden' : 'flex'
            )}
        >
            <div className="border-[16px] border-solid border-transparent border-t-primary rounded-full w-[10rem] h-[10rem] animate-spin"></div>
        </div>
    );
};

export default LoginBlock;

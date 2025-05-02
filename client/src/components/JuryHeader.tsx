import { useNavigate } from 'react-router-dom';
import { twMerge } from 'tailwind-merge';
import Cookies from 'universal-cookie';
import logoutButton from '../assets/logout.svg';

interface JuryHeaderProps {
    /* Whether to show the back to admin button */
    withBack?: boolean;

    /* Whether to show the logout button */
    withLogout?: boolean;

    /* Whether the user is an admin */
    isAdmin?: boolean;

    /* Custom back location */
    backLocation?: string;
}

const JuryHeader = (props: JuryHeaderProps) => {
    const navigate = useNavigate();
    const cookies = new Cookies();

    const logout = () => {
        const isJudge = window.location.pathname.startsWith('/judge');
        cookies.remove(isJudge ? 'token' : 'admin-pass', { path: '/' });
        navigate('/');
    };

    const backToAdmin = () => navigate(props.backLocation ?? '/admin');

    const adminCenter = props.isAdmin ? 'lg:text-center' : '';

    return (
        <div
            className={twMerge(
                'px-2 relative mx-auto pt-2 md:pt-6 w-full flex flex-col bg-background',
                props.isAdmin ? 'lg:items-center md:pt-2 lg:pt-6' : 'md:w-[30rem]'
            )}
        >
            <a
                href="/"
                className={twMerge(
                    'font-bold hover:text-primary duration-200 block max-w-fit z-10',
                    props.isAdmin ? 'lg:text-5xl text-4xl' : 'text-4xl',
                    adminCenter
                )}
            >
                {props.isAdmin ? 'Jury Admin' : 'Jury'}
            </a>
            <div
                className={twMerge(
                    'font-bold text-primary z-10',
                    props.isAdmin && 'lg:text-[1.5rem]',
                    adminCenter
                )}
            >
                {import.meta.env.VITE_JURY_NAME}
            </div>
            {props.withBack && (
                <div
                    className="lg:absolute lg:top-6 lg:left-6 mt-2 ml-2 flex items-center cursor-pointer border-none bg-transparent hover:scale-110 duration-200 text-light text-xl mr-2"
                    onClick={backToAdmin}
                >
                    ◂&nbsp;&nbsp;Back
                </div>
            )}
            {props.withLogout && (
                <div
                    className="absolute top-4 right-4 md:top-6 md:right-6 flex items-center cursor-pointer border-none bg-transparent hover:scale-110 duration-200 z-10"
                    onClick={logout}
                >
                    <div className="text-light lg:text-xl text-lg mr-2">Logout</div>
                    <img className="w-4 h-4" src={logoutButton} alt="logout icon" />
                </div>
            )}
        </div>
    );
};

export default JuryHeader;

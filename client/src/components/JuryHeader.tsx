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
}

const JuryHeader = (props: JuryHeaderProps) => {
    const navigate = useNavigate();
    const cookies = new Cookies();

    const logout = () => {
        const isJudge = window.location.pathname.startsWith('/judge');
        cookies.remove(isJudge ? 'token' : 'admin-pass', { path: '/' });
        navigate('/');
    };

    const backToAdmin = () => navigate('/admin');

    const adminCenter = props.isAdmin ? 'text-center' : '';

    return (
        <div
            className={twMerge(
                'md:px-2 px-4 relative mx-auto mt-6 w-full flex flex-col',
                props.isAdmin ? 'items-center' : 'md:w-[30rem]'
            )}
        >
            <a
                href="/"
                className={twMerge(
                    'font-bold hover:text-primary duration-200 block max-w-fit',
                    props.isAdmin ? 'text-5xl' : 'text-4xl',
                    adminCenter
                )}
            >
                {props.isAdmin ? 'Jury Admin' : 'Jury'}
            </a>
            <div
                className={twMerge(
                    'font-bold text-primary',
                    props.isAdmin && 'text-[1.5rem]',
                    adminCenter
                )}
            >
                {process.env.REACT_APP_JURY_NAME}
            </div>
            {props.withBack && (
                <div
                    className="absolute top-0 left-6 flex items-center cursor-pointer border-none bg-transparent hover:scale-110 duration-200 text-light text-xl mr-2"
                    onClick={backToAdmin}
                >
                    ◂&nbsp;&nbsp;Back
                </div>
            )}
            {props.withLogout && (
                <div
                    className="absolute top-0 right-6 flex items-center cursor-pointer border-none bg-transparent hover:scale-110 duration-200"
                    onClick={logout}
                >
                    <div className="text-light text-xl mr-2">Logout</div>
                    <img className="w-4 h-4" src={logoutButton} alt="logout icon" />
                </div>
            )}
        </div>
    );
};

export default JuryHeader;

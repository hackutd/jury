import { useNavigate } from 'react-router-dom';
import Cookies from 'universal-cookie';
import logoutButton from '../assets/logout.svg';

const JuryHeader = (props: { withLogout?: boolean }) => {
    const navigate = useNavigate();
    const cookies = new Cookies();

    const logout = () => {
        const isJudge = window.location.pathname.startsWith('/judge');
        cookies.remove(isJudge ? 'token' : 'admin-pass', { path: '/' });
        navigate('/');
    };

    return (
        <div className="px-2 relative md:w-[30rem] mx-auto mt-2">
            <div className="font-bold text-4xl">Jury</div>
            <div className="font-bold text-primary">{process.env.REACT_APP_JURY_NAME}</div>
            {props.withLogout && (
                <div
                    className="absolute top-4 right-4 flex items-center cursor-pointer border-none bg-transparent"
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

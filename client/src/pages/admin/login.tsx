import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import Cookies from 'universal-cookie';
import Button from '../../components/Button';
import Container from '../../components/Container';
import JuryHeader from '../../components/JuryHeader';
import Loading from '../../components/Loading';
import PasswordInput from '../../components/PasswordInput';
import { postRequest } from '../../api';
import { errorAlert } from '../../util';

const AdminLogin = () => {
    const [password, setPassword] = useState('');
    const [loginLock, setLoginLock] = useState(false);
    const [error, setError] = useState(false);
    const cookies = new Cookies();
    const navigate = useNavigate();

    // If token cookie is already defined and valid, redirect to admin page
    useEffect(() => {
        async function checkLoggedIn() {
            const cookies = new Cookies();
            if (cookies.get('admin-pass')) {
                const res = await postRequest<OkResponse>('/admin/auth', 'admin', null);
                if (res.status === 200 && res.data?.ok === 1) {
                    navigate('/admin');
                    return;
                }
            }

            // If invalid, delete the token
            cookies.remove('admin-pass');
        }

        checkLoggedIn();
    }, [navigate]);

    const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        setPassword(e.target.value);
    };

    const handleEnter = (e: React.KeyboardEvent<HTMLInputElement>) => {
        if (e.key === 'Enter') {
            login();
        }
    };

    const login = async () => {
        // Prevent spamming of login while waiting for fetch
        if (loginLock) return;
        setLoginLock(true);

        // Make async call to check code
        const res = await postRequest<OkResponse>('/admin/login', 'admin', { password });

        // Invalid code
        if (res.status === 400) {
            setError(true);
            setLoginLock(false);
            return;
        }

        // Internal server error
        if (res.status !== 200) {
            errorAlert(res.status);

            setError(true);
            setLoginLock(false);
            return;
        }
        
        // Wrong password
        if (res.data?.ok !== 1) {
            setError(true);
            setLoginLock(false);
            return;
        }

        // Correct code; save password as cookie
        cookies.set('admin-pass', password, {
            path: '/',
            sameSite: 'strict',
            secure: true,
            maxAge: 60 * 60 * 24,
        });

        // Redirect
        navigate('/admin');

        setLoginLock(false);
    };

    return (
        <>
            <JuryHeader />
            <Container>
                <PasswordInput
                    label="Enter the admin password"
                    placeholder="Admin password..."
                    onKeyPress={handleEnter}
                    onChange={handleChange}
                    error={error}
                    setError={setError}
                    errorMessage="Invalid admin password"
                    isHidden
                    className="w-4/5"
                />
                <div className="my-12" />
                <Button type="primary" onClick={login}>
                    Log In
                </Button>
            </Container>
            <Loading disabled={!loginLock} />
        </>
    );
};

export default AdminLogin;

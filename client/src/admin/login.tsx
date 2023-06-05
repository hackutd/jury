import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import Cookies from 'universal-cookie';
import Button from '../components/Button';
import Container from '../components/Container';
import JuryHeader from '../components/JuryHeader';
import LoginBlock from '../components/LoginBlock';
import PasswordInput from '../components/PasswordInput';

const AdminLogin = () => {
    const [password, setPassword] = useState('');
    const [loginLock, setLoginLock] = useState(false);
    const [error, setError] = useState(false);
    const cookies = new Cookies();
    const navigate = useNavigate();

    // If token cookie is already defined and valid, redirect to admin page
    useEffect(() => {
        const cookies = new Cookies();
        // TODO: Check if the cookie is valid lmao
        if (cookies.get('admin-pass')) {
            navigate('/admin')
            return;
        }

        // If invalid, delete the token
        cookies.remove('admin-pass');
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
        const res = await fetch(`${process.env.REACT_APP_JURY_URL}/admin/login`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ password }),
        });
        console.log(res);

        // Invalid code
        if (res.status === 400) {
            setError(true);
            setLoginLock(false);
            return;
        }

        // Internal server error
        if (res.status !== 200) {
            const err = await res.text();
            console.error(err);
            alert(err);

            setError(true);
            setLoginLock(false);
            return;
        }

        // Correct code; save token as cookie
        const token = await res.text();
        cookies.set('admin-pass', token, {
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
                    className='w-4/5'
                />
                <div className="my-12" />
                <Button type="primary" onClick={login}>
                    Log In
                </Button>
            </Container>
            <LoginBlock disabled={!loginLock} />
        </>
    );
};

export default AdminLogin;

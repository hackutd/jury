import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import Cookies from 'universal-cookie';
import Button from '../components/Button';
import Container from '../components/Container';
import JuryHeader from '../components/JuryHeader';
import PasswordInput from '../components/PasswordInput';
import LoginBlock from '../components/LoginBlock';

const JudgeLogin = () => {
    const [disabled, setDisabled] = useState(true);
    const [code, setCode] = useState('');
    const [loginLock, setLoginLock] = useState(false);
    const [error, setError] = useState(false);
    const navigate = useNavigate();

    // If token cookie is already defined and valid, redirect to judge page
    useEffect(() => {
        async function checkCookies() {
            const cookies = new Cookies();

            // If no token cookie, just ignore
            if (!cookies.get('token')) {
                return;
            }

            // If invalid, delete the token
            const res = await fetch(`${process.env.REACT_APP_JURY_URL}/judge/auth`, {
                method: 'POST',
                credentials: 'include',
            });
            if (!res.ok) {
                cookies.remove('token');
                return;
            }

            // If all valid, redirect to judge page
            navigate('/judge');
        }

        checkCookies();
    }, [navigate]);

    // Disable button if length of code is not 6
    useEffect(() => {
        setDisabled(code.length !== 6);
    }, [code]);

    const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        setCode(e.target.value);
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

        // Check for length of code
        if (code.length < 6) {
            setError(true);
            setLoginLock(false);
            return;
        }

        // Make async call to check code
        const res = await fetch(`${process.env.REACT_APP_JURY_URL}/judge/login`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ code }),
        });

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
        const cookies = new Cookies();
        cookies.set('token', token, {
            path: '/',
            sameSite: 'strict',
            secure: true,
            maxAge: 60 * 60 * 24,
        });

        // Redirect
        navigate('/judge');

        setLoginLock(false);
    };

    return (
        <>
            <JuryHeader />
            <Container>
                <PasswordInput
                    label="Enter your judging code"
                    maxLength={6}
                    placeholder="000000"
                    onKeyPress={handleEnter}
                    onChange={handleChange}
                    error={error}
                    setError={setError}
                    errorMessage="Invalid judging code"
                    className="text-8xl text-center w-4/5"
                />
                <p className="text-2xl text-light text-center mx-4 my-12">
                    If you did not get a code, check your email or contact an organizer.
                </p>
                <Button type="primary" disabled={disabled} onClick={login}>
                    Log In
                </Button>
            </Container>
            <LoginBlock disabled={!loginLock} />
        </>
    );
};

export default JudgeLogin;
